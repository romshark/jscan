<a href="https://pkg.go.dev/github.com/romshark/jscan/v2">
    <img src="https://godoc.org/github.com/romshark/jscan/v2?status.svg" alt="GoDoc">
</a>
<a href="https://goreportcard.com/report/github.com/romshark/jscan/v2">
    <img src="https://goreportcard.com/badge/github.com/romshark/jscan/v2" alt="GoReportCard">
</a>
<a href='https://coveralls.io/github/romshark/jscan?branch=main'>
    <img src='https://coveralls.io/repos/github/romshark/jscan/badge.svg?branch=main' alt='Coverage Status' />
</a>


# jscan
jscan provides high-performance zero-allocation JSON iterator and validator for Go. This module doesn't provide `Marshal`/`Unmarshal` capabilities *yet*, instead it focuses on highly efficient iteration over JSON data with on-the-fly validation.

An [experimental decoder](https://github.com/romshark/jscan-experimental-decoder) with backward compatibility to [encoding/json](https://pkg.go.dev/encoding/json) is WiP 🧪 and is expected to be introduced together with jscan v3.

jscan is tested against https://github.com/nst/JSONTestSuite, a comprehensive test suite for [RFC 8259](https://datatracker.ietf.org/doc/html/rfc8259) compliant JSON parsers.

See [jscan-benchmark](https://github.com/romshark/jscan-benchmark) for benchmark results 🏎️ 🏁.

## Example: Scan
https://go.dev/play/p/moP3l9EkebF

```go
package main

import (
	"fmt"

	"github.com/romshark/jscan/v2"
)

func main() {
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
		return false // Resume scanning
	})

	if err.IsErr() {
		fmt.Printf("ERR: %s\n", err)
		return
	}
}
```

## Example: Tokenizer

```go
package main

import (
	"fmt"
	"strconv"

	"github.com/romshark/jscan/v2"
)

func main() {
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
```
