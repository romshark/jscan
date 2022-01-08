package jscan_test

import (
	"encoding/json"
	"testing"

	"github.com/romshark/jscan"

	"github.com/stretchr/testify/require"
)

func TestScan(t *testing.T) {
	type Record struct {
		Level      int
		ValueType  jscan.ValueType
		Key        string
		Value      string
		ArrayIndex int
		Path       string
	}

	for _, tt := range []struct {
		name       string
		escapePath bool
		input      string
		expect     []Record
	}{
		{
			name:       "null",
			escapePath: true,
			input:      "null",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeNull,
					ArrayIndex: -1,
					Value:      "null",
				},
			},
		},
		{
			name:       "bool_true",
			escapePath: true,
			input:      "true",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeTrue,
					ArrayIndex: -1,
					Value:      "true",
				},
			},
		},
		{
			name:       "bool_false",
			escapePath: true,
			input:      "false",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeFalse,
					ArrayIndex: -1,
					Value:      "false",
				},
			},
		},
		{
			name:       "number_int",
			escapePath: true,
			input:      "42",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeNumber,
					Value:      "42",
					ArrayIndex: -1,
				},
			},
		},
		{
			name:       "number_decimal",
			escapePath: true,
			input:      "42.5",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeNumber,
					Value:      "42.5",
					ArrayIndex: -1,
				},
			},
		},
		{
			name:       "number_negative",
			escapePath: true,
			input:      "-42.5",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeNumber,
					Value:      "-42.5",
					ArrayIndex: -1,
				},
			},
		},
		{
			name:       "number_exponent",
			escapePath: true,
			input:      "2.99792458e8",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeNumber,
					Value:      "2.99792458e8",
					ArrayIndex: -1,
				},
			},
		},
		{
			name:       "string",
			escapePath: true,
			input:      `"42"`,
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeString,
					Value:      "42",
					ArrayIndex: -1,
				},
			},
		},
		{
			name:       "escaped unicode string",
			escapePath: true,
			input:      `"жш\"ц\\\\\""`,
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeString,
					Value:      `жш\"ц\\\\\"`,
					ArrayIndex: -1,
				},
			},
		},
		{
			name:       "empty array",
			escapePath: true,
			input:      "[]",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeArray,
					ArrayIndex: -1,
				},
			},
		},
		{
			name:       "empty object",
			escapePath: true,
			input:      "{}",
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeObject,
					ArrayIndex: -1,
				},
			},
		},
		{
			name:       "nested array",
			escapePath: true,
			input:      `[[null,[{"key":true}]],[]]`,
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeArray,
					ArrayIndex: -1,
				},
				{
					ValueType:  jscan.ValueTypeArray,
					Level:      1,
					ArrayIndex: 0,
					Path:       "[0]",
				},
				{
					ValueType:  jscan.ValueTypeNull,
					Level:      2,
					Value:      "null",
					ArrayIndex: 0,
					Path:       "[0][0]",
				},
				{
					ValueType:  jscan.ValueTypeArray,
					Level:      2,
					ArrayIndex: 1,
					Path:       "[0][1]",
				},
				{
					ValueType:  jscan.ValueTypeObject,
					Level:      3,
					ArrayIndex: 0,
					Path:       "[0][1][0]",
				},
				{
					ValueType:  jscan.ValueTypeTrue,
					Key:        "key",
					Value:      "true",
					Level:      4,
					ArrayIndex: -1,
					Path:       "[0][1][0].key",
				},
				{
					ValueType:  jscan.ValueTypeArray,
					Level:      1,
					ArrayIndex: 1,
					Path:       "[1]",
				},
			},
		},
		{
			name:       "escaped path",
			escapePath: true,
			input:      `{"[0]":[{"y.z":null},0]}`,
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeObject,
					ArrayIndex: -1,
				},
				{
					ValueType:  jscan.ValueTypeArray,
					Level:      1,
					Key:        "[0]",
					ArrayIndex: -1,
					Path:       `\[0\]`,
				},
				{
					ValueType:  jscan.ValueTypeObject,
					Level:      2,
					ArrayIndex: 0,
					Path:       `\[0\][0]`,
				},
				{
					ValueType:  jscan.ValueTypeNull,
					Key:        "y.z",
					Value:      "null",
					Level:      3,
					ArrayIndex: -1,
					Path:       `\[0\][0].y\.z`,
				},
				{
					ValueType:  jscan.ValueTypeNumber,
					Value:      "0",
					Level:      2,
					ArrayIndex: 1,
					Path:       `\[0\][1]`,
				},
			},
		},
		{
			name:       "unescaped path",
			escapePath: false,
			input:      `{"[0]":[{"y.z":null},0]}`,
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeObject,
					ArrayIndex: -1,
				},
				{
					ValueType:  jscan.ValueTypeArray,
					Level:      1,
					Key:        "[0]",
					ArrayIndex: -1,
					Path:       `[0]`,
				},
				{
					ValueType:  jscan.ValueTypeObject,
					Level:      2,
					ArrayIndex: 0,
					Path:       `[0][0]`,
				},
				{
					ValueType:  jscan.ValueTypeNull,
					Key:        "y.z",
					Value:      "null",
					Level:      3,
					ArrayIndex: -1,
					Path:       `[0][0].y.z`,
				},
				{
					ValueType:  jscan.ValueTypeNumber,
					Value:      "0",
					Level:      2,
					ArrayIndex: 1,
					Path:       `[0][1]`,
				},
			},
		},
		{
			name:       "nested object",
			escapePath: true,
			input: `{"s":"value","t":true,"f":false,"0":null,"n":-9.123e3,` +
				`"o0":{},"a0":[],"o":{"k":"\"v\"",` +
				`"a":[true,null,"item",-67.02e9,["foo"]]},"a3":[0,{"a3":8}]}`,
			expect: []Record{
				{
					ValueType:  jscan.ValueTypeObject,
					ArrayIndex: -1,
				},
				{
					Level:      1,
					ValueType:  jscan.ValueTypeString,
					Key:        "s",
					Value:      "value",
					ArrayIndex: -1,
					Path:       "s",
				},
				{
					Level:      1,
					ValueType:  jscan.ValueTypeTrue,
					Key:        "t",
					Value:      "true",
					ArrayIndex: -1,
					Path:       "t",
				},
				{
					Level:      1,
					ValueType:  jscan.ValueTypeFalse,
					Key:        "f",
					Value:      "false",
					ArrayIndex: -1,
					Path:       "f",
				},
				{
					Level:      1,
					ValueType:  jscan.ValueTypeNull,
					Key:        "0",
					Value:      "null",
					ArrayIndex: -1,
					Path:       "0",
				},
				{
					Level:      1,
					ValueType:  jscan.ValueTypeNumber,
					Key:        "n",
					Value:      "-9.123e3",
					ArrayIndex: -1,
					Path:       "n",
				},
				{
					Level:      1,
					ValueType:  jscan.ValueTypeObject,
					Key:        "o0",
					ArrayIndex: -1,
					Path:       "o0",
				},
				{
					Level:      1,
					ValueType:  jscan.ValueTypeArray,
					Key:        "a0",
					ArrayIndex: -1,
					Path:       "a0",
				},
				{
					Level:      1,
					ValueType:  jscan.ValueTypeObject,
					Key:        "o",
					ArrayIndex: -1,
					Path:       "o",
				},
				{
					Level:      2,
					ValueType:  jscan.ValueTypeString,
					Key:        "k",
					Value:      `\"v\"`,
					ArrayIndex: -1,
					Path:       "o.k",
				},
				{
					Level:      2,
					ValueType:  jscan.ValueTypeArray,
					Key:        "a",
					ArrayIndex: -1,
					Path:       "o.a",
				},
				{
					Level:      3,
					ValueType:  jscan.ValueTypeTrue,
					Value:      "true",
					ArrayIndex: 0,
					Path:       "o.a[0]",
				},
				{
					Level:      3,
					ValueType:  jscan.ValueTypeNull,
					Value:      "null",
					ArrayIndex: 1,
					Path:       "o.a[1]",
				},
				{
					Level:      3,
					ValueType:  jscan.ValueTypeString,
					Value:      "item",
					ArrayIndex: 2,
					Path:       "o.a[2]",
				},
				{
					Level:      3,
					ValueType:  jscan.ValueTypeNumber,
					Value:      "-67.02e9",
					ArrayIndex: 3,
					Path:       "o.a[3]",
				},
				{
					Level:      3,
					ValueType:  jscan.ValueTypeArray,
					ArrayIndex: 4,
					Path:       "o.a[4]",
				},
				{
					Level:      4,
					ValueType:  jscan.ValueTypeString,
					Value:      "foo",
					ArrayIndex: 0,
					Path:       "o.a[4][0]",
				},
				{
					Level:      1,
					ValueType:  jscan.ValueTypeArray,
					Key:        "a3",
					ArrayIndex: -1,
					Path:       "a3",
				},
				{
					Level:      2,
					ValueType:  jscan.ValueTypeNumber,
					Value:      "0",
					ArrayIndex: 0,
					Path:       "a3[0]",
				},
				{
					Level:      2,
					ValueType:  jscan.ValueTypeObject,
					ArrayIndex: 1,
					Path:       "a3[1]",
				},
				{
					Level:      3,
					ValueType:  jscan.ValueTypeNumber,
					Key:        "a3",
					Value:      "8",
					ArrayIndex: -1,
					Path:       "a3[1].a3",
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			q := require.New(t)
			q.True(json.Valid([]byte(tt.input)))

			j := 0
			check := func(t *testing.T) func(i *jscan.Iterator) bool {
				q := require.New(t)
				return func(i *jscan.Iterator) bool {
					if j >= len(tt.expect) {
						t.Errorf("unexpected value at %d", j)
						j++
						return false
					}
					e := tt.expect[j]
					q.Equal(e.ValueType, i.ValueType, "ValueType at %d", i)
					q.Equal(e.Level, i.Level, "Level at %d", i)
					q.Equal(e.Value, i.Value(), "Value at %d", i)
					q.Equal(e.Key, i.Key(), "Key at %d", i)
					q.Equal(e.ArrayIndex, i.ArrayIndex, "ArrayIndex at %d", i)
					q.Equal(e.Path, i.Path(), "Path at %d", i)
					j++
					return false
				}
			}

			t.Run("cachepath", func(t *testing.T) {
				err := jscan.Scan(jscan.Options{
					CachePath: true, EscapePath: tt.escapePath,
				}, tt.input, check(t))
				require.False(t, err.IsErr(), "unexpected error: %s", err)
			})
			t.Run("nocachepath", func(t *testing.T) {
				j = 0
				err := jscan.Scan(jscan.Options{
					CachePath: false, EscapePath: tt.escapePath,
				}, tt.input, check(t))
				require.False(t, err.IsErr(), "unexpected error: %s", err)
			})
		})
	}
}

