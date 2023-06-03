package jscan_test

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/romshark/jscan"
	"github.com/stretchr/testify/require"
)

func TestX(t *testing.T) {
	in := `[ `

	require.False(t, json.Valid([]byte(in)))

	t.Run("string", func(t *testing.T) {
		check := func(t *testing.T) func(i *jscan.Iterator[string]) bool {
			return func(i *jscan.Iterator[string]) (err bool) { return false }
		}

		t.Run("Valid", func(t *testing.T) {
			require.False(t, jscan.Valid(in))
		})
		t.Run("Scan", func(t *testing.T) {
			err := jscan.Scan(in, check(t))
			require.True(t, err.IsErr())
		})
	})

	t.Run("bytes", func(t *testing.T) {
		check := func(t *testing.T) func(i *jscan.Iterator[[]byte]) bool {
			return func(i *jscan.Iterator[[]byte]) (err bool) {
				return false
			}
		}

		t.Run("Valid", func(t *testing.T) {
			require.False(t, jscan.Valid([]byte(in)))
		})
		t.Run("Scan", func(t *testing.T) {
			err := jscan.Scan([]byte(in), check(t))
			require.True(t, err.IsErr())
		})
	})
}

//go:embed jsontestsuite
var fsSuite embed.FS

func TestJSONTestSuite(t *testing.T) {
	d, err := fsSuite.ReadDir("jsontestsuite")
	require.NoError(t, err)
	for _, f := range d {
		n := f.Name()
		t.Run(n, func(t *testing.T) {
			switch {
			case strings.HasPrefix(n, "i_"):
				testOKOrErr(t, readTestFileContents(t, n))
			case strings.HasPrefix(n, "y_"):
				testStrictOK(t, readTestFileContents(t, n))
			case strings.HasPrefix(n, "n_"):
				testStrictErr(t, readTestFileContents(t, n))
			default:
				t.Skip(n)
			}
		})
	}
}

func readTestFileContents(t *testing.T, name string) []byte {
	c, err := fsSuite.ReadFile(filepath.Join("jsontestsuite", name))
	require.NoError(t, err)
	return c
}

// testStrictOK runs tests with the "y_" prefix that parsers must accept.
func testStrictOK(t *testing.T, input []byte) {
	t.Run("Valid", func(t *testing.T) {
		require.True(t, jscan.Valid(string(input)))
	})
	t.Run("Valid_bytes", func(t *testing.T) {
		require.True(t, jscan.Valid(input))
	})
	t.Run("Scan", func(t *testing.T) {
		err := jscan.Scan(
			string(input),
			func(i *jscan.Iterator[string]) (err bool) { return false },
		)
		require.False(t, err.IsErr())
	})
	t.Run("ScanBytes", func(t *testing.T) {
		err := jscan.Scan(
			input,
			func(i *jscan.Iterator[[]byte]) (err bool) { return false },
		)
		require.False(t, err.IsErr())
	})
	t.Run("ScanOne", func(t *testing.T) {
		_, err := jscan.ScanOne(
			string(input),
			func(i *jscan.Iterator[string]) (err bool) { return false },
		)
		require.False(t, err.IsErr())
	})
	t.Run("ScanOneBytes", func(t *testing.T) {
		_, err := jscan.ScanOne(
			input,
			func(i *jscan.Iterator[[]byte]) (err bool) { return false },
		)
		require.False(t, err.IsErr())
	})
}

// testOKOrErr runs tests with the "i_" prefix that
// parsers are free to accept or reject.
func testOKOrErr(t *testing.T, input []byte) {
	t.Run("Valid", func(t *testing.T) {
		if !jscan.Valid(string(input)) {
			t.Skip("allowed to fail")
		}
	})
	t.Run("Valid_bytes", func(t *testing.T) {
		if !jscan.Valid(input) {
			t.Skip("allowed to fail")
		}
	})
}

// testStrictErr runs tests with the "n_" prefix that parsers must reject.
func testStrictErr(t *testing.T, input []byte) {
	t.Run("Valid", func(t *testing.T) {
		require.False(t, jscan.Valid(string(input)))
	})
	t.Run("Valid_bytes", func(t *testing.T) {
		require.False(t, jscan.Valid(input))
	})
	t.Run("Scan", func(t *testing.T) {
		err := jscan.Scan(
			string(input),
			func(i *jscan.Iterator[string]) (err bool) { return false },
		)
		require.True(t, err.IsErr())
	})
	t.Run("ScanBytes", func(t *testing.T) {
		err := jscan.Scan(
			input,
			func(i *jscan.Iterator[[]byte]) (err bool) { return false },
		)
		require.True(t, err.IsErr())
	})
}

type Record struct {
	Level      int
	ValueType  jscan.ValueType
	Key        string
	Value      string
	ArrayIndex int
	Pointer    string
}

