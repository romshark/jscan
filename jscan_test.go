package jscan_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"

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
				v := jscan.NewValidator[S](jscan.DefaultStackSizeValidator)
				a := v.Valid(input, jscan.Options{})
				require.True(t, a)
			})
			t.Run("Validate", func(t *testing.T) {
				v := jscan.NewValidator[S](jscan.DefaultStackSizeValidator)
				err := v.Validate(input, jscan.Options{})
				require.False(t, err.IsErr())
			})
			t.Run("ValidateOne", func(t *testing.T) {
				v := jscan.NewValidator[S](jscan.DefaultStackSizeValidator)
				_, err := v.ValidateOne(input, jscan.Options{})
				require.False(t, err.IsErr())
			})
		})
		t.Run("Valid", func(t *testing.T) {
			a := jscan.Valid[S](input, jscan.Options{})
			require.True(t, a)
		})
		t.Run("Scan", func(t *testing.T) {
			err := jscan.Scan(
				string(input), jscan.Options{},
				func(i *jscan.Iterator[string]) (err bool) { return false },
			)
			require.False(t, err.IsErr())
		})
		t.Run("ScanOne", func(t *testing.T) {
			_, err := jscan.ScanOne(
				string(input), jscan.Options{},
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
		if !jscan.Valid(string(input), jscan.Options{}) {
			t.Skip("allowed to fail")
		}
	})
	t.Run("Validate", func(t *testing.T) {
		err := jscan.Validate(string(input), jscan.Options{})
		if err.IsErr() {
			t.Skip("allowed to fail")
		}
	})
	t.Run("Validate_DisableUTF8Validation", func(t *testing.T) {
		err := jscan.Validate(string(input), jscan.Options{
			DisableUTF8Validation: true,
		})
		if err.IsErr() {
			t.Skip("allowed to fail")
		}
	})
	t.Run("Validator_Valid", func(t *testing.T) {
		v := jscan.NewValidator[string](jscan.DefaultStackSizeValidator)
		if !v.Valid(string(input), jscan.Options{}) {
			t.Skip("allowed to fail")
		}
	})
	t.Run("Validator_Validate", func(t *testing.T) {
		v := jscan.NewValidator[string](jscan.DefaultStackSizeValidator)
		err := v.Validate(string(input), jscan.Options{})
		if err.IsErr() {
			t.Skip("allowed to fail")
		}
	})
	t.Run("Validator_Validate_DisableUTF8Validation", func(t *testing.T) {
		v := jscan.NewValidator[string](jscan.DefaultStackSizeValidator)
		err := v.Validate(string(input), jscan.Options{
			DisableUTF8Validation: true,
		})
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
			p := jscan.NewParser[S](jscan.DefaultStackSizeParser)
			err := p.Scan(
				input, jscan.Options{},
				func(i *jscan.Iterator[S]) (err bool) { return false },
			)
			require.True(t, err.IsErr())
		})
		t.Run("Scan", func(t *testing.T) {
			err := jscan.Scan(
				input, jscan.Options{},
				func(i *jscan.Iterator[S]) (err bool) { return false },
			)
			require.True(t, err.IsErr())
		})
		t.Run("ValidatorValid", func(t *testing.T) {
			v := jscan.NewValidator[S](jscan.DefaultStackSizeValidator)
			a := v.Valid(input, jscan.Options{})
			require.False(t, a)
		})
		t.Run("ValidatorValidate", func(t *testing.T) {
			v := jscan.NewValidator[S](jscan.DefaultStackSizeValidator)
			err := v.Validate(input, jscan.Options{})
			require.True(t, err.IsErr())
		})
		t.Run("Valid", func(t *testing.T) {
			a := jscan.Valid[S](input, jscan.Options{})
			require.False(t, a)
		})
		t.Run("Validate", func(t *testing.T) {
			err := jscan.Validate[S](input, jscan.Options{})
			require.True(t, err.IsErr())
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
			name:  "escaped_unicode_string",
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
			name:  "empty_array",
			input: "[]",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeArray,
					ArrayIndex: -1,
				},
			},
		},
		{
			name:  "empty_object",
			input: "{}",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeObject,
					ArrayIndex: -1,
				},
			},
		},
		{
			name:  "nested_array",
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
			name:  "escaped_pointer",
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
			name: "nested_object",
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
			name:  "trailing_space",
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
			name:  "trailing_carriage_return",
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
			name:  "trailing_tab",
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
			name:  "trailing_line_break",
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
			a := jscan.Valid[S](S(td.input), jscan.Options{})
			require.True(t, a)
		})
		t.Run("ValidatorValid", func(t *testing.T) {
			v := jscan.NewValidator[S](0)
			a := v.Valid(S(td.input), jscan.Options{})
			require.True(t, a)
		})
		t.Run("Scan", func(t *testing.T) {
			j = 0
			err := jscan.Scan[S](S(td.input), jscan.Options{}, check(t))
			require.False(t, err.IsErr(), "unexpected error: %s", err)
		})
		t.Run("ParserScan", func(t *testing.T) {
			j = 0
			p := jscan.NewParser[S](0)
			err := p.Scan(S(td.input), jscan.Options{}, check(t))
			require.False(t, err.IsErr(), "unexpected error: %s", err)
		})
	})
}

