package jscan

import (
	"github.com/romshark/jscan/v2/internal/jsonnum"
	"github.com/romshark/jscan/v2/internal/strfind"
	"github.com/romshark/jscan/v2/internal/utf8"
)

// scan calls fn for every value encountered.
// Returns the remainder of i.src and an error if any is encountered.
func scan[S ~string | ~[]byte](
	i *Iterator[S], fn func(*Iterator[S]) (err bool), noUTF8Validation bool,
) (S, Error[S]) {
	var (
		ss     S // Used as fallback for error report and for UTF-8 validation
		s      = i.src
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

	{ // Invoke callback
		i.arrayIndex = -1
		if len(i.stack) != 0 &&
			i.stack[len(i.stack)-1].Type == stackNodeTypeArray {
			i.arrayIndex = i.stack[len(i.stack)-1].ArrLen
			i.stack[len(i.stack)-1].ArrLen++
		}
		if fn(i) {
			return s, i.getError(ErrorCodeCallback)
		}
		i.keyIndex = -1
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

	{ // Invoke callback
		i.arrayIndex = -1
		if len(i.stack) != 0 &&
			i.stack[len(i.stack)-1].Type == stackNodeTypeArray {
			i.arrayIndex = i.stack[len(i.stack)-1].ArrLen
			i.stack[len(i.stack)-1].ArrLen++
		}
		if fn(i) {
			return s, i.getError(ErrorCodeCallback)
		}
		i.keyIndex = -1
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
		{
			ss = s
			if s, b = jsonnum.ReadNumber(s); b {
				return s, getError(ErrorCodeMalformedNumber, i.src, ss)
			}
		}
		i.valueIndexEnd = len(i.src) - len(s)
		i.valueType = ValueTypeNumber

		{ // Invoke callback
			i.arrayIndex = -1
			if len(i.stack) != 0 &&
				i.stack[len(i.stack)-1].Type == stackNodeTypeArray {
				i.arrayIndex = i.stack[len(i.stack)-1].ArrLen
				i.stack[len(i.stack)-1].ArrLen++
			}
			if fn(i) {
				return s, i.getError(ErrorCodeCallback)
			}
			i.keyIndex = -1
		}
	}
	goto AFTER_VALUE

VALUE_STRING:
	s = s[1:]
	ss = s
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
				s = s[1:]
				return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
			}
			if lutEscape[s[1]] == 1 {
				s = s[2:]
				continue
			}
			if s[1] != 'u' {
				return s, getError(ErrorCodeInvalidEscape, i.src, s)
			}
			if len(s) < 6 ||
				lutSX[s[5]] != 2 ||
				lutSX[s[4]] != 2 ||
				lutSX[s[3]] != 2 ||
				lutSX[s[2]] != 2 {
				return s, getError(ErrorCodeInvalidEscape, i.src, s)
			}
			s = s[5:]
		case '"':
			s = s[1:]

			if noUTF8Validation {
				goto AFTER_VALUE
			}

			// The UTF-8 verification code was borrowed from utf8.ValidString
			// https://cs.opensource.google/go/go/+/refs/tags/go1.21.2:src/unicode/utf8/utf8.go;l=528
			// See LICENCES.md for more information.
			{
				ss = ss[:len(ss)-len(s)]
				// Fast path. Check for and skip 8 bytes of
				// ASCII characters per iteration.
				for len(ss) >= 8 {
					// Combining two 32 bit loads allows the same code to be used for 32 and 64 bit platforms.
					// The compiler can generate a 32bit load for first32 and second32 on many platforms.
					// See test/codegen/memcombine.go.
					first32 := uint32(ss[0]) |
						uint32(ss[1])<<8 |
						uint32(ss[2])<<16 |
						uint32(ss[3])<<24
					second32 := uint32(ss[4]) |
						uint32(ss[5])<<8 |
						uint32(ss[6])<<16 |
						uint32(ss[7])<<24
					if (first32|second32)&0x80808080 != 0 {
						// Found a non ASCII byte (>= RuneSelf).
						break
					}
					ss = ss[8:]
				}
				n := len(ss)
				for j := 0; j < n; {
					si := ss[j]
					if si < utf8.RuneSelf {
						j++
						continue
					}
					x := utf8.First[si]
					if x == utf8.XX {
						// Illegal starter byte.
						return s, getError(ErrorCodeInvalidUTF8, i.src, s)
					}
					size := int(x & 7)
					if j+size > n {
						// Short or invalid.
						return s, getError(ErrorCodeInvalidUTF8, i.src, s)
					}
					accept := utf8.AcceptRanges[x>>4]
					if c := ss[j+1]; c < accept.Lo || accept.Hi < c {
						return s, getError(ErrorCodeInvalidUTF8, i.src, s)
					} else if size == 2 {
					} else if c := ss[j+2]; c < utf8.Locb || utf8.Hicb < c {
						return s, getError(ErrorCodeInvalidUTF8, i.src, s)
					} else if size == 3 {
					} else if c := ss[j+3]; c < utf8.Locb || utf8.Hicb < c {
						return s, getError(ErrorCodeInvalidUTF8, i.src, s)
					}
					j += size
				}
			}

			i.valueIndexEnd = len(i.src) - len(s)
			i.valueType = ValueTypeString

			{ // Invoke callback
				i.arrayIndex = -1
				if len(i.stack) != 0 &&
					i.stack[len(i.stack)-1].Type == stackNodeTypeArray {
					i.arrayIndex = i.stack[len(i.stack)-1].ArrLen
					i.stack[len(i.stack)-1].ArrLen++
				}
				if fn(i) {
					return s, i.getError(ErrorCodeCallback)
				}
				i.keyIndex = -1
			}

			goto AFTER_VALUE
		default:
			if s[0] < 0x20 {
				return s, getError(ErrorCodeIllegalControlChar, i.src, s)
			}
			s = s[1:]
		}
	}

