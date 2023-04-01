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
		switch {
		case std != jscanStr:
			t.Fatalf(
				`Valid(%q): %t (std) != %t (jscanStr)`,
				data, std, jscanStr,
			)
		case std != jscanBytes:
			t.Fatalf(
				`Valid(%q): %t (std) != %t (jscanBytes)`,
				data, std, jscanBytes,
			)
		case std != jscanParserStr:
			t.Fatalf(
				`Valid(%q): %t (std) != %t (jscanParserStr)`,
				data, std, jscanParserStr,
			)
		case std != jscanParserBytes:
			t.Fatalf(
				`Valid(%q): %t (std) != %t (jscanParserBytes)`,
				data, std, jscanParserBytes,
			)
		}
	})
}
