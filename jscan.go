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
	DefaultStackSizeParser    = 64
	DefaultStackSizeValidator = 128
)

func newIterator[S ~string | ~[]byte]() *Iterator[S] {
	return &Iterator[S]{stack: make([]stackNode, 0, DefaultStackSizeParser)}
}

func newValidator[S ~string | ~[]byte]() *Validator[S] {
	return &Validator[S]{stack: make([]stackNodeType, 0, DefaultStackSizeValidator)}
}

var (
	iteratorPoolString  = sync.Pool{New: func() any { return newIterator[string]() }}
	iteratorPoolBytes   = sync.Pool{New: func() any { return newIterator[[]byte]() }}
	validatorPoolString = sync.Pool{New: func() any { return newValidator[string]() }}
	validatorPoolBytes  = sync.Pool{New: func() any { return newValidator[[]byte]() }}
)

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

// Valid returns true if s is a valid JSON value, otherwise returns false.
//
// Unlike (*Validator).Valid this function will take a validator instance
// from a global pool and can therefore be less efficient.
// Consider reusing a Validator instance instead.
func Valid[S ~string | ~[]byte](s S, o Options) bool {
	return !Validate(s, o).IsErr()
}

// ValidateOne scans one JSON value from s and returns an error if it's invalid
// and trailing as substring of s with the scanned value cut.
// In case of an error trailing will be a substring of s cut up until the index
// where the error was encountered.
//
// Unlike (*Validator).ValidateOne this function will take a validator instance
// from a global pool and can therefore be less efficient.
// Consider reusing a Validator instance instead.
//
// TIP: Explicitly cast s to string or []byte to use the global validator pools
// and avoid an unecessary validator allocation such as when dealing with
// json.RawMessage and similar types derived from string or []byte.
//
//	m := json.RawMessage(`1`)
//	jscan.ValidateOne([]byte(m), // Cast m to []byte to avoid allocation!
func ValidateOne[S ~string | ~[]byte](
	s S, o Options,
) (trailing S, err Error[S]) {
	var v *Validator[S]
	switch any(s).(type) {
	case string:
		x := validatorPoolString.Get()
		defer validatorPoolString.Put(x)
		v = x.(*Validator[S])
	case []byte:
		x := validatorPoolBytes.Get()
		defer validatorPoolBytes.Put(x)
		v = x.(*Validator[S])
	}
	v.stack = v.stack[:0]

	return validate(v.stack, s, o.DisableUTF8Validation)
}

// Validate returns an error if s is invalid JSON.
//
// Unlike (*Validator).Validate this function will take a validator instance
// from a global pool and can therefore be less efficient.
// Consider reusing a Validator instance instead.
//
// TIP: Explicitly cast s to string or []byte to use the global validator pools
// and avoid an unecessary validator allocation such as when dealing with
// json.RawMessage and similar types derived from string or []byte.
//
//	m := json.RawMessage(`1`)
//	jscan.Validate([]byte(m), // Cast m to []byte to avoid allocation!
func Validate[S ~string | ~[]byte](s S, o Options) Error[S] {
	var v *Validator[S]
	switch any(s).(type) {
	case string:
		x := validatorPoolString.Get()
		defer validatorPoolString.Put(x)
		v = x.(*Validator[S])
	case []byte:
		x := validatorPoolBytes.Get()
		defer validatorPoolBytes.Put(x)
		v = x.(*Validator[S])
	}
	v.stack = v.stack[:0]

	t, err := validate(v.stack, s, o.DisableUTF8Validation)
	if err.IsErr() {
		return err
	}
	var illegalChar bool
	t, illegalChar = strfind.EndOfWhitespaceSeq(t)
	if illegalChar {
		return getError(ErrorCodeIllegalControlChar, s, t)
	}
	if len(t) > 0 {
		return getError(ErrorCodeUnexpectedToken, s, t)
	}
	return Error[S]{}
}

// Validator is a reusable validator instance.
// The validator is more efficient than the parser at JSON validation.
// A validator instance can be more efficient than global Valid, Validate and ValidateOne
// function calls due to potential stack frame allocation avoidance.
type Validator[S ~string | ~[]byte] struct{ stack []stackNodeType }

