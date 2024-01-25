package jscan_test

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/romshark/jscan/v2"
	"github.com/stretchr/testify/assert"
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
		t.Run("TokenizerTokenize", func(t *testing.T) {
			k := jscan.NewTokenizer[string](64, 128)
			err := k.Tokenize(
				string(input),
				func(tokens []jscan.Token[string]) (err bool) { return false },
			)
			require.False(t, err.IsErr())
		})
		t.Run("TokenizerTokenizeOne", func(t *testing.T) {
			k := jscan.NewTokenizer[string](64, 128)
			_, err := k.TokenizeOne(
				string(input),
				func(tokens []jscan.Token[string]) (err bool) { return false },
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
	t.Run("TokenizerTokenize", func(t *testing.T) {
		v := jscan.NewTokenizer[string](1024, 4*1024)
		err := v.Tokenize(string(input), func(tokens []jscan.Token[string]) (err bool) {
			return false
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
			err := jscan.NewParser[S](1024).Scan(
				input, func(i *jscan.Iterator[S]) (err bool) { return false },
			)
			require.True(t, err.IsErr())
			require.NotEqual(t, err.Code, jscan.ErrorCodeCallback)
		})
		t.Run("Scan", func(t *testing.T) {
			err := jscan.Scan(
				input, func(i *jscan.Iterator[S]) (err bool) { return false },
			)
			require.True(t, err.IsErr())
			require.NotEqual(t, err.Code, jscan.ErrorCodeCallback)
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
		t.Run("TokenizerTokenize", func(t *testing.T) {
			err := jscan.NewTokenizer[S](1024, 4*1024).Tokenize(
				input, func(tokens []jscan.Token[S]) (err bool) { return false },
			)
			require.True(t, err.IsErr())
			require.NotEqual(t, err.Code, jscan.ErrorCodeCallback)
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

type ScanTest[S ~string | ~[]byte] struct {
	name         string
	input        S
	expect       []Record
	expectTokens []jscan.Token[S]
}

func TestParsingValid(t *testing.T) {
	for _, td := range []ScanTest[string]{
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
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeNull, Index: 0, End: 4},
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
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeTrue, Index: 0, End: 4},
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
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeFalse, Index: 0, End: 5},
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
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeInteger, Index: 0, End: 2},
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
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeNumber, Index: 0, End: 4},
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
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeNumber, Index: 0, End: 5},
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
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeNumber, Index: 0, End: 12},
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
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeString, Index: 0, End: 4},
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
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeString, Index: 0, End: len(`"жш\"ц\\\\\""`)},
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
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeArray, Index: 0, End: 1},
				{Type: jscan.TokenTypeArrayEnd, Index: 1, End: 0},
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
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeObject, Index: 0, End: 1},
				{Type: jscan.TokenTypeObjectEnd, Index: 1, End: 0},
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
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeArray, Index: 0, Elements: 2, End: 12}, // <────┐ 0
				{Type: jscan.TokenTypeArray, Index: 1, Elements: 2, End: 9},  // <───┐│ 1
				{Type: jscan.TokenTypeNull, Index: 2, End: 6},                //     ││ 2
				{Type: jscan.TokenTypeArray, Index: 7, Elements: 1, End: 8},  // <──┐││ 3
				{Type: jscan.TokenTypeObject, Index: 8, Elements: 1, End: 7}, // <─┐│││ 4
				{Type: jscan.TokenTypeKey, Index: 9, End: 14},                //   ││││ 5
				{Type: jscan.TokenTypeTrue, Index: 15, End: 19},              //   ││││ 6
				{Type: jscan.TokenTypeObjectEnd, Index: 19, End: 4},          // ──┘│││ 7
				{Type: jscan.TokenTypeArrayEnd, Index: 20, End: 3},           // ───┘││ 8
				{Type: jscan.TokenTypeArrayEnd, Index: 21, End: 1},           // ────┘│ 9
				{Type: jscan.TokenTypeArray, Index: 23, End: 11},             // <─┐  │ 10
				{Type: jscan.TokenTypeArrayEnd, Index: 24, End: 10},          // ──┘  │ 11
				{Type: jscan.TokenTypeArrayEnd, Index: 25, End: 0},           // ─────┘ 12
			},
		},
		{
			name:  "escaped_reverse_solidus_in_field_name",
			input: `{"\\":null}`,
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeObject,
					ArrayIndex: -1,
				},
				{
					ValueType:  jscan.ValueTypeNull,
					Key:        `"\\"`,
					Value:      "null",
					Level:      1,
					ArrayIndex: -1,
					Pointer:    `/\\`,
				},
			},
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeObject, Index: 0, Elements: 1, End: 3}, // <─┐ 0
				{Type: jscan.TokenTypeKey, Index: 1, End: 5},                 //   │ 1
				{Type: jscan.TokenTypeNull, Index: 6, End: 10},               //   │ 2
				{Type: jscan.TokenTypeObjectEnd, Index: 10, End: 0},          // ──┘ 3
			},
		},
		{
			name:  "escaped_quotes_in_field_name",
			input: `{"\"":null}`,
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeObject,
					ArrayIndex: -1,
				},
				{
					ValueType:  jscan.ValueTypeNull,
					Key:        `"\""`,
					Value:      "null",
					Level:      1,
					ArrayIndex: -1,
					Pointer:    `/\"`,
				},
			},
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeObject, Index: 0, Elements: 1, End: 3}, // <─┐ 0
				{Type: jscan.TokenTypeKey, Index: 1, End: 5},                 //   │ 1
				{Type: jscan.TokenTypeNull, Index: 6, End: 10},               //   │ 2
				{Type: jscan.TokenTypeObjectEnd, Index: 10, End: 0},          // ──┘ 3
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
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeObject, Index: 0, Elements: 1, End: 9}, // <───┐ 0
				{Type: jscan.TokenTypeKey, Index: 1, End: 4},                 //     │ 1
				{Type: jscan.TokenTypeArray, Index: 5, Elements: 2, End: 8},  // <──┐│ 2
				{Type: jscan.TokenTypeObject, Index: 6, Elements: 1, End: 6}, // <─┐││ 3
				{Type: jscan.TokenTypeKey, Index: 7, End: 10},                //   │││ 4
				{Type: jscan.TokenTypeNull, Index: 11, End: 15},              //   │││ 5
				{Type: jscan.TokenTypeObjectEnd, Index: 15, End: 3},          // ──┘││ 6
				{Type: jscan.TokenTypeInteger, Index: 17, End: 18},           //    ││ 7
				{Type: jscan.TokenTypeArrayEnd, Index: 18, End: 2},           // ───┘│ 8
				{Type: jscan.TokenTypeObjectEnd, Index: 19, End: 0},          // ────┘ 9
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
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeObject, Index: 0, Elements: 9, End: 41},   // <───┐ 0
				{Type: jscan.TokenTypeKey, Index: 6, End: 9},                    //     │ 1
				{Type: jscan.TokenTypeString, Index: 11, End: 18},               //     │ 2
				{Type: jscan.TokenTypeKey, Index: 24, End: 27},                  //     │ 3
				{Type: jscan.TokenTypeTrue, Index: 29, End: 33},                 //     │ 4
				{Type: jscan.TokenTypeKey, Index: 39, End: 42},                  //     │ 5
				{Type: jscan.TokenTypeFalse, Index: 44, End: 49},                //     │ 6
				{Type: jscan.TokenTypeKey, Index: 55, End: 58},                  //     │ 7
				{Type: jscan.TokenTypeNull, Index: 60, End: 64},                 //     │ 8
				{Type: jscan.TokenTypeKey, Index: 70, End: 73},                  //     │ 9
				{Type: jscan.TokenTypeNumber, Index: 75, End: 83},               //     │ 10
				{Type: jscan.TokenTypeKey, Index: 89, End: 93},                  //     │ 11
				{Type: jscan.TokenTypeObject, Index: 95, Elements: 0, End: 13},  // <┐  │ 12
				{Type: jscan.TokenTypeObjectEnd, Index: 96, End: 12},            // ─┘  │ 13
				{Type: jscan.TokenTypeKey, Index: 103, End: 107},                //     │ 14
				{Type: jscan.TokenTypeArray, Index: 109, Elements: 0, End: 16},  // <┐  │ 15
				{Type: jscan.TokenTypeArrayEnd, Index: 110, End: 15},            // ─┘  │ 16
				{Type: jscan.TokenTypeKey, Index: 117, End: 120},                //     │ 17
				{Type: jscan.TokenTypeObject, Index: 122, Elements: 2, End: 32}, // <──┐│ 18
				{Type: jscan.TokenTypeKey, Index: 129, End: 132},                //    ││ 19
				{Type: jscan.TokenTypeString, Index: 134, End: 141},             //    ││ 20
				{Type: jscan.TokenTypeKey, Index: 148, End: 151},                //    ││ 21
				{Type: jscan.TokenTypeArray, Index: 153, Elements: 6, End: 31},  // <─┐││ 22
				{Type: jscan.TokenTypeTrue, Index: 161, End: 165},               //   │││ 23
				{Type: jscan.TokenTypeFalse, Index: 173, End: 178},              //   │││ 24
				{Type: jscan.TokenTypeNull, Index: 186, End: 190},               //   │││ 25
				{Type: jscan.TokenTypeString, Index: 198, End: 204},             //   │││ 26
				{Type: jscan.TokenTypeNumber, Index: 212, End: 220},             //   │││ 27
				{Type: jscan.TokenTypeArray, Index: 228, Elements: 1, End: 30},  // <┐│││ 28
				{Type: jscan.TokenTypeString, Index: 229, End: 234},             //  ││││ 29
				{Type: jscan.TokenTypeArrayEnd, Index: 234, End: 28},            // ─┘│││ 30
				{Type: jscan.TokenTypeArrayEnd, Index: 241, End: 22},            // ──┘││ 31
				{Type: jscan.TokenTypeObjectEnd, Index: 247, End: 18},           // ───┘│ 32
				{Type: jscan.TokenTypeKey, Index: 254, End: 258},                //     │ 33
				{Type: jscan.TokenTypeArray, Index: 260, Elements: 2, End: 40},  // <─┐ │ 34
				{Type: jscan.TokenTypeInteger, Index: 262, End: 263},            //   │ │ 35
				{Type: jscan.TokenTypeObject, Index: 265, Elements: 1, End: 39}, // <┐│ │ 36
				{Type: jscan.TokenTypeKey, Index: 266, End: 270},                //  ││ │ 37
				{Type: jscan.TokenTypeInteger, Index: 272, End: 273},            //  ││ │ 38
				{Type: jscan.TokenTypeObjectEnd, Index: 273, End: 36},           // ─┘│ │ 39
				{Type: jscan.TokenTypeArrayEnd, Index: 275, End: 34},            // ──┘ │ 40
				{Type: jscan.TokenTypeObjectEnd, Index: 280, End: 0},            // ────┘ 41
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
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeNull, Index: 0, End: 4},
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
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeNull, Index: 0, End: 4},
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
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeNull, Index: 0, End: 4},
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
			expectTokens: []jscan.Token[string]{
				{Type: jscan.TokenTypeNull, Index: 0, End: 4},
			},
		},
	} {
		t.Run(td.name, func(t *testing.T) {
			require.True(t, json.Valid([]byte(td.input)))
			testParsingValid[string](t, td)
		})
	}
}

