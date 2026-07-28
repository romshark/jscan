package jsonnum_test

import (
	"strings"
	"testing"

	"github.com/romshark/jscan/v2/internal/jsonnum"
)

// refScan is a deliberately naive reference implementation of the RFC 8259
// number grammar used to cross-check the unrolled loops of ReadNumber.
// It returns the number of bytes consumed, whether the value is an integer
// (neither fraction nor exponent) and whether it's a valid number.
//
// Like ReadNumber it treats a '.' or 'e' that isn't followed by at least one
// digit as a hard error instead of backtracking to the last valid position,
// hence "1." is an error rather than the integer "1" followed by ".".
func refScan(s string) (n int, isInt, ok bool) {
	i := 0
	digit := func(j int) bool { return j < len(s) && s[j] >= '0' && s[j] <= '9' }

	if i < len(s) && s[i] == '-' {
		i++
	}
	if !digit(i) {
		return 0, false, false
	}
	if s[i] == '0' {
		// A leading zero terminates the integer part, anything
		// following it is left to the caller as trailing.
		i++
	} else {
		for digit(i) {
			i++
		}
	}
	isInt = true

	if i < len(s) && s[i] == '.' {
		i++
		if !digit(i) {
			return 0, false, false
		}
		for digit(i) {
			i++
		}
		isInt = false
	}

	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		i++
		if i < len(s) && (s[i] == '-' || s[i] == '+') {
			i++
		}
		if !digit(i) {
			return 0, false, false
		}
		for digit(i) {
			i++
		}
		isInt = false
	}
	return i, isInt, true
}

// FuzzReadNumber cross-checks ReadNumber against refScan.
//
// The distinction between [jsonnum.ReturnCodeInteger] and [jsonnum.ReturnCodeNumber] is
// deliberately asserted here because it isn't covered by FuzzValid in the root package:
// the classification doesn't affect whether a document is valid, yet it decides whether
// package jscan hands the bytes to internal/atoi, which assumes decimal digits only.
func FuzzReadNumber(f *testing.F) {
	for _, s := range []string{
		"0", "-0", "1", "42", "-42", "01", "-01",
		"0.0", "1.5", "-1.5", "1.",
		"0e0", "1e2", "1E2", "1e+2", "1e-2", "1e", "1e-", "1e+",
		"0.0e-1", "-1234567890.1234567890E-1234567890",
		"12345678", "123456789", "1234567890123456789012345",
		"1.12345678901234567890", "1e123456789012345678",
		"a", "-", ".5", "+1", "e1", "0x1",
		"1x", "1.5x", "1e2x", "12345678x",
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, data string) {
		// ReadNumber requires a non-empty input. Package jscan only ever
		// invokes it after having encountered '-' or a digit,
		// so an empty input can't occur in practice and isn't part of the contract.
		if len(data) < 1 {
			return
		}

		trailing, rc := jsonnum.ReadNumber(data)
		if !strings.HasSuffix(data, trailing) {
			t.Fatalf("ReadNumber(%q): trailing %q isn't a suffix of the input",
				data, trailing)
		}
		consumed := len(data) - len(trailing)

		// The []byte instantiation must agree with the string one.
		trailingBytes, rcBytes := jsonnum.ReadNumber([]byte(data))
		if rcBytes != rc || string(trailingBytes) != trailing {
			t.Fatalf("ReadNumber(%q): string returned (%q, %d) "+
				"but []byte returned (%q, %d)",
				data, trailing, rc, string(trailingBytes), rcBytes)
		}

		wantConsumed, wantIsInt, wantOK := refScan(data)

		if !wantOK {
			if rc != jsonnum.ReturnCodeErr {
				t.Fatalf("ReadNumber(%q): expected ReturnCodeErr, "+
					"received rc=%d (consumed %d)", data, rc, consumed)
			}
			return
		}
		if rc == jsonnum.ReturnCodeErr {
			t.Fatalf("ReadNumber(%q): unexpected ReturnCodeErr, expected to "+
				"consume %d byte(s) (isInt=%t)", data, wantConsumed, wantIsInt)
		}
		if consumed != wantConsumed {
			t.Fatalf("ReadNumber(%q): consumed %d byte(s), expected %d "+
				"(trailing=%q)", data, consumed, wantConsumed, trailing)
		}

		expect := jsonnum.ReturnCodeNumber
		if wantIsInt {
			expect = jsonnum.ReturnCodeInteger
		}
		if rc != expect {
			t.Fatalf("ReadNumber(%q): received rc=%d, expected %d",
				data, rc, expect)
		}
	})
}