func TestScan(t *testing.T) {
	for _, tt := range []struct {
		name   string
		input  string
		expect []Record
	}{
		{
			name:  "null",
			input: "null",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeNull,
					ArrayIndex: -1,
					Value:      "null",
				},
			},
		},
		{
			name:  "bool_true",
			input: "true",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeTrue,
					ArrayIndex: -1,
					Value:      "true",
				},
			},
		},
		{
			name:  "bool_false",
			input: "false",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeFalse,
					ArrayIndex: -1,
					Value:      "false",
				},
			},
		},
		{
			name:  "number_int",
			input: "42",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeNumber,
					Value:      "42",
					ArrayIndex: -1,
				},
			},
		},
		{
			name:  "number_decimal",
			input: "42.5",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeNumber,
					Value:      "42.5",
					ArrayIndex: -1,
				},
			},
		},
		{
			name:  "number_negative",
			input: "-42.5",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeNumber,
					Value:      "-42.5",
					ArrayIndex: -1,
				},
			},
		},
		{
			name:  "number_exponent",
			input: "2.99792458e8",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeNumber,
					Value:      "2.99792458e8",
					ArrayIndex: -1,
				},
			},
		},
		{
			name:  "string",
			input: `"42"`,
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeString,
					Value:      `"42"`,
					ArrayIndex: -1,
				},
			},
		},
		{
			name:  "escaped unicode string",
			input: `"жш\"ц\\\\\""`,
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeString,
					Value:      `"жш\"ц\\\\\""`,
					ArrayIndex: -1,
				},
			},
		},
		{
			name:  "empty array",
			input: "[]",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeArray,
					ArrayIndex: -1,
				},
			},
		},
		{
			name:  "empty object",
			input: "{}",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeObject,
					ArrayIndex: -1,
				},
			},
		},
		{
			name:  "nested array",
			input: `[[null,[{"key":true}]],[]]`,
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeArray,
					ArrayIndex: -1,
				},
				{
					ValueType:  jscan.ValueTypeArray,
					Level:      1,
					ArrayIndex: 0,
					Pointer:    "/0",
				},
				{
					ValueType:  jscan.ValueTypeNull,
					Level:      2,
					Value:      "null",
					ArrayIndex: 0,
					Pointer:    "/0/0",
				},
				{
					ValueType:  jscan.ValueTypeArray,
					Level:      2,
					ArrayIndex: 1,
					Pointer:    "/0/1",
				},
				{
					ValueType:  jscan.ValueTypeObject,
					Level:      3,
					ArrayIndex: 0,
					Pointer:    "/0/1/0",
				},
				{
					ValueType:  jscan.ValueTypeTrue,
					Key:        `"key"`,
					Value:      "true",
					Level:      4,
					ArrayIndex: -1,
					Pointer:    "/0/1/0/key",
				},
				{
					ValueType:  jscan.ValueTypeArray,
					Level:      1,
					ArrayIndex: 1,
					Pointer:    "/1",
				},
			},
		},
		{
			name:  "escaped pointer",
			input: `{"/":[{"~":null},0]}`,
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeObject,
					ArrayIndex: -1,
				},
				{
					ValueType:  jscan.ValueTypeArray,
					Level:      1,
					Key:        `"/"`,
					ArrayIndex: -1,
					Pointer:    `/~1`,
				},
				{
					ValueType:  jscan.ValueTypeObject,
					Level:      2,
					ArrayIndex: 0,
					Pointer:    `/~1/0`,
				},
				{
					ValueType:  jscan.ValueTypeNull,
					Key:        `"~"`,
					Value:      "null",
					Level:      3,
					ArrayIndex: -1,
					Pointer:    `/~1/0/~0`,
				},
				{
					ValueType:  jscan.ValueTypeNumber,
					Value:      "0",
					Level:      2,
					ArrayIndex: 1,
					Pointer:    `/~1/1`,
				},
			},
		},
		{
			name: "nested object",
			input: `{
				"s": "value",
				"t": true,
				"f": false,
				"0": null,
				"n": -9.123e3,
				"o0": {},
				"a0": [],
				"o": {
					"k": "\"v\"",
					"a": [
						true,
						false,
						null,
						"item",
						-67.02e9,
						["foo"]
					]
				},
				"a3": [ 0, {"a3": 8} ]
			}`,
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeObject,
					ArrayIndex: -1,
				},
				{
					Level:      1,
					ValueType:  jscan.ValueTypeString,
					Key:        `"s"`,
					Value:      `"value"`,
					ArrayIndex: -1,
					Pointer:    "/s",
				},
				{
					Level:      1,
					ValueType:  jscan.ValueTypeTrue,
					Key:        `"t"`,
					Value:      "true",
					ArrayIndex: -1,
					Pointer:    "/t",
				},
				{
					Level:      1,
					ValueType:  jscan.ValueTypeFalse,
					Key:        `"f"`,
					Value:      "false",
					ArrayIndex: -1,
					Pointer:    "/f",
				},
				{
					Level:      1,
					ValueType:  jscan.ValueTypeNull,
					Key:        `"0"`,
					Value:      "null",
					ArrayIndex: -1,
					Pointer:    "/0",
				},
				{
					Level:      1,
					ValueType:  jscan.ValueTypeNumber,
					Key:        `"n"`,
					Value:      "-9.123e3",
					ArrayIndex: -1,
					Pointer:    "/n",
				},
				{
					Level:      1,
					ValueType:  jscan.ValueTypeObject,
					Key:        `"o0"`,
					ArrayIndex: -1,
					Pointer:    "/o0",
				},
				{
					Level:      1,
					ValueType:  jscan.ValueTypeArray,
					Key:        `"a0"`,
					ArrayIndex: -1,
					Pointer:    "/a0",
				},
				{
					Level:      1,
					ValueType:  jscan.ValueTypeObject,
					Key:        `"o"`,
					ArrayIndex: -1,
					Pointer:    "/o",
				},
				{
					Level:      2,
					ValueType:  jscan.ValueTypeString,
					Key:        `"k"`,
					Value:      `"\"v\""`,
					ArrayIndex: -1,
					Pointer:    "/o/k",
				},
				{
					Level:      2,
					ValueType:  jscan.ValueTypeArray,
					Key:        `"a"`,
					ArrayIndex: -1,
					Pointer:    "/o/a",
				},
				{
					Level:      3,
					ValueType:  jscan.ValueTypeTrue,
					Value:      "true",
					ArrayIndex: 0,
					Pointer:    "/o/a/0",
				},
				{
					Level:      3,
					ValueType:  jscan.ValueTypeFalse,
					Value:      "false",
					ArrayIndex: 1,
					Pointer:    "/o/a/1",
				},
				{
					Level:      3,
					ValueType:  jscan.ValueTypeNull,
					Value:      "null",
					ArrayIndex: 2,
					Pointer:    "/o/a/2",
				},
				{
					Level:      3,
					ValueType:  jscan.ValueTypeString,
					Value:      `"item"`,
					ArrayIndex: 3,
					Pointer:    "/o/a/3",
				},
				{
					Level:      3,
					ValueType:  jscan.ValueTypeNumber,
					Value:      "-67.02e9",
					ArrayIndex: 4,
					Pointer:    "/o/a/4",
				},
				{
					Level:      3,
					ValueType:  jscan.ValueTypeArray,
					ArrayIndex: 5,
					Pointer:    "/o/a/5",
				},
				{
					Level:      4,
					ValueType:  jscan.ValueTypeString,
					Value:      `"foo"`,
					ArrayIndex: 0,
					Pointer:    "/o/a/5/0",
				},
				{
					Level:      1,
					ValueType:  jscan.ValueTypeArray,
					Key:        `"a3"`,
					ArrayIndex: -1,
					Pointer:    "/a3",
				},
				{
					Level:      2,
					ValueType:  jscan.ValueTypeNumber,
					Value:      "0",
					ArrayIndex: 0,
					Pointer:    "/a3/0",
				},
				{
					Level:      2,
					ValueType:  jscan.ValueTypeObject,
					ArrayIndex: 1,
					Pointer:    "/a3/1",
				},
				{
					Level:      3,
					ValueType:  jscan.ValueTypeNumber,
					Key:        `"a3"`,
					Value:      "8",
					ArrayIndex: -1,
					Pointer:    "/a3/1/a3",
				},
			},
		},
		{
			name:  "trailing space",
			input: "null ",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeNull,
					ArrayIndex: -1,
					Value:      "null",
				},
			},
		},
		{
			name:  "trailing carriage return",
			input: "null\r",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeNull,
					ArrayIndex: -1,
					Value:      "null",
				},
			},
		},
		{
			name:  "trailing tab",
			input: "null\t",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeNull,
					ArrayIndex: -1,
					Value:      "null",
				},
			},
		},
		{
			name:  "trailing line-break",
			input: "null\n",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeNull,
					ArrayIndex: -1,
					Value:      "null",
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, json.Valid([]byte(tt.input)))

			t.Run("string", func(t *testing.T) {
				j := 0
				check := func(t *testing.T) func(i *jscan.Iterator[string]) bool {
					q := require.New(t)
					return func(i *jscan.Iterator[string]) bool {
						if j >= len(tt.expect) {
							t.Errorf("unexpected value at %d", j)
							j++
							return false
						}
						e := tt.expect[j]
						q.Equal(
							e.ValueType, i.ValueType(),
							"ValueType at %d", i.ValueIndex(),
						)
						q.Equal(
							e.Level, i.Level(),
							"Level at %d", i.ValueIndex(),
						)
						q.Equal(
							e.Value, i.Value(),
							"Value at %d", i.ValueIndex(),
						)
						q.Equal(
							e.Key, i.Key(),
							"Key at %d", i.ValueIndex(),
						)
						q.Equal(
							e.ArrayIndex, i.ArrayIndex(),
							"ArrayIndex at %d", i.ValueIndex(),
						)
						q.Equal(
							e.Pointer, i.Pointer(),
							"Pointer at %d", i.ValueIndex(),
						)
						j++
						return false
					}
				}

				t.Run("valid", func(t *testing.T) {
					require.True(t, jscan.Valid(tt.input))
				})
				t.Run("Scan", func(t *testing.T) {
					j = 0
					err := jscan.Scan(tt.input, check(t))
					require.False(t, err.IsErr(), "unexpected error: %s", err)
				})
			})

			t.Run("bytes", func(t *testing.T) {
				j := 0
				check := func(t *testing.T) func(i *jscan.Iterator[[]byte]) bool {
					q := require.New(t)
					return func(i *jscan.Iterator[[]byte]) bool {
						if j >= len(tt.expect) {
							t.Errorf("unexpected value at %d", j)
							j++
							return false
						}
						e := tt.expect[j]
						q.Equal(
							e.ValueType, i.ValueType(),
							"ValueType at %d", i.ValueIndex(),
						)
						q.Equal(
							e.Level, i.Level(),
							"Level at %d", i.ValueIndex(),
						)
						q.Equal(
							e.Value, string(i.Value()),
							"Value at %d", i.ValueIndex(),
						)
						q.Equal(
							e.Key, string(i.Key()),
							"Key at %d", i.ValueIndex(),
						)
						q.Equal(
							e.ArrayIndex, i.ArrayIndex(),
							"ArrayIndex at %d", i.ValueIndex(),
						)
						q.Equal(
							e.Pointer, string(i.Pointer()),
							"Pointer at %d", i.ValueIndex(),
						)
						j++
						return false
					}
				}

				t.Run("valid", func(t *testing.T) {
					require.True(t, jscan.Valid([]byte(tt.input)))
				})
				t.Run("Scan", func(t *testing.T) {
					j = 0
					err := jscan.Scan([]byte(tt.input), check(t))
					require.False(t, err.IsErr(), "unexpected error: %s", err)
				})
			})
		})
	}
}