func testParsingValid[S ~string | ~[]byte](t *testing.T, td ScanTest[S]) {
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
		t.Run("TokenizerTokenize", func(t *testing.T) {
			k := jscan.NewTokenizer[S](64, 1024)
			var cp []jscan.Token[S]
			err := k.Tokenize(S(td.input), func(tokens []jscan.Token[S]) (err bool) {
				cp = make([]jscan.Token[S], len(tokens))
				copy(cp, tokens)
				return false
			})
			require.False(t, err.IsErr(), "unexpected error: %s", err)
			compareTokens[S](t, td.expectTokens, cp)
		})
	})
}

func compareTokens[S ~string | ~[]byte](t *testing.T, expected, actual []jscan.Token[S]) {
	t.Helper()
	assert.Len(t, actual, len(expected))
	for i, e := range expected {
		if i >= len(actual) {
			t.Errorf("missing index %d: %v", i, e)
			continue
		}
		assert.Equal(t, e.Type.String(), actual[i].Type.String(), "type at index %d", i)
		assert.Equal(t, e.Index, actual[i].Index, "index at index %d", i)
		assert.Equal(t, e.End, actual[i].End, "end at index %d", i)
		assert.Equal(t, e.Elements, actual[i].Elements, "elements at index %d", i)
	}
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
		t.Run("TokenizerTokenize", func(t *testing.T) {
			k := jscan.NewTokenizer[S](64, 128)
			err := k.Tokenize(
				S(td.input),
				func(tokens []jscan.Token[S]) (err bool) { return false },
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
		t.Run("TokenizerTokenize", func(t *testing.T) {
			k := jscan.NewTokenizer[S](64, 128)
			err := k.Tokenize(
				S(input),
				func(tokens []jscan.Token[S]) (err bool) { return false },
			)
			require.Equal(t, expectErr, err.Error())
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCodeIllegalControlChar, err.Code)
		})
		t.Run("TokenizerTokenizeOne", func(t *testing.T) {
			k := jscan.NewTokenizer[S](64, 128)
			_, err := k.TokenizeOne(
				S(input),
				func(tokens []jscan.Token[S]) (err bool) { return false },
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
		t.Run("TokenizerTokenize", func(t *testing.T) {
			k := jscan.NewTokenizer[S](64, 1024)
			c := 0
			err := k.Tokenize(input, func(tokens []jscan.Token[S]) (err bool) {
				require.Len(t, tokens, 1)
				require.Equal(t, 0, tokens[0].Index)
				require.Equal(t, jscan.TokenTypeString, tokens[0].Type)
				require.Equal(t, len(input), tokens[0].End)
				c++
				return false
			})
			require.False(t, err.IsErr())
			require.Equal(t, 1, c)
		})
		t.Run("TokenizerTokenizeOne", func(t *testing.T) {
			k := jscan.NewTokenizer[S](64, 1024)
			c := 0
			tail, err := k.TokenizeOne(input, func(tokens []jscan.Token[S]) (err bool) {
				require.Len(t, tokens, 1)
				require.Equal(t, 0, tokens[0].Index)
				require.Equal(t, jscan.TokenTypeString, tokens[0].Type)
				require.Equal(t, len(input), tokens[0].End)
				c++
				return false
			})
			require.False(t, err.IsErr())
			require.Equal(t, 1, c)
			require.Len(t, tail, 0)
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

func TestTokenBool(t *testing.T) {
	src := `null`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Bool(src)
		require.NoError(t, err)
		require.Zero(t, v)
	})

	src = `true`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Bool(src)
		require.NoError(t, err)
		require.Equal(t, true, v)
	})

	src = `false`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Bool(src)
		require.NoError(t, err)
		require.Equal(t, false, v)
	})

	src = `42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Bool(src)
		require.ErrorIs(t, err, jscan.ErrWrongType)
		require.Zero(t, v)
	})
}

func TestTokenString(t *testing.T) {
	src := `null`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.String(src)
		require.NoError(t, err)
		require.Zero(t, v)
	})

	src = `"text"`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.String(src)
		require.NoError(t, err)
		require.Equal(t, "text", v)
	})

	src = `""`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.String(src)
		require.NoError(t, err)
		require.Equal(t, "", v)
	})

	src = `42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.String(src)
		require.ErrorIs(t, err, jscan.ErrWrongType)
		require.Zero(t, v)
	})
}

