package unescape_test

import (
	"encoding/json"
	"os"
	"runtime"
	"strconv"
	"testing"

	"github.com/romshark/jscan/v2/internal/unescape"
	"github.com/stretchr/testify/require"
)

type test struct{ Name, Input, Expect string }

var tests = []test{
	{Name: "empty", Input: "", Expect: ""},
	{Name: "noescape_1", Input: "a", Expect: "a"},
	{Name: "noescape_2", Input: "ж", Expect: "ж"},
	{Name: "reverse_solidus", Input: `\\`, Expect: `\`},
	{Name: "solidus", Input: `\/`, Expect: `/`},
	{Name: `\"`, Input: `\"`, Expect: "\""},
	{Name: "\\b", Input: `\b`, Expect: "\b"},
	{Name: "\\f", Input: `\f`, Expect: "\f"},
	{Name: "\\n", Input: `\n`, Expect: "\n"},
	{Name: "\\r", Input: `\r`, Expect: "\r"},
	{Name: "\\t", Input: `\t`, Expect: "\t"},
	{Name: "multiple_escaped", Input: `\r\n\r\n\tabc\n.`, Expect: "\r\n\r\n\tabc\n."},
	{Name: "unicode_0436", Input: `\u0436`, Expect: `ж`},
	{Name: "unicode_ffa2", Input: `\uffa2`, Expect: `ﾢ`},
	{Name: "unicode_fc00", Input: `\ufc00`, Expect: `ﰀ`},
	{Name: "solidus_escaped", Input: `\u002f`, Expect: `/`},
	{Name: "unicode_uppercase", Input: `\uAAAA`, Expect: `ꪪ`},
	{Name: "unicode_lowercase", Input: `\uAAaa`, Expect: `ꪪ`},
	{
		Name:   "noescape_65",
		Input:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed doe.",
		Expect: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed doe.",
	},
	{
		Name:   "escaped_with_gaps1",
		Input:  `\ra\nb\tc\td\\e\\f\ng`,
		Expect: "\ra\nb\tc\td\\e\\f\ng",
	},
	{
		Name:   "escaped_with_gaps2",
		Input:  `\r12\n12\t12\t12\\12\\12\n12`,
		Expect: "\r12\n12\t12\t12\\12\\12\n12",
	},
	{
		Name:   "escaped_with_gaps3",
		Input:  `\r123\n123\t123\t123\\123\\123\n123`,
		Expect: "\r123\n123\t123\t123\\123\\123\n123",
	},
	{
		Name:   "escaped_with_gaps4",
		Input:  `\r1234\n1234\t1234\t1234\\1234\\1234\n1234`,
		Expect: "\r1234\n1234\t1234\t1234\\1234\\1234\n1234",
	},
	{
		Name:   "escaped_with_gaps5",
		Input:  `\r12345\n12345\t12345\t12345\\12345\\12345\n12345`,
		Expect: "\r12345\n12345\t12345\t12345\\12345\\12345\n12345",
	},
	{
		Name:   "escaped_gap6",
		Input:  `\r123456\n`,
		Expect: "\r123456\n",
	},
	{
		Name:   "escaped_gap7",
		Input:  `\r1234567\n`,
		Expect: "\r1234567\n",
	},
	{
		Name:   "escaped_gap8",
		Input:  `\r12345678\n`,
		Expect: "\r12345678\n",
	},
	{
		Name:   "noescape_8",
		Input:  `12345678`,
		Expect: "12345678",
	},
	{
		Name:   "noescape_9",
		Input:  `123456789`,
		Expect: "123456789",
	},
	{
		Name:   "suffix_9",
		Input:  `12345678\n123456789`,
		Expect: "12345678\n123456789",
	},
	{
		Name:   "prefix21",
		Input:  `Lorem ipsum dolor sit\n`,
		Expect: "Lorem ipsum dolor sit\n",
	},
	{
		Name:   "prefix64",
		Input:  `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed doe.\n`,
		Expect: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed doe.\n",
	},
	{
		Name:   "unicode_latin",
		Input:  `latin prefix アパート 飛行船 \r\n some latin suffix`,
		Expect: "latin prefix アパート 飛行船 \r\n some latin suffix",
	},
	func() test {
		c, err := os.ReadFile("testdata/lorem_ipsum_noescape_5k.txt")
		if err != nil {
			panic(err)
		}
		return test{Name: "lorem_ipsum_noescape_5k", Input: string(c), Expect: string(c)}
	}(),
}

func TestValid(t *testing.T) {
	for _, td := range tests {
		t.Run(td.Name+"/string", func(t *testing.T) {
			actual := unescape.Valid(td.Input)
			runtime.GC() // Make sure the GC is happy.
			require.Equal(t, td.Expect, actual)
			var std string
			err := json.Unmarshal([]byte(`"`+td.Input+`"`), &std)
			require.NoError(t, err)
			require.Equal(t, td.Expect, std, "deviation between strconv and jscan")

			if td.Input == `\/` {
				return // Escaped solidus is not supported by strconv.Unquote
			}
			strconvResult, err := strconv.Unquote(`"` + td.Input + `"`)
			require.NoError(t, err)
			require.Equal(t, td.Expect, strconvResult)
		})

		t.Run(td.Name+"/bytes", func(t *testing.T) {
			input := []byte(td.Input)
			actual := unescape.Valid(input)
			runtime.GC() // Make sure the GC is happy.
			require.Equal(t, td.Expect, actual)
			var std string
			err := json.Unmarshal([]byte(`"`+td.Input+`"`), &std)
			require.NoError(t, err)
			require.Equal(t, td.Expect, std, "deviation between strconv and jscan")
		})
	}
}

func TestValidErrRune(t *testing.T) {
	for _, td := range []struct{ Name, Input, Expect string }{
		{Name: "D800", Input: `\uD800`, Expect: "�"},
		{Name: "DFFF", Input: `\uDFFF`, Expect: "�"},
		{Name: "FFFF", Input: `\uFFFF`, Expect: "￿"},
	} {
		t.Run(td.Name, func(t *testing.T) {
			unescaped := unescape.Valid(td.Input)
			require.Equal(t, td.Expect, unescaped)

			strconvResult, err := strconv.Unquote(td.Input)
			require.Equal(t, strconv.ErrSyntax, err)
			require.Zero(t, strconvResult)
		})
	}
}

func BenchmarkValid(b *testing.B) {
	var result string
	var err error
	for _, td := range tests {
		b.Run(td.Name, func(b *testing.B) {
			quoted := []byte(`"` + td.Input + `"`)
			quotedStr := string(quoted)

			b.Run("json_unmarshal", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					if err = json.Unmarshal(quoted, &result); err != nil {
						b.Fatal(err)
					}
				}
			})

			if td.Input != `\/` { // Escaped solidus is not supported by strconv.Unquote
				b.Run("strconv", func(b *testing.B) {
					for n := 0; n < b.N; n++ {
						if result, err = strconv.Unquote(quotedStr); err != nil {
							b.Fatal(err)
						}
					}
				})
			}

			b.Run("jscan", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					result = unescape.Valid(td.Input)
				}
			})
		})
	}
	runtime.KeepAlive(result)
	runtime.KeepAlive(err)
}
