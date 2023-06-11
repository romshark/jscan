package jscan

import (
	"fmt"
	"strconv"
	"sync"
	"unicode/utf8"

	"github.com/romshark/jscan/internal/keyescape"
)

const (
	defaultIteratorStackSize  = 64
	defaultValidatorStackSize = 128
)

var iteratorPoolString = sync.Pool{
	New: func() any {
		return &Iterator[string]{
			stack: make([]stackNode, 0, defaultIteratorStackSize),
		}
	},
}

var iteratorPoolBytes = sync.Pool{
	New: func() any {
		return &Iterator[[]byte]{
			stack: make([]stackNode, 0, defaultIteratorStackSize),
		}
	},
}

var validatorPoolString = sync.Pool{
	New: func() any {
		return &Validator[string]{
			stack: make([]stackNodeType, 0, defaultValidatorStackSize),
		}
	},
}

var validatorPoolBytes = sync.Pool{
	New: func() any {
		return &Validator[[]byte]{
			stack: make([]stackNodeType, 0, defaultValidatorStackSize),
		}
	},
}

type stackNodeType int8

const (
	_                   stackNodeType = iota
	stackNodeTypeObject               = 1
	stackNodeTypeArray                = 2
)

type stackNode struct {
	ArrLen                int
	KeyIndex, KeyIndexEnd int
	Type                  stackNodeType
}

// Iterator provides access to the recently encountered value.
type Iterator[S ~string | ~[]byte] struct {
	stack   []stackNode
	src     S
	pointer []byte

	valueType             ValueType
	valueIndex            int
	valueIndexEnd         int
	level                 int
	keyIndex, keyIndexEnd int
	arrayIndex            int
}

// Level returns the depth level of the current value
func (i *Iterator[S]) Level() int { return len(i.stack) }

// ArrayIndex returns either the index of the value item in the array
// or -1 if the value item isn't inside an array.
func (i *Iterator[S]) ArrayIndex() int { return i.arrayIndex }

// ValueType returns the value type identifier.
func (i *Iterator[S]) ValueType() ValueType { return i.valueType }

// ValueIndex returns the index of the value in the source.
func (i *Iterator[S]) ValueIndex() int { return i.valueIndex }

// ValueIndexEnd returns the end index of the value in the source if any.
// Object and array values have a -1 end index because their end is unknown
// during traversal.
func (i *Iterator[S]) ValueIndexEnd() int { return i.valueIndexEnd }

// KeyIndex returns either the index of the key of the value in the source
// or -1 when the value isn't inside an object and hence doesn't have a key.
func (i *Iterator[S]) KeyIndex() int { return i.keyIndex }

// KeyIndexEnd returns either the end index of the key of the value in the source
// or -1 when the value isn't inside an object and hence doesn't have a key.
func (i *Iterator[S]) KeyIndexEnd() int { return i.keyIndexEnd }

// Key returns the object field key if any.
func (i *Iterator[S]) Key() (key S) {
	if i.keyIndex == -1 {
		return
	}
	return i.src[i.keyIndex:i.keyIndexEnd]
}

// Value returns the value if any.
func (i *Iterator[S]) Value() (value S) {
	if i.valueIndexEnd == -1 {
		return
	}
	return i.src[i.valueIndex:i.valueIndexEnd]
}

// ScanStack calls fn for every element in the stack.
// If keyIndex is != -1 then the element is a field value, otherwise
// arrayIndex indicates the index of the item in the underlying array.
func (i *Iterator[S]) ScanStack(fn func(keyIndex, keyEnd, arrayIndex int)) {
	for j := range i.stack {
		if i.stack[j].KeyIndex > -1 {
			fn(i.stack[j].KeyIndex, i.stack[j].KeyIndexEnd, -1)
		}
		if i.stack[j].Type == stackNodeTypeArray {
			fn(-1, -1, i.stack[j].ArrLen-1)
		}
	}
}

// Pointer returns the JSON pointer in RFC-6901 format.
func (i *Iterator[S]) Pointer() (s S) {
	i.ViewPointer(func(p []byte) {
		switch any(s).(type) {
		case string:
			s = S(p)
		case []byte:
			b := make([]byte, len(p))
			copy(b, p)
			s = S(b)
		}
	})
	return
}

// ViewPointer calls fn and provides the buffer holding the
// JSON pointer in RFC-6901 format.
// Consider using (*Iter[S]).Pointer() instead for safety and convenience.
//
// WARNING: do not use or alias p after fn returns,
// only reading and copying p are considered safe!
func (i *Iterator[S]) ViewPointer(fn func(p []byte)) {
	i.ScanStack(func(keyIndex, keyEnd, arrayIndex int) {
		if keyIndex != -1 {
			// Object key
			i.pointer = append(i.pointer, '/')
			i.pointer = keyescape.Append(i.pointer, i.src[keyIndex+1:keyEnd-1])
			return
		}
		// Array index
		i.pointer = append(i.pointer, '/')
		i.pointer = strconv.AppendInt(i.pointer, int64(arrayIndex), 10)
	})
	if i.keyIndex != -1 {
		i.pointer = append(i.pointer, '/')
		i.pointer = keyescape.Append(i.pointer, i.src[i.keyIndex+1:i.keyIndexEnd-1])
	}
	fn(i.pointer)
	i.pointer = i.pointer[:0]
}