// func TestScanError(t *testing.T) {
// 	for _, tt := range []struct {
// 		name   string
// 		input  string
// 		expect string
// 	}{
// 		{
// 			name:   "empty input",
// 			input:  "",
// 			expect: `error at index 0: unexpected EOF`,
// 		},
// 		{
// 			name:   "empty input",
// 			input:  " ",
// 			expect: `error at index 1: unexpected EOF`,
// 		},
// 		{
// 			name:   "empty input",
// 			input:  "\t\r\n ",
// 			expect: `error at index 4: unexpected EOF`,
// 		},
// 		{
// 			name:   "invalid literal",
// 			input:  "nul",
// 			expect: `error at index 0 ('n'): unexpected token`,
// 		},
// 		{
// 			name:   "invalid literal",
// 			input:  "fals",
// 			expect: `error at index 0 ('f'): unexpected token`,
// 		},
// 		{
// 			name:   "invalid literal",
// 			input:  "tru",
// 			expect: `error at index 0 ('t'): unexpected token`,
// 		},
// 		{
// 			name:   "invalid negative number",
// 			input:  "-",
// 			expect: `error at index 0 ('-'): malformed number`,
// 		},
// 		{
// 			name:   "invalid number fraction",
// 			input:  "0.",
// 			expect: `error at index 0 ('0'): malformed number`,
// 		},
// 		{
// 			name:   "invalid number exponent",
// 			input:  "0e",
// 			expect: `error at index 0 ('0'): malformed number`,
// 		},
// 		{
// 			name:   "invalid number exponent",
// 			input:  "1e-",
// 			expect: `error at index 0 ('1'): malformed number`,
// 		},
// 		{
// 			name:   "invalid number integer",
// 			input:  "e1",
// 			expect: `error at index 0 ('e'): unexpected token`,
// 		},
// 		{
// 			name:   "invalid escape sequence in string",
// 			input:  `"\0"`,
// 			expect: `error at index 2 ('0'): invalid escape sequence`,
// 		},
// 		{
// 			name:   "missing closing }",
// 			input:  `{"x":null`,
// 			expect: `error at index 9: unexpected EOF`,
// 		},
// 		{
// 			name:   "missing closing }",
// 			input:  `{"x":{`,
// 			expect: `error at index 6: unexpected EOF`,
// 		},
// 		{
// 			name:   "missing closing ]",
// 			input:  `[null`,
// 			expect: `error at index 5: unexpected EOF`,
// 		},
// 		{
// 			name:   "missing closing ]",
// 			input:  `[[null`,
// 			expect: `error at index 6: unexpected EOF`,
// 		},
// 		{
// 			name:   `missing closing quotes`,
// 			input:  `"string`,
// 			expect: `error at index 7: unexpected EOF`,
// 		},
// 		{
// 			name:   `missing closing quotes after escaped quotes`,
// 			input:  `"string\"`,
// 			expect: `error at index 9: unexpected EOF`,
// 		},
// 		{
// 			name:   `missing closing quotes after escaped sequences`,
// 			input:  `"string\\\"`,
// 			expect: `error at index 11: unexpected EOF`,
// 		},
// 		{
// 			name:   `unfinished key`,
// 			input:  `{"key`,
// 			expect: `error at index 5: unexpected EOF`,
// 		},
// 		{
// 			name:   `missing column`,
// 			input:  `{"key"}`,
// 			expect: `error at index 6 ('}'): unexpected token`,
// 		},
// 		{
// 			name:   `invalid content before column`,
// 			input:  `{"key"1 :}`,
// 			expect: `error at index 6 ('1'): unexpected token`,
// 		},
// 		{
// 			name:   `invalid column`,
// 			input:  `{"key";1}`,
// 			expect: `error at index 6 (';'): unexpected token`,
// 		},
// 		{
// 			name:   `missing field value`,
// 			input:  `{"okay":}`,
// 			expect: `error at index 8 ('}'): unexpected token`,
// 		},
// 		{
// 			name:   "unexpected object",
// 			input:  `{"key":12,{}}`,
// 			expect: `error at index 10 ('{'): unexpected token`,
// 		},
// 		{
// 			name:   `missing array item`,
// 			input:  `["okay",]`,
// 			expect: `error at index 8 (']'): unexpected token`,
// 		},
// 		{
// 			name:   `missing comma`,
// 			input:  `["okay"[`,
// 			expect: `error at index 7 ('['): unexpected token`,
// 		},
// 		{
// 			name:   `missing comma`,
// 			input:  `["okay"-12`,
// 			expect: `error at index 7 ('-'): unexpected token`,
// 		},
// 		{
// 			name:   `missing comma`,
// 			input:  `["okay"0`,
// 			expect: `error at index 7 ('0'): unexpected token`,
// 		},
// 		{
// 			name:   `missing comma`,
// 			input:  `["okay""not okay"]`,
// 			expect: `error at index 7 ('"'): unexpected token`,
// 		},
// 		{
// 			name:   `missing comma`,
// 			input:  `{"foo":"bar" "baz":"fuz"}`,
// 			expect: `error at index 13 ('"'): unexpected token`,
// 		},
// 		{
// 			name:   `missing comma`,
// 			input:  `[null false]`,
// 			expect: `error at index 6 ('f'): unexpected token`,
// 		},
// 		{
// 			name:   "expect EOF after number zero",
// 			input:  "01",
// 			expect: `error at index 1 ('1'): unexpected token`,
// 		},
// 		{
// 			name:   `expect EOF after negative number zero`,
// 			input:  `-00`,
// 			expect: `error at index 2 ('0'): unexpected token`,
// 		},
// 		{
// 			name:   `expect EOF after string get comma`,
// 			input:  `"okay",null`,
// 			expect: `error at index 6 (','): unexpected token`,
// 		},
// 		{
// 			name:   `expect EOF after string`,
// 			input:  `"str" "str"`,
// 			expect: `error at index 6 ('"'): unexpected token`,
// 		},
// 		{
// 			name:   `expect EOF after number`,
// 			input:  `0 0`,
// 			expect: `error at index 2 ('0'): unexpected token`,
// 		},
// 		{
// 			name:   `expect EOF after false`,
// 			input:  `false false`,
// 			expect: `error at index 6 ('f'): unexpected token`,
// 		},
// 		{
// 			name:   `expect EOF after true`,
// 			input:  `true true`,
// 			expect: `error at index 5 ('t'): unexpected token`,
// 		},
// 		{
// 			name:   `expect EOF after null`,
// 			input:  `null null`,
// 			expect: `error at index 5 ('n'): unexpected token`,
// 		},
// 		{
// 			name:   `expect EOF after array`,
// 			input:  `[] []`,
// 			expect: `error at index 3 ('['): unexpected token`,
// 		},
// 		{
// 			name:   `expect EOF after object`,
// 			input:  `{"k":0} {"k":0}`,
// 			expect: `error at index 8 ('{'): unexpected token`,
// 		},
// 	} {
// 		require.False(t, json.Valid([]byte(tt.input)))

// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Run("string", func(t *testing.T) {
// 				check := func(t *testing.T) func(i *jscan.Iterator) bool {
// 					return func(i *jscan.Iterator) (err bool) { return false }
// 				}

// 				t.Run("valid", func(t *testing.T) {
// 					require.False(t, jscan.Valid(tt.input))
// 				})

// 				t.Run("cachepath", func(t *testing.T) {
// 					err := jscan.Scan(jscan.Options{
// 						CachePath: true,
// 					}, tt.input, check(t))
// 					require.Equal(t, tt.expect, err.Error())
// 					require.True(t, err.IsErr())
// 				})
// 				t.Run("nocachepath", func(t *testing.T) {
// 					err := jscan.Scan(jscan.Options{}, tt.input, check(t))
// 					require.Equal(t, tt.expect, err.Error())
// 					require.True(t, err.IsErr())
// 				})
// 			})

// 			t.Run("bytes", func(t *testing.T) {
// 				check := func(t *testing.T) func(i *jscan.IteratorBytes) bool {
// 					return func(i *jscan.IteratorBytes) (err bool) {
// 						return false
// 					}
// 				}

// 				t.Run("valid", func(t *testing.T) {
// 					require.False(t, jscan.ValidBytes([]byte(tt.input)))
// 				})