func TestTokenFloat32(t *testing.T) {
	src := `null`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Float32(src)
		require.NoError(t, err)
		require.Zero(t, v)
	})

	src = `3.1415`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Float32(src)
		require.NoError(t, err)
		require.Equal(t, float32(3.1415), v)
	})

	src = `0`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Float32(src)
		require.NoError(t, err)
		require.Equal(t, float32(0), v)
	})

	src = `42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Float32(src)
		require.NoError(t, err)
		require.Equal(t, float32(42), v)
	})

	src = `-42.0`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Float32(src)
		require.NoError(t, err)
		require.Equal(t, float32(-42), v)
	})

	src = `123456e123456`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Float32(src)
		require.ErrorIs(t, err, strconv.ErrRange)
		require.Zero(t, v)
	})

	src = `false`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Float32(src)
		require.ErrorIs(t, err, jscan.ErrWrongType)
		require.Zero(t, v)
	})

	{ // Bytes
		src := []byte(`3.1415`)
		testValueToken(t, src, func(t *testing.T, token jscan.Token[[]byte]) {
			v, err := token.Float32(src)
			require.NoError(t, err)
			require.Equal(t, float32(3.1415), v)
		})
	}
}

func TestTokenFloat64(t *testing.T) {
	src := `null`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Float64(src)
		require.NoError(t, err)
		require.Zero(t, v)
	})

	src = `3.1415`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Float64(src)
		require.NoError(t, err)
		require.Equal(t, float64(3.1415), v)
	})

	src = `0`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Float64(src)
		require.NoError(t, err)
		require.Equal(t, float64(0), v)
	})

	src = `42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Float64(src)
		require.NoError(t, err)
		require.Equal(t, float64(42), v)
	})

	src = `-42.0`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Float64(src)
		require.NoError(t, err)
		require.Equal(t, float64(-42), v)
	})

	src = `123456e123456`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Float64(src)
		require.ErrorIs(t, err, strconv.ErrRange)
		require.Zero(t, v)
	})

	src = `false`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Float64(src)
		require.ErrorIs(t, err, jscan.ErrWrongType)
		require.Zero(t, v)
	})

	{ // Bytes
		src := []byte(`3.1415`)
		testValueToken(t, src, func(t *testing.T, token jscan.Token[[]byte]) {
			v, err := token.Float64(src)
			require.NoError(t, err)
			require.Equal(t, float64(3.1415), v)
		})
	}
}

