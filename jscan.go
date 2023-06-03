package jscan

import (
	"fmt"
	"strconv"
	"sync"
	"unicode/utf8"
	"unsafe"

	"github.com/romshark/jscan/internal/jsonnum"
	"github.com/romshark/jscan/internal/keyescape"
	"github.com/romshark/jscan/internal/strfind"
)

var itrPoolString = sync.Pool{
	New: func() any {
		return &Iterator[string]{
			stack: make([]stackNode, 0, 64),
		}
	},
}

var itrPoolBytes = sync.Pool{
	New: func() any {
		return &Iterator[[]byte]{
			stack: make([]stackNode, 0, 64),
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

func (i *Iterator[S]) callFn(fn func(i *Iterator[S]) (err bool)) (err bool) {
	i.arrayIndex = -1
	if len(i.stack) != 0 &&
		i.stack[len(i.stack)-1].Type == stackNodeTypeArray {
		i.arrayIndex = i.stack[len(i.stack)-1].ArrLen
		i.stack[len(i.stack)-1].ArrLen++
	}

	if fn(i) {
		return true
	}

	i.keyIndex = -1
	return false
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

// Parser wraps an iterator instance in a reusable instance.
// Using a parser instance is more efficient than global functions
// that rely on a global iterator pool.
type Parser[S ~string | ~[]byte] struct{ i *Iterator[S] }

// NewParser creates a new reusable parser instance.
func NewParser[S ~string | ~[]byte](stackCap int) *Parser[S] {
	i := &Iterator[S]{stack: make([]stackNode, stackCap)}
	reset(i)
	return &Parser[S]{i: i}
}

// ScanOne calls fn for every encountered value including objects and arrays.
// When an object or array is encountered fn will also be called for each of its
// field values and array items respectively.
// Unlike Scan, ScanOne doesn't return ErrorCodeUnexpectedToken when
// it encounters anything other than EOF after reading a valid JSON value.
// Returns s with the scanned value cut.
//
// WARNING: Don't use or alias *Iterator after fn returns!
func (p *Parser[S]) ScanOne(
	s S, fn func(*Iterator[S]) (err bool),
) (trailing S, err Error[S]) {
	reset(p.i)
	p.i.src = s
	return scan(p.i, fn)
}

// Scan calls fn for every encountered value including objects and arrays.
// When an object or array is encountered fn will also be called for each of its
// field values and array items respectively.
//
// WARNING: Don't use or alias *Iterator after fn returns!
func (p *Parser[S]) Scan(
	s S, fn func(*Iterator[S]) (err bool),
) Error[S] {
	reset(p.i)
	p.i.src = s

	trailing, err := scan(p.i, fn)
	if err.IsErr() {
		return err
	}
	s, illegalChar := strfind.EndOfWhitespaceSeq(trailing)
	if illegalChar {
		return getError(ErrorCodeIllegalControlChar, s, trailing)
	}
	if len(s) > 0 {
		return getError(ErrorCodeUnexpectedToken, s, trailing)
	}
	return Error[S]{}
}

// ScanOne calls fn for every encountered value including objects and arrays.
// When an object or array is encountered fn will also be called for each of its
// field values and array items respectively.
// Unlike Scan, ScanOne doesn't return ErrorCodeUnexpectedToken when
// it encounters anything other than EOF after reading a valid JSON value.
// Returns s with the scanned value cut.
// Unlike (*Parser).ScanOne this function will take an iterator instance
// from a global iterator pool and can therefore be less efficient.
// Consider reusing parser instances instead.
//
// WARNING: Don't use or alias *Iterator after fn returns!
func ScanOne[S ~string | ~[]byte](
	s S, fn func(*Iterator[S]) (err bool),
) (trailing S, err Error[S]) {
	var i *Iterator[S]
	switch any(s).(type) {
	case string:
		x := itrPoolString.Get()
		defer itrPoolString.Put(x)
		i = x.(*Iterator[S])
	case []byte:
		x := itrPoolBytes.Get()
		defer itrPoolBytes.Put(x)
		i = x.(*Iterator[S])
	}
	i.src = s
	reset(i)
	return scan(i, fn)
}

// Scan calls fn for every encountered value including objects and arrays.
// When an object or array is encountered fn will also be called for each of its
// field values and array items respectively.
// Unlike (*Parser).Scan this function will take an iterator instance
// from a global iterator pool and can therefore be less efficient.
// Consider reusing parser instances instead.
//
// WARNING: Don't use or alias *Iterator after fn returns!
func Scan[S ~string | ~[]byte](
	s S, fn func(*Iterator[S]) (err bool),
) (err Error[S]) {
	var i *Iterator[S]
	switch any(s).(type) {
	case string:
		x := itrPoolString.Get()
		defer itrPoolString.Put(x)
		i = x.(*Iterator[S])
	case []byte:
		x := itrPoolBytes.Get()
		defer itrPoolBytes.Put(x)
		i = x.(*Iterator[S])
	}
	i.src = s
	reset(i)
	trailing, err := scan(i, fn)
	if err.IsErr() {
		return err
	}
	s, illegalChar := strfind.EndOfWhitespaceSeq(trailing)
	if illegalChar {
		return getError(ErrorCodeIllegalControlChar, s, trailing)
	}
	if len(s) > 0 {
		return getError(ErrorCodeUnexpectedToken, s, trailing)
	}
	return Error[S]{}
}

// Valid returns true if s is a valid JSON value.
func Valid[S ~string | ~[]byte](s S) bool {
	return !Validate(s).IsErr()
}

// ValidateOne scans a JSON value from s and returns an error if it's invalid,
// otherwise returns s with the scanned value cut.
func ValidateOne[S ~string | ~[]byte](s S) (trailing S, err Error[S]) {
	return validate(s)
}

// Validate returns an error if s is invalid JSON.
func Validate[S ~string | ~[]byte](s S) Error[S] {
	trailing, err := validate(s)
	if err.IsErr() {
		return err
	}
	s, illegalChar := strfind.EndOfWhitespaceSeq(trailing)
	if illegalChar {
		return getError(ErrorCodeIllegalControlChar, s, trailing)
	}
	if len(s) > 0 {
		return getError(ErrorCodeUnexpectedToken, s, trailing)
	}
	return Error[S]{}
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

func unsafeB2S(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
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

// scan calls fn for every value encountered.
// Returns the remainder of i.src and an error if any is encountered.
func scan[S ~string | ~[]byte](
	i *Iterator[S], fn func(*Iterator[S]) (err bool),
) (S, Error[S]) {
	var (
		s      = i.src
		err    error
		b      bool
		ks, ke int
	)

VALUE:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, i.src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
		}
	}
	switch s[0] {
	case '{':
		goto VALUE_OBJECT
	case '[':
		goto VALUE_ARRAY
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		goto VALUE_NUMBER
	case '"':
		goto VALUE_STRING
	case 'n':
		goto VALUE_NULL
	case 'f':
		goto VALUE_FALSE
	case 't':
		goto VALUE_TRUE
	}
	if s[0] < 0x20 {
		return s, getError(ErrorCodeIllegalControlChar, i.src, s)
	}
	return s, getError(ErrorCodeUnexpectedToken, i.src, s)

VALUE_OBJECT:
	i.valueType = ValueTypeObject
	i.valueIndex, i.valueIndexEnd = len(i.src)-len(s), -1
	s = s[1:]
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, i.src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
		}
	}
	ks, ke = i.keyIndex, i.keyIndexEnd
	if i.callFn(fn) {
		return s, i.getError(ErrorCodeCallback)
	}
	if s[0] == '}' {
		s = s[1:]
		goto AFTER_VALUE
	}
	i.stack = append(i.stack, stackNode{
		Type:        stackNodeTypeObject,
		KeyIndex:    ks,
		KeyIndexEnd: ke,
	})
	goto OBJ_KEY

VALUE_ARRAY:
	i.valueType = ValueTypeArray
	i.valueIndex, i.valueIndexEnd = len(i.src)-len(s), -1
	s = s[1:]
	ks, ke = i.keyIndex, i.keyIndexEnd
	if i.callFn(fn) {
		return s, i.getError(ErrorCodeCallback)
	}
	i.stack = append(i.stack, stackNode{
		Type:        stackNodeTypeArray,
		KeyIndex:    ks,
		KeyIndexEnd: ke,
	})
	goto VALUE_OR_ARR_TERM

VALUE_NUMBER:
	{
		i.valueIndex = len(i.src) - len(s)
		if s, b = jsonnum.ReadNumber(s); b {
			return s, getError(ErrorCodeMalformedNumber, i.src, s)
		}
		i.valueIndexEnd = len(i.src) - len(s)
		i.valueType = ValueTypeNumber
		switch src := any(i.src).(type) {
		case string:
			_, err = strconv.ParseFloat(
				src[i.valueIndex:i.valueIndexEnd], 64,
			)
		case []byte:
			_, err = strconv.ParseFloat(
				unsafeB2S(src[i.valueIndex:i.valueIndexEnd]), 64,
			)
		}
		if err != nil {
			return s, getError(ErrorCodeMalformedNumber, i.src, s)
		}
		if i.callFn(fn) {
			return s, i.getError(ErrorCodeCallback)
		}
	}
	goto AFTER_VALUE

VALUE_STRING:
	s, b = s[1:], false
	goto STRING

VALUE_NULL:
	if len(s) < 4 || string(s[:4]) != "null" {
		return s, getError(ErrorCodeUnexpectedToken, i.src, s)
	}
	i.valueType = ValueTypeNull
	i.valueIndex = len(i.src) - len(s)
	i.valueIndexEnd = len(i.src) - len(s) + len("null")
	s = s[len("null"):]
	if i.callFn(fn) {
		return s, i.getError(ErrorCodeCallback)
	}
	goto AFTER_VALUE

VALUE_FALSE:
	if len(s) < 5 || string(s[:5]) != "false" {
		return s, getError(ErrorCodeUnexpectedToken, i.src, s)
	}
	i.valueType = ValueTypeFalse
	i.valueIndex = len(i.src) - len(s)
	i.valueIndexEnd = len(i.src) - len(s) + len("false")
	s = s[len("false"):]
	if i.callFn(fn) {
		return s, i.getError(ErrorCodeCallback)
	}
	goto AFTER_VALUE

VALUE_TRUE:
	if s := s; len(s) < 4 || string(s[:4]) != "true" {
		return s, getError(ErrorCodeUnexpectedToken, i.src, s)
	}
	i.valueType = ValueTypeTrue
	i.valueIndex = len(i.src) - len(s)
	i.valueIndexEnd = len(i.src) - len(s) + len("true")
	s = s[len("true"):]
	if i.callFn(fn) {
		return s, i.getError(ErrorCodeCallback)
	}
	goto AFTER_VALUE

OBJ_KEY:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, i.src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
		}
	}
	if s[0] != '"' {
		if s[0] < 0x20 {
			return s, getError(ErrorCodeIllegalControlChar, i.src, s)
		}
		return s, getError(ErrorCodeUnexpectedToken, i.src, s)
	}

	s = s[1:]

	b = true
	goto STRING
AFTER_OBJ_KEY_STRING:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, i.src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
		}
	}
	if s[0] != ':' {
		if s[0] < 0x20 {
			return s, getError(ErrorCodeIllegalControlChar, i.src, s)
		}
		return s, getError(ErrorCodeUnexpectedToken, i.src, s)
	}
	s = s[1:]
	goto VALUE