// 				t.Run("cachepath", func(t *testing.T) {
// 					err := jscan.ScanBytes(jscan.Options{
// 						CachePath: true,
// 					}, []byte(tt.input), check(t))
// 					require.Equal(t, tt.expect, err.Error())
// 					require.True(t, err.IsErr())
// 				})
// 				t.Run("nocachepath", func(t *testing.T) {
// 					err := jscan.ScanBytes(
// 						jscan.Options{}, []byte(tt.input), check(t),
// 					)
// 					require.Equal(t, tt.expect, err.Error())
// 					require.True(t, err.IsErr())
// 				})
// 			})
// 		})
// 	}
// }

func TestControlCharacters(t *testing.T) {
	test := func(t *testing.T, in, expectErr string) {
		t.Run("Validate", func(t *testing.T) {
			err := jscan.Validate(in)
			require.Equal(t, expectErr, err.Error())
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
		})

		t.Run("Validate_bytes", func(t *testing.T) {
			err := jscan.Validate([]byte(in))
			require.Equal(t, expectErr, err.Error(), "IN: %q", in)
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
		})

		t.Run("Validate", func(t *testing.T) {
			t.Run("string", func(t *testing.T) {
				err := jscan.Validate(in)
				require.Equal(t, expectErr, err.Error(), "IN: %q", in)
				require.True(t, err.IsErr())
				require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
			})
			t.Run("bytes", func(t *testing.T) {
				err := jscan.Validate([]byte(in))
				require.Equal(t, expectErr, err.Error(), "IN: %q", in)
				require.True(t, err.IsErr())
				require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
			})
		})

		t.Run("ValidateOne", func(t *testing.T) {
			t.Run("string", func(t *testing.T) {
				t.Logf("in:%q", in)
				for s := in; len(s) > 0; {
					trailing, err := jscan.ValidateOne(s)
					require.True(t, err.IsErr())
					require.Equal(t, expectErr, err.Error())
					require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
					require.Equal(t, err.Src[err.Index:], trailing)
					s = trailing
					if err.IsErr() {
						break
					}
				}
			})
			t.Run("bytes", func(t *testing.T) {
				for s := []byte(in); len(s) > 0; {
					trailing, err := jscan.ValidateOne(s)
					require.True(t, err.IsErr())
					require.Equal(t, expectErr, err.Error())
					require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
					require.Equal(t, err.Src[err.Index:], trailing)
					s = trailing
					if err.IsErr() {
						break
					}
				}
			})
		})

		t.Run("Valid", func(t *testing.T) {
			require.False(t, jscan.Valid(in))
		})

		t.Run("Valid_bytes", func(t *testing.T) {
			require.False(t, jscan.Valid([]byte(in)))
		})

		t.Run("ScanOne", func(t *testing.T) {
			_, err := jscan.ScanOne(
				in,
				func(i *jscan.Iterator[string]) (err bool) { return false },
			)
			require.Equal(t, expectErr, err.Error())
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
		})

		t.Run("ScanOne_bytes", func(t *testing.T) {
			_, err := jscan.ScanOne(
				[]byte(in),
				func(i *jscan.Iterator[[]byte]) (err bool) { return false },
			)
			require.Equal(t, expectErr, err.Error())
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
		})

		t.Run("Scan", func(t *testing.T) {
			err := jscan.Scan(
				in,
				func(i *jscan.Iterator[string]) (err bool) { return false },
			)
			require.Equal(t, expectErr, err.Error())
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
		})

		t.Run("Scan_bytes", func(t *testing.T) {
			err := jscan.Scan(
				[]byte(in),
				func(i *jscan.Iterator[[]byte]) (err bool) { return false },
			)
			require.Equal(t, expectErr, err.Error())
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
		})
	}

	// Control characters in a string value
	t.Run("value_string", func(t *testing.T) {
		ForASCIIControlChars(t, func(t *testing.T, b byte) {
			var buf bytes.Buffer
			buf.WriteByte('"')
			buf.WriteByte(b)
			buf.WriteByte('"')
			test(t, buf.String(), fmt.Sprintf(
				"error at index 1 (0x%x): illegal control character", b,
			))
		})
	})

	// Control characters in a string field value
	t.Run("field_string", func(t *testing.T) {
		ForASCIIControlChars(t, func(t *testing.T, b byte) {
			var buf bytes.Buffer
			buf.WriteString(`{"x":"`)
			buf.WriteByte(b)
			buf.WriteString(`"}`)
			test(t, buf.String(), fmt.Sprintf(
				"error at index 6 (0x%x): illegal control character", b,
			))
		})
	})

	// Control characters in an array item string
	t.Run("array_item_string", func(t *testing.T) {
		ForASCIIControlChars(t, func(t *testing.T, b byte) {
			var buf bytes.Buffer
			buf.WriteString(`["`)
			buf.WriteByte(b)
			buf.WriteString(`"]`)
			test(t, buf.String(), fmt.Sprintf(
				"error at index 2 (0x%x): illegal control character", b,
			))
		})
	})

	// Control characters in head
	t.Run("head", func(t *testing.T) {
		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			var buf bytes.Buffer
			buf.WriteByte('\t')
			buf.WriteByte(b)
			buf.WriteString(`[]`)
			test(t, buf.String(), fmt.Sprintf(
				"error at index 1 (0x%x): illegal control character", b,
			))
		})
	})

	// Control characters outside of strings
	t.Run("nonstring", func(t *testing.T) {
		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			var buf bytes.Buffer
			buf.WriteByte(b)
			buf.WriteString("null")
			test(t, buf.String(), fmt.Sprintf(
				"error at index 0 (0x%x): illegal control character", b,
			))
		})

		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			var buf bytes.Buffer
			buf.WriteString("[")
			buf.WriteByte(b)
			buf.WriteString("]")
			test(t, buf.String(), fmt.Sprintf(
				"error at index 1 (0x%x): illegal control character", b,
			))
		})

		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			var buf bytes.Buffer
			buf.WriteString("[\n")
			buf.WriteByte(b)
			buf.WriteString("]")
			test(t, buf.String(), fmt.Sprintf(
				"error at index 2 (0x%x): illegal control character", b,
			))
		})

		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			var buf bytes.Buffer
			buf.WriteString("{")
			buf.WriteByte(b)
			buf.WriteString("}")
			test(t, buf.String(), fmt.Sprintf(
				"error at index 1 (0x%x): illegal control character", b,
			))
		})

		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			var buf bytes.Buffer
			buf.WriteString(`{"foo":`)
			buf.WriteByte(b)
			buf.WriteString(`false}`)
			test(t, buf.String(), fmt.Sprintf(
				"error at index 7 (0x%x): illegal control character", b,
			))
		})
	})
}