func TestTokenInt(t *testing.T) {
	src := `null`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int(src)
		require.NoError(t, err)
		require.Zero(t, v)
	})

	src = `0`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int(src)
		require.NoError(t, err)
		require.Equal(t, int(0), v)
	})

	src = `42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int(src)
		require.NoError(t, err)
		require.Equal(t, int(42), v)
	})

	src = `-42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int(src)
		require.NoError(t, err)
		require.Equal(t, int(-42), v)
	})

	src = `9223372036854775808`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int(src)
		require.ErrorIs(t, err, jscan.ErrOverflow)
		require.Zero(t, v)
	})

	src = `-9223372036854775809`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int(src)
		require.ErrorIs(t, err, jscan.ErrOverflow)
		require.Zero(t, v)
	})

	src = `3.14`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int(src)
		require.ErrorIs(t, err, jscan.ErrWrongType)
		require.Zero(t, v)
	})
}

func TestTokenInt8(t *testing.T) {
	src := `null`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int8(src)
		require.NoError(t, err)
		require.Zero(t, v)
	})

	src = `0`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int8(src)
		require.NoError(t, err)
		require.Equal(t, int8(0), v)
	})

	src = `42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int8(src)
		require.NoError(t, err)
		require.Equal(t, int8(42), v)
	})

	src = `-42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int8(src)
		require.NoError(t, err)
		require.Equal(t, int8(-42), v)
	})

	src = fmt.Sprintf("%d", math.MaxInt8+1)
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int8(src)
		require.ErrorIs(t, err, jscan.ErrOverflow)
		require.Zero(t, v)
	})

	src = fmt.Sprintf("%d", math.MinInt8-1)
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int8(src)
		require.ErrorIs(t, err, jscan.ErrOverflow)
		require.Zero(t, v)
	})

	src = `3.14`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int8(src)
		require.ErrorIs(t, err, jscan.ErrWrongType)
		require.Zero(t, v)
	})
}

