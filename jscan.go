package jscan

import (
	"fmt"
	"unsafe"
)

// ErrorCode defines the error type.
type ErrorCode int8

// All error codes
const (
	_ ErrorCode = iota
	ErrorCodeInvalidEscapeSeq
	ErrorCodeIllegalControlChar
	ErrorCodeUnexpectedEOF
	ErrorCodeUnexpectedToken
	ErrorCodeMalformedNumber
	ErrorCodeCallback
)

type Options struct {
	// When CachePath == true paths are generated on the fly and cached
	// reducing their performance penalty when read. If you don't require
	// paths then setting CachePath to false will improve performance.
	CachePath  bool
	EscapePath bool
}

type expectation int8

const (
	_ expectation = iota
	expectVal
	expectCommaOrObjTerm
	expectCommaOrArrTerm
	expectKeyOrObjTerm
	expectValOrArrTerm
	expectKey
)

// ValueType defines a JSON value type
type ValueType int

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
		if atIndex < 32 {
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
	if len(b) < 1 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}
