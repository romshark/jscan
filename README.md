<a href="https://github.com/romshark/jscan/actions?query=workflow%3ACI">
    <img src="https://github.com/romshark/jscan/workflows/CI/badge.svg" alt="GitHub Actions: CI">
</a>
<a href="https://coveralls.io/github/romshark/jscan">
    <img src="https://coveralls.io/repos/github/romshark/jscan/badge.svg" alt="Coverage Status" />
</a>
<a href="https://goreportcard.com/report/github.com/romshark/jscan">
    <img src="https://goreportcard.com/badge/github.com/romshark/jscan" alt="GoReportCard">
</a>
<a href="https://pkg.go.dev/github.com/romshark/jscan/v2">
    <img src="https://pkg.go.dev/badge/github.com/romshark/jscan/v2.svg" alt="Go Reference">
</a>

# jscan
jscan provides high-performance zero-allocation JSON iterator and validator for Go. This module doesn't provide `Marshal`/`Unmarshal` capabilities, instead it focuses on highly efficient iteration over JSON data with on-the-fly validation.

jscan is tested against https://github.com/nst/JSONTestSuite, a comprehensive test suite for [RFC 8259](https://datatracker.ietf.org/doc/html/rfc8259) compliant JSON parsers.

## Example
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

## Benchmark Results

### Apple M1 - macOS

<details>

