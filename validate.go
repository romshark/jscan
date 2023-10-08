package jscan

import (
	"github.com/romshark/jscan/v2/internal/jsonnum"
	"github.com/romshark/jscan/v2/internal/strfind"
)

// Valid returns true if s is a valid JSON value, otherwise returns false.
//
// Unlike (*Validator).Valid this function will take a validator instance
// from a global pool and can therefore be less efficient.
// Consider reusing a Validator instance instead.
func Valid[S ~string | ~[]byte](s S) bool {
	return !Validate(s).IsErr()
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
func ValidateOne[S ~string | ~[]byte](s S) (trailing S, err Error[S]) {
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
	default:
		v = newValidator[S]()
	}
	v.stack = v.stack[:0]

	return validate(v.stack, s, false)
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
func Validate[S ~string | ~[]byte](s S) Error[S] {
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
	default:
		v = newValidator[S]()
	}
	v.stack = v.stack[:0]

	t, err := validate(v.stack, s, false)
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
type Validator[S ~string | ~[]byte] struct {
	stack                 []stackNodeType
	disableUTF8Validation bool
}

// ValidatorOptions configures a validator instance.
type ValidatorOptions struct {
	PreallocStackFrames int

	// DisableUTF8Validation disables UTF-8 validation which significantly
	// improves performance at the cost of RFC8259 compliance which states that
	// JSON strings must not contain illegal UTF-8 sequences.
	DisableUTF8Validation bool
}

// NewValidator creates a new reusable validator instance.
// A higher preallocStackFrames value implies greater memory usage but also reduces
// the chance of dynamic memory allocations if the JSON depth surpasses the stack size.
// preallocStackFrames of 1024 is equivalent to ~1KiB of memory usage (1 frame = 1 byte).
// Use DefaultStackSizeValidator when not sure.
func NewValidator[S ~string | ~[]byte](o ValidatorOptions) *Validator[S] {
	if o.PreallocStackFrames == 0 {
		o.PreallocStackFrames = DefaultStackSizeValidator
	}
	return &Validator[S]{
		stack:                 make([]stackNodeType, 0, o.PreallocStackFrames),
		disableUTF8Validation: o.DisableUTF8Validation,
	}
}

// Valid returns true if s is a valid JSON value, otherwise returns false.
func (v *Validator[S]) Valid(s S) bool {
	return !v.Validate(s).IsErr()
}

// ValidateOne scans one JSON value from s and returns an error if it's invalid
// and trailing as substring of s with the scanned value cut.
// In case of an error trailing will be a substring of s cut up until the index
// where the error was encountered.
func (v *Validator[S]) ValidateOne(s S) (trailing S, err Error[S]) {
	return validate(v.stack, s, v.disableUTF8Validation)
}

// Validate returns an error if s is invalid JSON,
// otherwise returns a zero value of Error[S].
func (v *Validator[S]) Validate(s S) Error[S] {
	t, err := validate(v.stack, s, v.disableUTF8Validation)
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

// validate returns the remainder of i.src and an error if any is encountered.
func validate[S ~string | ~[]byte](st []stackNodeType, s S, noUTF8validation bool) (S, Error[S]) {
	var (
		rollback S // Used as fallback for error report
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
		rollback = s
		if s, b = jsonnum.ReadNumber(s); b {
			return s, getError(ErrorCodeMalformedNumber, src, rollback)
		}
	}
	goto AFTER_VALUE

VALUE_STRING:
	s = s[1:]
	{
		ss := s
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
					s = s[1:]
					return s, getError(ErrorCodeUnexpectedEOF, src, s)
				}
				if lutEscape[s[1]] == 1 {
					s = s[2:]
					continue
				}
				if s[1] != 'u' {
					return s, getError(ErrorCodeInvalidEscape, src, s)
				}
				if len(s) < 6 ||
					lutSX[s[5]] != 2 ||
					lutSX[s[4]] != 2 ||
					lutSX[s[3]] != 2 ||
					lutSX[s[2]] != 2 {
					return s, getError(ErrorCodeInvalidEscape, src, s)
				}
				s = s[5:]
			case '"':
				s = s[1:]

				if noUTF8validation {
					goto AFTER_VALUE
				}

				// The UTF-8 verification code was borrowed from utf8.ValidString
				// https://cs.opensource.google/go/go/+/refs/tags/go1.21.2:src/unicode/utf8/utf8.go;l=528
				{
					sv := ss[:len(ss)-len(s)]
					// Fast path. Check for and skip 8 bytes of
					// ASCII characters per iteration.
					for len(sv) >= 8 {
						// Combining two 32 bit loads allows the same code to be used for 32 and 64 bit platforms.
						// The compiler can generate a 32bit load for first32 and second32 on many platforms.
						// See test/codegen/memcombine.go.
						first32 := uint32(sv[0]) |
							uint32(sv[1])<<8 |
							uint32(sv[2])<<16 |
							uint32(sv[3])<<24
						second32 := uint32(sv[4]) |
							uint32(sv[5])<<8 |
							uint32(sv[6])<<16 |
							uint32(sv[7])<<24
						if (first32|second32)&0x80808080 != 0 {
							// Found a non ASCII byte (>= RuneSelf).
							break
						}
						sv = sv[8:]
					}
					n := len(sv)
					for j := 0; j < n; {
						si := sv[j]
						if si < utf8RuneSelf {
							j++
							continue
						}
						x := utf8First[si]
						if x == utf8xx {
							// Illegal starter byte.
							return s, getError(ErrorCodeInvalidUTF8, src, s)
						}
						size := int(x & 7)
						if j+size > n {
							// Short or invalid.
							return s, getError(ErrorCodeInvalidUTF8, src, s)
						}
						accept := utf8AcceptRanges[x>>4]
						if c := sv[j+1]; c < accept.lo || accept.hi < c {
							return s, getError(ErrorCodeInvalidUTF8, src, s)
						} else if size == 2 {
						} else if c := sv[j+2]; c < utf8locb || utf8hicb < c {
							return s, getError(ErrorCodeInvalidUTF8, src, s)
						} else if size == 3 {
						} else if c := sv[j+3]; c < utf8locb || utf8hicb < c {
							return s, getError(ErrorCodeInvalidUTF8, src, s)
						}
						j += size
					}
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
	{
		ss := s
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
				return s, getError(ErrorCodeUnexpectedEOF, src, s)
			}
			switch s[0] {
			case '\\':
				if len(s) < 2 {
					s = s[1:]
					return s, getError(ErrorCodeUnexpectedEOF, src, s)
				}
				if lutEscape[s[1]] == 1 {
					s = s[2:]
					continue
				}
				if s[1] != 'u' {
					return s, getError(ErrorCodeInvalidEscape, src, s)
				}
				if len(s) < 6 ||
					lutSX[s[5]] != 2 ||
					lutSX[s[4]] != 2 ||
					lutSX[s[3]] != 2 ||
					lutSX[s[2]] != 2 {
					return s, getError(ErrorCodeInvalidEscape, src, s)
				}
				s = s[5:]
			case '"':
				s = s[1:]

				if noUTF8validation {
					goto AFTER_OBJ_KEY_STRING
				}

				// The UTF-8 verification code was borrowed from utf8.ValidString
				// https://cs.opensource.google/go/go/+/refs/tags/go1.21.2:src/unicode/utf8/utf8.go;l=528
				{
					sv := ss[:len(ss)-len(s)]
					// Fast path. Check for and skip 8 bytes of
					// ASCII characters per iteration.
					for len(sv) >= 8 {
						// Combining two 32 bit loads allows the same code to be used for 32 and 64 bit platforms.
						// The compiler can generate a 32bit load for first32 and second32 on many platforms.
						// See test/codegen/memcombine.go.
						first32 := uint32(sv[0]) |
							uint32(sv[1])<<8 |
							uint32(sv[2])<<16 |
							uint32(sv[3])<<24
						second32 := uint32(sv[4]) |
							uint32(sv[5])<<8 |
							uint32(sv[6])<<16 |
							uint32(sv[7])<<24
						if (first32|second32)&0x80808080 != 0 {
							// Found a non ASCII byte (>= RuneSelf).
							break
						}
						sv = sv[8:]
					}
					n := len(sv)
					for j := 0; j < n; {
						si := sv[j]
						if si < utf8RuneSelf {
							j++
							continue
						}
						x := utf8First[si]
						if x == utf8xx {
							// Illegal starter byte.
							return s, getError(ErrorCodeInvalidUTF8, src, s)
						}
						size := int(x & 7)
						if j+size > n {
							// Short or invalid.
							return s, getError(ErrorCodeInvalidUTF8, src, s)
						}
						accept := utf8AcceptRanges[x>>4]
						if c := sv[j+1]; c < accept.lo || accept.hi < c {
							return s, getError(ErrorCodeInvalidUTF8, src, s)
						} else if size == 2 {
						} else if c := sv[j+2]; c < utf8locb || utf8hicb < c {
							return s, getError(ErrorCodeInvalidUTF8, src, s)
						} else if size == 3 {
						} else if c := sv[j+3]; c < utf8locb || utf8hicb < c {
							return s, getError(ErrorCodeInvalidUTF8, src, s)
						}
						j += size
					}
				}

				goto AFTER_OBJ_KEY_STRING
			default:
				if s[0] < 0x20 {
					return s, getError(ErrorCodeIllegalControlChar, src, s)
				}
				s = s[1:]
			}
		}
	}

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
	if s[0] < 0x20 {
		return s, getError(ErrorCodeIllegalControlChar, src, s)
	}
	return s, getError(ErrorCodeUnexpectedToken, src, s)
}
