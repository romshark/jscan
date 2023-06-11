package jscan_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/romshark/jscan/v2"
	"github.com/stretchr/testify/require"
)

func TestJSONTestSuite(t *testing.T) {
	fileContents := func(t *testing.T, name string) []byte {
		c, err := os.ReadFile(filepath.Join("testdata/jsontestsuite", name))
		require.NoError(t, err)
		return c
	}

	d, err := os.ReadDir("testdata/jsontestsuite")
	require.NoError(t, err)
	for _, f := range d {
		n := f.Name()
		t.Run(n, func(t *testing.T) {
			switch {
			case strings.HasPrefix(n, "i_"):
				c := fileContents(t, n)
				testOKOrErr(t, []byte(c))
				testOKOrErr(t, string(c))
			case strings.HasPrefix(n, "y_"):
				c := fileContents(t, n)
				testStrictOK(t, []byte(c))
				testStrictOK(t, string(c))
			case strings.HasPrefix(n, "n_"):
				c := fileContents(t, n)
				testStrictErr(t, []byte(c))
				testStrictErr(t, string(c))
			default:
				t.Skip(n)
			}
		})
	}
}

// testStrictOK runs tests with the "y_" prefix that parsers must accept.
func testStrictOK[S ~string | ~[]byte](t *testing.T, input S) {
	t.Run(testDataType(input), func(t *testing.T) {
		t.Run("Validator", func(t *testing.T) {
			t.Run("Valid", func(t *testing.T) {
				require.True(t, jscan.NewValidator[S](1024).Valid(input))
			})
			t.Run("Validate", func(t *testing.T) {
				err := jscan.NewValidator[S](1024).Validate(input)
				require.False(t, err.IsErr())
			})
			t.Run("ValidateOne", func(t *testing.T) {
				_, err := jscan.NewValidator[S](1024).ValidateOne(input)
				require.False(t, err.IsErr())
			})
		})
		t.Run("Valid", func(t *testing.T) {
			require.True(t, jscan.Valid[S](input))
		})
		t.Run("Scan", func(t *testing.T) {
			err := jscan.Scan(
				string(input),
				func(i *jscan.Iterator[string]) (err bool) { return false },
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
	})
}

// testOKOrErr runs tests with the "i_" prefix that
// parsers are free to accept or reject.
func testOKOrErr[S ~string | ~[]byte](t *testing.T, input S) {
	t.Run("Valid", func(t *testing.T) {
		if !jscan.Valid(string(input)) {
			t.Skip("allowed to fail")
		}
	})
	t.Run("Validate", func(t *testing.T) {
		err := jscan.Validate(string(input))
		if err.IsErr() {
			t.Skip("allowed to fail")
		}
	})
	t.Run("ValidatorValid", func(t *testing.T) {
		v := jscan.NewValidator[string](1024)
		if !v.Valid(string(input)) {
			t.Skip("allowed to fail")
		}
	})
	t.Run("ValidatorValidate", func(t *testing.T) {
		v := jscan.NewValidator[string](1024)
		err := v.Validate(string(input))
		if err.IsErr() {
			t.Skip("allowed to fail")
		}
	})
}

// testStrictErr runs tests with the "n_" prefix that parsers must reject.
// (*Validator).ValidateOne, (*Parser).ScanOne, jscan.ValidateOne and jscan.ScanOne
// aren't tested because https://github.com/nst/JSONTestSuite
// assumes EOF after a valid JSON value.
func testStrictErr[S ~string | ~[]byte](t *testing.T, input S) {
	t.Run(testDataType(input), func(t *testing.T) {
		t.Run("ParserScan", func(t *testing.T) {
			err := jscan.NewParser[S](1024).Scan(
				input, func(i *jscan.Iterator[S]) (err bool) { return false },
			)
			require.True(t, err.IsErr())
		})
		t.Run("Scan", func(t *testing.T) {
			err := jscan.Scan(
				input, func(i *jscan.Iterator[S]) (err bool) { return false },
			)
			require.True(t, err.IsErr())
		})
		t.Run("ValidatorValid", func(t *testing.T) {
			require.False(t, jscan.NewValidator[S](1024).Valid(input))
		})
		t.Run("ValidatorValidate", func(t *testing.T) {
			err := jscan.NewValidator[S](1024).Validate(input)
			require.True(t, err.IsErr())
		})
		t.Run("Valid", func(t *testing.T) {
			require.False(t, jscan.Valid[S](input))
		})
		t.Run("Validate", func(t *testing.T) {
			require.True(t, jscan.Validate[S](input).IsErr())
		})
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

type ScanTest struct {
	name   string
	input  string
	expect []Record
}

func TestScan(t *testing.T) {
	for _, td := range []ScanTest{
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
		t.Run(td.name, func(t *testing.T) {
			require.True(t, json.Valid([]byte(td.input)))
			testScan[string](t, td)
			testScan[[]byte](t, td)
		})
	}
}

func testScan[S ~string | ~[]byte](t *testing.T, td ScanTest) {
	t.Run(testDataType(td.input), func(t *testing.T) {
		j := 0
		check := func(t *testing.T) func(i *jscan.Iterator[S]) bool {
			q := require.New(t)
			return func(i *jscan.Iterator[S]) bool {
				if j >= len(td.expect) {
					t.Errorf("unexpected value at %d", j)
					j++
					return false
				}
				e := td.expect[j]
				q.Equal(e.ValueType, i.ValueType(), "ValueType at %d", i.ValueIndex())
				q.Equal(e.Level, i.Level(), "Level at %d", i.ValueIndex())
				if e.Value == "" {
					q.Len(S(e.Value), 0, "Value at %d", i.ValueIndex())
				} else {
					q.Equal(S(e.Value), i.Value(), "Value at %d", i.ValueIndex())
				}
				if e.Key == "" {
					q.Len(i.Key(), 0, "Key at %d", i.ValueIndex())
				} else {
					q.Equal(S(e.Key), i.Key(), "Key at %d", i.ValueIndex())
				}
				q.Equal(e.ArrayIndex, i.ArrayIndex(), "ArrayIndex at %d", i.ValueIndex())
				if e.Pointer == "" {
					q.Len(i.Pointer(), 0, "Pointer at %d", i.ValueIndex())
				} else {
					q.Equal(S(e.Pointer), i.Pointer(), "Pointer at %d", i.ValueIndex())
				}
				j++
				return false
			}
		}

		t.Run("Valid", func(t *testing.T) {
			require.True(t, jscan.Valid[S](S(td.input)))
		})
		t.Run("ValidatorValid", func(t *testing.T) {
			v := jscan.NewValidator[S](64)
			require.True(t, v.Valid(S(td.input)))
		})
		t.Run("Scan", func(t *testing.T) {
			j = 0
			err := jscan.Scan[S](S(td.input), check(t))
			require.False(t, err.IsErr(), "unexpected error: %s", err)
		})
		t.Run("ParserScan", func(t *testing.T) {
			j = 0
			p := jscan.NewParser[S](64)
			err := p.Scan(S(td.input), check(t))
			require.False(t, err.IsErr(), "unexpected error: %s", err)
		})
	})
}

type ErrorTest struct {
	name   string
	input  string
	expect string
}

func TestError(t *testing.T) {
	for _, td := range []ErrorTest{
		{
			name:   "empty input",
			input:  "",
			expect: `error at index 0: unexpected EOF`,
		},
		{
			name:   "empty input",
			input:  " ",
			expect: `error at index 1: unexpected EOF`,
		},
		{
			name:   "empty input",
			input:  "\t\r\n ",
			expect: `error at index 4: unexpected EOF`,
		},
		{
			name:   "invalid literal",
			input:  "nul",
			expect: `error at index 0 ('n'): unexpected token`,
		},
		{
			name:   "invalid literal",
			input:  "fals",
			expect: `error at index 0 ('f'): unexpected token`,
		},
		{
			name:   "invalid literal",
			input:  "tru",
			expect: `error at index 0 ('t'): unexpected token`,
		},
		{
			name:   "invalid negative number",
			input:  "-",
			expect: `error at index 0 ('-'): malformed number`,
		},
		{
			name:   "invalid number fraction",
			input:  "0.",
			expect: `error at index 0 ('0'): malformed number`,
		},
		{
			name:   "invalid number exponent",
			input:  "0e",
			expect: `error at index 0 ('0'): malformed number`,
		},
		{
			name:   "invalid number exponent",
			input:  "1e-",
			expect: `error at index 0 ('1'): malformed number`,
		},
		{
			name:   "invalid number integer",
			input:  "e1",
			expect: `error at index 0 ('e'): unexpected token`,
		},
		{
			name:   "invalid escape sequence in string",
			input:  `"\0"`,
			expect: `error at index 1 ('\'): invalid escape sequence`,
		},
		{
			name:   "missing closing }",
			input:  `{"x":null`,
			expect: `error at index 9: unexpected EOF`,
		},
		{
			name:   "missing closing }",
			input:  `{"x":{`,
			expect: `error at index 6: unexpected EOF`,
		},
		{
			name:   "missing closing ] after space",
			input:  `[ `,
			expect: `error at index 2: unexpected EOF`,
		},
		{
			name:   "missing closing ] after null",
			input:  `[null`,
			expect: `error at index 5: unexpected EOF`,
		},
		{
			name:   "missing closing ] after null in array",
			input:  `[[null`,
			expect: `error at index 6: unexpected EOF`,
		},
		{
			name:   `missing closing quotes`,
			input:  `"string`,
			expect: `error at index 7: unexpected EOF`,
		},
		{
			name:   `missing closing quotes after escaped quotes`,
			input:  `"string\"`,
			expect: `error at index 9: unexpected EOF`,
		},
		{
			name:   `missing closing quotes after escaped sequences`,
			input:  `"string\\\"`,
			expect: `error at index 11: unexpected EOF`,
		},
		{
			name:   `unfinished key`,
			input:  `{"key`,
			expect: `error at index 5: unexpected EOF`,
		},
		{
			name:   `missing column`,
			input:  `{"key"}`,
			expect: `error at index 6 ('}'): unexpected token`,
		},
		{
			name:   `invalid content before column`,
			input:  `{"key"1 :}`,
			expect: `error at index 6 ('1'): unexpected token`,
		},
		{
			name:   `invalid column`,
			input:  `{"key";1}`,
			expect: `error at index 6 (';'): unexpected token`,
		},
		{
			name:   `missing field value`,
			input:  `{"okay":}`,
			expect: `error at index 8 ('}'): unexpected token`,
		},
		{
			name:   "unexpected object",
			input:  `{"key":12,{}}`,
			expect: `error at index 10 ('{'): unexpected token`,
		},
		{
			name:   `missing array item`,
			input:  `["okay",]`,
			expect: `error at index 8 (']'): unexpected token`,
		},
		{
			name:   `missing comma`,
			input:  `["okay"[`,
			expect: `error at index 7 ('['): unexpected token`,
		},
		{
			name:   `missing comma`,
			input:  `["okay"-12`,
			expect: `error at index 7 ('-'): unexpected token`,
		},
		{
			name:   `missing comma`,
			input:  `["okay"0`,
			expect: `error at index 7 ('0'): unexpected token`,
		},
		{
			name:   `missing comma`,
			input:  `["okay""not okay"]`,
			expect: `error at index 7 ('"'): unexpected token`,
		},
		{
			name:   `missing comma`,
			input:  `{"foo":"bar" "baz":"fuz"}`,
			expect: `error at index 13 ('"'): unexpected token`,
		},
		{
			name:   `missing comma`,
			input:  `[null false]`,
			expect: `error at index 6 ('f'): unexpected token`,
		},
		{
			name:   "expect EOF after number zero",
			input:  "01",
			expect: `error at index 1 ('1'): unexpected token`,
		},
		{
			name:   `expect EOF after negative number zero`,
			input:  `-00`,
			expect: `error at index 2 ('0'): unexpected token`,
		},
		{
			name:   `expect EOF after string get comma`,
			input:  `"okay",null`,
			expect: `error at index 6 (','): unexpected token`,
		},
		{
			name:   `expect EOF after string`,
			input:  `"str" "str"`,
			expect: `error at index 6 ('"'): unexpected token`,
		},
		{
			name:   `expect EOF after number`,
			input:  `0 0`,
			expect: `error at index 2 ('0'): unexpected token`,
		},
		{
			name:   `expect EOF after false`,
			input:  `false false`,
			expect: `error at index 6 ('f'): unexpected token`,
		},
		{
			name:   `expect EOF after true`,
			input:  `true true`,
			expect: `error at index 5 ('t'): unexpected token`,
		},
		{
			name:   `expect EOF after null`,
			input:  `null null`,
			expect: `error at index 5 ('n'): unexpected token`,
		},
		{
			name:   `expect EOF after array`,
			input:  `[] []`,
			expect: `error at index 3 ('['): unexpected token`,
		},
		{
			name:   `expect EOF after object`,
			input:  `{"k":0} {"k":0}`,
			expect: `error at index 8 ('{'): unexpected token`,
		},
	} {
		require.False(t, json.Valid([]byte(td.input)))

		t.Run(td.name, func(t *testing.T) {
			testError[string](t, td)
			testError[[]byte](t, td)
		})
	}
}

func testError[S ~string | ~[]byte](t *testing.T, td ErrorTest) {
	t.Run(testDataType(td.input), func(t *testing.T) {
		t.Run("Valid", func(t *testing.T) {
			require.False(t, jscan.Valid[S](S(td.input)))
		})
		t.Run("ValidatorValid", func(t *testing.T) {
			v := jscan.NewValidator[S](64)
			require.False(t, v.Valid(S(td.input)))
		})
		t.Run("Scan", func(t *testing.T) {
			err := jscan.Scan[S](
				S(td.input),
				func(i *jscan.Iterator[S]) (err bool) { return false },
			)
			require.Equal(t, td.expect, err.Error())
			require.True(t, err.IsErr())
		})
		t.Run("ParserScan", func(t *testing.T) {
			p := jscan.NewParser[S](64)
			err := p.Scan(
				S(td.input),
				func(i *jscan.Iterator[S]) (err bool) { return false },
			)
			require.Equal(t, td.expect, err.Error())
			require.True(t, err.IsErr())
		})
	})
}

func TestControlCharacters(t *testing.T) {
	// Control characters in a string value
	t.Run("value_string", func(t *testing.T) {
		ForASCIIControlChars(t, func(t *testing.T, b byte) {
			var buf bytes.Buffer
			buf.WriteByte('"')
			buf.WriteByte(b)
			buf.WriteByte('"')
			testControlCharacters(t, buf.String(), fmt.Sprintf(
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
			testControlCharacters(t, buf.String(), fmt.Sprintf(
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
			testControlCharacters(t, buf.String(), fmt.Sprintf(
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
			testControlCharacters(t, buf.String(), fmt.Sprintf(
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
			testControlCharacters(t, buf.String(), fmt.Sprintf(
				"error at index 0 (0x%x): illegal control character", b,
			))
		})

		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			var buf bytes.Buffer
			buf.WriteString("[")
			buf.WriteByte(b)
			buf.WriteString("]")
			testControlCharacters(t, buf.String(), fmt.Sprintf(
				"error at index 1 (0x%x): illegal control character", b,
			))
		})

		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			var buf bytes.Buffer
			buf.WriteString("[\n")
			buf.WriteByte(b)
			buf.WriteString("]")
			testControlCharacters(t, buf.String(), fmt.Sprintf(
				"error at index 2 (0x%x): illegal control character", b,
			))
		})

		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			var buf bytes.Buffer
			buf.WriteString("{")
			buf.WriteByte(b)
			buf.WriteString("}")
			testControlCharacters(t, buf.String(), fmt.Sprintf(
				"error at index 1 (0x%x): illegal control character", b,
			))
		})

		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			var buf bytes.Buffer
			buf.WriteString(`{"foo":`)
			buf.WriteByte(b)
			buf.WriteString(`false}`)
			testControlCharacters(t, buf.String(), fmt.Sprintf(
				"error at index 7 (0x%x): illegal control character", b,
			))
		})
	})
}

func testControlCharacters[S ~string | ~[]byte](t *testing.T, input S, expectErr string) {
	t.Run(testDataType(expectErr), func(t *testing.T) {
		t.Run("ValidatorValidate", func(t *testing.T) {
			v := jscan.NewValidator[S](64)
			err := v.Validate(S(input))
			require.Equal(t, expectErr, err.Error())
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
		})
		t.Run("ValidateValidateOne", func(t *testing.T) {
			for s := S(input); len(s) > 0; {
				v := jscan.NewValidator[S](64)
				trailing, err := v.ValidateOne(s)
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
		t.Run("ValidatorValid", func(t *testing.T) {
			v := jscan.NewValidator[S](64)
			require.False(t, v.Valid(S(input)))
		})
		t.Run("Validate", func(t *testing.T) {
			err := jscan.Validate[S](S(input))
			require.Equal(t, expectErr, err.Error())
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
		})
		t.Run("ValidateOne", func(t *testing.T) {
			for s := S(input); len(s) > 0; {
				trailing, err := jscan.ValidateOne[S](s)
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
		t.Run("Valid", func(t *testing.T) {
			require.False(t, jscan.Valid[S](S(input)))
		})
		t.Run("ParserScanOne", func(t *testing.T) {
			p := jscan.NewParser[S](64)
			_, err := p.ScanOne(
				S(input),
				func(i *jscan.Iterator[S]) (err bool) { return false },
			)
			require.Equal(t, expectErr, err.Error())
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
		})
		t.Run("ParserScan", func(t *testing.T) {
			p := jscan.NewParser[S](64)
			err := p.Scan(
				S(input),
				func(i *jscan.Iterator[S]) (err bool) { return false },
			)
			require.Equal(t, expectErr, err.Error())
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
		})
		t.Run("ScanOne", func(t *testing.T) {
			_, err := jscan.ScanOne[S](
				S(input),
				func(i *jscan.Iterator[S]) (err bool) { return false },
			)
			require.Equal(t, expectErr, err.Error())
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
		})
		t.Run("Scan", func(t *testing.T) {
			err := jscan.Scan[S](
				S(input),
				func(i *jscan.Iterator[S]) (err bool) { return false },
			)
			require.Equal(t, expectErr, err.Error())
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
		})
	})
}

func TestReturnErrorTrue(t *testing.T) {
	input := `{"foo":"bar","baz":null}`
	testReturnErrorTrue(t, string(input))
	testReturnErrorTrue(t, []byte(input))
}

func testReturnErrorTrue[S ~string | ~[]byte](t *testing.T, input S) {
	t.Run(testDataType(input), func(t *testing.T) {
		j := 0
		err := jscan.Scan(
			input,
			func(i *jscan.Iterator[S]) (err bool) {
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

type TypeValuePair struct {
	Type  jscan.ValueType
	Value string
}

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
	expect := []TypeValuePair{
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
	testScanOne[string](t, expect, inputs)
	testScanOne[[]byte](t, expect, inputs)
}

func testScanOne[S ~string | ~[]byte](
	t *testing.T, expect []TypeValuePair, inputs []string,
) {
	s := S(strings.Join(inputs, ""))
	t.Run(testDataType(s), func(t *testing.T) {
		for i, c, s := 0, 1, s; len(s) > 0; c++ {
			trailing, err := jscan.ScanOne(
				s,
				func(itr *jscan.Iterator[S]) (err bool) {
					require.Equal(t, expect[i], TypeValuePair{
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
		testStrings(t, string(td))
		testStrings(t, []byte(td))
	}
}

func testStrings[S ~string | ~[]byte](t *testing.T, input S) {
	t.Run(testDataType(input), func(t *testing.T) {
		inputObject := S(fmt.Sprintf(`{%s:%s}`, input, input))
		t.Run("ParserScan", func(t *testing.T) {
			p := jscan.NewParser[S](64)
			c := 0
			err := p.Scan(inputObject, func(i *jscan.Iterator[S]) (err bool) {
				if c < 1 {
					c++
					return false
				}
				require.Equal(t, jscan.ValueTypeString, i.ValueType())
				require.Equal(t, input, i.Value())
				require.Equal(t, input, i.Key())
				return false
			})
			require.False(t, err.IsErr())
		})
		t.Run("ParserScanOne", func(t *testing.T) {
			p := jscan.NewParser[S](64)
			c := 0
			trailing, err := p.ScanOne(
				inputObject,
				func(i *jscan.Iterator[S]) (err bool) {
					if c < 1 {
						c++
						return false
					}
					require.Equal(t, jscan.ValueTypeString, i.ValueType())
					require.Equal(t, input, i.Value())
					require.Equal(t, input, i.Key())
					return false
				},
			)
			require.False(t, err.IsErr())
			require.Len(t, trailing, 0)
		})
		t.Run("Scan", func(t *testing.T) {
			c := 0
			err := jscan.Scan[S](
				inputObject,
				func(i *jscan.Iterator[S]) (err bool) {
					if c < 1 {
						c++
						return false
					}
					require.Equal(t, jscan.ValueTypeString, i.ValueType())
					require.Equal(t, input, i.Value())
					require.Equal(t, input, i.Key())
					return false
				},
			)
			require.False(t, err.IsErr())
		})
		t.Run("ScanOne", func(t *testing.T) {
			c := 0
			trailing, err := jscan.ScanOne[S](
				inputObject,
				func(i *jscan.Iterator[S]) (err bool) {
					if c < 1 {
						c++
						return false
					}
					require.Equal(t, jscan.ValueTypeString, i.ValueType())
					require.Equal(t, input, i.Value())
					require.Equal(t, input, i.Key())
					return false
				},
			)
			require.False(t, err.IsErr())
			require.Len(t, trailing, 0)
		})
		t.Run("Valid", func(t *testing.T) {
			require.True(t, jscan.Valid[S](input))
		})
		t.Run("Validate", func(t *testing.T) {
			err := jscan.Validate[S](input)
			require.False(t, err.IsErr())
		})
		t.Run("ValidateOne", func(t *testing.T) {
			_, err := jscan.ValidateOne[S](input)
			require.False(t, err.IsErr())
		})
	})
}

func testDataType[S ~string | ~[]byte](input S) string {
	if _, ok := any(input).([]byte); ok {
		return "bytes"
	}
	return "string"
}

func TestIndexEnd(t *testing.T) {
	c := 0
	err := jscan.Scan(`{"x":["y"]}`, func(i *jscan.Iterator[string]) (err bool) {
		switch c {
		case 0:
			require.Equal(t, jscan.ValueTypeObject, i.ValueType())
			require.Equal(t, -1, i.ValueIndexEnd())
		case 1:
			require.Equal(t, jscan.ValueTypeArray, i.ValueType())
			require.Equal(t, -1, i.ValueIndexEnd())
		case 2:
			require.Equal(t, jscan.ValueTypeString, i.ValueType())
			require.Equal(t, len(`{"x":["y"`), i.ValueIndexEnd())
		default:
			t.Fatal("unexpected branch")
		}
		c++
		return false
	})
	require.False(t, err.IsErr())
}
