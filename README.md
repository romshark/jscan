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
jscan provides a high-performance zero-allocation JSON iterator for Go. It's **not** compatible with [encoding/json](https://pkg.go.dev/encoding/json) and doesn't provide the usual `Marshal`/`Unmarshal` capabilities, instead it focuses on fast and efficient scanning over JSON strings with on-the-fly validation and error reporting.

jscan is tested against https://github.com/nst/JSONTestSuite, a comprehensive test suite for RFC 8259 compliant JSON parsers.

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

The following results were recorded on an Apple M1 Max MBP running macOS 13.2.1

```
goos: darwin
goarch: arm64
```

|package|version|
|-|-|
|pkg.go.dev/encoding/json|[go1.20.2](https://pkg.go.dev/encoding/json)|
|github.com/go-faster/jx|[v1.0.0](https://github.com/go-faster/jx/releases/tag/v1.0.0)|
|github.com/go-faster/jx|[v1.0.0](https://github.com/go-faster/jx/releases/tag/v1.0.0)|
|github.com/json-iterator/go|[v1.1.12](https://github.com/json-iterator/go/releases/tag/v1.1.12)|
|github.com/sinhashubham95/jsonic|[v1.1.0](https://github.com/sinhashubham95/jsonic/releases/tag/v1.1.0)|
|github.com/tidwall/gjson|[v1.14.4](https://github.com/tidwall/gjson/releases/tag/v1.14.4)|
|github.com/valyala/fastjson|[v1.6.4](https://github.com/valyala/fastjson/releases/tag/v1.6.4)|

Calculating statistics for a tiny JSON document (`{"x":0}`):

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 46.39 ns/op | 0 B/op | 0 allocs/op |
| encoding-json | 46.39 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 85.98 ns/op | 160 B/op | 2 allocs/op |
| gofaster-jx | 54.21 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 50.17 ns/op | 0 B/op | 0 allocs/op |


Calculating statistics for a small JSON document (335 bytes):

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 486.3 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 768.7 ns/op | 224 B/op | 12 allocs/op |
| gofaster-jx | 561.2 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 553.0 ns/op | 0 B/op | 0 allocs/op |

Calculating statistics for a large JSON document (26.1 MB):

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 24468140 ns/op | 6 B/op | 0 allocs/op |
| jsoniter | 54457668 ns/op | 32851588 B/op | 1108519 allocs/op |
| gofaster-jx | 28043963 ns/op | 8 B/op | 0 allocs/op |
| valyala-fastjson | 28241526 ns/op | 7 B/op | 0 allocs/op |

Calculating statistics for an object that contains a key and a string value consisting entirely of escape sequences (~3KB):

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 1891 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 7862 ns/op | 2208 B/op | 16 allocs/op |
| gofaster-jx | 6677 ns/op | 504 B/op | 6 allocs/op |
| valyala-fastjson | 11292 ns/op | 0 B/op | 0 allocs/op |


Array of 1024 integers:

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 20105 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 38419 ns/op | 16528 B/op | 1025 allocs/op |
| gofaster-jx | 29924 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 21763 ns/op | 0 B/op | 0 allocs/op |

Array of 1024 floating point numbers:

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 19196 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 43730 ns/op | 16528 B/op | 1025 allocs/op |
| gofaster-jx | 37036 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 21320 ns/op | 6 B/op | 0 allocs/op |

Array of 1024 strings:

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 267332 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 585148 ns/op | 670313 B/op | 1019 allocs/op |
| gofaster-jx | 166966 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 61301 ns/op | 50 B/op | 0 allocs/op |

Array of 1024 nullable booleans:

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 11592 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 21736 ns/op | 144 B/op | 1 allocs/op |
| gofaster-jx | 33209 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 10913 ns/op | 0 B/op | 0 allocs/op |

Get by path:

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 229.7 ns/op | 16 B/op | 2 allocs/op |
| jsoniter | 815.8 ns/op | 496 B/op | 19 allocs/op |
| tidwallgjson | 148.9 ns/op | 16 B/op | 2 allocs/op |
| valyalafastjson | 108.9 ns/op | 0 B/op | 0 allocs/op |
| sinhashubham95jsonic | 194.1 ns/op | 96 B/op | 1 allocs/op |

Validation:

Tiny

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 35.13 ns/op | 0 B/op | 0 allocs/op |
| encoding-json | 42.93 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 45.14 ns/op | 0 B/op | 0 allocs/op |
| gofaster-jx | 39.19 ns/op | 0 B/op | 0 allocs/op |
| tidwallgjson | 15.85 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 19.25 ns/op | 0 B/op | 0 allocs/op |

Small

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 391.5 ns/op | 0 B/op | 0 allocs/op |
| encoding-json | 903.8 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 717.3 ns/op | 56 B/op | 7 allocs/op |
| gofaster-jx | 392.2 ns/op | 0 B/op | 0 allocs/op |
| tidwallgjson | 336.0 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 372.0 ns/op | 0 B/op | 0 allocs/op |

Large

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 21481452 ns/op | 7 B/op | 0 allocs/op |
| encoding-json | 68749373 ns/op | 29 B/op | 0 allocs/op |
| jsoniter | 44151891 ns/op | 13583525 B/op | 644362 allocs/op |
| gofaster-jx | 20625554 ns/op | 24 B/op | 0 allocs/op |
| tidwallgjson | 27378477 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 25695158 ns/op | 0 B/op | 0 allocs/op |

Validating an object that contains a key and a string value consisting entirely of escape sequences (~3KB):

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 1879 ns/op | 0 B/op | 0 allocs/op |
| encoding-json | 9313 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 7858 ns/op | 2064 B/op | 15 allocs/op |
| gofaster-jx | 5897 ns/op | 0 B/op | 0 allocs/op |
| tidwallgjson | 2989 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 9062 ns/op | 0 B/op | 0 allocs/op |

Unwinding Stack

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 2614 ns/op | 0 B/op | 0 allocs/op |
| encoding-json | 5133 ns/op | 24 B/op | 1 allocs/op |
| jsoniter | 65817 ns/op | 33149 B/op | 1033 allocs/op |
| gofaster-jx | 398530 ns/op | 65687 B/op | 1026 allocs/op |
| tidwallgjson | 14033 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 4876865 ns/op | 52431864 B/op | 4134 allocs/op |

Array of 1024 integers

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 17170 ns/op | 0 B/op | 0 allocs/op |
| encoding-json | 34104 ns/op  | 0 B/op | 0 allocs/op |
| jsoniter | 21317 ns/op | 0 B/op | 0 allocs/op |
| gofaster-jx | 17843 ns/op | 0 B/op | 0 allocs/op |
| tidwallgjson | 13525 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 14741 ns/op | 0 B/op | 0 allocs/op |

Array of 1024 floats

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 12925 ns/op | 0 B/op | 0 allocs/op |
| encoding-json | 36390 ns/op  | 0 B/op | 0 allocs/op |
| jsoniter | 66119 ns/op | 8755 B/op | 547 allocs/op |
| gofaster-jx | 20111 ns/op | 0 B/op | 0 allocs/op |
| tidwallgjson | 12176 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 17346 ns/op | 0 B/op | 0 allocs/op |

Array of 1024 nullable booleans

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 4968 ns/op | 0 B/op | 0 allocs/op |
| encoding-json | 20239 ns/op  | 0 B/op | 0 allocs/op |
| jsoniter | 16538 ns/op | 0 B/op | 0 allocs/op |
| gofaster-jx | 20987 ns/op | 0 B/op | 0 allocs/op |
| tidwallgjson | 5081 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 4891 ns/op | 0 B/op | 0 allocs/op |

Array of 1024 strings

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 264281 ns/op | 0 B/op | 0 allocs/op |
| encoding-json | 1391270 ns/op  | 0 B/op | 0 allocs/op |
| jsoniter | 507315 ns/op | 0 B/op | 0 allocs/op |
| gofaster-jx | 152812 ns/op | 0 B/op | 0 allocs/op |
| tidwallgjson | 501928 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 253798 ns/op | 0 B/op | 0 allocs/op |
