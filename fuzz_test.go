//go:build go1.18
// +build go1.18

package jscan

import (
	"encoding/json"
	"testing"
)

func FuzzValidBytes(f *testing.F) {
	for _, s := range []string{
		"{}",
		`{"foo": "bar"}`,
		``,
		`"foo"`,
		`"{"`,
		`"{}"`,
	} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		var (
			std   = json.Valid(data)
			jscan = ValidBytes(data)
		)
		if std != jscan {
			t.Fatalf(`Valid(%#v): %v (std) != %v (jscan)`, string(data), std, jscan)
		}
	})
}
