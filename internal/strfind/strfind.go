package strfind

import (
	"bytes"
	"strings"
)

// IndexTerm returns either -1 or the index of the string value terminator.
func IndexTerm(s string, i int) int {
	for {
		x := strings.IndexByte(s[i:], '"')
		if x < 0 {
			return -1
		}
		x += i

		bs := 0
		for j := x - 1; j >= 0 && s[j] == '\\'; j-- {
			bs++
		}
		if bs%2 > 0 {
			i++
			continue
		}
		return x
	}
}

// IndexTermBytes returns either -1 or the index of the string value terminator.
func IndexTermBytes(s []byte, i int) int {
	for {
		x := bytes.IndexByte(s[i:], '"')
		if x < 0 {
			return -1
		}
		x += i

		bs := 0
		for j := x - 1; j >= 0 && s[j] == '\\'; j-- {
			bs++
		}
		if bs%2 > 0 {
			i++
			continue
		}
		return x
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

// EndOfWhitespaceSeq returns the index of the end of
// the whitespace sequence.
// If the returned stoppedAtIllegalChar == true then index points at an
// illegal character that was encountered during the scan.
func EndOfWhitespaceSeq(s string) (index int, hasIllegalChars bool) {
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