```
goos: darwin
goarch: arm64
pkg: github.com/romshark/jscan/v2
BenchmarkCalcStats/miniscule_1b__________/jscan___________-10         	58520883	        20.44 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/miniscule_1b__________/jsoniter________-10         	48171327	        25.08 ns/op	      16 B/op	       1 allocs/op
BenchmarkCalcStats/miniscule_1b__________/gofaster-jx_____-10         	64921742	        18.07 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/miniscule_1b__________/valyala-fastjson-10         	72884450	        16.44 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/tiny_8b_______________/jscan___________-10         	41999097	        28.57 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/tiny_8b_______________/jsoniter________-10         	26395453	        44.93 ns/op	      16 B/op	       1 allocs/op
BenchmarkCalcStats/tiny_8b_______________/gofaster-jx_____-10         	27551818	        43.31 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/tiny_8b_______________/valyala-fastjson-10         	29055865	        41.18 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/small_336b____________/jscan___________-10         	 3680330	       326.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/small_336b____________/jsoniter________-10         	 1751287	       685.1 ns/op	      80 B/op	      11 allocs/op
BenchmarkCalcStats/small_336b____________/gofaster-jx_____-10         	 2176081	       553.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/small_336b____________/valyala-fastjson-10         	 2186758	       548.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/large_26m_____________/jscan___________-10         	      84	  13997098 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/large_26m_____________/jsoniter________-10         	      20	  53716958 ns/op	32851291 B/op	 1108518 allocs/op
BenchmarkCalcStats/large_26m_____________/gofaster-jx_____-10         	      42	  27925808 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/large_26m_____________/valyala-fastjson-10         	      37	  29352441 ns/op	 9104579 B/op	    8944 allocs/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/jscan___________-10         	   10000	    116999 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/jsoniter________-10         	    3429	    344725 ns/op	  144473 B/op	    7357 allocs/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/gofaster-jx_____-10         	    5198	    229932 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/valyala-fastjson-10         	    3528	    336446 ns/op	     671 B/op	       1 allocs/op
BenchmarkCalcStats/escaped_3k____________/jscan___________-10         	  855004	      1400 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/escaped_3k____________/jsoniter________-10         	  151801	      7801 ns/op	    2064 B/op	      15 allocs/op
BenchmarkCalcStats/escaped_3k____________/gofaster-jx_____-10         	  181478	      6578 ns/op	     504 B/op	       6 allocs/op
BenchmarkCalcStats/escaped_3k____________/valyala-fastjson-10         	  106310	     11305 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_int_1024_12k____/jscan___________-10         	   86882	     13755 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_int_1024_12k____/jsoniter________-10         	   31441	     37991 ns/op	   16384 B/op	    1024 allocs/op
BenchmarkCalcStats/array_int_1024_12k____/gofaster-jx_____-10         	   39762	     30128 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_int_1024_12k____/valyala-fastjson-10         	   61540	     19450 ns/op	       5 B/op	       0 allocs/op
BenchmarkCalcStats/array_dec_1024_10k____/jscan___________-10         	   85881	     12686 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_dec_1024_10k____/jsoniter________-10         	   27861	     42949 ns/op	   16384 B/op	    1024 allocs/op
BenchmarkCalcStats/array_dec_1024_10k____/gofaster-jx_____-10         	   31263	     37542 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_dec_1024_10k____/valyala-fastjson-10         	   50748	     23808 ns/op	       7 B/op	       0 allocs/op
BenchmarkCalcStats/array_nullbool_1024_5k/jscan___________-10         	  169617	      7114 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_nullbool_1024_5k/jsoniter________-10         	   56238	     21339 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_nullbool_1024_5k/gofaster-jx_____-10         	   36418	     32872 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_nullbool_1024_5k/valyala-fastjson-10         	  114373	     10516 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_str_1024_639k___/jscan___________-10         	    8172	    146423 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_str_1024_639k___/jsoniter________-10         	    1940	    600455 ns/op	  670172 B/op	    1018 allocs/op
BenchmarkCalcStats/array_str_1024_639k___/gofaster-jx_____-10         	    7315	    165087 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_str_1024_639k___/valyala-fastjson-10         	   18740	     63804 ns/op	      52 B/op	       0 allocs/op
BenchmarkValid/deeparray_____________/jscan___________-10             	75256148	        16.42 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/deeparray_____________/encoding-json___-10             	 8976027	       133.4 ns/op	     104 B/op	       5 allocs/op
BenchmarkValid/deeparray_____________/jsoniter________-10             	 3476221	       345.9 ns/op	     352 B/op	       9 allocs/op
BenchmarkValid/deeparray_____________/gofaster-jx_____-10             	 4368445	       274.1 ns/op	      80 B/op	       2 allocs/op
BenchmarkValid/deeparray_____________/tidwallgjson____-10             	297245341	         4.033 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/deeparray_____________/valyala-fastjson-10             	 1000000	      1026 ns/op	    1184 B/op	      11 allocs/op
BenchmarkValid/deeparray_____________/goccy-go-json___-10             	   15310	     78688 ns/op	   49295 B/op	    2062 allocs/op
BenchmarkValid/deeparray_____________/bytedance-sonic_-10             	 8948800	       133.6 ns/op	     104 B/op	       5 allocs/op
BenchmarkValid/unwind_stack__________/jscan___________-10             	  579363	      2068 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/unwind_stack__________/encoding-json___-10             	  232671	      5129 ns/op	      24 B/op	       1 allocs/op
BenchmarkValid/unwind_stack__________/jsoniter________-10             	   18024	     66312 ns/op	   33150 B/op	    1033 allocs/op
BenchmarkValid/unwind_stack__________/gofaster-jx_____-10             	    2992	    398218 ns/op	   65664 B/op	    1026 allocs/op
BenchmarkValid/unwind_stack__________/tidwallgjson____-10             	   85342	     14015 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/unwind_stack__________/valyala-fastjson-10             	     194	   5847743 ns/op	52443035 B/op	    4143 allocs/op
BenchmarkValid/unwind_stack__________/goccy-go-json___-10             	    7741	    148331 ns/op	  102298 B/op	    4105 allocs/op
BenchmarkValid/unwind_stack__________/bytedance-sonic_-10             	  232706	      5132 ns/op	      24 B/op	       1 allocs/op
BenchmarkValid/miniscule_1b__________/jscan___________-10             	100000000	        11.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/encoding-json___-10             	66638437	        18.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/jsoniter________-10             	21767868	        54.35 ns/op	      16 B/op	       1 allocs/op
BenchmarkValid/miniscule_1b__________/gofaster-jx_____-10             	87644018	        13.66 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/tidwallgjson____-10             	214682797	         5.585 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/valyala-fastjson-10             	138021925	         8.688 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/goccy-go-json___-10             	 5881270	       202.5 ns/op	     704 B/op	       5 allocs/op
BenchmarkValid/miniscule_1b__________/bytedance-sonic_-10             	64177986	        18.69 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/jscan___________-10             	67845481	        17.69 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/encoding-json___-10             	25996744	        45.27 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/jsoniter________-10             	26885408	        44.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/gofaster-jx_____-10             	40406250	        29.59 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/tidwallgjson____-10             	74256657	        16.15 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/valyala-fastjson-10             	60297847	        19.86 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/goccy-go-json___-10             	 3338076	       358.1 ns/op	    1072 B/op	       9 allocs/op
BenchmarkValid/tiny_8b_______________/bytedance-sonic_-10             	25624326	        46.74 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/jscan___________-10             	 4849407	       247.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/encoding-json___-10             	 1325954	       904.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/jsoniter________-10             	 1749552	       680.3 ns/op	      56 B/op	       7 allocs/op
BenchmarkValid/small_336b____________/gofaster-jx_____-10             	 3150200	       380.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/tidwallgjson____-10             	 3563703	       336.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/valyala-fastjson-10             	 3201638	       375.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/goccy-go-json___-10             	  478029	      2498 ns/op	    2866 B/op	      61 allocs/op
BenchmarkValid/small_336b____________/bytedance-sonic_-10             	 1324634	       914.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/jscan___________-10             	     100	  11160137 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/encoding-json___-10             	      16	  68620914 ns/op	      92 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/jsoniter________-10             	      25	  43679673 ns/op	13582690 B/op	  644360 allocs/op
BenchmarkValid/large_26m_____________/gofaster-jx_____-10             	      57	  20582050 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/tidwallgjson____-10             	      43	  27191413 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/valyala-fastjson-10             	      45	  25724133 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/goccy-go-json___-10             	       1	7218929625 ns/op	144669928 B/op	 2338258 allocs/op
BenchmarkValid/large_26m_____________/bytedance-sonic_-10             	      16	  68641424 ns/op	      80 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/jscan___________-10             	   13716	     87631 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/encoding-json___-10             	    3349	    357536 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/jsoniter________-10             	    4935	    237420 ns/op	   69236 B/op	    2121 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/gofaster-jx_____-10             	    8491	    140075 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/tidwallgjson____-10             	    9327	    128191 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/valyala-fastjson-10             	    4172	    286538 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/goccy-go-json___-10             	     418	   2920799 ns/op	  780737 B/op	   20801 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/bytedance-sonic_-10             	    3338	    357336 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/jscan___________-10             	  863568	      1388 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/encoding-json___-10             	  128864	      9283 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/jsoniter________-10             	  149144	      7823 ns/op	    2064 B/op	      15 allocs/op
BenchmarkValid/escaped_3k____________/gofaster-jx_____-10             	  204566	      5828 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/tidwallgjson____-10             	  400201	      2992 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/valyala-fastjson-10             	  131952	      9066 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/goccy-go-json___-10             	   80352	     14817 ns/op	    4480 B/op	      13 allocs/op
BenchmarkValid/escaped_3k____________/bytedance-sonic_-10             	  128760	      9284 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/jscan___________-10             	  122386	      9518 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/encoding-json___-10             	   35683	     33596 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/jsoniter________-10             	   57631	     20813 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/gofaster-jx_____-10             	   67642	     17537 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/tidwallgjson____-10             	   88839	     13472 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/valyala-fastjson-10             	   80936	     14705 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/goccy-go-json___-10             	   10000	    100703 ns/op	   73470 B/op	    2057 allocs/op
BenchmarkValid/array_int_1024_12k____/bytedance-sonic_-10             	   35652	     33584 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/jscan___________-10             	  136593	      8692 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/encoding-json___-10             	   34257	     35244 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/jsoniter________-10             	   17952	     66589 ns/op	    8755 B/op	     547 allocs/op
BenchmarkValid/array_dec_1024_10k____/gofaster-jx_____-10             	   52272	     23103 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/tidwallgjson____-10             	  104120	     11481 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/valyala-fastjson-10             	   75606	     15445 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/goccy-go-json___-10             	   10000	    105016 ns/op	   73466 B/op	    2057 allocs/op
BenchmarkValid/array_dec_1024_10k____/bytedance-sonic_-10             	   32086	     34394 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/jscan___________-10             	  323750	      3467 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/encoding-json___-10             	   58966	     20212 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/jsoniter________-10             	   69603	     17280 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/gofaster-jx_____-10             	   57130	     20754 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/tidwallgjson____-10             	  204843	      5619 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/valyala-fastjson-10             	  237021	      4961 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/goccy-go-json___-10             	   26042	     45734 ns/op	   48909 B/op	    1036 allocs/op
BenchmarkValid/array_nullbool_1024_5k/bytedance-sonic_-10             	   58744	     20296 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/jscan___________-10             	    8374	    143025 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/encoding-json___-10             	     859	   1390473 ns/op	       1 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/jsoniter________-10             	    2352	    506790 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/gofaster-jx_____-10             	    7866	    152462 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/tidwallgjson____-10             	    2384	    501218 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/valyala-fastjson-10             	    4720	    253760 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/goccy-go-json___-10             	    1370	    841067 ns/op	 2817342 B/op	    3081 allocs/op
BenchmarkValid/array_str_1024_639k___/bytedance-sonic_-10             	     860	   1390300 ns/op	       1 B/op	       0 allocs/op
PASS
ok  	github.com/romshark/jscan/v2	192.046s
```