func TestTokenInt16(t *testing.T) {
	src := `null`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int16(src)
		require.NoError(t, err)
		require.Zero(t, v)
	})

	src = `0`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int16(src)
		require.NoError(t, err)
		require.Equal(t, int16(0), v)
	})

	src = `42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int16(src)
		require.NoError(t, err)
		require.Equal(t, int16(42), v)
	})

	src = `-42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int16(src)
		require.NoError(t, err)
		require.Equal(t, int16(-42), v)
	})

	src = fmt.Sprintf("%d", math.MaxInt16+1)
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int16(src)
		require.ErrorIs(t, err, jscan.ErrOverflow)
		require.Zero(t, v)
	})

	src = fmt.Sprintf("%d", math.MinInt16-1)
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int16(src)
		require.ErrorIs(t, err, jscan.ErrOverflow)
		require.Zero(t, v)
	})

	src = `3.14`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int16(src)
		require.ErrorIs(t, err, jscan.ErrWrongType)
		require.Zero(t, v)
	})
}

func TestTokenInt32(t *testing.T) {
	src := `null`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int32(src)
		require.NoError(t, err)
		require.Zero(t, v)
	})

	src = `0`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int32(src)
		require.NoError(t, err)
		require.Equal(t, int32(0), v)
	})

	src = `42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int32(src)
		require.NoError(t, err)
		require.Equal(t, int32(42), v)
	})

	src = `-42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int32(src)
		require.NoError(t, err)
		require.Equal(t, int32(-42), v)
	})

	src = fmt.Sprintf("%d", math.MaxInt32+1)
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int32(src)
		require.ErrorIs(t, err, jscan.ErrOverflow)
		require.Zero(t, v)
	})

	src = fmt.Sprintf("%d", math.MinInt32-1)
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int32(src)
		require.ErrorIs(t, err, jscan.ErrOverflow)
		require.Zero(t, v)
	})

	src = `3.14`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int32(src)
		require.ErrorIs(t, err, jscan.ErrWrongType)
		require.Zero(t, v)
	})
}

