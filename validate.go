package jscan

import (
	"github.com/romshark/jscan/v2/internal/jsonnum"
	"github.com/romshark/jscan/v2/internal/strfind"
)

// Valid returns true if s is a valid JSON value, otherwise returns false.
//
// Unlike [Validator.Valid] this function will take a validator instance
// from a global pool and can therefore be less efficient.
// Consider reusing a [Validator] instance instead.
func Valid[S string | []byte](s S) bool {
	return !Validate(s).IsErr()
}

// ValidateOne scans one JSON value from s and returns an error if it's invalid
// and trailing as substring of s with the scanned value cut.
// In case of an error trailing will be a substring of s cut up until the index
// where the error was encountered.
//
// Unlike [Validator.ValidateOne] this function will take a validator instance
// from a global pool and can therefore be less efficient.
// Consider reusing a [Validator] instance instead.
//
// NOTE: Types derived from string or []byte such as [encoding/json.RawMessage]
// must be converted explicitly.
//
//	m := json.RawMessage(`1`)
//	jscan.ValidateOne([]byte(m), // Convert m to []byte.
func ValidateOne[S string | []byte](s S) (trailing S, err Error[S]) {
	st := validatorStackPool.Get().(*[]stackNodeType)
	defer validatorStackPool.Put(st)

	t, e := validate((*st)[:0], toStr(s))
	return fromStr[S](t), Error[S]{Src: s, Index: e.Index, Code: e.Code}
}

// Validate returns an error if s is invalid JSON.
//
// Unlike [Validator.Validate] this function will take a validator instance
// from a global pool and can therefore be less efficient.
// Consider reusing a [Validator] instance instead.
//
// NOTE: Types derived from string or []byte such as [encoding/json.RawMessage]
// must be converted explicitly.
//
//	m := json.RawMessage(`1`)
//	jscan.Validate([]byte(m), // Convert m to []byte.
func Validate[S string | []byte](s S) Error[S] {
	st := validatorStackPool.Get().(*[]stackNodeType)
	defer validatorStackPool.Put(st)

	return validateAll((*st)[:0], s)
}

// NewValidator creates a new reusable validator instance.
// A higher preallocStackFrames value implies greater memory usage but also reduces
// the chance of dynamic memory allocations if the JSON depth surpasses the stack size.
// preallocStackFrames of 1024 is equivalent to ~1KiB of memory usage (1 frame = 1 byte).
// Use [DefaultStackSizeValidator] when not sure.
func NewValidator[S string | []byte](preallocStackFrames int) *Validator[S] {
	return &Validator[S]{
		stack: make([]stackNodeType, 0, preallocStackFrames),
	}
}

// Validator is a reusable validator instance.
// The validator is more efficient than the scanner at JSON validation.
// A validator instance can be more efficient than global [Valid], [Validate] and
// [ValidateOne] function calls because it avoids the global stack pool.
type Validator[S string | []byte] struct{ stack []stackNodeType }

// Valid returns true if s is a valid JSON value, otherwise returns false.
func (v *Validator[S]) Valid(s S) bool {
	return !v.Validate(s).IsErr()
}

// ValidateOne scans one JSON value from s and returns an error if it's invalid
// and trailing as substring of s with the scanned value cut.
// In case of an error trailing will be a substring of s cut up until the index
// where the error was encountered.
func (v *Validator[S]) ValidateOne(s S) (trailing S, err Error[S]) {
	t, e := validate(v.stack, toStr(s))
	return fromStr[S](t), Error[S]{Src: s, Index: e.Index, Code: e.Code}
}

// Validate returns an error if s is invalid JSON,
// otherwise returns a zero value of [Error].
func (v *Validator[S]) Validate(s S) Error[S] {
	return validateAll(v.stack, s)
}