</details>

### AMD Ryzen 5 3600 - Linux

<details>

```
goos: linux
goarch: amd64
pkg: github.com/romshark/jscan/v2
cpu: AMD Ryzen 5 3600 6-Core Processor
BenchmarkCalcStats/miniscule_1b__________/jscan___________-12         	29536280	        38.06 ns/op
BenchmarkCalcStats/miniscule_1b__________/jsoniter________-12         	13069328	        89.97 ns/op
BenchmarkCalcStats/miniscule_1b__________/gofaster-jx_____-12         	55930816	        20.88 ns/op
BenchmarkCalcStats/miniscule_1b__________/valyala-fastjson-12         	45276549	        25.97 ns/op
BenchmarkCalcStats/tiny_8b_______________/jscan___________-12         	22317489	        51.13 ns/op
BenchmarkCalcStats/tiny_8b_______________/jsoniter________-12         	 8729518	       145.3 ns/op
BenchmarkCalcStats/tiny_8b_______________/gofaster-jx_____-12         	22108088	        52.41 ns/op
BenchmarkCalcStats/tiny_8b_______________/valyala-fastjson-12         	15440025	        74.07 ns/op
BenchmarkCalcStats/small_336b____________/jscan___________-12         	 2083824	       553.0 ns/op
BenchmarkCalcStats/small_336b____________/jsoniter________-12         	  746019	      1751 ns/op
BenchmarkCalcStats/small_336b____________/gofaster-jx_____-12         	 1492614	       780.9 ns/op
BenchmarkCalcStats/small_336b____________/valyala-fastjson-12         	 1474159	       803.9 ns/op
BenchmarkCalcStats/large_26m_____________/jscan___________-12         	      51	  21989699 ns/op
BenchmarkCalcStats/large_26m_____________/jsoniter________-12         	      14	 106006860 ns/op
BenchmarkCalcStats/large_26m_____________/gofaster-jx_____-12         	      32	  34525350 ns/op
BenchmarkCalcStats/large_26m_____________/valyala-fastjson-12         	      22	  47573176 ns/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/jscan___________-12         	    6600	    183015 ns/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/jsoniter________-12         	    1416	    993701 ns/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/gofaster-jx_____-12         	    4070	    288317 ns/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/valyala-fastjson-12         	    3777	    311561 ns/op
BenchmarkCalcStats/escaped_3k____________/jscan___________-12         	  592038	      1951 ns/op
BenchmarkCalcStats/escaped_3k____________/jsoniter________-12         	   62140	     19350 ns/op
BenchmarkCalcStats/escaped_3k____________/gofaster-jx_____-12         	   82678	     12159 ns/op
BenchmarkCalcStats/escaped_3k____________/valyala-fastjson-12         	   86558	     13867 ns/op
BenchmarkCalcStats/array_int_1024_12k____/jscan___________-12         	   53991	     19322 ns/op
BenchmarkCalcStats/array_int_1024_12k____/jsoniter________-12         	   10000	    112766 ns/op
BenchmarkCalcStats/array_int_1024_12k____/gofaster-jx_____-12         	   34392	     32628 ns/op
BenchmarkCalcStats/array_int_1024_12k____/valyala-fastjson-12         	   41499	     25589 ns/op
BenchmarkCalcStats/array_dec_1024_10k____/jscan___________-12         	   49053	     21098 ns/op
BenchmarkCalcStats/array_dec_1024_10k____/jsoniter________-12         	   10000	    119097 ns/op
BenchmarkCalcStats/array_dec_1024_10k____/gofaster-jx_____-12         	   25051	     44707 ns/op
BenchmarkCalcStats/array_dec_1024_10k____/valyala-fastjson-12         	   37179	     29858 ns/op
BenchmarkCalcStats/array_nullbool_1024_5k/jscan___________-12         	  131551	      9163 ns/op
BenchmarkCalcStats/array_nullbool_1024_5k/jsoniter________-12         	   41341	     25579 ns/op
BenchmarkCalcStats/array_nullbool_1024_5k/gofaster-jx_____-12         	   34100	     32163 ns/op
BenchmarkCalcStats/array_nullbool_1024_5k/valyala-fastjson-12         	   81313	     15249 ns/op
BenchmarkCalcStats/array_str_1024_639k___/jscan___________-12         	    5541	    203461 ns/op
BenchmarkCalcStats/array_str_1024_639k___/jsoniter________-12         	    1015	   1201872 ns/op
BenchmarkCalcStats/array_str_1024_639k___/gofaster-jx_____-12         	    4980	    235475 ns/op
BenchmarkCalcStats/array_str_1024_639k___/valyala-fastjson-12         	   16449	     70953 ns/op
BenchmarkValid/deeparray_____________/jscan___________-12             	35657382	        33.19 ns/op
BenchmarkValid/deeparray_____________/encoding-json___-12             	 2833788	       427.1 ns/op
BenchmarkValid/deeparray_____________/jsoniter________-12             	  947551	      1172 ns/op
BenchmarkValid/deeparray_____________/gofaster-jx_____-12             	 1220472	       981.8 ns/op
BenchmarkValid/deeparray_____________/tidwallgjson____-12             	219869212	         5.017 ns/op
BenchmarkValid/deeparray_____________/valyala-fastjson-12             	  369481	      3046 ns/op
BenchmarkValid/deeparray_____________/goccy-go-json___-12             	    4824	    235055 ns/op
BenchmarkValid/deeparray_____________/bytedance-sonic_-12             	48642750	        24.55 ns/op
BenchmarkValid/unwind_stack__________/jscan___________-12             	  564374	      2019 ns/op
BenchmarkValid/unwind_stack__________/encoding-json___-12             	  188680	      6377 ns/op
BenchmarkValid/unwind_stack__________/jsoniter________-12             	    6768	    157657 ns/op
BenchmarkValid/unwind_stack__________/gofaster-jx_____-12             	     963	   1147099 ns/op
BenchmarkValid/unwind_stack__________/tidwallgjson____-12             	  115238	     10386 ns/op
BenchmarkValid/unwind_stack__________/valyala-fastjson-12             	      66	  17407013 ns/op
BenchmarkValid/unwind_stack__________/goccy-go-json___-12             	    2392	    479779 ns/op
BenchmarkValid/unwind_stack__________/bytedance-sonic_-12             	  297740	      3919 ns/op
BenchmarkValid/miniscule_1b__________/jscan___________-12             	41000984	        28.06 ns/op
BenchmarkValid/miniscule_1b__________/encoding-json___-12             	34084412	        36.02 ns/op
BenchmarkValid/miniscule_1b__________/jsoniter________-12             	 5943859	       194.2 ns/op
BenchmarkValid/miniscule_1b__________/gofaster-jx_____-12             	68613285	        16.76 ns/op
BenchmarkValid/miniscule_1b__________/tidwallgjson____-12             	153976021	         7.833 ns/op
BenchmarkValid/miniscule_1b__________/valyala-fastjson-12             	114110142	         9.776 ns/op
BenchmarkValid/miniscule_1b__________/goccy-go-json___-12             	 1682247	       707.7 ns/op
BenchmarkValid/miniscule_1b__________/bytedance-sonic_-12             	36644451	        31.71 ns/op
BenchmarkValid/tiny_8b_______________/jscan___________-12             	39070484	        29.76 ns/op
BenchmarkValid/tiny_8b_______________/encoding-json___-12             	17939155	        67.83 ns/op
BenchmarkValid/tiny_8b_______________/jsoniter________-12             	19664109	        57.81 ns/op
BenchmarkValid/tiny_8b_______________/gofaster-jx_____-12             	27540657	        39.81 ns/op
BenchmarkValid/tiny_8b_______________/tidwallgjson____-12             	41461370	        27.44 ns/op
BenchmarkValid/tiny_8b_______________/valyala-fastjson-12             	42391431	        28.77 ns/op
BenchmarkValid/tiny_8b_______________/goccy-go-json___-12             	  897928	      1336 ns/op
BenchmarkValid/tiny_8b_______________/bytedance-sonic_-12             	20842724	        54.44 ns/op
BenchmarkValid/small_336b____________/jscan___________-12             	 3121490	       352.0 ns/op
BenchmarkValid/small_336b____________/encoding-json___-12             	  958030	      1221 ns/op
BenchmarkValid/small_336b____________/jsoniter________-12             	  796832	      1404 ns/op
BenchmarkValid/small_336b____________/gofaster-jx_____-12             	 2120688	       531.7 ns/op
BenchmarkValid/small_336b____________/tidwallgjson____-12             	 2512659	       453.9 ns/op
BenchmarkValid/small_336b____________/valyala-fastjson-12             	 2375257	       481.2 ns/op
BenchmarkValid/small_336b____________/goccy-go-json___-12             	  134289	      9373 ns/op
BenchmarkValid/small_336b____________/bytedance-sonic_-12             	 1870857	       647.2 ns/op
BenchmarkValid/large_26m_____________/jscan___________-12             	      63	  18528704 ns/op
BenchmarkValid/large_26m_____________/encoding-json___-12             	      15	  72678689 ns/op
BenchmarkValid/large_26m_____________/jsoniter________-12             	      15	  75510329 ns/op
BenchmarkValid/large_26m_____________/gofaster-jx_____-12             	      43	  26559201 ns/op
BenchmarkValid/large_26m_____________/tidwallgjson____-12             	      37	  30608782 ns/op
BenchmarkValid/large_26m_____________/valyala-fastjson-12             	      36	  33153089 ns/op
BenchmarkValid/large_26m_____________/goccy-go-json___-12             	       1	25530262976 ns/op
BenchmarkValid/large_26m_____________/bytedance-sonic_-12             	      56	  19892450 ns/op
BenchmarkValid/nasa_SxSW_2016_125k___/jscan___________-12             	    8973	    129394 ns/op
BenchmarkValid/nasa_SxSW_2016_125k___/encoding-json___-12             	    2739	    421571 ns/op
BenchmarkValid/nasa_SxSW_2016_125k___/jsoniter________-12             	    3480	    634091 ns/op
BenchmarkValid/nasa_SxSW_2016_125k___/gofaster-jx_____-12             	    6297	    179343 ns/op
BenchmarkValid/nasa_SxSW_2016_125k___/tidwallgjson____-12             	    7483	    156944 ns/op
BenchmarkValid/nasa_SxSW_2016_125k___/valyala-fastjson-12             	    4473	    269180 ns/op
BenchmarkValid/nasa_SxSW_2016_125k___/goccy-go-json___-12             	     212	   5740529 ns/op
BenchmarkValid/nasa_SxSW_2016_125k___/bytedance-sonic_-12             	    7290	    163337 ns/op
BenchmarkValid/escaped_3k____________/jscan___________-12             	  639826	      1836 ns/op
BenchmarkValid/escaped_3k____________/encoding-json___-12             	  116306	     10565 ns/op
BenchmarkValid/escaped_3k____________/jsoniter________-12             	   60562	     19184 ns/op
BenchmarkValid/escaped_3k____________/gofaster-jx_____-12             	  181603	      6546 ns/op
BenchmarkValid/escaped_3k____________/tidwallgjson____-12             	  387321	      2884 ns/op
BenchmarkValid/escaped_3k____________/valyala-fastjson-12             	  137722	      8484 ns/op
BenchmarkValid/escaped_3k____________/goccy-go-json___-12             	   33592	     35131 ns/op
BenchmarkValid/escaped_3k____________/bytedance-sonic_-12             	 4022890	       269.0 ns/op
BenchmarkValid/array_int_1024_12k____/jscan___________-12             	   88560	     13372 ns/op
BenchmarkValid/array_int_1024_12k____/encoding-json___-12             	   34699	     35938 ns/op
BenchmarkValid/array_int_1024_12k____/jsoniter________-12             	   48398	     23861 ns/op
BenchmarkValid/array_int_1024_12k____/gofaster-jx_____-12             	   51304	     20024 ns/op
BenchmarkValid/array_int_1024_12k____/tidwallgjson____-12             	   81090	     15038 ns/op
BenchmarkValid/array_int_1024_12k____/valyala-fastjson-12             	   59917	     17180 ns/op
BenchmarkValid/array_int_1024_12k____/goccy-go-json___-12             	    4555	    337608 ns/op
BenchmarkValid/array_int_1024_12k____/bytedance-sonic_-12             	   50977	     19822 ns/op
BenchmarkValid/array_dec_1024_10k____/jscan___________-12             	   81597	     14911 ns/op
BenchmarkValid/array_dec_1024_10k____/encoding-json___-12             	   26898	     43262 ns/op
BenchmarkValid/array_dec_1024_10k____/jsoniter________-12             	   10000	    186437 ns/op
BenchmarkValid/array_dec_1024_10k____/gofaster-jx_____-12             	   38944	     27930 ns/op
BenchmarkValid/array_dec_1024_10k____/tidwallgjson____-12             	   71010	     17063 ns/op
BenchmarkValid/array_dec_1024_10k____/valyala-fastjson-12             	   47227	     22618 ns/op
BenchmarkValid/array_dec_1024_10k____/goccy-go-json___-12             	    3338	    349254 ns/op
BenchmarkValid/array_dec_1024_10k____/bytedance-sonic_-12             	   43161	     24447 ns/op
BenchmarkValid/array_nullbool_1024_5k/jscan___________-12             	  291025	      4224 ns/op
BenchmarkValid/array_nullbool_1024_5k/encoding-json___-12             	   48828	     21422 ns/op
BenchmarkValid/array_nullbool_1024_5k/jsoniter________-12             	   53788	     18773 ns/op
BenchmarkValid/array_nullbool_1024_5k/gofaster-jx_____-12             	   62220	     20396 ns/op
BenchmarkValid/array_nullbool_1024_5k/tidwallgjson____-12             	  161696	      7254 ns/op
BenchmarkValid/array_nullbool_1024_5k/valyala-fastjson-12             	  145792	      8197 ns/op
BenchmarkValid/array_nullbool_1024_5k/goccy-go-json___-12             	    7735	    158039 ns/op
BenchmarkValid/array_nullbool_1024_5k/bytedance-sonic_-12             	   71770	     13938 ns/op
BenchmarkValid/array_str_1024_639k___/jscan___________-12             	    5163	    228909 ns/op
BenchmarkValid/array_str_1024_639k___/encoding-json___-12             	     817	   1435954 ns/op
BenchmarkValid/array_str_1024_639k___/jsoniter________-12             	    2497	    488499 ns/op
BenchmarkValid/array_str_1024_639k___/gofaster-jx_____-12             	    5104	    223810 ns/op
BenchmarkValid/array_str_1024_639k___/tidwallgjson____-12             	    2422	    481406 ns/op
BenchmarkValid/array_str_1024_639k___/valyala-fastjson-12             	    3362	    366067 ns/op
BenchmarkValid/array_str_1024_639k___/goccy-go-json___-12             	     543	   2236849 ns/op
BenchmarkValid/array_str_1024_639k___/bytedance-sonic_-12             	   26715	     42715 ns/op
PASS
ok  	github.com/romshark/jscan/v2	257.283s
```