VALUE_OR_OBJ_TERM:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, i.src, s)
			}
			goto VALUE_OR_OBJ_TERM
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
		}
	}
	switch s[0] {
	case '}':
		s = s[1:]
		i.stack = i.stack[:len(i.stack)-1]
		goto AFTER_VALUE
	case '{':
		goto VALUE_OBJECT
	case '[':
		goto VALUE_ARRAY
	case '"':
		goto VALUE_STRING
	case 't':
		goto VALUE_TRUE
	case 'f':
		goto VALUE_FALSE
	case 'n':
		goto VALUE_NULL
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		goto VALUE_NUMBER
	}
	if s[0] < 0x20 {
		return s, getError(ErrorCodeIllegalControlChar, i.src, s)
	}
	return s, getError(ErrorCodeUnexpectedToken, i.src, s)

VALUE_OR_ARR_TERM:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, i.src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
		}
	}
	switch s[0] {
	case ']':
		s = s[1:]
		i.stack = i.stack[:len(i.stack)-1]
		goto AFTER_VALUE
	case '{':
		goto VALUE_OBJECT
	case '[':
		goto VALUE_ARRAY
	case '"':
		goto VALUE_STRING
	case 't':
		goto VALUE_TRUE
	case 'f':
		goto VALUE_FALSE
	case 'n':
		goto VALUE_NULL
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		goto VALUE_NUMBER
	}
	if s[0] < 0x20 {
		return s, getError(ErrorCodeIllegalControlChar, i.src, s)
	}
	return s, getError(ErrorCodeUnexpectedToken, i.src, s)