VALUE_NULL:
	if len(s) < 4 || string(s[:4]) != "null" {
		return s, getError(ErrorCodeUnexpectedToken, i.src, s)
	}
	i.valueType = ValueTypeNull
	i.valueIndex = len(i.src) - len(s)
	i.valueIndexEnd = len(i.src) - len(s) + len("null")
	s = s[len("null"):]

	{ // Invoke callback
		i.arrayIndex = -1
		if len(i.stack) != 0 &&
			i.stack[len(i.stack)-1].Type == stackNodeTypeArray {
			i.arrayIndex = i.stack[len(i.stack)-1].ArrLen
			i.stack[len(i.stack)-1].ArrLen++
		}
		if fn(i) {
			return s, i.getError(ErrorCodeCallback)
		}
		i.keyIndex = -1
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

	{ // Invoke callback
		i.arrayIndex = -1
		if len(i.stack) != 0 &&
			i.stack[len(i.stack)-1].Type == stackNodeTypeArray {
			i.arrayIndex = i.stack[len(i.stack)-1].ArrLen
			i.stack[len(i.stack)-1].ArrLen++
		}
		if fn(i) {
			return s, i.getError(ErrorCodeCallback)
		}
		i.keyIndex = -1
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

	{ // Invoke callback
		i.arrayIndex = -1
		if len(i.stack) != 0 &&
			i.stack[len(i.stack)-1].Type == stackNodeTypeArray {
			i.arrayIndex = i.stack[len(i.stack)-1].ArrLen
			i.stack[len(i.stack)-1].ArrLen++
		}
		if fn(i) {
			return s, i.getError(ErrorCodeCallback)
		}
		i.keyIndex = -1
	}

	goto AFTER_VALUE

OBJ_KEY:
	ss = s
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

	i.valueIndex = len(i.src) - len(s) - 1
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
			return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
		}
		switch s[0] {
		case '\\':
			if len(s) < 2 {
				s = s[1:]
				return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
			}
			if lutEscape[s[1]] == 1 {
				s = s[2:]
				continue
			}
			if s[1] != 'u' {
				return s, getError(ErrorCodeInvalidEscape, i.src, s)
			}
			if len(s) < 6 ||
				lutSX[s[5]] != 2 ||
				lutSX[s[4]] != 2 ||
				lutSX[s[3]] != 2 ||
				lutSX[s[2]] != 2 {
				return s, getError(ErrorCodeInvalidEscape, i.src, s)
			}
			s = s[5:]
		case '"':
			s = s[1:]

			if noUTF8Validation {
				goto AFTER_OBJ_KEY_STRING
			}

			// The UTF-8 verification code was borrowed from utf8.ValidString
			// https://cs.opensource.google/go/go/+/refs/tags/go1.21.2:src/unicode/utf8/utf8.go;l=528
			// See LICENCES.md for more information.
			{ // Verify UTF-8 in string value
				ss = ss[:len(ss)-len(s)]
				// Fast path. Check for and skip 8 bytes of
				// ASCII characters per iteration.
				for len(ss) >= 8 {
					// Combining two 32 bit loads allows the same code to be used for 32 and 64 bit platforms.
					// The compiler can generate a 32bit load for first32 and second32 on many platforms.
					// See test/codegen/memcombine.go.
					first32 := uint32(ss[0]) |
						uint32(ss[1])<<8 |
						uint32(ss[2])<<16 |
						uint32(ss[3])<<24
					second32 := uint32(ss[4]) |
						uint32(ss[5])<<8 |
						uint32(ss[6])<<16 |
						uint32(ss[7])<<24
					if (first32|second32)&0x80808080 != 0 {
						// Found a non ASCII byte (>= RuneSelf).
						break
					}
					ss = ss[8:]
				}
				n := len(ss)
				for j := 0; j < n; {
					si := ss[j]
					if si < utf8.RuneSelf {
						j++
						continue
					}
					x := utf8.First[si]
					if x == utf8.XX {
						// Illegal starter byte.
						return s, getError(ErrorCodeInvalidUTF8, i.src, s)
					}
					size := int(x & 7)
					if j+size > n {
						// Short or invalid.
						return s, getError(ErrorCodeInvalidUTF8, i.src, s)
					}
					accept := utf8.AcceptRanges[x>>4]
					if c := ss[j+1]; c < accept.Lo || accept.Hi < c {
						return s, getError(ErrorCodeInvalidUTF8, i.src, s)
					} else if size == 2 {
					} else if c := ss[j+2]; c < utf8.Locb || utf8.Hicb < c {
						return s, getError(ErrorCodeInvalidUTF8, i.src, s)
					} else if size == 3 {
					} else if c := ss[j+3]; c < utf8.Locb || utf8.Hicb < c {
						return s, getError(ErrorCodeInvalidUTF8, i.src, s)
					}
					j += size
				}
			}

			i.keyIndex, i.keyIndexEnd = i.valueIndex, len(i.src)-len(s)
			goto AFTER_OBJ_KEY_STRING
		default:
			if s[0] < 0x20 {
				return s, getError(ErrorCodeIllegalControlChar, i.src, s)
			}
			s = s[1:]
		}
	}

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
	if s[0] < 0x20 {
		return s, getError(ErrorCodeIllegalControlChar, i.src, s)
	}
	return s, getError(ErrorCodeUnexpectedToken, i.src, s)
}