type ErrorTest struct {
	name   string
	input  string
	expect string
}

func TestErrorUnexpectedEOF(t *testing.T) {
	for _, td := range []ErrorTest{
		{
			name:   "before_value",
			input:  "",
			expect: `error at index 0: unexpected EOF`,
		},
		{
			name:   "before_value_after_space",
			input:  " ",
			expect: `error at index 1: unexpected EOF`,
		},
		{
			name:   "before_value_after_trn",
			input:  "\t\r\n ",
			expect: `error at index 4: unexpected EOF`,
		},
		{
			name:   "after_opening_curlbrack",
			input:  `{`,
			expect: `error at index 1: unexpected EOF`,
		},
		{
			name:   "after_opening_curlbrack_after_space",
			input:  `{ `,
			expect: `error at index 2: unexpected EOF`,
		},
		{
			name:   "after_key",
			input:  `{"x"`,
			expect: `error at index 4: unexpected EOF`,
		},
		{
			name:   "after_key_space",
			input:  `{"x" `,
			expect: `error at index 5: unexpected EOF`,
		},
		{
			name:   "after_value_after_key",
			input:  `{"x":null`,
			expect: `error at index 9: unexpected EOF`,
		},
		{
			name:   "after_key_after_value_after_space",
			input:  `{"x":null `,
			expect: `error at index 10: unexpected EOF`,
		},
		{
			name:   "after_field_after_comma",
			input:  `{"x":null,`,
			expect: `error at index 10: unexpected EOF`,
		},
		{
			name:   "after_field_after_comma_after_space",
			input:  `{"x":null, `,
			expect: `error at index 11: unexpected EOF`,
		},
		{
			name:   "after_opening_squarebrack",
			input:  `[`,
			expect: `error at index 1: unexpected EOF`,
		},
		{
			name:   "after_opening_squarebrack_after_space",
			input:  `[ `,
			expect: `error at index 2: unexpected EOF`,
		},
		{
			name:   "before_comma_in_array",
			input:  `[null`,
			expect: `error at index 5: unexpected EOF`,
		},
		{
			name:   "before_comma_in_array_after_space",
			input:  `[null `,
			expect: `error at index 6: unexpected EOF`,
		},
		{
			name:   "after_arrayitem_after_comma",
			input:  `[null,`,
			expect: `error at index 6: unexpected EOF`,
		},
		{
			name:   "after_arrayitem_after_comma_after_space",
			input:  `[null, `,
			expect: `error at index 7: unexpected EOF`,
		},
		{
			name:   "before_closing_squarebrack",
			input:  `[[null`,
			expect: `error at index 6: unexpected EOF`,
		},
		{
			name:   `before_closing_quotes`,
			input:  `"string`,
			expect: `error at index 7: unexpected EOF`,
		},
		{
			name:   `before_closing_quotes_after_escaped_quotes`,
			input:  `"string\"`,
			expect: `error at index 9: unexpected EOF`,
		},
		{
			name:   `before_closing_quotes_after_escaped_sequences`,
			input:  `"string\\\"`,
			expect: `error at index 11: unexpected EOF`,
		},
		{
			name:   `after_revsolidus_in_string`,
			input:  `{"key\`,
			expect: `error at index 6: unexpected EOF`,
		},
		{
			name:   `before_closing_quotes_in_key`,
			input:  `{"key`,
			expect: `error at index 5: unexpected EOF`,
		},
		{
			name:   `after_revsolidus_in_key`,
			input:  `{"key\`,
			expect: `error at index 6: unexpected EOF`,
		},
	} {
		require.False(t, json.Valid([]byte(td.input)))

		t.Run(td.name, func(t *testing.T) {
			testError[string](t, td)
			testError[[]byte](t, td)
		})
	}
}

func TestErrorInvalidUTF8(t *testing.T) {
	testFn := func(t *testing.T, input string, beforeInvalidUTF8 string) {
		t.Helper()
		require.False(t, utf8.ValidString(input))

		t.Run("Valid", func(t *testing.T) {
			a := jscan.Valid(input, jscan.Options{})
			require.False(t, a)
		})

		t.Run("Validate", func(t *testing.T) {
			err := jscan.Validate(input, jscan.Options{})
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeInvalidUTF8, err.Code)
			expectMsg := fmt.Sprintf(
				"error at index %d: invalid UTF-8", len(beforeInvalidUTF8),
			)
			require.Equal(t, expectMsg, err.Error())
		})

		t.Run("Validator_Valid", func(t *testing.T) {
			v := jscan.NewValidator[string](jscan.DefaultStackSizeValidator)
			a := v.Valid(input, jscan.Options{})
			require.False(t, a)
		})

		t.Run("Validator_Validate", func(t *testing.T) {
			v := jscan.NewValidator[string](jscan.DefaultStackSizeValidator)
			err := v.Validate(input, jscan.Options{})
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeInvalidUTF8, err.Code)
			expectMsg := fmt.Sprintf(
				"error at index %d: invalid UTF-8", len(beforeInvalidUTF8),
			)
			require.Equal(t, expectMsg, err.Error())
		})

		t.Run("Scan", func(t *testing.T) {
			err := jscan.Scan(
				input, jscan.Options{},
				func(i *jscan.Iterator[string]) (err bool) { return false },
			)
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeInvalidUTF8, err.Code)
			expectMsg := fmt.Sprintf(
				"error at index %d: invalid UTF-8", len(beforeInvalidUTF8),
			)
			require.Equal(t, expectMsg, err.Error())
		})

		t.Run("ScanOne", func(t *testing.T) {
			_, err := jscan.ScanOne(
				input, jscan.Options{},
				func(i *jscan.Iterator[string]) (err bool) { return false },
			)
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeInvalidUTF8, err.Code)
			expectMsg := fmt.Sprintf(
				"error at index %d: invalid UTF-8", len(beforeInvalidUTF8),
			)
			require.Equal(t, expectMsg, err.Error())
		})

		t.Run("Parser_Scan", func(t *testing.T) {
			p := jscan.NewParser[string](jscan.DefaultStackSizeParser)
			err := p.Scan(
				input, jscan.Options{},
				func(i *jscan.Iterator[string]) (err bool) { return false },
			)
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeInvalidUTF8, err.Code)
			expectMsg := fmt.Sprintf(
				"error at index %d: invalid UTF-8", len(beforeInvalidUTF8),
			)
			require.Equal(t, expectMsg, err.Error())
		})

		t.Run("Parser_ScanOne", func(t *testing.T) {
			p := jscan.NewParser[string](jscan.DefaultStackSizeParser)
			_, err := p.ScanOne(
				input, jscan.Options{},
				func(i *jscan.Iterator[string]) (err bool) { return false },
			)
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeInvalidUTF8, err.Code)
			expectMsg := fmt.Sprintf(
				"error at index %d: invalid UTF-8", len(beforeInvalidUTF8),
			)
			require.Equal(t, expectMsg, err.Error())
		})
	}

	invalidUTF8Strings := []struct{ Name, Str string }{
		// Lone start byte must be followed by
		// a continuation byte in the range 0x80 to 0xBF.
		{Name: "lone_start_byte", Str: "\xC3"},

		// Start byte followed by incorrect number of continuation bytes.
		// This 3-byte sequence start byte (E2) should be followed by
		// 2 continuation bytes, but there's only 1.
		{Name: "incorrect_continuation", Str: "\xE2\x82"},

		// This represents the null byte, which should be just 0x00 in UTF-8.
		// Using two bytes is an overlong encoding is invalid.
		{Name: "overlong_encoding", Str: "\xC0\x80"},

		// FE is not a valid start byte in UTF-8.
		{Name: "invalid_start_byte", Str: "\xFE"},

		// These are continuation bytes without a proper preceding start byte.
		{Name: "continuation_bytes_without_preceding_start_byte", Str: "\x80\x81"},

		// Surrogate pairs (U+D800 to U+DFFF) are not valid Unicode points and hence
		// must not appear in a UTF-8 encoded string.
		// This is an encoding for U+D800 which is a high surrogate.
		{Name: "surrogate_pairs", Str: "\xED\xA0\x80"},
	}

	t.Run("field_name", func(t *testing.T) {
		for _, td := range invalidUTF8Strings {
			t.Run(td.Name, func(t *testing.T) {
				in := `{"long_prefix_` + td.Str + `_long_suffix_string":0}`
				testFn(t, in, `{"long_prefix_`)
			})
		}
	})

	t.Run("string_value", func(t *testing.T) {
		for _, td := range invalidUTF8Strings {
			t.Run(td.Name, func(t *testing.T) {
				in := `[" ` + td.Str + `"]`
				testFn(t, in, `[" `)
			})
		}
	})
}