func TestReturnErrorTrue(t *testing.T) {
	input := `{"foo":"bar","baz":null}`

	t.Run("string", func(t *testing.T) {
		j := 0
		err := jscan.Scan(
			input,
			func(i *jscan.Iterator[string]) (err bool) {
				require.Equal(t, jscan.ValueTypeObject, i.ValueType())
				j++
				return true // Expect immediate return
			},
		)
		require.Equal(t, 1, j)
		require.True(t, err.IsErr())
		require.Equal(t, jscan.ErrorCodeCallback, err.Code)
		require.Equal(
			t, "error at index 0 ('{'): callback error", err.Error(),
		)
	})

	t.Run("bytes", func(t *testing.T) {
		j := 0
		err := jscan.Scan(
			[]byte(input),
			func(i *jscan.Iterator[[]byte]) (err bool) {
				require.Equal(t, jscan.ValueTypeObject, i.ValueType())
				j++
				return true // Expect immediate return
			},
		)
		require.Equal(t, 1, j)
		require.True(t, err.IsErr())
		require.Equal(t, jscan.ErrorCodeCallback, err.Code)
		require.Equal(
			t, "error at index 0 ('{'): callback error", err.Error(),
		)
	})
}

// func TestGet(t *testing.T) {
// 	for _, tt := range []struct {
// 		json       string
// 		path       string
// 		escapePath bool
// 		expect     Record
// 	}{
// 		{`{"key":null}`, "key", true, Record{
// 			Level:      1,
// 			ValueType:  jscan.ValueTypeNull,
// 			Key:        "key",
// 			Value:      "null",
// 			ArrayIndex: -1,
// 			Path:       "key",
// 		}},
// 		{`{"foo.bar":false}`, `foo\.bar`, true, Record{
// 			Level:      1,
// 			ValueType:  jscan.ValueTypeFalse,
// 			Key:        `foo.bar`,
// 			Value:      "false",
// 			ArrayIndex: -1,
// 			Path:       `foo\.bar`,
// 		}},
// 		{`{"foo.bar":false}`, `foo.bar`, false, Record{
// 			Level:      1,
// 			ValueType:  jscan.ValueTypeFalse,
// 			Key:        `foo.bar`,
// 			Value:      "false",
// 			ArrayIndex: -1,
// 			Path:       `foo.bar`,
// 		}},
// 		{`[true]`, `[0]`, true, Record{
// 			Level:      1,
// 			ValueType:  jscan.ValueTypeTrue,
// 			Value:      "true",
// 			ArrayIndex: 0,
// 			Path:       `[0]`,
// 		}},
// 		{`[false,[[2, 42]]]`, `[1][0][1]`, true, Record{
// 			Level:      3,
// 			ValueType:  jscan.ValueTypeNumber,
// 			Value:      "42",
// 			ArrayIndex: 1,
// 			Path:       `[1][0][1]`,
// 		}},
// 		{
// 			`[false,[[2, {"[foo]":[{"bar-baz":"fuz"}]}]]]`,
// 			`[1][0][1].\[foo\][0].bar-baz`, true,
// 			Record{
// 				Level:      6,
// 				ValueType:  jscan.ValueTypeString,
// 				Key:        "bar-baz",
// 				Value:      "fuz",
// 				ArrayIndex: -1,
// 				Path:       `[1][0][1].\[foo\][0].bar-baz`,
// 			},
// 		},
// 	} {
// 		t.Run("", func(t *testing.T) {
// 			t.Run("string", func(t *testing.T) {
// 				c := 0
// 				err := jscan.Get(
// 					tt.json, tt.path, tt.escapePath,
// 					func(i *jscan.Iterator) {
// 						c++
// 						require.Equal(t, tt.expect, Record{
// 							Level:      i.Level,
// 							ValueType:  i.ValueType,
// 							Key:        i.Key(),
// 							Value:      i.Value(),
// 							ArrayIndex: i.ArrayIndex,
// 							Path:       i.Path(),
// 						})
// 					},
// 				)
// 				require.False(t, err.IsErr(), "unexpected error: %s", err)
// 				require.Equal(t, 1, c)
// 			})

