package jscan

import (
	"encoding/json"
	"testing"
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
			jscanBytes = ValidBytes([]byte(data))
			jscanStr   = Valid(data)
		)
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
