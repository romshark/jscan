<a href="https://github.com/romshark/jscan/actions?query=workflow%3ACI">
    <img src="https://github.com/romshark/jscan/workflows/CI/badge.svg" alt="GitHub Actions: CI">
</a>
<a href="https://coveralls.io/github/romshark/jscan">
    <img src="https://coveralls.io/repos/github/romshark/jscan/badge.svg" alt="Coverage Status" />
</a>
<a href="https://goreportcard.com/report/github.com/romshark/jscan">
    <img src="https://goreportcard.com/badge/github.com/romshark/jscan" alt="GoReportCard">
</a>
<a href="https://pkg.go.dev/github.com/romshark/jscan">
    <img src="https://pkg.go.dev/badge/github.com/romshark/jscan.svg" alt="Go Reference">
</a>

# jscan
[jscan](https://github.com/romshark/jscan) provides a high-performance zero-allocation JSON iterator for Go. It's **not** compatible with [encoding/json](https://pkg.go.dev/encoding/json) and doesn't provide the usual Marshal/Unmarshal capabilities, instead it focuses on fast and efficient scanning over JSON strings with on-the-fly validation.

## Example
https://go.dev/play/p/v-VeiMO2fsJ

```go
package main

import (
	"fmt"

	"github.com/romshark/jscan"
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

}
```

## Benchmark Results

### jscan vs [json-iterator/go](https://github.com/json-iterator/go)

```
goos: darwin
goarch: arm64
pkg: github.com/romshark/jscan
```

Tiny JSON document (`{"x":0}`):

```
BenchmarkCalcStats/jscan/tiny-10                17627972        57.67 ns/op        0 B/op          0 allocs/op
BenchmarkCalcStats/jsoniter/tiny-10.            10658056       111.8 ns/op       168 B/op          3 allocs/op
BenchmarkCalcStats/jscan_withpath/tiny-10.      15655423        76.27 ns/op        0 B/op          0 allocs/op
BenchmarkCalcStats/jsoniter_withpath/tiny-10.    9757730       122.3 ns/op       168 B/op          3 allocs/op
```

Small JSON document (335 bytes):

```
BenchmarkCalcStats/jscan/small-10               1541143	       777.8 ns/op         0 B/op          0 allocs/op
BenchmarkCalcStats/jsoniter/small-10            1436154	       828.3 ns/op       576 B/op         13 allocs/op
BenchmarkCalcStats/jscan_withpath/small-10.     1313886	       903.3 ns/op         0 B/op          0 allocs/op
BenchmarkCalcStats/jsoniter_withpath/small-10.   978420	      1228 ns/op         640 B/op         22 allocs/op
```

Large JSON document (26.1 MB):

```
BenchmarkCalcStats/jscan/large-10.              27    42554798 ns/op          47 B/op          0 allocs/op
BenchmarkCalcStats/jsoniter/large-10.           18    60359799 ns/op    59209093 B/op    1108612 allocs/op
BenchmarkCalcStats/jscan_withpath/large-10.     24    45549453 ns/op         177 B/op          3 allocs/op
BenchmarkCalcStats/jsoniter_withpath/large-10.  13    85553199 ns/op    81954372 B/op    1757457 allocs/op
```

Get by path:

```
BenchmarkGet/jsoniter-10	1253300	       960.2 ns/op       512 B/op         21 allocs/op
BenchmarkGet/jscan-10		3115696	       386.1 ns/op       128 B/op         10 allocs/op
```