// Options specifies parser options. Use zero value for default settings.
type Options struct {
	// DisableUTF8Validation disables UTF-8 validation which improves performance
	// at the cost of RFC8259 compliance, see "8.1. Character Encoding"
	// (https://datatracker.ietf.org/doc/html/rfc8259#section-8.1).
	DisableUTF8Validation bool
}

// NewValidator creates a new reusable validator instance.
//
// preallocStackFrames determines how many stack frames will be preallocated.
// If 0, then DefaultStackSizeValidator is applied by default.
// A higher value implies greater memory usage but also reduces the chance of
// dynamic memory allocations if the JSON depth surpasses the stack size.
// 1024 is equivalent to ~1KiB of memory usage (1 frame = 1 byte).
func NewValidator[S ~string | ~[]byte](preallocStackFrames int) *Validator[S] {
	if preallocStackFrames == 0 {
		preallocStackFrames = DefaultStackSizeValidator
	}
	return &Validator[S]{stack: make([]stackNodeType, 0, preallocStackFrames)}
}

// Valid returns true if s is a valid JSON value, otherwise returns false.
func (v *Validator[S]) Valid(s S, o Options) bool {
	return !v.Validate(s, o).IsErr()
}

// ValidateOne scans one JSON value from s and returns an error if it's invalid
// and trailing as substring of s with the scanned value cut.
// In case of an error trailing will be a substring of s cut up until the index
// where the error was encountered.
func (v *Validator[S]) ValidateOne(s S, o Options) (trailing S, err Error[S]) {
	return validate(v.stack, s, o.DisableUTF8Validation)
}

// Validate returns an error if s is invalid JSON,
// otherwise returns a zero value of Error[S].
func (v *Validator[S]) Validate(s S, o Options) Error[S] {
	t, err := validate(v.stack, s, o.DisableUTF8Validation)
	if err.IsErr() {
		return err
	}
	var illegalChar bool
	t, illegalChar = strfind.EndOfWhitespaceSeq(t)
	if illegalChar {
		return getError(ErrorCodeIllegalControlChar, s, t)
	}
	if len(t) > 0 {
		return getError(ErrorCodeUnexpectedToken, s, t)
	}
	return Error[S]{}
}

// ScanOne calls fn for every encountered value including objects and arrays.
// When an object or array is encountered fn will also be called for each of its
// member and element values.
//
// Unlike Scan, ScanOne doesn't return ErrorCodeUnexpectedToken when
// it encounters anything other than EOF after reading a valid JSON value.
// Returns an error if any and trailing as substring of s with the scanned value cut.
// In case of an error trailing will be a substring of s cut up until the index
// where the error was encountered.
//
// Unlike (*Parser).ScanOne this function will take an iterator instance
// from a global iterator pool and can therefore be less efficient.
// Consider reusing a Parser instance instead.
//
// TIP: Explicitly cast s to string or []byte to use the global iterator pools
// and avoid an unecessary iterator allocation such as when dealing with
// json.RawMessage and similar types derived from string or []byte.
//
//	m := json.RawMessage(`1`)
//	jscan.ScanOne([]byte(m), // Cast m to []byte to avoid allocation!
//
// WARNING: Don't use or alias *Iterator[S] after fn returns!
func ScanOne[S ~string | ~[]byte](
	s S, o Options, fn func(*Iterator[S]) (err bool),
) (trailing S, err Error[S]) {
	var i *Iterator[S]
	switch any(s).(type) {
	case string:
		x := iteratorPoolString.Get()
		defer iteratorPoolString.Put(x)
		i = x.(*Iterator[S])
	case []byte:
		x := iteratorPoolBytes.Get()
		defer iteratorPoolBytes.Put(x)
		i = x.(*Iterator[S])
	}
	i.src = s
	reset(i)
	return scan(i, fn, o.DisableUTF8Validation)
}

