package strfind

import "bytes"

type ErrCode int

const (
	ErrCodeOK ErrCode = iota
	ErrCodeInvalidEscapeSeq
	ErrCodeIllegalControlChar
	ErrCodeUnexpectedEOF
)

// safeCharSet maps 0 to all inherently safe ASCII characters.
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

// IndexTerm returns either -1 or the index of the string value terminator.
func IndexTerm[S []byte | string](s S, i int) (indexEnd int, errCode ErrCode) {
	for j := 0; i < len(s); i++ {
		if i+8 < len(s) {
			if safeCharSet[s[i]] != 0 {
				goto CHECK
			}
			i++
			if safeCharSet[s[i]] != 0 {
				goto CHECK
			}
			i++
			if safeCharSet[s[i]] != 0 {
				goto CHECK
			}
			i++
			if safeCharSet[s[i]] != 0 {
				goto CHECK
			}
			i++
			if safeCharSet[s[i]] != 0 {
				goto CHECK
			}
			i++
			if safeCharSet[s[i]] != 0 {
				goto CHECK
			}
			i++
			if safeCharSet[s[i]] != 0 {
				goto CHECK
			}
			i++
			if safeCharSet[s[i]] != 0 {
				goto CHECK
			}
			i++
		}

	CHECK:
		switch s[i] {
		case '\\':
			i++
			if i >= len(s) {
				return i, ErrCodeUnexpectedEOF
			}
			if escapableChars[s[i]] == 1 {
				j = 0
			} else if s[i] == 'u' {
				if i+4 >= len(s) ||
					charMap[s[i+4]] != 2 ||
					charMap[s[i+3]] != 2 ||
					charMap[s[i+2]] != 2 ||
					charMap[s[i+1]] != 2 {
					return i, ErrCodeInvalidEscapeSeq
				}
				i, j = i+4, 0
			} else {
				return i, ErrCodeInvalidEscapeSeq
			}
		case '"':
			if j%2 == 0 {
				return i, ErrCodeOK
			}
		default:
			if s[i] < 0x20 {
				return i, ErrCodeIllegalControlChar
			}
		}
	}
	return i, ErrCodeUnexpectedEOF
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

// charMap maps space characters such as whitespace, tab
// line-break and carriage return to 1,
// valid hex digits to 2 and
// all other ASCII characters to 0
var charMap = [256]byte{
	' ': 1, '\n': 1, '\t': 1, '\r': 1,

	'0': 2, '1': 2, '2': 2, '3': 2, '4': 2, '5': 2, '6': 2, '7': 2, '8': 2, '9': 2,
	'a': 2, 'b': 2, 'c': 2, 'd': 2, 'e': 2, 'f': 2,
	'A': 2, 'B': 2, 'C': 2, 'D': 2, 'E': 2, 'F': 2,
}

// EndOfWhitespaceSeq returns the index of the end of
// the whitespace sequence.
// If the returned stoppedAtIllegalChar == true then index points at an
// illegal character that was encountered during the scan.
func EndOfWhitespaceSeq[S []byte | string](s S) (index int, stoppedAtIllegalChar bool) {
	if len(s) == 0 || s[0] > 0x20 {
		// Fast path
		return 0, false
	}
	i := 0
	for i < len(s) {
		if i+7 < len(s) {
			if charMap[s[i]] != 1 {
				if s[i] < 0x20 {
					return i, true
				}
				break
			}
			i++
			if charMap[s[i]] != 1 {
				if s[i] < 0x20 {
					return i, true
				}
				break
			}
			i++
			if charMap[s[i]] != 1 {
				if s[i] < 0x20 {
					return i, true
				}
				break
			}
			i++
			if charMap[s[i]] != 1 {
				if s[i] < 0x20 {
					return i, true
				}
				break
			}
			i++
			if charMap[s[i]] != 1 {
				if s[i] < 0x20 {
					return i, true
				}
				break
			}
			i++
			if charMap[s[i]] != 1 {
				if s[i] < 0x20 {
					return i, true
				}
				break
			}
			i++
			if charMap[s[i]] != 1 {
				if s[i] < 0x20 {
					return i, true
				}
				break
			}
			i++
			if charMap[s[i]] != 1 {
				if s[i] < 0x20 {
					return i, true
				}
				break
			}
			i++
			continue
		}
		if charMap[s[i]] != 1 {
			if s[i] < 0x20 {
				return i, true
			}
			break
		}
		i++
	}
	return i, false
}
