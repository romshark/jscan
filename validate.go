package jscan

import (
	"github.com/romshark/jscan/v2/internal/jsonnum"
	"github.com/romshark/jscan/v2/internal/strfind"
)

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
