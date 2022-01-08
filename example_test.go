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

	err := jscan.Scan(jscan.Options{
		CachePath:  true,
		EscapePath: true,
	}, j, func(i *jscan.Iterator) (err bool) {
		fmt.Printf("| value:\n")
		fmt.Printf("|  level:      %d\n", i.Level)
		if k := i.Key(); k != "" {
			fmt.Printf("|  key:        %q\n", i.Key())
		}
		fmt.Printf("|  valueType:  %s\n", i.ValueType)
		if v := i.Value(); v != "" {
			fmt.Printf("|  value:      %q\n", i.Value())
		}
		fmt.Printf("|  arrayIndex: %d\n", i.ArrayIndex)
		fmt.Printf("|  path:       '%s'\n", i.Path())
		return false // No Error, resume scanning
	})

	if err.IsErr() {
		fmt.Printf("ERR: %s\n", err)
		return
	}

	// Output:
	// | value:
	// |  level:      0
	// |  valueType:  object
	// |  arrayIndex: -1
	// |  path:       ''
	// | value:
	// |  level:      1
	// |  key:        "s"
	// |  valueType:  string
	// |  value:      "value"
	// |  arrayIndex: -1
	// |  path:       's'
	// | value:
	// |  level:      1
	// |  key:        "t"
	// |  valueType:  true
	// |  value:      "true"
	// |  arrayIndex: -1
	// |  path:       't'
	// | value:
	// |  level:      1
	// |  key:        "f"
	// |  valueType:  false
	// |  value:      "false"
	// |  arrayIndex: -1
	// |  path:       'f'
	// | value:
	// |  level:      1
	// |  key:        "0"
	// |  valueType:  null
	// |  value:      "null"
	// |  arrayIndex: -1
	// |  path:       '0'
	// | value:
	// |  level:      1
	// |  key:        "n"
	// |  valueType:  number
	// |  value:      "-9.123e3"
	// |  arrayIndex: -1
	// |  path:       'n'
	// | value:
	// |  level:      1
	// |  key:        "o0"
	// |  valueType:  object
	// |  arrayIndex: -1
	// |  path:       'o0'
	// | value:
	// |  level:      1
	// |  key:        "a0"
	// |  valueType:  array
	// |  arrayIndex: -1
	// |  path:       'a0'
	// | value:
	// |  level:      1
	// |  key:        "o"
	// |  valueType:  object
	// |  arrayIndex: -1
	// |  path:       'o'
	// | value:
	// |  level:      2
	// |  key:        "k"
	// |  valueType:  string
	// |  value:      "\\\"v\\\""
	// |  arrayIndex: -1
	// |  path:       'o.k'
	// | value:
	// |  level:      2
	// |  key:        "a"
	// |  valueType:  array
	// |  arrayIndex: -1
	// |  path:       'o.a'
	// | value:
	// |  level:      3
	// |  valueType:  true
	// |  value:      "true"
	// |  arrayIndex: 0
	// |  path:       'o.a[0]'
	// | value:
	// |  level:      3
	// |  valueType:  null
	// |  value:      "null"
	// |  arrayIndex: 1
	// |  path:       'o.a[1]'
	// | value:
	// |  level:      3
	// |  valueType:  string
	// |  value:      "item"
	// |  arrayIndex: 2
	// |  path:       'o.a[2]'
	// | value:
	// |  level:      3
	// |  valueType:  number
	// |  value:      "-67.02e9"
	// |  arrayIndex: 3
	// |  path:       'o.a[3]'
	// | value:
	// |  level:      3
	// |  valueType:  array
	// |  arrayIndex: 4
	// |  path:       'o.a[4]'
	// | value:
	// |  level:      4
	// |  valueType:  string
	// |  value:      "foo"
	// |  arrayIndex: 0
	// |  path:       'o.a[4][0]'
	// | value:
	// |  level:      1
	// |  key:        "a3"
	// |  valueType:  array
	// |  arrayIndex: -1
	// |  path:       'a3'
	// | value:
	// |  level:      2
	// |  valueType:  number
	// |  value:      "0"
	// |  arrayIndex: 0
	// |  path:       'a3[0]'
	// | value:
	// |  level:      2
	// |  valueType:  object
	// |  arrayIndex: 1
	// |  path:       'a3[1]'
	// | value:
	// |  level:      3
	// |  key:        "a3.a3"
	// |  valueType:  number
	// |  value:      "8"
	// |  arrayIndex: -1
	// |  path:       'a3[1].a3\.a3'
}

func ExampleScan_error_handling() {
	j := `"something...`

	err := jscan.Scan(jscan.Options{}, j, func(i *jscan.Iterator) (err bool) {
		fmt.Println("This shall never be executed")
		return false // No Error, resume scanning
	})

	if err.IsErr() {
		fmt.Printf("ERR: %s\n", err)
		return
	}

	// Output:
	// ERR: error at index 0 ('"'): unexpected EOF
}

func ExampleGet() {
	j := `[false,[[2, {"[escaped]":[{"test-key":"string value"}]}]]]`

	if err := jscan.Get(
		j, `[1][0][1].\[escaped\][0].test-key`,
		true, func(i *jscan.Iterator) {
			fmt.Println(i.Value())
		},
	); err.IsErr() {
		fmt.Printf("ERR: %s\n", err)
		return
	}

	// Output:
	// string value
}
