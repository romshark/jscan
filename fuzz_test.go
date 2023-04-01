package jscan_test

import (
	"encoding/json"
	"testing"

	"github.com/romshark/jscan"
)

func FuzzValid(f *testing.F) {
	p, pb := jscan.New(64), jscan.NewBytes(64)

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
			std              = json.Valid([]byte(data))
			jscanParserBytes = pb.Valid([]byte(data))
			jscanParserStr   = p.Valid(data)
			jscanBytes       = jscan.ValidBytes([]byte(data))
			jscanStr         = jscan.Valid(data)
		)
		if std != jscanStr {
			t.Errorf(
				`Valid(%q): %t (std) != %t (jscanStr)`,
				data, std, jscanStr,
			)
		}
		if std != jscanBytes {
			t.Errorf(
				`Valid(%q): %t (std) != %t (jscanBytes)`,
				data, std, jscanBytes,
			)
		}
		if std != jscanParserStr {
			t.Errorf(
				`Valid(%q): %t (std) != %t (jscanParserStr)`,
				data, std, jscanParserStr,
			)
		}
		if std != jscanParserBytes {
			t.Errorf(
				`Valid(%q): %t (std) != %t (jscanParserBytes)`,
				data, std, jscanParserBytes,
			)
		}
	})
}