func TestTokenInt64(t *testing.T) {
	src := `null`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int64(src)
		require.NoError(t, err)
		require.Zero(t, v)
	})

	src = `0`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int64(src)
		require.NoError(t, err)
		require.Equal(t, int64(0), v)
	})

	src = `42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int64(src)
		require.NoError(t, err)
		require.Equal(t, int64(42), v)
	})

	src = `-42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int64(src)
		require.NoError(t, err)
		require.Equal(t, int64(-42), v)
	})

	src = `9223372036854775808`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int64(src)
		require.ErrorIs(t, err, jscan.ErrOverflow)
		require.Zero(t, v)
	})

	src = `-9223372036854775809`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int64(src)
		require.ErrorIs(t, err, jscan.ErrOverflow)
		require.Zero(t, v)
	})

	src = `3.14`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Int64(src)
		require.ErrorIs(t, err, jscan.ErrWrongType)
		require.Zero(t, v)
	})
}

func TestTokenUint(t *testing.T) {
	src := `null`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint(src)
		require.NoError(t, err)
		require.Zero(t, v)
	})

	src = `0`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint(src)
		require.NoError(t, err)
		require.Equal(t, uint(0), v)
	})

	src = `42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint(src)
		require.NoError(t, err)
		require.Equal(t, uint(42), v)
	})

	src = fmt.Sprintf("%d", math.MaxUint32)
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint(src)
		require.NoError(t, err)
		require.Equal(t, uint(math.MaxUint32), v)
	})

	src = `18446744073709551616`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint(src)
		require.ErrorIs(t, err, jscan.ErrOverflow)
		require.Zero(t, v)
	})

	src = `-42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint(src)
		require.ErrorIs(t, err, jscan.ErrWrongType)
		require.Zero(t, v)
	})
}

