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
// the whitespace sequence
func EndOfWhitespaceSeq(s string) int {
	if len(s) == 0 {
		return 0
	}
	switch s[0] {
	case ' ', '\n', '\t', '\r':
	default:
		return 0
	}
	i := 1
	for ; i < len(s); i++ {
		switch s[i] {
		case ' ', '\n', '\t', '\r':
		default:
			return i
		}
	}
	return i
}

// EndOfWhitespaceSeqBytes returns the index of the end of
// the whitespace sequence
func EndOfWhitespaceSeqBytes(s []byte) int {
	if len(s) == 0 {
		return 0
	}
	switch s[0] {
	case ' ', '\n', '\t', '\r':
	default:
		return 0
	}
	i := 1
	for ; i < len(s); i++ {
		switch s[i] {
		case ' ', '\n', '\t', '\r':
		default:
			return i
		}
	}
	return i
}