// 			t.Run("bytes", func(t *testing.T) {
// 				c := 0
// 				err := jscan.GetBytes(
// 					[]byte(tt.json), []byte(tt.path), tt.escapePath,
// 					func(i *jscan.IteratorBytes) {
// 						c++
// 						require.Equal(t, tt.expect, Record{
// 							Level:      i.Level,
// 							ValueType:  i.ValueType,
// 							Key:        string(i.Key()),
// 							Value:      string(i.Value()),
// 							ArrayIndex: i.ArrayIndex,
// 							Path:       string(i.Path()),
// 						})
// 					},
// 				)
// 				require.False(t, err.IsErr(), "unexpected error: %s", err)
// 				require.Equal(t, 1, c)
// 			})
// 		})
// 	}
// }

// func TestGetNotFound(t *testing.T) {
// 	for _, tt := range []struct {
// 		json       string
// 		path       string
// 		escapePath bool
// 	}{
// 		{`{"key":null}`, "non-existing-key", true},
// 		{`{"foo.bar":false}`, `foo.bar`, true},
// 		{`{"foo.bar":false}`, `foo\.bar`, false},
// 		{`[true]`, `[1]`, true},
// 		{`[false,[[2, 42]]]`, `[1][0][2]`, true},
// 		{
// 			`[false,[[2, {"[foo]":[{"bar-baz":"fuz"}]}]]]`,
// 			`[1][0][1].\[foo\][0].bar-`, true,
// 		},
// 	} {
// 		t.Run("", func(t *testing.T) {
// 			t.Run("string", func(t *testing.T) {
// 				c := 0
// 				err := jscan.Get(
// 					tt.json, tt.path, tt.escapePath,
// 					func(i *jscan.Iterator) {
// 						c++
// 					},
// 				)
// 				require.False(t, err.IsErr(), "unexpected error: %s", err)
// 				require.Zero(t, c, "unexpected call")
// 			})

// 			t.Run("bytes", func(t *testing.T) {
// 				c := 0
// 				err := jscan.GetBytes(
// 					[]byte(tt.json), []byte(tt.path), tt.escapePath,
// 					func(i *jscan.IteratorBytes) {
// 						c++
// 					},
// 				)
// 				require.False(t, err.IsErr(), "unexpected error: %s", err)
// 				require.Zero(t, c, "unexpected call")
// 			})
// 		})
// 	}
// }

func TestScanOne(t *testing.T) {
	inputs := []string{
		`-120.4`,
		`"string"`,
		`{"key":"value"}`,
		`[0,1]`,
		`true`,
		`false`,
		`null`,
	}
	s := strings.Join(inputs, "")

	type V struct {
		Type  jscan.ValueType
		Value string
	}
	expect := []V{
		{jscan.ValueTypeNumber, "-120.4"},
		{jscan.ValueTypeString, `"string"`},
		{jscan.ValueTypeObject, ""},
		{jscan.ValueTypeString, `"value"`},
		{jscan.ValueTypeArray, ""},
		{jscan.ValueTypeNumber, "0"},
		{jscan.ValueTypeNumber, "1"},
		{jscan.ValueTypeTrue, "true"},
		{jscan.ValueTypeFalse, "false"},
		{jscan.ValueTypeNull, "null"},
	}

	t.Run("string", func(t *testing.T) {
		for i, c, s := 0, 1, s; s != ""; c++ {
			trailing, err := jscan.ScanOne(
				s,
				func(itr *jscan.Iterator[string]) (err bool) {
					require.Equal(t, expect[i], V{
						Type:  itr.ValueType(),
						Value: string(itr.Value()),
					}, "unexpected value at index %d", i)
					i++
					return false
				},
			)
			require.False(t, err.IsErr())
			require.Equal(t, trailing, strings.Join(inputs[c:], ""))
			s = trailing
		}
	})

	t.Run("bytes", func(t *testing.T) {
		s := []byte(s)
		for i, c, s := 0, 1, s; len(s) > 0; c++ {
			trailing, err := jscan.ScanOne(
				s,
				func(itr *jscan.Iterator[[]byte]) (err bool) {
					require.Equal(t, expect[i], V{
						Type:  itr.ValueType(),
						Value: string(itr.Value()),
					}, "unexpected value at index %d", i)
					i++
					return false
				},
			)
			require.False(t, err.IsErr())
			require.Equal(t, string(trailing), strings.Join(inputs[c:], ""))
			s = trailing
		}
	})
}

