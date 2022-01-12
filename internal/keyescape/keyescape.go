package keyescape

import "strings"

var keyEscapeReplacer = strings.NewReplacer(
	".", `\.`,
	"[", `\[`,
	"]", `\]`,
)

var escapedPeriod = []byte(`\.`)
var escapedBracketLeft = []byte(`\[`)
var escapedBracketRight = []byte(`\]`)

func Escape(s string) string {
	return keyEscapeReplacer.Replace(s)
}

func EscapeAppend(dst, src []byte) []byte {
	if len(src) < 1 {
		return nil
	}
	for len(src) > 0 {
		switch src[0] {
		case '.':
			dst = append(dst, escapedPeriod...)
		case '[':
			dst = append(dst, escapedBracketLeft...)
		case ']':
			dst = append(dst, escapedBracketRight...)
		default:
			dst = append(dst, src[0])
		}
		src = src[1:]
	}
	return dst
}
