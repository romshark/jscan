package jscan

import (
	"fmt"
	"strconv"
	"sync"
	"unicode/utf8"

	"github.com/romshark/jscan/v2/internal/keyescape"
	"github.com/romshark/jscan/v2/internal/strfind"
)

// Default stack sizes
const (
	DefaultStackSizeIterator  = 64
	DefaultStackSizeValidator = 128
)

var iteratorPoolString = sync.Pool{
	New: func() any {
		return &Iterator[string]{
			stack: make([]stackNode, 0, DefaultStackSizeIterator),
		}
	},
}

var iteratorPoolBytes = sync.Pool{
	New: func() any {
		return &Iterator[[]byte]{
			stack: make([]stackNode, 0, DefaultStackSizeIterator),
		}
	},
}

var validatorPoolString = sync.Pool{
	New: func() any {
		return &Validator[string]{
			stack: make([]stackNodeType, 0, DefaultStackSizeValidator),
		}
	},
}

var validatorPoolBytes = sync.Pool{
	New: func() any {
		return &Validator[[]byte]{
			stack: make([]stackNodeType, 0, DefaultStackSizeValidator),
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

// ArrayLen returns the number of items in the current array,
// or -1 if the current value isn't a valid array value.
func (i *Iterator[S]) ArrayLen() (count int) {
	if i.valueType != ValueTypeArray {
		return -1
	}
	s := i.src[i.valueIndex+1:]
	if len(s) < 1 {
		return -1
	}
	switch s[0] {
	case ' ', '\t', '\r', '\n':
		var notOK bool
		s, notOK = strfind.EndOfWhitespaceSeq(s)
		if notOK {
			return -1
		}
	}
	if len(s) < 1 {
		return -1
	}
	if s[0] != ']' {
		count++
	}

	stack := 0

MAIN_LOOP:
	for {
		for ; len(s) > 15; s = s[16:] {
			if lutCount[s[0]] == 1 {
				goto CHECK_CHAR
			}
			if lutCount[s[1]] == 1 {
				s = s[1:]
				goto CHECK_CHAR
			}
			if lutCount[s[2]] == 1 {
				s = s[2:]
				goto CHECK_CHAR
			}
			if lutCount[s[3]] == 1 {
				s = s[3:]
				goto CHECK_CHAR
			}
			if lutCount[s[4]] == 1 {
				s = s[4:]
				goto CHECK_CHAR
			}
			if lutCount[s[5]] == 1 {
				s = s[5:]
				goto CHECK_CHAR
			}
			if lutCount[s[6]] == 1 {
				s = s[6:]
				goto CHECK_CHAR
			}
			if lutCount[s[7]] == 1 {
				s = s[7:]
				goto CHECK_CHAR
			}
			if lutCount[s[8]] == 1 {
				s = s[8:]
				goto CHECK_CHAR
			}
			if lutCount[s[9]] == 1 {
				s = s[9:]
				goto CHECK_CHAR
			}
			if lutCount[s[10]] == 1 {
				s = s[10:]
				goto CHECK_CHAR
			}
			if lutCount[s[11]] == 1 {
				s = s[11:]
				goto CHECK_CHAR
			}
			if lutCount[s[12]] == 1 {
				s = s[12:]
				goto CHECK_CHAR
			}
			if lutCount[s[13]] == 1 {
				s = s[13:]
				goto CHECK_CHAR
			}
			if lutCount[s[14]] == 1 {
				s = s[14:]
				goto CHECK_CHAR
			}
			if lutCount[s[15]] == 1 {
				s = s[15:]
				goto CHECK_CHAR
			}
			s = s[16:]
		}

	CHECK_CHAR:
		if len(s) < 1 {
			return -1
		}
		switch s[0] {
		case ',':
			if stack < 1 {
				count++
			}
		case ']', '}':
			if stack == 0 {
				return count
			}
			stack--
		case '[', '{':
			stack++
		case '"':
			goto SKIP_STRING
		}
		s = s[1:]
	}

SKIP_STRING:
	s = s[1:]
	for {
		for ; len(s) > 15; s = s[16:] {
			if lutStr[s[0]] != 0 {
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[1]] != 0 {
				s = s[1:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[2]] != 0 {
				s = s[2:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[3]] != 0 {
				s = s[3:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[4]] != 0 {
				s = s[4:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[5]] != 0 {
				s = s[5:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[6]] != 0 {
				s = s[6:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[7]] != 0 {
				s = s[7:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[8]] != 0 {
				s = s[8:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[9]] != 0 {
				s = s[9:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[10]] != 0 {
				s = s[10:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[11]] != 0 {
				s = s[11:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[12]] != 0 {
				s = s[12:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[13]] != 0 {
				s = s[13:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[14]] != 0 {
				s = s[14:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[15]] != 0 {
				s = s[15:]
				goto CHECK_STRING_CHARACTER
			}
			continue
		}

	CHECK_STRING_CHARACTER:
		if len(s) < 1 {
			return -1
		}
		switch s[0] {
		case '\\':
			if len(s) < 2 {
				return -1
			}
			if lutEscape[s[1]] == 1 {
				s = s[2:]
				continue
			}
			if s[1] != 'u' {
				return -1
			}
			if len(s) < 6 ||
				lutSX[s[5]] != 2 ||
				lutSX[s[4]] != 2 ||
				lutSX[s[3]] != 2 ||
				lutSX[s[2]] != 2 {
				return -1
			}
			s = s[5:]
		case '"':
			s = s[1:]
			goto MAIN_LOOP
		default:
			if s[0] < 0x20 {
				return -1
			}
			s = s[1:]
		}
	}
}

// Level returns the depth level of the current value.
//
// For example in the following JSON: `[1,2,3]` the array is situated at level 0
// while the integers inside are situated at level 1.
func (i *Iterator[S]) Level() int { return len(i.stack) }

// ArrayIndex returns either the index of the element value in the array
// or -1 if the value isn't inside an array.
func (i *Iterator[S]) ArrayIndex() int { return i.arrayIndex }

// ValueType returns the value type identifier.
func (i *Iterator[S]) ValueType() ValueType { return i.valueType }

// ValueIndex returns the start index of the value in the source.
func (i *Iterator[S]) ValueIndex() int { return i.valueIndex }

// ValueIndexEnd returns the end index of the value in the source if any.
// Object and array values have a -1 end index because their end is unknown
// during traversal.
func (i *Iterator[S]) ValueIndexEnd() int { return i.valueIndexEnd }

// KeyIndex returns either the start index of the member key string in the source
// or -1 when the value isn't a member of an object and hence doesn't have a key.
func (i *Iterator[S]) KeyIndex() int { return i.keyIndex }

// KeyIndexEnd returns either the end index of the member key string in the source
// or -1 when the value isn't a member of an object and hence doesn't have a key.
func (i *Iterator[S]) KeyIndexEnd() int { return i.keyIndexEnd }

// Key returns either the object member key or "" when the value
// isn't a member of an object and hence doesn't have a key.
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
// If keyIndex is != -1 then the element is a member value, otherwise
// arrayIndex indicates the index of the element in the underlying array.
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
// Consider using (*Iterator[S]).Pointer() instead for safety and convenience.
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

func (i *Iterator[S]) getError(c ErrorCode) Error[S] {
	return Error[S]{
		Code:  c,
		Src:   i.src,
		Index: i.valueIndex,
	}
}

// Error is a syntax error encountered during validation or iteration.
// The only exception is ErrorCodeCallback which indicates a callback
// explicitly breaking by returning true instead of a syntax error.
// (Error).IsErr() returning false is equivalent to err == nil.
type Error[S ~string | ~[]byte] struct {
	// Src refers to the original source.
	Src S

	// Index points to the error start index in the source.
	Index int

	// Code indicates the type of the error.
	Code ErrorCode
}

var _ error = Error[string]{}

// IsErr returns true if there is an error, otherwise returns false.
func (e Error[S]) IsErr() bool { return e.Code != 0 }

// Error stringifies the error implementing the built-in error interface.
// Calling Error should be avoided in performance-critical code as it
// relies on dynamic memory allocation.
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

	// ErrorCodeInvalidEscape indicates the encounter of an invalid escape sequence.
	ErrorCodeInvalidEscape

	// ErrorCodeIllegalControlChar indicates the encounter of
	// an illegal control character in the source.
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
		return fmt.Sprintf("error at index %d: unexpected EOF", index)
	case ErrorCodeInvalidEscape:
		errMsg = "invalid escape"
	case ErrorCodeIllegalControlChar:
		errMsg = "illegal control character"
	case ErrorCodeCallback:
		errMsg = "callback error"
	default:
		return ""
	}
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

// lutCount maps all characters interesting.
var lutCount = [256]byte{
	',': 1, '[': 1, ']': 1, '{': 1, '}': 1, '"': 1,
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
