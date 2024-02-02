package unescape

import (
	"strings"
	"unsafe"
)

// Valid returns the unescaped version of str relying on str to be valid.
// Don't use this function if str isn't guaranteed to contain no
// invalid escape sequences.
func Valid[S ~[]byte | ~string](str S) string {
	if len(str) < 1 {
		return ""
	}
	var s string
	switch in := any(str).(type) {
	case string:
		s = in
	case []byte:
		// Avoid copying str to a string, treat the bytes as read-only instead
		// since str is guaranteed to remain immutable.
		s = unsafe.String(unsafe.SliceData(in), len(in))
	default:
		s = string(str)
	}
	i := strings.IndexByte(s, '\\')
	if i < 0 {
		return string(s)
	}
	buf := make([]byte, 0, len(s))
	for {
		buf = append(buf, s[:i]...)
		s = s[i:]
		if len(s) < 1 {
			break
		}
	ESCAPE_SEQUENCE:
		if s[1] == 'u' {
			switch r := uint32(lut[s[2]]<<12 + lut[s[3]]<<8 + lut[s[4]]<<4 + lut[s[5]]); {
			case r <= rune1Max:
				buf = append(buf, byte(r))
			case r <= rune2Max:
				buf = append(buf, t2|byte(r>>6), tx|byte(r)&maskx)
			case r > maxRune, surrogateMin <= r && r <= surrogateMax:
				buf = append(buf, RuneError...) // Error placeholder rune
			default:
				buf = append(buf, t3|byte(r>>12), tx|byte(r>>6)&maskx, tx|byte(r)&maskx)
			}
			s = s[6:]
		} else {
			buf, s = append(buf, lutReplace[s[1]]), s[2:]
		}
		const batchSize = 8
		if len(s) < batchSize {
			for i = 0; i < len(s) && s[i] != '\\'; i++ {
			}
			continue
		}
		if s[0] == '\\' {
			// Return directly to the escape sequence handler.
			goto ESCAPE_SEQUENCE
		}
		if s[1] == '\\' {
			i = 1
			continue
		}
		if s[2] == '\\' {
			i = 2
			continue
		}
		if s[3] == '\\' {
			i = 3
			continue
		}
		if s[4] == '\\' {
			i = 4
			continue
		}
		if s[5] == '\\' {
			i = 5
			continue
		}
		if s[6] == '\\' {
			i = 6
			continue
		}
		if s[7] == '\\' {
			i = 7
			continue
		}
		i = strings.IndexByte(s[batchSize:], '\\')
		if i < 0 {
			buf = append(buf, s...)
			break
		}
		i += batchSize
		continue
	}
	// Avoid copying the buffer, return it as a string instead
	// because it's guaranteed to never be referenced and mutated.
	return unsafe.String(unsafe.SliceData(buf), len(buf))
}

var lutReplace = [256]byte{
	'\\': '\\',
	'/':  '/',
	'"':  '"',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
}

var lut = [256]int{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4,
	'5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15,
	'a': 10, 'b': 11, 'c': 12, 'd': 13, 'e': 14, 'f': 15,
}

const (
	// Code points in the surrogate range are not valid for UTF-8.
	surrogateMin, surrogateMax = 0xD800, 0xDFFF

	tx, t2, t3, t4               = 0b10000000, 0b11000000, 0b11100000, 0b11110000
	maskx                        = 0b00111111
	rune1Max, rune2Max, rune3Max = 1<<7 - 1, 1<<11 - 1, 1<<16 - 1

	// maxRune is the maximum valid Unicode code point.
	maxRune = '\U0010FFFF'

	// RuneError is the "error" rune or "Unicode replacement character".
	RuneError = "\uFFFD"
)
