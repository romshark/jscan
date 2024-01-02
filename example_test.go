package jscan_test

import (
	"fmt"
	"strconv"

	"github.com/romshark/jscan/v2"
)

func ExampleScan() {
	j := `{
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
				null,
				"item",
				-67.02e9,
				["foo"]
			]
		},
		"a3": [
			0,
			{
				"a3.a3":8
			}
		]
	}`

	err := jscan.Scan(j, func(i *jscan.Iterator[string]) (err bool) {
		fmt.Printf("%q:\n", i.Pointer())
		fmt.Printf("├─ valueType:  %s\n", i.ValueType().String())
		if k := i.Key(); k != "" {
			fmt.Printf("├─ key:        %q\n", k[1:len(k)-1])
		}
		if ai := i.ArrayIndex(); ai != -1 {
			fmt.Printf("├─ arrayIndex: %d\n", ai)
		}
		if v := i.Value(); v != "" {
			fmt.Printf("├─ value:      %q\n", v)
		}
		fmt.Printf("└─ level:      %d\n", i.Level())
		return false // No Error, resume scanning
	})

	if err.IsErr() {
		fmt.Printf("ERR: %s\n", err)
		return
	}

	// Output:
	// "":
	// ├─ valueType:  object
	// └─ level:      0
	// "/s":
	// ├─ valueType:  string
	// ├─ key:        "s"
	// ├─ value:      "\"value\""
	// └─ level:      1
	// "/t":
	// ├─ valueType:  true
	// ├─ key:        "t"
	// ├─ value:      "true"
	// └─ level:      1
	// "/f":
	// ├─ valueType:  false
	// ├─ key:        "f"
	// ├─ value:      "false"
	// └─ level:      1
	// "/0":
	// ├─ valueType:  null
	// ├─ key:        "0"
	// ├─ value:      "null"
	// └─ level:      1
	// "/n":
	// ├─ valueType:  number
	// ├─ key:        "n"
	// ├─ value:      "-9.123e3"
	// └─ level:      1
	// "/o0":
	// ├─ valueType:  object
	// ├─ key:        "o0"
	// └─ level:      1
	// "/a0":
	// ├─ valueType:  array
	// ├─ key:        "a0"
	// └─ level:      1
	// "/o":
	// ├─ valueType:  object
	// ├─ key:        "o"
	// └─ level:      1
	// "/o/k":
	// ├─ valueType:  string
	// ├─ key:        "k"
	// ├─ value:      "\"\\\"v\\\"\""
	// └─ level:      2
	// "/o/a":
	// ├─ valueType:  array
	// ├─ key:        "a"
	// └─ level:      2
	// "/o/a/0":
	// ├─ valueType:  true
	// ├─ arrayIndex: 0
	// ├─ value:      "true"
	// └─ level:      3
	// "/o/a/1":
	// ├─ valueType:  null
	// ├─ arrayIndex: 1
	// ├─ value:      "null"
	// └─ level:      3
	// "/o/a/2":
	// ├─ valueType:  string
	// ├─ arrayIndex: 2
	// ├─ value:      "\"item\""
	// └─ level:      3
	// "/o/a/3":
	// ├─ valueType:  number
	// ├─ arrayIndex: 3
	// ├─ value:      "-67.02e9"
	// └─ level:      3
	// "/o/a/4":
	// ├─ valueType:  array
	// ├─ arrayIndex: 4
	// └─ level:      3
	// "/o/a/4/0":
	// ├─ valueType:  string
	// ├─ arrayIndex: 0
	// ├─ value:      "\"foo\""
	// └─ level:      4
	// "/a3":
	// ├─ valueType:  array
	// ├─ key:        "a3"
	// └─ level:      1
	// "/a3/0":
	// ├─ valueType:  number
	// ├─ arrayIndex: 0
	// ├─ value:      "0"
	// └─ level:      2
	// "/a3/1":
	// ├─ valueType:  object
	// ├─ arrayIndex: 1
	// └─ level:      2
	// "/a3/1/a3.a3":
	// ├─ valueType:  number
	// ├─ key:        "a3.a3"
	// ├─ value:      "8"
	// └─ level:      3
}

func ExampleScan_error_handling() {
	j := `"something...`

	err := jscan.Scan(j, func(i *jscan.Iterator[string]) (err bool) {
		fmt.Println("This shall never be executed")
		return false // No Error, resume scanning
	})

	if err.IsErr() {
		fmt.Printf("ERR: %s\n", err)
		return
	}

	// Output:
	// ERR: error at index 13: unexpected EOF
}