func TestErrorInvalidEscapeSequence(t *testing.T) {
	for _, td := range []ErrorTest{
		{
			name:   "invalid_escape_sequence_in_string",
			input:  `"\0"`,
			expect: `error at index 1 ('\'): invalid escape`,
		},
		{
			name:   "invalid_escape_sequence_in_string",
			input:  `"\u000m"`,
			expect: `error at index 1 ('\'): invalid escape`,
		},
		{
			name:   "invalid_escape_sequence_in_fieldname",
			input:  `{"\0":true}`,
			expect: `error at index 2 ('\'): invalid escape`,
		},
		{
			name:   "invalid_escape_sequence_in_string",
			input:  `{"\u000m":true}`,
			expect: `error at index 2 ('\'): invalid escape`,
		},
	} {
		require.False(t, json.Valid([]byte(td.input)))

		t.Run(td.name, func(t *testing.T) {
			testError[string](t, td)
			testError[[]byte](t, td)
		})
	}
}

func TestErrorUnexpectedToken(t *testing.T) {
	for _, td := range []ErrorTest{
		{
			name:   "before_value_invalid_literal_null",
			input:  "nul",
			expect: `error at index 0 ('n'): unexpected token`,
		},
		{
			name:   "before_value_invalid_literal_false",
			input:  "fals",
			expect: `error at index 0 ('f'): unexpected token`,
		},
		{
			name:   "before_value_invalid_literal_true",
			input:  "tru",
			expect: `error at index 0 ('t'): unexpected token`,
		},
		{
			name:   "before_value_invalid_literal_number",
			input:  "e1",
			expect: `error at index 0 ('e'): unexpected token`,
		},
		{
			name:   `after_key_closing_curlybrack`,
			input:  `{"key"}`,
			expect: `error at index 6 ('}'): unexpected token`,
		},
		{
			name:   `after_key_number`,
			input:  `{"key"1 :}`,
			expect: `error at index 6 ('1'): unexpected token`,
		},
		{
			name:   `after_key_semicolon`,
			input:  `{"key";1}`,
			expect: `error at index 6 (';'): unexpected token`,
		},
		{
			name:   `after_key_closing_curlybrack`,
			input:  `{"okay":}`,
			expect: `error at index 8 ('}'): unexpected token`,
		},
		{
			name:   "after_field_comma_empty_object",
			input:  `{"key":12,{}}`,
			expect: `error at index 10 ('{'): unexpected token`,
		},
		{
			name:   `after_fieldvalue_squarebrack`,
			input:  `{"f":""]`,
			expect: `error at index 7 (']'): unexpected token`,
		},
		{
			name:   `after_element_curlybrack`,
			input:  `[null}`,
			expect: `error at index 5 ('}'): unexpected token`,
		},
		{
			name:   `after_element_comma`,
			input:  `["okay",]`,
			expect: `error at index 8 (']'): unexpected token`,
		},
		{
			name:   `after_element_squarebrack`,
			input:  `["okay"[`,
			expect: `error at index 7 ('['): unexpected token`,
		},
		{
			name:   `after_element_number`,
			input:  `["okay"-12`,
			expect: `error at index 7 ('-'): unexpected token`,
		},
		{
			name:   `after_element_number_zero`,
			input:  `["okay"0`,
			expect: `error at index 7 ('0'): unexpected token`,
		},
		{
			name:   `after_element_string`,
			input:  `["okay""not okay"]`,
			expect: `error at index 7 ('"'): unexpected token`,
		},
		{
			name:   `after_field_string`,
			input:  `{"foo":"bar" "baz":"fuz"}`,
			expect: `error at index 13 ('"'): unexpected token`,
		},
		{
			name:   `after_element_false`,
			input:  `[null false]`,
			expect: `error at index 6 ('f'): unexpected token`,
		},
		{
			name:   `after_element_true`,
			input:  `[null true]`,
			expect: `error at index 6 ('t'): unexpected token`,
		},
		{
			name:   "after_number_zero_number",
			input:  "01",
			expect: `error at index 1 ('1'): unexpected token`,
		},
		{
			name:   `after_number_negzero_number`,
			input:  `-00`,
			expect: `error at index 2 ('0'): unexpected token`,
		},
		{
			name:   `after_string_comma`,
			input:  `"okay",null`,
			expect: `error at index 6 (','): unexpected token`,
		},
		{
			name:   `after_string_space_string`,
			input:  `"str" "str"`,
			expect: `error at index 6 ('"'): unexpected token`,
		},
		{
			name:   `after_zero_space_zero`,
			input:  `0 0`,
			expect: `error at index 2 ('0'): unexpected token`,
		},
		{
			name:   `after_false_space_false`,
			input:  `false false`,
			expect: `error at index 6 ('f'): unexpected token`,
		},
		{
			name:   `after_true_space_true`,
			input:  `true true`,
			expect: `error at index 5 ('t'): unexpected token`,
		},
		{
			name:   `after_null_space_null`,
			input:  `null null`,
			expect: `error at index 5 ('n'): unexpected token`,
		},
		{
			name:   `after_array_space_array`,
			input:  `[] []`,
			expect: `error at index 3 ('['): unexpected token`,
		},
		{
			name:   `after_object_space_object`,
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

func TestErrorMalformedNumber(t *testing.T) {
	for _, td := range []ErrorTest{
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
			a := jscan.Valid[S](S(td.input), jscan.Options{})
			require.False(t, a)
		})
		t.Run("ValidatorValid", func(t *testing.T) {
			v := jscan.NewValidator[S](jscan.DefaultStackSizeValidator)
			a := v.Valid(S(td.input), jscan.Options{})
			require.False(t, a)
		})
		t.Run("Scan", func(t *testing.T) {
			err := jscan.Scan[S](
				S(td.input), jscan.Options{},
				func(i *jscan.Iterator[S]) (err bool) { return false },
			)
			require.Equal(t, td.expect, err.Error())
			require.True(t, err.IsErr())
		})
		t.Run("ParserScan", func(t *testing.T) {
			p := jscan.NewParser[S](jscan.DefaultStackSizeParser)
			err := p.Scan(
				S(td.input), jscan.Options{},
				func(i *jscan.Iterator[S]) (err bool) { return false },
			)
			require.Equal(t, td.expect, err.Error())
			require.True(t, err.IsErr())
		})
	})
}