// ForASCIIControlChars calls fn with each possible ASCII control character
func ForASCIIControlChars(t *testing.T, fn func(t *testing.T, b byte)) {
	for b := byte(0); b < 0x20; b++ {
		t.Run(fmt.Sprintf("%U", b), func(t *testing.T) {
			fn(t, b)
		})
	}
}

// ForASCIIControlCharsExceptTRN calls fn with each possible ASCII
// control character except '\t', '\r', '\n'
func ForASCIIControlCharsExceptTRN(t *testing.T, fn func(t *testing.T, b byte)) {
	for b := byte(0); b < 0x20; b++ {
		switch b {
		case '\t', '\r', '\n':
			continue
		}
		t.Run(fmt.Sprintf("%U", b), func(t *testing.T) {
			fn(t, b)
		})
	}
}

func TestStrings(t *testing.T) {
	for _, td := range []string{
		`""`,
		`"a\""`,
		`"ab\""`,
		`"abc\""`,
		`"abcd\""`,
		`"abcde\""`,
		`"abcdef\""`,
		`"abcdefg\""`,
		`"abcdefgh\""`,
		`"abcdefghi\""`,
		`"abcdefghik\""`,
		`"abcdefghikl\""`,
		`"abcdefghiklm\""`,
		`"abcdefghiklmn\""`,
		`"abcdefghiklmno\""`,
		`"abcdefghiklmnop\""`,
		`"abcdefghiklmnopq\""`,
		`"abcdefghiklmnopqr\""`,
		`"abcdefghiklmnopqrs\""`,
		`"abcdefghiklmnopqrst\""`,
		`"abcdefghiklmnopqrstu\""`,
		`"abcdefghiklmnopqrstuv\""`,
		`"abcdefghiklmnopqrstuvw\""`,
		`"abcdefghiklmnopqrstuvwx\""`,
		`"abcdefghiklmnopqrstuvwxy\""`,
		`"abcdefghiklmnopqrstuvwxyz\""`,
		`"a1234567812345678\""`,
		`"ab1234567812345678\""`,
		`"abc1234567812345678\""`,
		`"abcd1234567812345678\""`,
		`"abcde1234567812345678\""`,
		`"abcdef1234567812345678\""`,
		`"abcdefg1234567812345678\""`,
		`"abcdefgh1234567812345678\""`,
		`"abcdefghi1234567812345678\""`,
		`"abcdefghik1234567812345678\""`,
		`"abcdefghikl1234567812345678\""`,
		`"abcdefghiklm1234567812345678\""`,
		`"abcdefghiklmn1234567812345678\""`,
		`"abcdefghiklmno1234567812345678\""`,
		`"abcdefghiklmnop1234567812345678\""`,
		`"abcdefghiklmnopq1234567812345678\""`,
		`"abcdefghiklmnopqr1234567812345678\""`,
		`"abcdefghiklmnopqrs1234567812345678\""`,
		`"abcdefghiklmnopqrst1234567812345678\""`,
		`"abcdefghiklmnopqrstu1234567812345678\""`,
		`"abcdefghiklmnopqrstuv1234567812345678\""`,
		`"abcdefghiklmnopqrstuvw1234567812345678\""`,
		`"abcdefghiklmnopqrstuvwx1234567812345678\""`,
		`"abcdefghiklmnopqrstuvwxy1234567812345678\""`,
		`"abcdefghiklmnopqrstuvwxyz1234567812345678\""`,
	} {
		t.Run("ParserScan", func(t *testing.T) {
			p := jscan.NewParser[string](64)
			err := p.Scan(td, func(i *jscan.Iterator[string]) (err bool) {
				require.Equal(t, jscan.ValueTypeString, i.ValueType())
				require.Equal(t, 0, i.Level())
				require.Equal(t, 0, i.ValueIndex())
				require.Equal(t, len(td), i.ValueIndexEnd())
				require.Zero(t, i.Pointer())
				return false
			})
			require.False(t, err.IsErr())
		})
		t.Run("ParserScan_bytes", func(t *testing.T) {
			p := jscan.NewParser[[]byte](64)
			err := p.Scan([]byte(td), func(i *jscan.Iterator[[]byte]) (err bool) {
				require.Equal(t, jscan.ValueTypeString, i.ValueType())
				require.Equal(t, 0, i.Level())
				require.Equal(t, 0, i.ValueIndex())
				require.Equal(t, len(td), i.ValueIndexEnd())
				require.Equal(t, []byte(""), i.Pointer())
				return false
			})
			require.False(t, err.IsErr())
		})
		t.Run("Scan", func(t *testing.T) {
			err := jscan.Scan(td, func(i *jscan.Iterator[string]) (err bool) {
				require.Equal(t, jscan.ValueTypeString, i.ValueType())
				require.Equal(t, 0, i.Level())
				require.Equal(t, 0, i.ValueIndex())
				require.Equal(t, len(td), i.ValueIndexEnd())
				require.Zero(t, i.Pointer())
				return false
			})
			require.False(t, err.IsErr())
		})
		t.Run("Scan_bytes", func(t *testing.T) {
			err := jscan.Scan([]byte(td), func(i *jscan.Iterator[[]byte]) (err bool) {
				require.Equal(t, jscan.ValueTypeString, i.ValueType())
				require.Equal(t, 0, i.Level())
				require.Equal(t, 0, i.ValueIndex())
				require.Equal(t, len(td), i.ValueIndexEnd())
				require.Equal(t, []byte(""), i.Pointer())
				return false
			})
			require.False(t, err.IsErr())
		})
	}
}