func ExampleValidateOne() {
	s := `-120.4` +
		`"string"` +
		`{"key":"value"}` +
		`[0,1]` +
		`true` +
		`false` +
		`null`

	for offset, x := 0, s; x != ""; offset = len(s) - len(x) {
		var err jscan.Error[string]
		if x, err = jscan.ValidateOne(x); err.IsErr() {
			panic(fmt.Errorf("unexpected error: %w", err))
		}
		fmt.Println(s[offset : len(s)-len(x)])
	}

	// Output:
	// -120.4
	// "string"
	// {"key":"value"}
	// [0,1]
	// true
	// false
	// null
}

func ExampleScan_decode2DIntArray() {
	j := `[[1,2,34,567],[8901,2147483647,-1,42]]`

	s := [][]int{}
	currentIndex := 0
	err := jscan.Scan(j, func(i *jscan.Iterator[string]) (err bool) {
		switch i.Level() {
		case 0: // Root array
			return i.ValueType() != jscan.ValueTypeArray
		case 1: // Sub-array
			if i.ValueType() != jscan.ValueTypeArray {
				return true
			}
			currentIndex = len(s)
			s = append(s, []int{})
			return false
		}
		if i.ValueType() != jscan.ValueTypeNumber {
			// Unexpected array element type
			return true
		}
		vi, errp := strconv.ParseInt(i.Value(), 10, 32)
		if errp != nil {
			// Not a valid 32-bit signed integer
			return true
		}
		s[currentIndex] = append(s[currentIndex], int(vi))
		return false
	})
	if err.IsErr() {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(s)

	// Output:
	// [[1 2 34 567] [8901 2147483647 -1 42]]
}

func ExampleTokenizer_decodeVector3DArray() {
	src := `[
		{"x": 12,   "y": 24,   "z": 12},
		{"x": 10.3, "y": 0.42, "z": 0.5},
		{"x": 0,    "y": 0.2,  "z": 10.275}
	]`

	// Initialize reusable tokenizer.
	tokenizer := jscan.NewTokenizer[string](
		jscan.DefaultStackSizeTokenizer,
		jscan.DefaultTokenBufferSize,
	)

	type Vector3D struct{ X, Y, Z float64 }
	var data []Vector3D

	var err error
	errTokenizer := tokenizer.Tokenize(src, func(tokens []jscan.Token) (errTok bool) {
		if tokens[0].Type != jscan.TokenTypeArray {
			err = fmt.Errorf("expected array at index %d", tokens[0].Index)
			return true
		}
		tokens = tokens[1 : len(tokens)-1]
		// Preallocate slice since we know the number of objects in advance.
		data = make([]Vector3D, tokens[0].Elements)

		mustParseField := func(defined bool, val jscan.Token) (float64, error) {
			if defined {
				return 0, fmt.Errorf("duplicated field at index %d", tokens[0].Index)
			}
			if val.Type != jscan.TokenTypeNumber && val.Type != jscan.TokenTypeInteger {
				return 0, fmt.Errorf("expected number at index %d", tokens[0].Index)
			}
			v, errParse := strconv.ParseFloat(src[val.Index:val.End], 64)
			if errParse != nil {
				return 0, fmt.Errorf("parsing number at index %d: %v",
					tokens[0].Index, err)
			}
			return v, nil
		}

		for i := 0; i < len(data); i++ {
			if tokens[0].Type != jscan.TokenTypeObject {
				err = fmt.Errorf("expected object at index %d", tokens[0].Index)
				return true
			}
			tokens = tokens[1:] // Skip object start token
			hasX, hasY, hasZ := false, false, false

			for k := 0; k < 3; k++ {
				fieldName := src[tokens[0].Index:tokens[0].End]
				switch string(fieldName) {
				case `"x"`:
					if data[i].X, err = mustParseField(hasX, tokens[1]); err != nil {
						return true
					}
					hasX, tokens = true, tokens[2:] // Skip key and value
				case `"y"`:
					if data[i].Y, err = mustParseField(hasY, tokens[1]); err != nil {
						return true
					}
					hasY, tokens = true, tokens[2:]
				case `"z"`:
					if data[i].Z, err = mustParseField(hasZ, tokens[1]); err != nil {
						return true
					}
					hasZ, tokens = true, tokens[2:]
				default:
					err = fmt.Errorf("unknown field %q at index %d",
						string(fieldName), tokens[0].Index)
					return true
				}
			}

			if tokens[0].Type != jscan.TokenTypeObjectEnd {
				err = fmt.Errorf("unknown extra field %q in object at index %d",
					string(src[tokens[0].Index:tokens[0].End]), tokens[0].Index)
			}
			tokens = tokens[1:] // Skip object end
		}

		return false
	})
	if errTokenizer.IsErr() {
		if errTokenizer.Code == jscan.ErrorCodeCallback {
			fmt.Printf("ERR: %v\n", err)
			return
		}
		fmt.Printf("ERR: %v\n", errTokenizer)
		return
	}

	fmt.Println(data)
	// Output:
	// [{12 24 12} {10.3 0.42 0.5} {0 0.2 10.275}]
}