AFTER_VALUE:
	if len(i.stack) == 0 {
		return s, Error[S]{}
	}
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, i.src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
		}
	}
	switch s[0] {
	case ',':
		s = s[1:]
		if i.stack[len(i.stack)-1].Type == stackNodeTypeArray {
			goto VALUE
		}
		goto OBJ_KEY
	case '}':
		if i.stack[len(i.stack)-1].Type != stackNodeTypeObject {
			return s, getError(ErrorCodeUnexpectedToken, i.src, s)
		}
		s = s[1:]
		i.stack = i.stack[:len(i.stack)-1]
		i.keyIndex, i.keyIndexEnd = -1, -1
		goto AFTER_VALUE
	case ']':
		if i.stack[len(i.stack)-1].Type != stackNodeTypeArray {
			return s, getError(ErrorCodeUnexpectedToken, i.src, s)
		}
		s = s[1:]
		i.stack = i.stack[:len(i.stack)-1]
		i.keyIndex, i.keyIndexEnd = -1, -1
		goto AFTER_VALUE
	}
	return s, getError(ErrorCodeUnexpectedToken, i.src, s)

STRING:
	i.valueIndex = len(i.src) - len(s) - 1
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
			return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
		}
		switch s[0] {
		case '\\':
			if len(s) < 2 {
				return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
			}
			if lutEscape[s[1]] == 1 {
				s = s[2:]
				continue
			}
			if s[1] != 'u' {
				return s, getError(ErrorCodeInvalidEscapeSeq, i.src, s)
			}
			if len(s) < 6 ||
				lutSX[s[5]] != 2 ||
				lutSX[s[4]] != 2 ||
				lutSX[s[3]] != 2 ||
				lutSX[s[2]] != 2 {
				return s, getError(ErrorCodeInvalidEscapeSeq, i.src, s)
			}
			s = s[5:]
		case '"':
			s = s[1:]
			if b {
				i.keyIndex, i.keyIndexEnd = i.valueIndex, len(i.src)-len(s)
				goto AFTER_OBJ_KEY_STRING
			}
			i.valueIndexEnd = len(i.src) - len(s)
			i.valueType = ValueTypeString
			if i.callFn(fn) {
				return s, i.getError(ErrorCodeCallback)
			}
			goto AFTER_VALUE
		default:
			if s[0] < 0x20 {
				return s, getError(ErrorCodeIllegalControlChar, i.src, s)
			}
			s = s[1:]
		}
	}
}

