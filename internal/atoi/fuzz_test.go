package atoi_test

import (
	"strconv"
	"testing"

	"github.com/romshark/jscan/v2/internal/atoi"
)

// Package atoi assumes canonical decimal input as produced by the jscan
// tokenizer: optional '-', no '+', no leading zeros. All helpers below
// therefore only ever generate canonical strings.

func checkU(t *testing.T, fn, in string, gotN uint64, gotOverflow bool, bits int) {
	t.Helper()
	want, err := strconv.ParseUint(in, 10, bits)
	wantOverflow := err != nil
	if gotOverflow != wantOverflow {
		t.Fatalf("%s(%q): overflow=%t, expected %t (value %d)",
			fn, in, gotOverflow, wantOverflow, gotN)
	}
	if !wantOverflow && gotN != want {
		t.Fatalf("%s(%q): received %d, expected %d", fn, in, gotN, want)
	}
	if gotOverflow && gotN != 0 {
		t.Fatalf("%s(%q): expected 0 on overflow, received %d", fn, in, gotN)
	}
}

func checkI(t *testing.T, fn, in string, gotN int64, gotOverflow bool, bits int) {
	t.Helper()
	want, err := strconv.ParseInt(in, 10, bits)
	wantOverflow := err != nil
	if gotOverflow != wantOverflow {
		t.Fatalf("%s(%q): overflow=%t, expected %t (value %d)",
			fn, in, gotOverflow, wantOverflow, gotN)
	}
	if !wantOverflow && gotN != want {
		t.Fatalf("%s(%q): received %d, expected %d", fn, in, gotN, want)
	}
	if gotOverflow && gotN != 0 {
		t.Fatalf("%s(%q): expected 0 on overflow, received %d", fn, in, gotN)
	}
}

// checkAll runs every parser against s, which must be a canonical
// unsigned decimal string, both as-is and negated.
func checkAll(t *testing.T, s string) {
	t.Helper()
	neg := "-" + s

	u8, o := atoi.U8(s)
	checkU(t, "U8", s, uint64(u8), o, 8)
	u16, o := atoi.U16(s)
	checkU(t, "U16", s, uint64(u16), o, 16)
	u32, o := atoi.U32(s)
	checkU(t, "U32", s, uint64(u32), o, 32)
	u64, o := atoi.U64(s)
	checkU(t, "U64", s, u64, o, 64)

	i8, o := atoi.I8(s)
	checkI(t, "I8", s, int64(i8), o, 8)
	i8, o = atoi.I8(neg)
	checkI(t, "I8", neg, int64(i8), o, 8)

	i16, o := atoi.I16(s)
	checkI(t, "I16", s, int64(i16), o, 16)
	i16, o = atoi.I16(neg)
	checkI(t, "I16", neg, int64(i16), o, 16)

	i32, o := atoi.I32(s)
	checkI(t, "I32", s, int64(i32), o, 32)
	i32, o = atoi.I32(neg)
	checkI(t, "I32", neg, int64(i32), o, 32)

	i64, o := atoi.I64(s)
	checkI(t, "I64", s, i64, o, 64)
	i64, o = atoi.I64(neg)
	checkI(t, "I64", neg, i64, o, 64)
}

// TestExhaustive8And16 checks every value that the 8- and 16-bit parsers
// can be handed, including the digit lengths at which they must report
// overflow. It's the regression test for the wraparound bug in which
// I8, I16 and I32 mistook a value that wrapped more than once for an
// in-range result (I8("300") returned 44 instead of overflowing).
func TestExhaustive8And16(t *testing.T) {
	for v := 0; v <= 99999; v++ {
		checkAll(t, strconv.Itoa(v))
	}
}

// TestOverflowBoundaries checks the exact edges of every parser.
func TestOverflowBoundaries(t *testing.T) {
	for _, v := range []uint64{
		0, 1, 9, 10, 99, 100,
		127, 128, 129, 255, 256, 257, 999,
		32767, 32768, 32769, 65535, 65536, 65537, 99999,
		2147483647, 2147483648, 2147483649,
		4294967295, 4294967296, 4294967297,
		5000000000, 6442450943, 6442450944,
		8589934591, 8589934592, 9000000000, 9999999999,
		9223372036854775807, 9223372036854775808, 9223372036854775809,
		18446744073709551614, 18446744073709551615,
	} {
		checkAll(t, strconv.FormatUint(v, 10))
	}
	// Digit lengths beyond what any parser accepts.
	for _, s := range []string{
		"11111111111111111111",    // 20 digits
		"111111111111111111111",   // 21 digits
		"99999999999999999999999", // 23 digits
	} {
		checkAll(t, s)
	}
}

// FuzzAtoi cross-checks every parser against strconv.
func FuzzAtoi(f *testing.F) {
	for _, s := range []string{
		"0", "1", "9", "10", "127", "128", "255", "256", "300", "999",
		"32767", "32768", "65535", "65536", "70000", "98303",
		"2147483647", "2147483648", "4294967295", "4294967296",
		"5000000000", "9999999999",
		"9223372036854775807", "9223372036854775808",
		"18446744073709551615", "18446744073709551616",
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, data string) {
		// The parsers assume canonical decimal input, anything else is
		// outside their contract and would be rejected by the tokenizer
		// long before reaching them.
		if len(data) < 1 || len(data) > 32 {
			return
		}
		for i := 0; i < len(data); i++ {
			if data[i] < '0' || data[i] > '9' {
				return
			}
		}
		if len(data) > 1 && data[0] == '0' {
			return // leading zero
		}
		checkAll(t, data)
	})
}
