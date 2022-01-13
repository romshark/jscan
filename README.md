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

The following results were recorded on an Apple M1 Max MBP running macOS 12.1

```
goos: darwin
goarch: arm64
```

Tiny JSON document (`{"x":0}`):

```
BenchmarkCalcStats/jscan/tiny-10                   17627972     57.67 ns/op      0 B/op    0 allocs/op
BenchmarkCalcStats/jsoniter/tiny-10                10658056    111.8 ns/op     168 B/op    3 allocs/op
BenchmarkCalcStats/gofaster-jx/tiny-10              9148058    130.0 ns/op      40 B/op    2 allocs/op

BenchmarkCalcStats/jscan_withpath/tiny-10          15655423     76.27 ns/op      0 B/op    0 allocs/op
BenchmarkCalcStats/jsoniter_withpath/tiny-10        9757730    122.3 ns/op     168 B/op    3 allocs/op
BenchmarkCalcStats/gofaster-jx_withpath/tiny-10     8676548    138.9 ns/op      40 B/op    2 allocs/op
```

Small JSON document (335 bytes):

```
BenchmarkCalcStats/jscan/small-10                   1541143     777.8 ns/op      0 B/op     0 allocs/op
BenchmarkCalcStats/jsoniter/small-10                1436154     828.3 ns/op    576 B/op    13 allocs/op
BenchmarkCalcStats/gofaster-jx/small-10             1000000    1002 ns/op       80 B/op     8 allocs/op

BenchmarkCalcStats/jscan_withpath/small-10          1313886     903.3 ns/op      0 B/op     0 allocs/op
BenchmarkCalcStats/jsoniter_withpath/small-10        978420    1228 ns/op      640 B/op    22 allocs/op
BenchmarkCalcStats/gofaster-jx_withpath/small-10     855382    1407 ns/op      144 B/op    17 allocs/op
```

Large JSON document (26.1 MB):

```
BenchmarkCalcStats/jscan/large-10                   27     42554798 ns/op          47 B/op          0 allocs/op
BenchmarkCalcStats/jsoniter/large-10                18     60359799 ns/op    59209093 B/op    1108612 allocs/op
BenchmarkCalcStats/gofaster-jx/large-10             14     82380768 ns/op    35167540 B/op    1117362 allocs/op

BenchmarkCalcStats/jscan_withpath/large-10          24     45549453 ns/op         177 B/op          3 allocs/op
BenchmarkCalcStats/jsoniter_withpath/large-10       13     85553199 ns/op    81954372 B/op    1757457 allocs/op
BenchmarkCalcStats/gofaster-jx_withpath/large-10    10    109315817 ns/op    57892535 B/op    1766207 allocs/op
```

Get by path:

```
BenchmarkGet/jscan-10           4125115    287.7 ns/op     16 B/op     2 allocs/op
BenchmarkGet/jsoniter-10        1283280    932.2 ns/op    496 B/op    19 allocs/op
BenchmarkGet/tidwallgjson-10    6280434    190.4 ns/op     16 B/op     2 allocs/op
```

Validation:

```
BenchmarkValid/tiny
BenchmarkValid/tiny/jscan-10                       23285302          51.39 ns/op           0 B/op         0 allocs/op
BenchmarkValid/tiny/jsoniter-10                    19864123          60.38 ns/op           0 B/op         0 allocs/op
BenchmarkValid/tiny/gofaster-jx-10                 15638914          76.54 ns/op           0 B/op         0 allocs/op
BenchmarkValid/tiny/encoding-json-10               26987947          44.10 ns/op           0 B/op         0 allocs/op
BenchmarkValid/tiny/valyala-fastjson-10            41702926          28.73 ns/op           0 B/op         0 allocs/op

BenchmarkValid/small/jscan-10                       1540806         777.4 ns/op            0 B/op         0 allocs/op
BenchmarkValid/small/jsoniter-10                    1551792         772.8 ns/op           56 B/op         7 allocs/op
BenchmarkValid/small/gofaster-jx-10                 1404654         853.8 ns/op           16 B/op         2 allocs/op
BenchmarkValid/small/encoding-json-10               1000000        1007 ns/op              0 B/op         0 allocs/op
BenchmarkValid/small/valyala-fastjson-10            2757339         435.6 ns/op            0 B/op         0 allocs/op

BenchmarkValid/large/jscan-10                            27    41861076 ns/op            148 B/op         2 allocs/op
BenchmarkValid/large/jsoniter-10                         24    46747641 ns/op       13791527 B/op    644453 allocs/op
BenchmarkValid/large/gofaster-jx-10                      25    45416868 ns/op             54 B/op         0 allocs/op
BenchmarkValid/large/encoding-json-10                    16    70370383 ns/op             92 B/op         0 allocs/op
BenchmarkValid/large/valyala-fastjson-10                 40    29149781 ns/op              0 B/op         0 allocs/op

BenchmarkValid/unwind_stack/jscan-10                 169131        6998 ns/op              0 B/op         0 allocs/op
BenchmarkValid/unwind_stack/jsoniter-10               15855       75387 ns/op          33145 B/op      1033 allocs/op
BenchmarkValid/unwind_stack/gofaster-jx-10             1569      762668 ns/op         131117 B/op      2048 allocs/op
BenchmarkValid/unwind_stack/encoding-json-10         212806        5525 ns/op             24 B/op         1 allocs/op
BenchmarkValid/unwind_stack/valyala-fastjson-10         218     5456765 ns/op       52468093 B/op      4132 allocs/op
```