// validateAll validates s expecting it to contain
// exactly one JSON value and nothing but whitespace after it.
func validateAll[S string | []byte](st []stackNodeType, s S) Error[S] {
	src := toStr(s)
	t, e := validate(st, src)
	if e.IsErr() {
		return Error[S]{Src: s, Index: e.Index, Code: e.Code}
	}
	var illegalChar bool
	t, illegalChar = strfind.EndOfWhitespaceSeq(t)
	if illegalChar {
		return Error[S]{
			Src:   s,
			Index: len(src) - len(t),
			Code:  ErrorCodeIllegalControlChar,
		}
	}
	if len(t) > 0 {
		return Error[S]{
			Src:   s,
			Index: len(src) - len(t),
			Code:  ErrorCodeUnexpectedToken,
		}
	}
	return Error[S]{}
}

// validate returns the remainder of s and an error if any is encountered.
func validate(st []stackNodeType, s string) (string, srcErr) {
	var (
		rollback string // Used as fallback for error report
		src      = s
		top      stackNodeType
		b        bool
	)

	stPop := func() { st = st[:len(st)-1] }
	stTop := func() {
		if len(st) < 1 {
			top = 0
			return
		}
		top = st[len(st)-1]
	}
	stPush := func(t stackNodeType) { st = append(st, t) }

VALUE:
	if len(s) < 1 {
		return s, errAt(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, errAt(ErrorCodeIllegalControlChar, src, s)
			}
		}
		if len(s) < 1 {
			return s, errAt(ErrorCodeUnexpectedEOF, src, s)
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
		return s, errAt(ErrorCodeIllegalControlChar, src, s)
	}
	return s, errAt(ErrorCodeUnexpectedToken, src, s)

VALUE_OBJECT:
	s = s[1:]
	if len(s) < 1 {
		return s, errAt(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, errAt(ErrorCodeIllegalControlChar, src, s)
			}
		}
		if len(s) < 1 {
			return s, errAt(ErrorCodeUnexpectedEOF, src, s)
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
		rollback = s
		var rc jsonnum.ReturnCode
		if s, rc = jsonnum.ReadNumber(s); rc == jsonnum.ReturnCodeErr {
			return s, errAt(ErrorCodeMalformedNumber, src, rollback)
		}
	}
	goto AFTER_VALUE

VALUE_STRING:
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
			return s, errAt(ErrorCodeUnexpectedEOF, src, s)
		}
		switch s[0] {
		case '\\':
			if len(s) < 2 {
				s = s[1:]
				return s, errAt(ErrorCodeUnexpectedEOF, src, s)
			}
			if lutEscape[s[1]] == 1 {
				s = s[2:]
				continue
			}
			if s[1] != 'u' {
				return s, errAt(ErrorCodeInvalidEscape, src, s)
			}
			if len(s) < 6 ||
				lutSX[s[5]] != 2 ||
				lutSX[s[4]] != 2 ||
				lutSX[s[3]] != 2 ||
				lutSX[s[2]] != 2 {
				return s, errAt(ErrorCodeInvalidEscape, src, s)
			}
			s = s[5:]
		case '"':
			s = s[1:]
			goto AFTER_VALUE
		default:
			if s[0] < 0x20 {
				return s, errAt(ErrorCodeIllegalControlChar, src, s)
			}
			s = s[1:]
		}
	}

VALUE_NULL:
	if len(s) < 4 || string(s[:4]) != "null" {
		return s, errAt(ErrorCodeUnexpectedToken, src, s)
	}
	s = s[len("null"):]
	goto AFTER_VALUE

VALUE_FALSE:
	if len(s) < 5 || string(s[:5]) != "false" {
		return s, errAt(ErrorCodeUnexpectedToken, src, s)
	}
	s = s[len("false"):]
	goto AFTER_VALUE

VALUE_TRUE:
	if s := s; len(s) < 4 || string(s[:4]) != "true" {
		return s, errAt(ErrorCodeUnexpectedToken, src, s)
	}
	s = s[len("true"):]
	goto AFTER_VALUE