func TestControlCharacters(t *testing.T) {
	for i := 0; i <= 24; i++ {
		t.Run(fmt.Sprintf("value_string_%d", i), func(t *testing.T) {
			var b strings.Builder
			b.WriteString(`"`)
			for x := 0; x < i; x++ {
				b.WriteString(`x`)
			}
			b.WriteString("\x00")
			b.WriteString("1234567812345678\"")
			testControlCharacters(t, b.String(), fmt.Sprintf(
				"error at index %d (0x0): illegal control character", i+1,
			))
		})
	}

	for i := 0; i <= 24; i++ {
		t.Run(fmt.Sprintf("fieldname_string_%d", i), func(t *testing.T) {
			var b strings.Builder
			b.WriteString(`{"`)
			for x := 0; x < i; x++ {
				b.WriteString(`x`)
			}
			b.WriteString("\x00\":\"1234567812345678\"}")
			testControlCharacters(t, b.String(), fmt.Sprintf(
				"error at index %d (0x0): illegal control character", i+2,
			))
		})
	}

	t.Run("array_item_string", func(t *testing.T) {
		ForASCIIControlChars(t, func(t *testing.T, b byte) {
			s := `["` + string(b) + `"]`
			testControlCharacters(t, s, fmt.Sprintf(
				"error at index 2 (0x%x): illegal control character", b,
			))
		})
	})

	t.Run("before_key_after_space", func(t *testing.T) {
		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			s := `{ ` + string(b) + `:null}`
			testControlCharacters(t, s, fmt.Sprintf(
				"error at index 2 (0x%x): illegal control character", b,
			))
		})
	})

	t.Run("before_key_after_space_after_comma", func(t *testing.T) {
		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			s := `{"k1":null, ` + string(b) + ` "k2":null}`
			testControlCharacters(t, s, fmt.Sprintf(
				"error at index 12 (0x%x): illegal control character", b,
			))
		})
	})

	t.Run("before_colon", func(t *testing.T) {
		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			s := `{"key"` + string(b) + `:null}`
			testControlCharacters(t, s, fmt.Sprintf(
				"error at index 6 (0x%x): illegal control character", b,
			))
		})
	})

	t.Run("before_colon_after_space", func(t *testing.T) {
		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			s := `{"key" ` + string(b) + `:null}`
			testControlCharacters(t, s, fmt.Sprintf(
				"error at index 7 (0x%x): illegal control character", b,
			))
		})
	})

	t.Run("head", func(t *testing.T) {
		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			s := "\t" + string(b) + `[]`
			testControlCharacters(t, s, fmt.Sprintf(
				"error at index 1 (0x%x): illegal control character", b,
			))
		})
	})

	t.Run("before_key", func(t *testing.T) {
		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			s := "{" + string(b) + "}"
			testControlCharacters(t, s, fmt.Sprintf(
				"error at index 1 (0x%x): illegal control character", b,
			))
		})
	})

	t.Run("before_array_item_after_br", func(t *testing.T) {
		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			s := "[\n" + string(b) + "]"
			testControlCharacters(t, s, fmt.Sprintf(
				"error at index 2 (0x%x): illegal control character", b,
			))
		})
	})

	t.Run("after_colon", func(t *testing.T) {
		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			s := `{"foo":` + string(b) + `false}`
			testControlCharacters(t, s, fmt.Sprintf(
				"error at index 7 (0x%x): illegal control character", b,
			))
		})
	})

	t.Run("after_value", func(t *testing.T) {
		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			s := `[null` + string(b) + `,null]`
			testControlCharacters(t, s, fmt.Sprintf(
				"error at index 5 (0x%x): illegal control character", b,
			))
		})
	})

	t.Run("after_value_after_space", func(t *testing.T) {
		ForASCIIControlCharsExceptTRN(t, func(t *testing.T, b byte) {
			s := `[null ` + string(b) + `,null]`
			testControlCharacters(t, s, fmt.Sprintf(
				"error at index 6 (0x%x): illegal control character", b,
			))
		})
	})
}