</details>

### Intel i7-3930K - Linux

<details>

```
goos: linux
goarch: amd64
pkg: github.com/romshark/jscan/v2
cpu: Intel(R) Core(TM) i7-3930K CPU @ 3.20GHz
BenchmarkCalcStats/miniscule_1b__________/jscan___________-12         	2380541        45.38 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/miniscule_1b__________/jsoniter________-12         	 701319       167.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkCalcStats/miniscule_1b__________/gofaster-jx_____-12         	3282648        34.36 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/miniscule_1b__________/valyala-fastjson-12         	3235749        34.87 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/tiny_8b_______________/jscan___________-12         	17619775	        69.18 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/tiny_8b_______________/jsoniter________-12         	 5898264	       263.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkCalcStats/tiny_8b_______________/gofaster-jx_____-12         	14033302	        82.65 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/tiny_8b_______________/valyala-fastjson-12         	13171170	        92.52 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/small_336b____________/jscan___________-12         	 1684648	       649.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/small_336b____________/jsoniter________-12         	  610396	      2299 ns/op	      80 B/op	      11 allocs/op
BenchmarkCalcStats/small_336b____________/gofaster-jx_____-12         	  994642	      1036 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/small_336b____________/valyala-fastjson-12         	 1061858	       973.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/large_26m_____________/jscan___________-12         	      37	  28861004 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/large_26m_____________/jsoniter________-12         	       8	 151542467 ns/op	32851282 B/op	 1108518 allocs/op
BenchmarkCalcStats/large_26m_____________/gofaster-jx_____-12         	      22	  49226281 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/large_26m_____________/valyala-fastjson-12         	      15	  70521173 ns/op	22457962 B/op	   22063 allocs/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/jscan___________-12         	    4712	    242357 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/jsoniter________-12         	     609	   2042643 ns/op	  144472 B/op	    7357 allocs/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/gofaster-jx_____-12         	    2742	    410905 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/valyala-fastjson-12         	    2718	    443195 ns/op	     871 B/op	       1 allocs/op
BenchmarkCalcStats/escaped_3k____________/jscan___________-12         	  381399	      2921 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/escaped_3k____________/jsoniter________-12         	   30962	     39746 ns/op	    2064 B/op	      15 allocs/op
BenchmarkCalcStats/escaped_3k____________/gofaster-jx_____-12         	   54873	     18626 ns/op	     504 B/op	       6 allocs/op
BenchmarkCalcStats/escaped_3k____________/valyala-fastjson-12         	   60566	     18513 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_int_1024_12k____/jscan___________-12         	   38097	     26954 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_int_1024_12k____/jsoniter________-12         	    4539	    226520 ns/op	   16384 B/op	    1024 allocs/op
BenchmarkCalcStats/array_int_1024_12k____/gofaster-jx_____-12         	   22924	     48433 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_int_1024_12k____/valyala-fastjson-12         	   28704	     37840 ns/op	      12 B/op	       0 allocs/op
BenchmarkCalcStats/array_dec_1024_10k____/jscan___________-12         	   33211	     31693 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_dec_1024_10k____/jsoniter________-12         	    4372	    244692 ns/op	   16384 B/op	    1024 allocs/op
BenchmarkCalcStats/array_dec_1024_10k____/gofaster-jx_____-12         	   17073	     65184 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_dec_1024_10k____/valyala-fastjson-12         	   26547	     41406 ns/op	      13 B/op	       0 allocs/op
BenchmarkCalcStats/array_nullbool_1024_5k/jscan___________-12         	   81963	     14344 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_nullbool_1024_5k/jsoniter________-12         	   30834	     35080 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_nullbool_1024_5k/gofaster-jx_____-12         	   25578	     45238 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_nullbool_1024_5k/valyala-fastjson-12         	   52755	     21375 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_str_1024_639k___/jscan___________-12         	    4482	    247600 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_str_1024_639k___/jsoniter________-12         	     378	   3391742 ns/op	  670172 B/op	    1018 allocs/op
BenchmarkCalcStats/array_str_1024_639k___/gofaster-jx_____-12         	    3843	    291812 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_str_1024_639k___/valyala-fastjson-12         	    6870	    150559 ns/op	     143 B/op	       0 allocs/op
BenchmarkValid/deeparray_____________/jscan___________-12             	31938502	        34.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/deeparray_____________/encoding-json___-12             	 1376350	       900.2 ns/op	     104 B/op	       5 allocs/op
BenchmarkValid/deeparray_____________/jsoniter________-12             	  428743	      2544 ns/op	     352 B/op	       9 allocs/op
BenchmarkValid/deeparray_____________/gofaster-jx_____-12             	  748923	      1692 ns/op	      80 B/op	       2 allocs/op
BenchmarkValid/deeparray_____________/tidwallgjson____-12             	123089697	         8.947 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/deeparray_____________/valyala-fastjson-12             	  160580	      7041 ns/op	    1184 B/op	      11 allocs/op
BenchmarkValid/deeparray_____________/goccy-go-json___-12             	    2145	    514302 ns/op	   49327 B/op	    2062 allocs/op
BenchmarkValid/deeparray_____________/bytedance-sonic_-12             	37683284	        27.90 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/unwind_stack__________/jscan___________-12             	  301344	      3537 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/unwind_stack__________/encoding-json___-12             	  112620	     10373 ns/op	      24 B/op	       1 allocs/op
BenchmarkValid/unwind_stack__________/jsoniter________-12             	    3033	    346668 ns/op	   33159 B/op	    1033 allocs/op
BenchmarkValid/unwind_stack__________/gofaster-jx_____-12             	     596	   1971132 ns/op	   65664 B/op	    1026 allocs/op
BenchmarkValid/unwind_stack__________/tidwallgjson____-12             	  142334	      8537 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/unwind_stack__________/valyala-fastjson-12             	      22	  53200125 ns/op	52453560 B/op	    4141 allocs/op
BenchmarkValid/unwind_stack__________/goccy-go-json___-12             	     837	   1354517 ns/op	  102342 B/op	    4105 allocs/op
BenchmarkValid/unwind_stack__________/bytedance-sonic_-12             	  216874	      5249 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/jscan___________-12             	41986868	        29.68 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/encoding-json___-12             	28363678	        44.30 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/jsoniter________-12             	 5517650	       321.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkValid/miniscule_1b__________/gofaster-jx_____-12             	40494812	        27.49 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/tidwallgjson____-12             	98816259	        12.14 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/valyala-fastjson-12             	75405946	        16.62 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/goccy-go-json___-12             	  729242	      1436 ns/op	     704 B/op	       5 allocs/op
BenchmarkValid/miniscule_1b__________/bytedance-sonic_-12             	31713258	        35.96 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/jscan___________-12             	25472898	        42.94 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/encoding-json___-12             	15659935	        79.36 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/jsoniter________-12             	14986729	        78.35 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/gofaster-jx_____-12             	19895846	        59.32 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/tidwallgjson____-12             	37944624	        33.31 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/valyala-fastjson-12             	28123518	        39.87 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/goccy-go-json___-12             	  554565	      2614 ns/op	    1072 B/op	       9 allocs/op
BenchmarkValid/tiny_8b_______________/bytedance-sonic_-12             	18913327	        60.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/jscan___________-12             	 2559782	       492.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/encoding-json___-12             	  725848	      1426 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/jsoniter________-12             	  609745	      2231 ns/op	      56 B/op	       7 allocs/op
BenchmarkValid/small_336b____________/gofaster-jx_____-12             	 1461477	       765.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/tidwallgjson____-12             	 2018919	       624.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/valyala-fastjson-12             	 1856380	       620.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/goccy-go-json___-12             	   72609	     18546 ns/op	    2867 B/op	      61 allocs/op
BenchmarkValid/small_336b____________/bytedance-sonic_-12             	 1116440	       987.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/jscan___________-12             	      45	  23841660 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/encoding-json___-12             	       9	 116953436 ns/op	     171 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/jsoniter________-12             	      10	 104225305 ns/op	13582885 B/op	  644361 allocs/op
BenchmarkValid/large_26m_____________/gofaster-jx_____-12             	      30	  37663437 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/tidwallgjson____-12             	      22	  50724139 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/valyala-fastjson-12             	      26	  45830408 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/goccy-go-json___-12             	       1	29805696498 ns/op	144651488 B/op	 2338192 allocs/op
BenchmarkValid/large_26m_____________/bytedance-sonic_-12             	      36	  33382804 ns/op	    1180 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/jscan___________-12             	    6804	    185507 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/encoding-json___-12             	    2020	    588064 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/jsoniter________-12             	    1066	   1297357 ns/op	   69247 B/op	    2121 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/gofaster-jx_____-12             	    4688	    254605 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/tidwallgjson____-12             	    4382	    252564 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/valyala-fastjson-12             	    3228	    347860 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/goccy-go-json___-12             	      58	  22098508 ns/op	  780453 B/op	   20800 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/bytedance-sonic_-12             	    5460	    202111 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/jscan___________-12             	  379418	      2827 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/encoding-json___-12             	   83638	     12846 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/jsoniter________-12             	   27159	     38728 ns/op	    2065 B/op	      15 allocs/op
BenchmarkValid/escaped_3k____________/gofaster-jx_____-12             	  111859	     10460 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/tidwallgjson____-12             	  219693	      5327 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/valyala-fastjson-12             	   93456	     11917 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/goccy-go-json___-12             	   14906	     85299 ns/op	    4480 B/op	      13 allocs/op
BenchmarkValid/escaped_3k____________/bytedance-sonic_-12             	 2185236	       472.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/jscan___________-12             	   60303	     20032 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/encoding-json___-12             	   21495	     47921 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/jsoniter________-12             	   29016	     37251 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/gofaster-jx_____-12             	   42692	     29489 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/tidwallgjson____-12             	   51464	     24469 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/valyala-fastjson-12             	   42097	     27093 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/goccy-go-json___-12             	    1686	    696244 ns/op	   73502 B/op	    2057 allocs/op
BenchmarkValid/array_int_1024_12k____/bytedance-sonic_-12             	   49515	     22477 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/jscan___________-12             	   50834	     22318 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/encoding-json___-12             	   21704	     52112 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/jsoniter________-12             	    4500	    266575 ns/op	    8756 B/op	     547 allocs/op
BenchmarkValid/array_dec_1024_10k____/gofaster-jx_____-12             	   23377	     46317 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/tidwallgjson____-12             	   43099	     29381 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/valyala-fastjson-12             	   35184	     35693 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/goccy-go-json___-12             	    1640	    687252 ns/op	   73544 B/op	    2057 allocs/op
BenchmarkValid/array_dec_1024_10k____/bytedance-sonic_-12             	   44972	     27603 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/jscan___________-12             	  142741	      8794 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/encoding-json___-12             	   39296	     25989 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/jsoniter________-12             	   40736	     24746 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/gofaster-jx_____-12             	   36211	     28741 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/tidwallgjson____-12             	  107116	     10878 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/valyala-fastjson-12             	  100891	     11356 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/goccy-go-json___-12             	    3442	    318392 ns/op	   48944 B/op	    1036 allocs/op
BenchmarkValid/array_nullbool_1024_5k/bytedance-sonic_-12             	   80179	     13957 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/jscan___________-12             	    4534	    256575 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/encoding-json___-12             	     439	   2519619 ns/op	       3 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/jsoniter________-12             	    1292	    880043 ns/op	       1 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/gofaster-jx_____-12             	    4042	    274603 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/tidwallgjson____-12             	    1291	    869398 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/valyala-fastjson-12             	    2434	    458109 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/goccy-go-json___-12             	     181	   6589450 ns/op	 2817336 B/op	    3080 allocs/op
BenchmarkValid/array_str_1024_639k___/bytedance-sonic_-12             	   12112	     90231 ns/op	       3 B/op	       0 allocs/op
PASS
ok  	github.com/romshark/jscan/v2	263.004s
```

</details>

|package|version|
|-|-|
|pkg.go.dev/encoding/json|[go1.20.5](https://pkg.go.dev/encoding/json)|
|github.com/go-faster/jx|[v1.0.0](https://github.com/go-faster/jx/releases/tag/v1.0.0)|
|github.com/json-iterator/go|[v1.1.12](https://github.com/json-iterator/go/releases/tag/v1.1.12)|
|github.com/tidwall/gjson|[v1.14.4](https://github.com/tidwall/gjson/releases/tag/v1.14.4)|
|github.com/valyala/fastjson|[v1.6.4](https://github.com/valyala/fastjson/releases/tag/v1.6.4)|
|github.com/goccy/go-json|[v0.10.2](https://github.com/goccy/go-json/releases/tag/v0.10.2)|
|github.com/bytedance/sonic|[v1.9.1](https://github.com/bytedance/sonic/releases/tag/v1.9.1)|