// validate returns the remainder of i.src and an error if any is encountered.
func validate[S ~string | ~[]byte](s S) (S, Error[S]) {
	var (
		src = s
		top stackNodeType
		err error
		b   bool
		st  = make([]stackNodeType, 0, 128)
	)

	stPop := func() {
		if len(st) < 1 {
			return
		}
		st = st[:len(st)-1]
	}
	stTop := func() {
		if len(st) < 1 {
			top = 0
			return
		}
		top = st[len(st)-1]
	}
	stPush := func(t stackNodeType) {
		st = append(st, t)
	}

VALUE:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
	}
	switch s[0] {
	case '{': // Object
		goto VALUE_OBJECT
	case '[': // Array
		goto VALUE_ARRAY
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		goto VALUE_NUMBER
	case '"': // String
		goto VALUE_STRING
	case 'n': // Null
		goto VALUE_NULL
	case 'f': // False
		goto VALUE_FALSE
	case 't': // True
		goto VALUE_TRUE
	}
	if s[0] < 0x20 {
		return s, getError(ErrorCodeIllegalControlChar, src, s)
	}
	return s, getError(ErrorCodeUnexpectedToken, src, s)

VALUE_OBJECT:
	s = s[1:]
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
	}
	if s[0] == '}' {
		s = s[1:]
		goto AFTER_VALUE
	}
	stPush(stackNodeTypeObject)
	goto OBJ_KEY

VALUE_ARRAY:
	stPush(stackNodeTypeArray)
	s = s[1:]
	goto VALUE_OR_ARR_TERM

VALUE_NUMBER:
	{
		before := s
		if s, b = jsonnum.ReadNumber(s); b {
			return s, getError(ErrorCodeMalformedNumber, src, s)
		}
		switch before := any(before).(type) {
		case string:
			_, err = strconv.ParseFloat(
				before[:len(before)-len(s)], 64,
			)
		case []byte:
			_, err = strconv.ParseFloat(
				unsafeB2S(before[:len(before)-len(s)]), 64,
			)
		}
		if err != nil {
			return s, getError(ErrorCodeMalformedNumber, src, s)
		}
	}
	goto AFTER_VALUE

VALUE_STRING:
	s, b = s[1:], false
	goto STRING

