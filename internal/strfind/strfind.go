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
var charMap = [256]byte{' ': 1, '\n': 1, '\t': 1, '\r': 1}

// EndOfWhitespaceSeq returns the index of the end of
// the whitespace sequence.
// If the returned ctrlChar == true then index points at an
// illegal character that was encountered during the scan.
func EndOfWhitespaceSeq[S ~string | ~[]byte](s S) (trailing S) {
	for ; len(s) > 15; s = s[16:] {
		if charMap[s[0]] != 1 {
			return s
		}
		if charMap[s[1]] != 1 {
			s = s[1:]
			return s
		}
		if charMap[s[2]] != 1 {
			s = s[2:]
			return s
		}
		if charMap[s[3]] != 1 {
			s = s[3:]
			return s
		}
		if charMap[s[4]] != 1 {
			s = s[4:]
			return s
		}
		if charMap[s[5]] != 1 {
			s = s[5:]
			return s
		}
		if charMap[s[6]] != 1 {
			s = s[6:]
			return s
		}
		if charMap[s[7]] != 1 {
			s = s[7:]
			return s
		}
		if charMap[s[8]] != 1 {
			s = s[8:]
			return s
		}
		if charMap[s[9]] != 1 {
			s = s[9:]
			return s
		}
		if charMap[s[10]] != 1 {
			s = s[10:]
			return s
		}
		if charMap[s[11]] != 1 {
			s = s[11:]
			return s
		}
		if charMap[s[12]] != 1 {
			s = s[12:]
			return s
		}
		if charMap[s[13]] != 1 {
			s = s[13:]
			return s
		}
		if charMap[s[14]] != 1 {
			s = s[14:]
			return s
		}
		if charMap[s[15]] != 1 {
			s = s[15:]
			return s
		}
	}
	for ; len(s) > 0; s = s[1:] {
		if charMap[s[0]] != 1 {
			return s
		}
	}
	return s
}