func testControlCharacters[S ~string | ~[]byte](t *testing.T, input S, expectErr string) {
	t.Run(testDataType(expectErr), func(t *testing.T) {
		t.Run("ValidatorValidate", func(t *testing.T) {
			v := jscan.NewValidator[S](jscan.DefaultStackSizeValidator)
			err := v.Validate(S(input), jscan.Options{})
			require.Equal(t, expectErr, err.Error())
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
		})
		t.Run("ValidateValidateOne", func(t *testing.T) {
			for s := S(input); len(s) > 0; {
				v := jscan.NewValidator[S](jscan.DefaultStackSizeValidator)
				trailing, err := v.ValidateOne(s, jscan.Options{})
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
			v := jscan.NewValidator[S](jscan.DefaultStackSizeValidator)
			a := v.Valid(S(input), jscan.Options{})
			require.False(t, a)
		})
		t.Run("Validate", func(t *testing.T) {
			err := jscan.Validate[S](S(input), jscan.Options{})
			require.Equal(t, expectErr, err.Error())
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
		})
		t.Run("ValidateOne", func(t *testing.T) {
			for s := S(input); len(s) > 0; {
				trailing, err := jscan.ValidateOne[S](s, jscan.Options{})
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
			a := jscan.Valid[S](S(input), jscan.Options{})
			require.False(t, a)
		})
		t.Run("ParserScanOne", func(t *testing.T) {
			p := jscan.NewParser[S](jscan.DefaultStackSizeParser)
			_, err := p.ScanOne(
				S(input), jscan.Options{},
				func(i *jscan.Iterator[S]) (err bool) { return false },
			)
			require.Equal(t, expectErr, err.Error())
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
		})
		t.Run("ParserScan", func(t *testing.T) {
			p := jscan.NewParser[S](jscan.DefaultStackSizeParser)
			err := p.Scan(
				S(input), jscan.Options{},
				func(i *jscan.Iterator[S]) (err bool) { return false },
			)
			require.Equal(t, expectErr, err.Error())
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
		})
		t.Run("ScanOne", func(t *testing.T) {
			_, err := jscan.ScanOne[S](
				S(input), jscan.Options{},
				func(i *jscan.Iterator[S]) (err bool) { return false },
			)
			require.Equal(t, expectErr, err.Error())
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
		})
		t.Run("Scan", func(t *testing.T) {
			err := jscan.Scan[S](
				S(input), jscan.Options{},
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
			input, jscan.Options{},
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
				s, jscan.Options{},
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
			p := jscan.NewParser[S](jscan.DefaultStackSizeParser)
			c := 0
			err := p.Scan(
				inputObject, jscan.Options{},
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
		t.Run("ParserScanOne", func(t *testing.T) {
			p := jscan.NewParser[S](jscan.DefaultStackSizeParser)
			c := 0
			trailing, err := p.ScanOne(
				inputObject, jscan.Options{},
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
				inputObject, jscan.Options{},
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
				inputObject, jscan.Options{},
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
			a := jscan.Valid[S](input, jscan.Options{})
			require.True(t, a)
		})
		t.Run("Validate", func(t *testing.T) {
			err := jscan.Validate[S](input, jscan.Options{})
			require.False(t, err.IsErr())
		})
		t.Run("ValidateOne", func(t *testing.T) {
			_, err := jscan.ValidateOne[S](input, jscan.Options{})
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
	err := jscan.Scan(
		`{"x":["y"]}`,
		jscan.Options{},
		func(i *jscan.Iterator[string]) (err bool) {
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
		},
	)
	require.False(t, err.IsErr())
}

// TestDerivedTypes tests types derived from string and []byte as inputs
func TestDerivedTypes(t *testing.T) {
	t.Run("byte_slice_derivative", func(t *testing.T) {
		input := json.RawMessage(`0`)

		t.Run("Scan", func(t *testing.T) {
			err := jscan.Scan(
				input,
				func(i *jscan.Iterator[json.RawMessage]) (err bool) { return false },
			)
			require.False(t, err.IsErr())
		})
		t.Run("ScanOne", func(t *testing.T) {
			_, err := jscan.ScanOne(
				input,
				func(i *jscan.Iterator[json.RawMessage]) (err bool) { return false },
			)
			require.False(t, err.IsErr())
		})
		t.Run("Validate", func(t *testing.T) {
			err := jscan.Validate(input)
			require.False(t, err.IsErr())
		})
		t.Run("ValidateOne", func(t *testing.T) {
			_, err := jscan.ValidateOne(input)
			require.False(t, err.IsErr())
		})
		t.Run("Valid", func(t *testing.T) {
			isValid := jscan.Valid(input)
			require.True(t, isValid)
		})
	})

	t.Run("string_derivative", func(t *testing.T) {
		type String string

		input := String(`0`)

		t.Run("Scan", func(t *testing.T) {
			err := jscan.Scan(
				input,
				func(i *jscan.Iterator[String]) (err bool) { return false },
			)
			require.False(t, err.IsErr())
		})
		t.Run("ScanOne", func(t *testing.T) {
			_, err := jscan.ScanOne(
				input,
				func(i *jscan.Iterator[String]) (err bool) { return false },
			)
			require.False(t, err.IsErr())
		})
		t.Run("Validate", func(t *testing.T) {
			err := jscan.Validate(input)
			require.False(t, err.IsErr())
		})
		t.Run("ValidateOne", func(t *testing.T) {
			_, err := jscan.ValidateOne(input)
			require.False(t, err.IsErr())
		})
		t.Run("Valid", func(t *testing.T) {
			isValid := jscan.Valid(input)
			require.True(t, isValid)
		})
	})
}