VALUE_NULL:
	if len(s) < 4 || string(s[:4]) != "null" {
		return s, getError(ErrorCodeUnexpectedToken, src, s)
	}
	s = s[len("null"):]
	goto AFTER_VALUE

VALUE_FALSE:
	if len(s) < 5 || string(s[:5]) != "false" {
		return s, getError(ErrorCodeUnexpectedToken, src, s)
	}
	s = s[len("false"):]
	goto AFTER_VALUE

VALUE_TRUE:
	if s := s; len(s) < 4 || string(s[:4]) != "true" {
		return s, getError(ErrorCodeUnexpectedToken, src, s)
	}
	s = s[len("true"):]
	goto AFTER_VALUE

OBJ_KEY:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
	}
	if s[0] != '"' {
		if s[0] < 0x20 {
			return s, getError(ErrorCodeIllegalControlChar, src, s)
		}
		return s, getError(ErrorCodeUnexpectedToken, src, s)
	}

	s = s[1:]

	b = true
	goto STRING
AFTER_OBJ_KEY_STRING:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
	}
	if s[0] != ':' {
		if s[0] < 0x20 {
			return s, getError(ErrorCodeIllegalControlChar, src, s)
		}
		return s, getError(ErrorCodeUnexpectedToken, src, s)
	}
	s = s[1:]
	goto VALUE

VALUE_OR_OBJ_TERM:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, src, s)
			}
			goto VALUE_OR_OBJ_TERM
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
	}
	switch s[0] {
	case '}':
		s = s[1:]
		stPop()
		goto AFTER_VALUE
	case '{':
		goto VALUE_OBJECT
	case '[':
		goto VALUE_ARRAY
	case '"':
		goto VALUE_STRING
	case 't':
		goto VALUE_TRUE
	case 'f':
		goto VALUE_FALSE
	case 'n':
		goto VALUE_NULL
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		goto VALUE_NUMBER
	}
	if s[0] < 0x20 {
		return s, getError(ErrorCodeIllegalControlChar, src, s)
	}
	return s, getError(ErrorCodeUnexpectedToken, src, s)

VALUE_OR_ARR_TERM:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
	}
	switch s[0] {
	case ']':
		s = s[1:]
		stPop()
		goto AFTER_VALUE
	case '{':
		goto VALUE_OBJECT
	case '[':
		goto VALUE_ARRAY
	case '"':
		goto VALUE_STRING
	case 't':
		goto VALUE_TRUE
	case 'f':
		goto VALUE_FALSE
	case 'n':
		goto VALUE_NULL
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		goto VALUE_NUMBER
	}
	if s[0] < 0x20 {
		return s, getError(ErrorCodeIllegalControlChar, src, s)
	}
	return s, getError(ErrorCodeUnexpectedToken, src, s)

AFTER_VALUE:
	stTop()
	if top == 0 {
		return s, Error[S]{}
	}
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
	}
	switch s[0] {
	case ',':
		s = s[1:]
		if top == stackNodeTypeArray {
			goto VALUE
		}
		goto OBJ_KEY
	case '}':
		if top != stackNodeTypeObject {
			return s, getError(ErrorCodeUnexpectedToken, src, s)
		}
		s = s[1:]
		stPop()
		goto AFTER_VALUE
	case ']':
		if top != stackNodeTypeArray {
			return s, getError(ErrorCodeUnexpectedToken, src, s)
		}
		s = s[1:]
		stPop()
		goto AFTER_VALUE
	}
	return s, getError(ErrorCodeUnexpectedToken, src, s)

STRING:
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
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
		switch s[0] {
		case '\\':
			if len(s) < 2 {
				return s, getError(ErrorCodeUnexpectedEOF, src, s)
			}
			if lutEscape[s[1]] == 1 {
				s = s[2:]
				continue
			}
			if s[1] != 'u' {
				return s, getError(ErrorCodeInvalidEscapeSeq, src, s)
			}
			if len(s) < 6 ||
				lutSX[s[5]] != 2 ||
				lutSX[s[4]] != 2 ||
				lutSX[s[3]] != 2 ||
				lutSX[s[2]] != 2 {
				return s, getError(ErrorCodeInvalidEscapeSeq, src, s)
			}
			s = s[5:]
		case '"':
			s = s[1:]
			if b {
				goto AFTER_OBJ_KEY_STRING
			}
			goto AFTER_VALUE
		default:
			if s[0] < 0x20 {
				return s, getError(ErrorCodeIllegalControlChar, src, s)
			}
			s = s[1:]
		}
	}
}
