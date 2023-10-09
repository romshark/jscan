package jscan_test

import (
	"encoding/json"
	"testing"
	"unicode/utf8"

	"github.com/romshark/jscan/v2"
)

func FuzzValid(f *testing.F) {
	for _, s := range []string{
		``,
		`0`,
		`42`,
		`3.14159`,
		`null`,
		`false`,
		`true`,
		`{"\"escaped_key\"":"\"escaped_value\""}`,
		`[]`,
		`[1]`,
		`[null, 1,  3 , 3.14159]`,
		"{}",
		`{"foo": "bar"}`,
		`"foo"`,
		`"{"`,
		`"{}"`,
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, data string) {
		var (
			std        = json.Valid([]byte(data))
			jscanBytes = jscan.Valid([]byte(data), jscan.Options{})
			jscanStr   = jscan.Valid(data, jscan.Options{})
		)
		if !utf8.ValidString(data) {
			if jscanBytes || jscanStr {
				t.Fatalf(
					`Invalid UTF8: jscanBytes: %t; jscanStr: %t`,
					jscanBytes, jscanStr,
				)
			}
			// no need to check against encoding/json as it doesn't validate UTF-8.
			return
		}
		if std != jscanStr {
			t.Fatalf(
				`Valid(%q): %t (std) != %t (jscanStr)`,
				data, std, jscanStr,
			)
		} else if std != jscanBytes {
			t.Fatalf(
				`Valid(%q): %t (std) != %t (jscanBytes)`,
				data, std, jscanBytes,
			)
		}
	})
}
