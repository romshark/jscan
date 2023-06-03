package keyescape

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var GI int

func repeat(s string, n int) string {
	b := make([]byte, 0, n*len(s))
	for i := 0; i < n; i++ {
		b = append(b, s...)
	}
	return string(b)
}

var tests = []struct {
	name   string
	input  string
	expect string
}{
	{
		"empty",
		"",
		"",
	},
	{
		"tilde",
		"~",
		"~0",
	},
	{
		"slash",
		"/",
		"~1",
	},
	{
		"space",
		" ",
		" ",
	},
	{
		"text3",
		"abc",
		"abc",
	},
	{
		"text9",
		"abc-def-1",
		"abc-def-1",
	},
	{
		"text12",
		"012345678912",
		"012345678912",
	},
	{
		"text36",
		"012345670123456701234567012345670123",
		"012345670123456701234567012345670123",
	},
	{
		"prefix1_tilde",
		"x~",
		"x~0",
	},
	{
		"prefix2_tilde",
		repeat("x", 2) + "~",
		repeat("x", 2) + "~0",
	},
	{
		"prefix3_tilde",
		repeat("x", 3) + "~",
		repeat("x", 3) + "~0",
	},
	{
		"prefix4_tilde",
		repeat("x", 4) + "~",
		repeat("x", 4) + "~0",
	},
	{
		"prefix5_tilde",
		repeat("x", 5) + "~",
		repeat("x", 5) + "~0",
	},
	{
		"prefix6_tilde",
		repeat("x", 6) + "~",
		repeat("x", 6) + "~0",
	},
	{
		"prefix7_tilde",
		repeat("x", 7) + "~",
		repeat("x", 7) + "~0",
	},

	{
		"slash_postfix6",
		"/" + repeat("t", 6),
		"~1" + repeat("t", 6),
	},
	{
		"prefix1_slash_postfix6",
		"x/" + repeat("t", 6),
		"x~1" + repeat("t", 6),
	},

	{
		"prefix0_slash_postfix4",
		"/" + repeat("t", 4),
		"~1" + repeat("t", 4),
	},
	{
		"prefix0_tilde_postfix4",
		"~" + repeat("t", 4),
		"~0" + repeat("t", 4),
	},
	{
		"prefix1_tilde_postfix4",
		"x~" + repeat("t", 4),
		"x~0" + repeat("t", 4),
	},
	{
		"prefix2_tilde_postfix4",
		repeat("x", 2) + "~" + repeat("t", 4),
		repeat("x", 2) + "~0" + repeat("t", 4),
	},
	{
		"prefix3_tilde_postfix4",
		repeat("x", 3) + "~" + repeat("t", 4),
		repeat("x", 3) + "~0" + repeat("t", 4),
	},
	{
		"prefix4_tilde_postfix4",
		repeat("x", 4) + "~" + repeat("t", 4),
		repeat("x", 4) + "~0" + repeat("t", 4),
	},
	{
		"prefix5_tilde_postfix4",
		repeat("x", 5) + "~" + repeat("t", 4),
		repeat("x", 5) + "~0" + repeat("t", 4),
	},
	{
		"prefix6_tilde_postfix4",
		repeat("x", 6) + "~" + repeat("t", 4),
		repeat("x", 6) + "~0" + repeat("t", 4),
	},
	{
		"prefix7_tilde_postfix4",
		repeat("x", 7) + "~" + repeat("t", 4),
		repeat("x", 7) + "~0" + repeat("t", 4),
	},

	{
		"prefix0_slash_postfix12",
		"/" + repeat("t", 12),
		"~1" + repeat("t", 12),
	},
	{
		"prefix0_tilde_postfix12",
		"~" + repeat("t", 12),
		"~0" + repeat("t", 12),
	},
	{
		"prefix1_tilde_postfix12",
		"x~" + repeat("t", 12),
		"x~0" + repeat("t", 12),
	},
	{
		"prefix2_tilde_postfix12",
		repeat("x", 2) + "~" + repeat("t", 12),
		repeat("x", 2) + "~0" + repeat("t", 12),
	},
	{
		"prefix3_tilde_postfix12",
		repeat("x", 3) + "~" + repeat("t", 12),
		repeat("x", 3) + "~0" + repeat("t", 12),
	},
	{
		"prefix4_tilde_postfix12",
		repeat("x", 4) + "~" + repeat("t", 12),
		repeat("x", 4) + "~0" + repeat("t", 12),
	},
	{
		"prefix5_tilde_postfix12",
		repeat("x", 5) + "~" + repeat("t", 12),
		repeat("x", 5) + "~0" + repeat("t", 12),
	},
	{
		"prefix6_tilde_postfix12",
		repeat("x", 6) + "~" + repeat("t", 12),
		repeat("x", 6) + "~0" + repeat("t", 12),
	},
	{
		"prefix7_tilde_postfix12",
		repeat("x", 7) + "~" + repeat("t", 12),
		repeat("x", 7) + "~0" + repeat("t", 12),
	},

	{
		"text1024",
		repeat("x", 1024),
		repeat("x", 1024),
	},
	{
		"multiple",
		"~abc/def/g~ ",
		"~0abc~1def~1g~0 ",
	},
	{
		"prefix1024",
		repeat("x", 1024) + "~abc/def/g~ ",
		repeat("x", 1024) + "~0abc~1def~1g~0 ",
	},
}