func TestScanError(t *testing.T) {
	for _, tt := range []struct {
		name   string
		input  string
		expect string
	}{
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
			name:   "invalid number",
			input:  "01",
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
			name:   "missing closing ]",
			input:  `[null`,
			expect: `error at index 5: unexpected EOF`,
		},
		{
			name:   "missing closing ]",
			input:  `[[null`,
			expect: `error at index 6: unexpected EOF`,
		},
		{
			name:   `missing closing quotes`,
			input:  `"string`,
			expect: `error at index 0 ('"'): unexpected EOF`,
		},
		{
			name:   `missing closing quotes`,
			input:  `"string\"`,
			expect: `error at index 0 ('"'): unexpected EOF`,
		},
		{
			name:   `missing value`,
			input:  `{"key"}`,
			expect: `error at index 6 ('}'): unexpected token`,
		},
		{
			name:   `missing column`,
			input:  `{"key"1 :}`,
			expect: `error at index 6 ('1'): unexpected token`,
		},
		{
			name:   `missing column`,
			input:  `{"key";1}`,
			expect: `error at index 6 (';'): unexpected token`,
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
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.False(t, json.Valid([]byte(tt.input)))

			check := func(t *testing.T) func(i *jscan.Iterator) bool {
				return func(i *jscan.Iterator) (err bool) { return false }
			}

			t.Run("cachepath", func(t *testing.T) {
				err := jscan.Scan(jscan.Options{
					CachePath: true,
				}, tt.input, check(t))
				require.True(t, err.IsErr())
				require.Equal(t, tt.expect, err.Error())
			})
			t.Run("nocachepath", func(t *testing.T) {
				err := jscan.Scan(jscan.Options{}, tt.input, check(t))
				require.True(t, err.IsErr())
				require.Equal(t, tt.expect, err.Error())
			})
		})
	}
}

