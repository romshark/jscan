package strfind

type ErrCode int

const (
	ErrCodeOK ErrCode = iota
	ErrCodeInvalidEscapeSeq
	ErrCodeIllegalControlChar
	ErrCodeUnexpectedEOF
)

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