// Scan calls fn for every encountered value including objects and arrays.
// When an object or array is encountered fn will also be called for each of its
// member and element values.
//
// Unlike (*Parser).Scan this function will take an iterator instance
// from a global iterator pool and can therefore be less efficient.
// Consider reusing a Parser instance instead.
//
// TIP: Explicitly cast s to string or []byte to use the global iterator pools
// and avoid an unecessary iterator allocation such as when dealing with
// json.RawMessage and similar types derived from string or []byte.
//
//	m := json.RawMessage(`1`)
//	jscan.Scan([]byte(m), // Cast m to []byte to avoid allocation!
//
// WARNING: Don't use or alias *Iterator[S] after fn returns!
func Scan[S ~string | ~[]byte](
	s S, o Options, fn func(*Iterator[S]) (err bool),
) (err Error[S]) {
	var i *Iterator[S]
	switch any(s).(type) {
	case string:
		x := iteratorPoolString.Get()
		defer iteratorPoolString.Put(x)
		i = x.(*Iterator[S])
	case []byte:
		x := iteratorPoolBytes.Get()
		defer iteratorPoolBytes.Put(x)
		i = x.(*Iterator[S])
	}
	i.src = s
	reset(i)
	t, err := scan(i, fn, o.DisableUTF8Validation)
	if err.IsErr() {
		return err
	}
	var illegalChar bool
	t, illegalChar = strfind.EndOfWhitespaceSeq(t)
	if illegalChar {
		return getError(ErrorCodeIllegalControlChar, s, t)
	}
	if len(t) > 0 {
		return getError(ErrorCodeUnexpectedToken, s, t)
	}
	return Error[S]{}
}

// Parser wraps an iterator in a reusable instance.
// Reusing a parser instance is more efficient than global functions
// that rely on a global iterator pool.
type Parser[S ~string | ~[]byte] struct{ i *Iterator[S] }

// NewParser creates a new reusable parser instance.
//
// preallocStackFrames determines how many stack frames will be preallocated.
// If 0, then DefaultStackSizeParser is applied by default.
// A higher value implies greater memory usage but also reduces the chance of
// dynamic memory allocations if the JSON depth surpasses the stack size.
// 32 is equivalent to ~1KiB of memory usage on 64-bit systems (1 frame = ~32 bytes).
func NewParser[S ~string | ~[]byte](preallocStackFrames int) *Parser[S] {
	if preallocStackFrames == 0 {
		preallocStackFrames = DefaultStackSizeValidator
	}
	i := &Iterator[S]{stack: make([]stackNode, preallocStackFrames)}
	reset(i)
	return &Parser[S]{i: i}
}

// ScanOne calls fn for every encountered value including objects and arrays.
// When an object or array is encountered fn will also be called for each of its
// member and element values.
//
// Unlike Scan, ScanOne doesn't return ErrorCodeUnexpectedToken when
// it encounters anything other than EOF after reading a valid JSON value.
// Returns an error if any and trailing as substring of s with the scanned value cut.
// In case of an error trailing will be a substring of s cut up until the index
// where the error was encountered.
//
// WARNING: Don't use or alias *Iterator[S] after fn returns!
func (p *Parser[S]) ScanOne(
	s S, o Options, fn func(*Iterator[S]) (err bool),
) (trailing S, err Error[S]) {
	reset(p.i)
	p.i.src = s
	return scan(p.i, fn, o.DisableUTF8Validation)
}

// Scan calls fn for every encountered value including objects and arrays.
// When an object or array is encountered fn will also be called for each of its
// member and element values.
//
// WARNING: Don't use or alias *Iterator[S] after fn returns!
func (p *Parser[S]) Scan(
	s S, o Options, fn func(*Iterator[S]) (err bool),
) Error[S] {
	reset(p.i)
	p.i.src = s

	t, err := scan(p.i, fn, o.DisableUTF8Validation)
	if err.IsErr() {
		return err
	}
	var illegalChar bool
	t, illegalChar = strfind.EndOfWhitespaceSeq(t)
	if illegalChar {
		return getError(ErrorCodeIllegalControlChar, s, t)
	}
	if len(t) > 0 {
		return getError(ErrorCodeUnexpectedToken, s, t)
	}
	return Error[S]{}
}

// Iterator provides access to the recently encountered value.
type Iterator[S ~string | ~[]byte] struct {
	stack   []stackNode
	src     S
	pointer []byte

	valueType             ValueType
	valueIndex            int
	valueIndexEnd         int
	keyIndex, keyIndexEnd int
	arrayIndex            int
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

	// ErrorCodeInvalidUTF8 indicates invalid UTF-8.
	ErrorCodeInvalidUTF8
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
	case ErrorCodeInvalidUTF8:
		errMsg = "invalid UTF-8"
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