func TestTokenUint8(t *testing.T) {
	src := `null`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint8(src)
		require.NoError(t, err)
		require.Zero(t, v)
	})

	src = `0`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint8(src)
		require.NoError(t, err)
		require.Equal(t, uint8(0), v)
	})

	src = `42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint8(src)
		require.NoError(t, err)
		require.Equal(t, uint8(42), v)
	})

	src = fmt.Sprintf("%d", uint8(math.MaxUint8))
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint8(src)
		require.NoError(t, err)
		require.Equal(t, uint8(math.MaxUint8), v)
	})

	src = fmt.Sprintf("%d", math.MaxUint8+1)
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint8(src)
		require.ErrorIs(t, err, jscan.ErrOverflow)
		require.Zero(t, v)
	})

	src = `-42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint8(src)
		require.ErrorIs(t, err, jscan.ErrWrongType)
		require.Zero(t, v)
	})
}

func TestTokenUint16(t *testing.T) {
	src := `null`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint16(src)
		require.NoError(t, err)
		require.Zero(t, v)
	})

	src = `0`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint16(src)
		require.NoError(t, err)
		require.Equal(t, uint16(0), v)
	})

	src = `42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint16(src)
		require.NoError(t, err)
		require.Equal(t, uint16(42), v)
	})

	src = fmt.Sprintf("%d", uint16(math.MaxUint16))
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint16(src)
		require.NoError(t, err)
		require.Equal(t, uint16(math.MaxUint16), v)
	})

	src = fmt.Sprintf("%d", math.MaxUint16+1)
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint16(src)
		require.ErrorIs(t, err, jscan.ErrOverflow)
		require.Zero(t, v)
	})

	src = `-42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint16(src)
		require.ErrorIs(t, err, jscan.ErrWrongType)
		require.Zero(t, v)
	})
}

func TestTokenUint32(t *testing.T) {
	src := `null`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint32(src)
		require.NoError(t, err)
		require.Zero(t, v)
	})

	src = `0`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint32(src)
		require.NoError(t, err)
		require.Equal(t, uint32(0), v)
	})

	src = `42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint32(src)
		require.NoError(t, err)
		require.Equal(t, uint32(42), v)
	})

	src = fmt.Sprintf("%d", uint32(math.MaxUint32))
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint32(src)
		require.NoError(t, err)
		require.Equal(t, uint32(math.MaxUint32), v)
	})

	src = fmt.Sprintf("%d", math.MaxUint32+1)
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint32(src)
		require.ErrorIs(t, err, jscan.ErrOverflow)
		require.Zero(t, v)
	})

	src = `-42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint32(src)
		require.ErrorIs(t, err, jscan.ErrWrongType)
		require.Zero(t, v)
	})
}

func TestTokenUint64(t *testing.T) {
	src := `null`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint64(src)
		require.NoError(t, err)
		require.Zero(t, v)
	})

	src = `0`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint64(src)
		require.NoError(t, err)
		require.Equal(t, uint64(0), v)
	})

	src = `42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint64(src)
		require.NoError(t, err)
		require.Equal(t, uint64(42), v)
	})

	src = fmt.Sprintf("%d", uint64(math.MaxUint64))
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint64(src)
		require.NoError(t, err)
		require.Equal(t, uint64(math.MaxUint64), v)
	})

	src = `18446744073709551616`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint64(src)
		require.ErrorIs(t, err, jscan.ErrOverflow)
		require.Zero(t, v)
	})

	src = `-42`
	testValueToken(t, src, func(t *testing.T, token jscan.Token[string]) {
		v, err := token.Uint64(src)
		require.ErrorIs(t, err, jscan.ErrWrongType)
		require.Zero(t, v)
	})
}

func testValueToken[S ~string | ~[]byte](
	t *testing.T, input S, check func(t *testing.T, token jscan.Token[S]),
) {
	t.Helper()
	tok := jscan.NewTokenizer[S](4, 16)
	err := tok.Tokenize(input, func(tokens []jscan.Token[S]) (err bool) {
		require.Len(t, tokens, 1)
		check(t, tokens[0])
		return false
	})
	require.False(t, err.IsErr())
}