OBJ_KEY:
	if len(s) < 1 {
		return s, errAt(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, errAt(ErrorCodeIllegalControlChar, src, s)
			}
		}
		if len(s) < 1 {
			return s, errAt(ErrorCodeUnexpectedEOF, src, s)
		}
	}
	if s[0] != '"' {
		if s[0] < 0x20 {
			return s, errAt(ErrorCodeIllegalControlChar, src, s)
		}
		return s, errAt(ErrorCodeUnexpectedToken, src, s)
	}

	s = s[1:]
	for {
		for ; len(s) > 15; s = s[16:] {
			if lutStr[s[0]] != 0 {
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[1]] != 0 {
				s = s[1:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[2]] != 0 {
				s = s[2:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[3]] != 0 {
				s = s[3:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[4]] != 0 {
				s = s[4:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[5]] != 0 {
				s = s[5:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[6]] != 0 {
				s = s[6:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[7]] != 0 {
				s = s[7:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[8]] != 0 {
				s = s[8:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[9]] != 0 {
				s = s[9:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[10]] != 0 {
				s = s[10:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[11]] != 0 {
				s = s[11:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[12]] != 0 {
				s = s[12:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[13]] != 0 {
				s = s[13:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[14]] != 0 {
				s = s[14:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[15]] != 0 {
				s = s[15:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			continue
		}

	CHECK_FIELDNAME_STRING_CHARACTER:
		if len(s) < 1 {
			return s, errAt(ErrorCodeUnexpectedEOF, src, s)
		}
		switch s[0] {
		case '\\':
			if len(s) < 2 {
				s = s[1:]
				return s, errAt(ErrorCodeUnexpectedEOF, src, s)
			}
			if lutEscape[s[1]] == 1 {
				s = s[2:]
				continue
			}
			if s[1] != 'u' {
				return s, errAt(ErrorCodeInvalidEscape, src, s)
			}
			if len(s) < 6 ||
				lutSX[s[5]] != 2 ||
				lutSX[s[4]] != 2 ||
				lutSX[s[3]] != 2 ||
				lutSX[s[2]] != 2 {
				return s, errAt(ErrorCodeInvalidEscape, src, s)
			}
			s = s[5:]
		case '"':
			s = s[1:]
			goto AFTER_OBJ_KEY_STRING
		default:
			if s[0] < 0x20 {
				return s, errAt(ErrorCodeIllegalControlChar, src, s)
			}
			s = s[1:]
		}
	}
AFTER_OBJ_KEY_STRING:
	if len(s) < 1 {
		return s, errAt(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, errAt(ErrorCodeIllegalControlChar, src, s)
			}
		}
		if len(s) < 1 {
			return s, errAt(ErrorCodeUnexpectedEOF, src, s)
		}
	}
	if s[0] != ':' {
		if s[0] < 0x20 {
			return s, errAt(ErrorCodeIllegalControlChar, src, s)
		}
		return s, errAt(ErrorCodeUnexpectedToken, src, s)
	}
	s = s[1:]
	goto VALUE

VALUE_OR_ARR_TERM:
	if len(s) < 1 {
		return s, errAt(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, errAt(ErrorCodeIllegalControlChar, src, s)
			}
		}
		if len(s) < 1 {
			return s, errAt(ErrorCodeUnexpectedEOF, src, s)
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
		return s, errAt(ErrorCodeIllegalControlChar, src, s)
	}
	return s, errAt(ErrorCodeUnexpectedToken, src, s)

AFTER_VALUE:
	stTop()
	if top == 0 {
		return s, srcErr{}
	}
	if len(s) < 1 {
		return s, errAt(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, errAt(ErrorCodeIllegalControlChar, src, s)
			}
		}
		if len(s) < 1 {
			return s, errAt(ErrorCodeUnexpectedEOF, src, s)
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
			return s, errAt(ErrorCodeUnexpectedToken, src, s)
		}
		s = s[1:]
		stPop()
		goto AFTER_VALUE
	case ']':
		if top != stackNodeTypeArray {
			return s, errAt(ErrorCodeUnexpectedToken, src, s)
		}
		s = s[1:]
		stPop()
		goto AFTER_VALUE
	}
	if s[0] < 0x20 {
		return s, errAt(ErrorCodeIllegalControlChar, src, s)
	}
	return s, errAt(ErrorCodeUnexpectedToken, src, s)
}
