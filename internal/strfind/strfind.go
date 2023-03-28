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
			switch s[i] {
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
				j = 0
			case 'u':
				if i+4 >= len(s) ||
					!isValidHexDigits(s[i+4]) ||
					!isValidHexDigits(s[i+3]) ||
					!isValidHexDigits(s[i+2]) ||
					!isValidHexDigits(s[i+1]) {
					return i, ErrCodeInvalidEscapeSeq
				}
				i, j = i+4, 0
			default:
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

func isValidHexDigits(b byte) bool {
	return (b >= '0' && b <= '9') ||
		(b >= 'A' && b <= 'F') ||
		(b >= 'a' && b <= 'f')
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

// EndOfWhitespaceSeq returns the index of the end of
// the whitespace sequence.
// If the returned stoppedAtIllegalChar == true then index points at an
// illegal character that was encountered during the scan.
func EndOfWhitespaceSeq(s string) (index int, stoppedAtIllegalChar bool) {
	if len(s) == 0 || s[0] > 32 {
		return 0, false
	}
	i := 0
	for ; i < len(s); i++ {
		switch s[i] {
		case ' ', '\n', '\t', '\r':
		default:
			if s[i] < 0x20 {
				return i, true
			}
			return i, false
		}
	}
	return i, false
}

// EndOfWhitespaceSeqBytes returns the index of the end of
// the whitespace sequence.
// If the returned stoppedAtIllegalChar == true then index points at an
// illegal character that was encountered during the scan.
func EndOfWhitespaceSeqBytes(s []byte) (index int, stoppedAtIllegalChar bool) {
	if len(s) == 0 || s[0] > 32 {
		return 0, false
	}
	i := 0
	for ; i < len(s); i++ {
		switch s[i] {
		case ' ', '\n', '\t', '\r':
		default:
			if s[i] < 0x20 {
				return i, true
			}
			return i, false
		}
	}
	return i, false
}
