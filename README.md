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

|package|version|
|-|-|
|github.com/go-faster/jx|[v1.0.0](https://github.com/go-faster/jx/releases/tag/v1.0.0)|
|github.com/json-iterator/go|[v1.1.12](https://github.com/json-iterator/go/releases/tag/v1.1.12)|
|github.com/sinhashubham95/jsonic|[v1.1.0](https://github.com/sinhashubham95/jsonic/releases/tag/v1.1.0)|
|github.com/tidwall/gjson|[v1.14.4](https://github.com/tidwall/gjson/releases/tag/v1.14.4)|
|github.com/valyala/fastjson|[v1.6.4](https://github.com/valyala/fastjson/releases/tag/v1.6.4)|

Tiny JSON document (`{"x":0}`):

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 46.22 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 85.98 ns/op | 160 B/op | 2 allocs/op |
| gofaster-jx | 54.21 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 50.17 ns/op | 0 B/op | 0 allocs/op |
| jscan_withpath | 61.20 ns/op | 0 B/op | 0 allocs/op |
| jsoniter_withpath | 92.64 ns/op | 160 B/op | 2 allocs/op |
| gofaster-jx_withpath | 64.24 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson_withpath | 73.10 ns/op | 8 B/op | 1 allocs/op |


Small JSON document (335 bytes):

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 506.7 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 768.7 ns/op | 224 B/op | 12 allocs/op |
| gofaster-jx | 561.2 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 553.0 ns/op | 0 B/op | 0 allocs/op |
| jscan_withpath | 670.7 ns/op | 0 B/op | 0 allocs/op |
| jsoniter_withpath | 1091 ns/op | 288 B/op | 21 allocs/op |
| gofaster-jx_withpath | 934.6 ns/op | 80 B/op | 13 allocs/op |
| valyala-fastjson_withpath | 784.7 ns/op | 88 B/op | 10 allocs/op |

Large JSON document (26.1 MB):

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan/ | 24938751 ns/op | 28 B/op | 0 allocs/op |
| jsoniter/ | 54856708 ns/op | 32851643 B/op | 1108519 allocs/op |
| gofaster-jx/ | 27957946 ns/op | 30 B/op | 0 allocs/op |
| valyala-fastjson/ | 28625101 ns/op | 35 B/op | 0 allocs/op |
| jscan_withpath/ | 30114974 ns/op | 33 B/op | 0 allocs/op |
| jsoniter_withpath/ | 74731725 ns/op | 55597093 B/op | 1757365 allocs/op |
| gofaster-jx_withpath/ | 61522048 ns/op | 51623144 B/op | 1533007 allocs/op |
| valyala-fastjson_withpath/ | 37039634 ns/op | 13135410 B/op | 325005 allocs/op |

Array of 1024 integers:

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 20441 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 38419 ns/op | 16528 B/op | 1025 allocs/op |
| gofaster-jx | 29924 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 21763 ns/op | 0 B/op | 0 allocs/op |
| jscan_withpath | 30549 ns/op | 0 B/op | 0 allocs/op |
| jsoniter_withpath | 82523 ns/op | 24496 B/op | 2973 allocs/op |
| gofaster-jx_withpath | 75776 ns/op | 7970 B/op | 1948 allocs/op |
| valyala-fastjson_withpath | 42226 ns/op | 8207 B/op | 1024 allocs/op |

Array of 1024 floating point numbers:

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 19628 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 43730 ns/op | 16528 B/op | 1025 allocs/op |
| gofaster-jx | 37036 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 21320 ns/op | 6 B/op | 0 allocs/op |
| jscan_withpath | 34130 ns/op | 0 B/op | 0 allocs/op |
| jsoniter_withpath | 88569 ns/op | 24496 B/op | 2973 allocs/op |
| gofaster-jx_withpath | 91708 ns/op | 7970 B/op | 1948 allocs/op |
| valyala-fastjson_withpath | 42265 ns/op | 8207 B/op | 1024 allocs/op |

Array of 1024 strings:

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 277372 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 585148 ns/op | 670313 B/op | 1019 allocs/op |
| gofaster-jx | 166966 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 61301 ns/op | 50 B/op | 0 allocs/op |
| jscan_withpath | 287013 ns/op | 0 B/op | 0 allocs/op |
| jsoniter_withpath | 629999 ns/op | 678286 B/op | 2967 allocs/op |
| gofaster-jx_withpath | 281957 ns/op | 677639 B/op | 2927 allocs/op |
| valyala-fastjson_withpath | 80256 ns/op | 8261 B/op | 1024 allocs/op |

Array of 1024 nullable booleans:

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 12308 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 21736 ns/op | 144 B/op | 1 allocs/op |
| gofaster-jx | 33209 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 10913 ns/op | 0 B/op | 0 allocs/op |
| jscan_withpath | 22391 ns/op | 0 B/op | 0 allocs/op |
| jsoniter_withpath | 62427 ns/op | 8112 B/op | 1949 allocs/op |
| gofaster-jx_withpath | 65697 ns/op | 7970 B/op | 1948 allocs/op |
| valyala-fastjson_withpath | 31945 ns/op | 8194 B/op | 1024 allocs/op |

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
| jscan | 35.17 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 45.14 ns/op | 0 B/op | 0 allocs/op |
| gofaster-jx | 39.19 ns/op | 0 B/op | 0 allocs/op |
| encoding-json | 42.53 ns/op | 0 B/op | 0 allocs/op |
| tidwallgjson | 15.85 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 19.25 ns/op | 0 B/op | 0 allocs/op |

Small

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 411.0 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 717.3 ns/op | 56 B/op | 7 allocs/op |
| gofaster-jx | 392.2 ns/op | 0 B/op | 0 allocs/op |
| encoding-json | 904.0 ns/op | 0 B/op | 0 allocs/op |
| tidwallgjson | 336.0 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 372.0 ns/op | 0 B/op | 0 allocs/op |

Large

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 22007394 ns/op | 23 B/op | 0 allocs/op |
| jsoniter | 44151891 ns/op | 13583525 B/op | 644362 allocs/op |
| gofaster-jx | 20625554 ns/op | 24 B/op | 0 allocs/op |
| encoding-json | 68774883 ns/op | 92 B/op | 0 allocs/op |
| tidwallgjson | 27378477 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 25695158 ns/op | 0 B/op | 0 allocs/op |

Unwinding Stack

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 2612 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 65817 ns/op | 33149 B/op | 1033 allocs/op |
| gofaster-jx | 398530 ns/op | 65687 B/op | 1026 allocs/op |
| encoding-json | 5133 ns/op | 24 B/op | 1 allocs/op |
| tidwallgjson | 14033 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 4876865 ns/op | 52431864 B/op | 4134 allocs/op |

Array of 1024 integers

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 17942 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 21317 ns/op | 0 B/op | 0 allocs/op |
| gofaster-jx | 17843 ns/op | 0 B/op | 0 allocs/op |
| encoding-json | 33607 ns/op | 0 B/op | 0 allocs/op |
| tidwallgjson | 13525 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 14741 ns/op | 0 B/op | 0 allocs/op |

Array of 1024 floats

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 12755 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 66119 ns/op | 8755 B/op | 547 allocs/op |
| gofaster-jx | 20111 ns/op | 0 B/op | 0 allocs/op |
| encoding-json | 35978 ns/op | 0 B/op | 0 allocs/op |
| tidwallgjson | 12176 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 17346 ns/op | 0 B/op | 0 allocs/op |

Array of 1024 nullable booleans

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 5483 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 16538 ns/op | 0 B/op | 0 allocs/op |
| gofaster-jx | 20987 ns/op | 0 B/op | 0 allocs/op |
| encoding-json | 20309 ns/op | 0 B/op | 0 allocs/op |
| tidwallgjson | 5081 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 4891 ns/op | 0 B/op | 0 allocs/op |

Array of 1024 strings

|implementation|ns/op|B/op|allocs/op|
|-|-|-|-|
| jscan | 275319 ns/op | 0 B/op | 0 allocs/op |
| jsoniter | 507315 ns/op | 0 B/op | 0 allocs/op |
| gofaster-jx | 152812 ns/op | 0 B/op | 0 allocs/op |
| encoding-json | 1391323 ns/op | 1 B/op | 0 allocs/op |
| tidwallgjson | 501928 ns/op | 0 B/op | 0 allocs/op |
| valyala-fastjson | 253798 ns/op | 0 B/op | 0 allocs/op |