func TestCachePathUnescaped(t *testing.T) {
	expect := []string{"", "x.y[0]", "x.y[0].z"}

	for _, tt := range []struct {
		name      string
		cachePath bool
	}{
		{"cachepath_escape", true},
		{"nocachepath_escape", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			j := 0
			err := jscan.Scan(jscan.Options{
				CachePath:  tt.cachePath,
				EscapePath: false,
			}, `{"x.y[0]":{"z":null}}`, func(i *jscan.Iterator) (err bool) {
				if j >= len(expect) {
					t.Fatalf("unexpected call: %d", j)
					return true
				}
				require.Equal(t, expect[j], i.Path())
				j++
				return false
			})
			require.False(t, err.IsErr(), "unexpected error: %s", err)
		})
	}
}

func TestReturnErrorTrue(t *testing.T) {
	ForPossibleOptions(func(name string, o jscan.Options) {
		t.Run(name, func(t *testing.T) {
			j := 0
			err := jscan.Scan(
				o,
				`{"foo":"bar","baz":null}`,
				func(i *jscan.Iterator) (err bool) {
					require.Equal(t, jscan.ValueTypeObject, i.ValueType)
					j++
					return true // Expect immediate return
				},
			)
			require.Equal(t, 1, j)
			require.True(t, err.IsErr())
			require.Equal(t, jscan.ErrorCallback, err.Code)
			require.Equal(
				t, "error at index 0 ('{'): callback error", err.Error(),
			)
		})
	})
}

// ForPossibleOptions calls fn with all possible option configurations
func ForPossibleOptions(fn func(name string, o jscan.Options)) {
	fn("cachepath_escaped", jscan.Options{
		CachePath:  true,
		EscapePath: true,
	})
	fn("cachepath_unescaped", jscan.Options{
		CachePath:  true,
		EscapePath: false,
	})
	fn("nocachepath_escaped", jscan.Options{
		CachePath:  false,
		EscapePath: true,
	})
	fn("nocachepath_unescaped", jscan.Options{
		CachePath:  false,
		EscapePath: false,
	})
}