func Test(t *testing.T) {
	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			t.Run("string", func(t *testing.T) {
				t.Run("Append", func(t *testing.T) {
					a := Append(nil, td.input)
					require.Equal(t, td.expect, string(a))
				})

				t.Run("checkAndReplace", func(t *testing.T) {
					a := variantCheckAndReplace(nil, td.input)
					require.Equal(t, td.expect, string(a))
				})

				t.Run("checkAndReplaceUnrolled", func(t *testing.T) {
					a := variantCheckAndReplaceUnrolled(nil, td.input)
					require.Equal(t, td.expect, string(a))
				})

				t.Run("stdReplacer", func(t *testing.T) {
					r := strings.NewReplacer("~", "~0", "/", "~1")
					a := variantStdReplacer(r, nil, td.input)
					require.Equal(t, td.expect, string(a))
				})
			})

			t.Run("bytes", func(t *testing.T) {
				t.Run("Append", func(t *testing.T) {
					a := Append(nil, []byte(td.input))
					require.Equal(t, td.expect, string(a))
				})

				t.Run("checkAndReplace", func(t *testing.T) {
					a := variantCheckAndReplace(nil, []byte(td.input))
					require.Equal(t, td.expect, string(a))
				})

				t.Run("checkAndReplaceUnrolled", func(t *testing.T) {
					a := variantCheckAndReplaceUnrolled(nil, []byte(td.input))
					require.Equal(t, td.expect, string(a))
				})

				t.Run("stdReplacer", func(t *testing.T) {
					r := strings.NewReplacer("~", "~0", "/", "~1")
					a := variantStdReplacer(r, nil, []byte(td.input))
					require.Equal(t, td.expect, string(a))
				})
			})
		})
	}
}

func Benchmark(b *testing.B) {
	buffer := make([]byte, 0, 1024*8)
	for _, td := range tests {
		b.Run(td.name, func(b *testing.B) {
			b.Run("Append", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					buffer = Append(buffer, td.input)
					buffer = buffer[:0]
				}
			})

			b.Run("checkAndReplace", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					buffer = variantCheckAndReplace(buffer, td.input)
					buffer = buffer[:0]
				}
			})

			b.Run("checkAndReplaceUnrolled", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					buffer = variantCheckAndReplaceUnrolled(buffer, td.input)
					buffer = buffer[:0]
				}
			})

			b.Run("stdReplacer", func(b *testing.B) {
				r := strings.NewReplacer("~", "~0", "/", "~1")
				b.ResetTimer()
				for n := 0; n < b.N; n++ {
					buffer = variantStdReplacer(r, buffer, td.input)
					buffer = buffer[:0]
				}
			})
		})
	}
}
