package jscan_test

import (
	"fmt"

	"github.com/romshark/jscan"
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

// func ExampleGet() {
// 	j := `[false,[[2, {"[escaped]":[{"test-key":"string value"}]}]]]`

// 	if err := jscan.Get(
// 		j, `[1][0][1].\[escaped\][0].test-key`,
// 		true, func(i *jscan.Iterator) {
// 			fmt.Println(i.Value())
// 		},
// 	); err.IsErr() {
// 		fmt.Printf("ERR: %s\n", err)
// 		return
// 	}

// 	// Output:
// 	// string value
// }

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
