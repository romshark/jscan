package strfind

import "bytes"

type ErrCode int

const (
	ErrCodeOK ErrCode = iota
	ErrCodeInvalidEscapeSeq
	ErrCodeIllegalControlChar
	ErrCodeUnexpectedEOF
)

// safeCharSet maps 0 to all ASCII characters that don't require checking
// during string traversal.
// 1 is mapped to control, quotation mark (") and reverse solidus ("\").
var safeCharSet = [256]byte{
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	'"': 1, '\\': 1,
}

// escapableChars maps escapable characters to 1,
// all other ASCII characters are mapped to 0.
var escapableChars = [256]byte{
	'"':  1,
	'\\': 1,
	'/':  1,
	'b':  1,
	'f':  1,
	'n':  1,
	'r':  1,
	't':  1,
}

// ReadString returns the index of the string value terminator in a JSON string
// with respect to escape sequences.
func ReadString[S ~string | ~[]byte](s S) (trailing S, errCode ErrCode) {
	for {
		for ; len(s) > 15; s = s[16:] {
			if safeCharSet[s[0]] != 0 {
				goto CHECK
			}
			if safeCharSet[s[1]] != 0 {
				s = s[1:]
				goto CHECK
			}
			if safeCharSet[s[2]] != 0 {
				s = s[2:]
				goto CHECK
			}
			if safeCharSet[s[3]] != 0 {
				s = s[3:]
				goto CHECK
			}
			if safeCharSet[s[4]] != 0 {
				s = s[4:]
				goto CHECK
			}
			if safeCharSet[s[5]] != 0 {
				s = s[5:]
				goto CHECK
			}
			if safeCharSet[s[6]] != 0 {
				s = s[6:]
				goto CHECK
			}
			if safeCharSet[s[7]] != 0 {
				s = s[7:]
				goto CHECK
			}
			if safeCharSet[s[8]] != 0 {
				s = s[8:]
				goto CHECK
			}
			if safeCharSet[s[9]] != 0 {
				s = s[9:]
				goto CHECK
			}
			if safeCharSet[s[10]] != 0 {
				s = s[10:]
				goto CHECK
			}
			if safeCharSet[s[11]] != 0 {
				s = s[11:]
				goto CHECK
			}
			if safeCharSet[s[12]] != 0 {
				s = s[12:]
				goto CHECK
			}
			if safeCharSet[s[13]] != 0 {
				s = s[13:]
				goto CHECK
			}
			if safeCharSet[s[14]] != 0 {
				s = s[14:]
				goto CHECK
			}
			if safeCharSet[s[15]] != 0 {
				s = s[15:]
				goto CHECK
			}
			continue
		}

	CHECK:
		if len(s) < 1 {
			return s, ErrCodeUnexpectedEOF
		}
		switch s[0] {
		case '\\':
			if len(s) < 2 {
				return s, ErrCodeUnexpectedEOF
			}
			if escapableChars[s[1]] == 1 {
				s = s[2:]
				continue
			}
			if s[1] != 'u' {
				return s, ErrCodeInvalidEscapeSeq
			}
			if len(s) < 6 ||
				charMap[s[5]] != 2 ||
				charMap[s[4]] != 2 ||
				charMap[s[3]] != 2 ||
				charMap[s[2]] != 2 {
				return s, ErrCodeInvalidEscapeSeq
			}
			s = s[5:]
		case '"':
			return s, ErrCodeOK
		default:
			if s[0] < 0x20 {
				return s, ErrCodeIllegalControlChar
			}
			s = s[1:]
		}
	}
}

func LastIndexUnescaped(path []byte, b byte) (i int) {
MAIN:
	for i = len(path); i >= 0; {
		path = path[:i]
		i = bytes.LastIndexByte(path, b)
		if i < 0 || i == 0 {
			return
		} else if path[i-1] != '\\' {
			return
		}
		for x := i - 1; ; x-- {
			if x == -1 || path[x] != '\\' {
				if z := x + 1; (i-(z))%2 > 0 {
					// Escaped, continue search
					i = z
					break
				}
				break MAIN
			}
		}
	}
	return
}

// charMap maps space characters such as whitespace, tab, line-break and
// carriage-return to 1, valid hex digits to 2 and all other ASCII characters to 0.
var charMap = [256]byte{
	' ': 1, '\n': 1, '\t': 1, '\r': 1,

	'0': 2, '1': 2, '2': 2, '3': 2, '4': 2, '5': 2, '6': 2, '7': 2, '8': 2, '9': 2,
	'a': 2, 'b': 2, 'c': 2, 'd': 2, 'e': 2, 'f': 2,
	'A': 2, 'B': 2, 'C': 2, 'D': 2, 'E': 2, 'F': 2,
}

// EndOfWhitespaceSeq returns the index of the end of
// the whitespace sequence.
// If the returned ok == false then index points at an
// illegal character that was encountered during the scan.
func EndOfWhitespaceSeq[S ~string | ~[]byte](s S) (trailing S, ok bool) {
	for ; len(s) > 15; s = s[16:] {
		if charMap[s[0]] != 1 {
			goto NONSPACE
		}
		if charMap[s[1]] != 1 {
			s = s[1:]
			goto NONSPACE
		}
		if charMap[s[2]] != 1 {
			s = s[2:]
			goto NONSPACE
		}
		if charMap[s[3]] != 1 {
			s = s[3:]
			goto NONSPACE
		}
		if charMap[s[4]] != 1 {
			s = s[4:]
			goto NONSPACE
		}
		if charMap[s[5]] != 1 {
			s = s[5:]
			goto NONSPACE
		}
		if charMap[s[6]] != 1 {
			s = s[6:]
			goto NONSPACE
		}
		if charMap[s[7]] != 1 {
			s = s[7:]
			goto NONSPACE
		}
		if charMap[s[8]] != 1 {
			s = s[8:]
			goto NONSPACE
		}
		if charMap[s[9]] != 1 {
			s = s[9:]
			goto NONSPACE
		}
		if charMap[s[10]] != 1 {
			s = s[10:]
			goto NONSPACE
		}
		if charMap[s[11]] != 1 {
			s = s[11:]
			goto NONSPACE
		}
		if charMap[s[12]] != 1 {
			s = s[12:]
			goto NONSPACE
		}
		if charMap[s[13]] != 1 {
			s = s[13:]
			goto NONSPACE
		}
		if charMap[s[14]] != 1 {
			s = s[14:]
			goto NONSPACE
		}
		if charMap[s[15]] != 1 {
			s = s[15:]
			goto NONSPACE
		}
	}
	for ; len(s) > 0; s = s[1:] {
		if charMap[s[0]] != 1 {
			goto NONSPACE
		}
	}
	return s, false

NONSPACE:
	return s, s[0] <= 0x20
}