// getError returns the stringified error, if any.
func (i *Iterator[S]) getError(c ErrorCode) Error[S] {
	return Error[S]{
		Code:  c,
		Src:   i.src,
		Index: i.valueIndex,
	}
}

// Error represents an error encountered during validation or iteration.
type Error[S ~string | ~[]byte] struct {
	Src   S
	Index int
	Code  ErrorCode
}

// IsErr returns true if there is an error, otherwise returns false.
func (e Error[S]) IsErr() bool { return e.Code != 0 }

// Error stringifies the error.
// Calling Error should be avoided in performance-critical code
// as it relies on dynamic memory allocation.
func (e Error[S]) Error() string {
	if e.Index < len(e.Src) {
		var r rune
		switch x := any(e.Src).(type) {
		case string:
			r, _ = utf8.DecodeRuneInString(x[e.Index:])
		case []byte:
			r, _ = utf8.DecodeRune(x[e.Index:])
		}
		return errorMessage(e.Code, e.Index, r)
	}
	return errorMessage(e.Code, e.Index, 0)
}

func reset[S ~string | ~[]byte](i *Iterator[S]) {
	i.stack = i.stack[:0]
	i.pointer = i.pointer[:0]
	i.valueType = 0
	i.level = 0
	i.keyIndex, i.keyIndexEnd = -1, -1
	i.valueIndexEnd = -1
	i.arrayIndex = 0
}

// ErrorCode defines the error type.
type ErrorCode int8

const (
	_ ErrorCode = iota

	// ErrorCodeInvalidEscapeSeq indicates the encounter of
	// an invalid escape sequence.
	ErrorCodeInvalidEscapeSeq

	// ErrorCodeIllegalControlChar indicates the presence of
	// a control character in the source.
	ErrorCodeIllegalControlChar

	// ErrorCodeUnexpectedEOF indicates the encounter an unexpected end of file.
	ErrorCodeUnexpectedEOF

	// ErrorCodeUnexpectedToken indicates the encounter of an unexpected token.
	ErrorCodeUnexpectedToken

	// ErrorCodeMalformedNumber indicates the encounter of a malformed number.
	ErrorCodeMalformedNumber

	// ErrorCodeCallback indicates return of true from the callback function.
	ErrorCodeCallback
)

// ValueType defines a JSON value type
type ValueType int8

// JSON value types
const (
	_ ValueType = iota
	ValueTypeObject
	ValueTypeArray
	ValueTypeNull
	ValueTypeFalse
	ValueTypeTrue
	ValueTypeString
	ValueTypeNumber
)

func (t ValueType) String() string {
	switch t {
	case ValueTypeObject:
		return "object"
	case ValueTypeArray:
		return "array"
	case ValueTypeNull:
		return "null"
	case ValueTypeFalse:
		return "false"
	case ValueTypeTrue:
		return "true"
	case ValueTypeString:
		return "string"
	case ValueTypeNumber:
		return "number"
	}
	return ""
}

func errorMessage(c ErrorCode, index int, atIndex rune) string {
	errMsg := ""
	switch c {
	case ErrorCodeUnexpectedToken:
		errMsg = "unexpected token"
	case ErrorCodeMalformedNumber:
		errMsg = "malformed number"
	case ErrorCodeUnexpectedEOF:
		if atIndex == 0 {
			return fmt.Sprintf(
				"error at index %d: unexpected EOF",
				index,
			)
		}
		errMsg = "unexpected EOF"
	case ErrorCodeInvalidEscapeSeq:
		errMsg = "invalid escape sequence"
	case ErrorCodeIllegalControlChar:
		errMsg = "illegal control character"
	case ErrorCodeCallback:
		errMsg = "callback error"
	default:
		return ""
	}
	if index >= 0 {
		if atIndex < 0x20 {
			return fmt.Sprintf(
				"error at index %d (0x%x): %s",
				index, atIndex, errMsg,
			)
		}
		return fmt.Sprintf(
			"error at index %d ('%s'): %s",
			index, string(atIndex), errMsg,
		)
	}
	return fmt.Sprintf(
		"error at index %d: %s",
		index, errMsg,
	)
}

// lutSX maps space characters such as whitespace, tab, line-break and
// carriage-return to 1, valid hex digits to 2 and others to 0.
var lutSX = [256]byte{
	' ': 1, '\n': 1, '\t': 1, '\r': 1,

	'0': 2, '1': 2, '2': 2, '3': 2, '4': 2, '5': 2, '6': 2, '7': 2, '8': 2, '9': 2,
	'a': 2, 'b': 2, 'c': 2, 'd': 2, 'e': 2, 'f': 2,
	'A': 2, 'B': 2, 'C': 2, 'D': 2, 'E': 2, 'F': 2,
}

// lutStr maps 0 to all bytes that don't require checking during string traversal.
// 1 is mapped to control, quotation mark (") and reverse solidus ("\").
var lutStr = [256]byte{
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	'"': 1, '\\': 1,
}

// lutEscape maps escapable characters to 1,
// all other ASCII characters are mapped to 0.
var lutEscape = [256]byte{
	'"':  1,
	'\\': 1,
	'/':  1,
	'b':  1,
	'f':  1,
	'n':  1,
	'r':  1,
	't':  1,
}

// getError returns the stringified error, if any.
func getError[S ~string | ~[]byte](c ErrorCode, src S, s S) Error[S] {
	return Error[S]{
		Code:  c,
		Src:   src,
		Index: len(src) - len(s),
	}
}
