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

### Apple M1 - macOS

<details>

```
goos: darwin
goarch: arm64
pkg: github.com/romshark/jscan/v2
BenchmarkCalcStats/miniscule_1b__________/jscan___________-10         	59269868	        20.30 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/miniscule_1b__________/jsoniter________-10         	47673040	        24.97 ns/op	      16 B/op	       1 allocs/op
BenchmarkCalcStats/miniscule_1b__________/gofaster-jx_____-10         	66381014	        18.34 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/miniscule_1b__________/valyala-fastjson-10         	72915637	        16.45 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/tiny_8b_______________/jscan___________-10         	43010366	        27.75 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/tiny_8b_______________/jsoniter________-10         	26192532	        45.33 ns/op	      16 B/op	       1 allocs/op
BenchmarkCalcStats/tiny_8b_______________/gofaster-jx_____-10         	27700084	        43.41 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/tiny_8b_______________/valyala-fastjson-10         	29007373	        41.36 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/small_336b____________/jscan___________-10         	 3803469	       315.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/small_336b____________/jsoniter________-10         	 1719595	       698.5 ns/op	      80 B/op	      11 allocs/op
BenchmarkCalcStats/small_336b____________/gofaster-jx_____-10         	 2169474	       554.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/small_336b____________/valyala-fastjson-10         	 2194273	       546.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/large_26m_____________/jscan___________-10         	      84	  13955350 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/large_26m_____________/jsoniter________-10         	      20	  54410746 ns/op	32851290 B/op	 1108518 allocs/op
BenchmarkCalcStats/large_26m_____________/gofaster-jx_____-10         	      42	  28031878 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/large_26m_____________/valyala-fastjson-10         	      37	  29620818 ns/op	 9104579 B/op	    8944 allocs/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/jscan___________-10         	   10000	    115466 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/jsoniter________-10         	    3243	    351182 ns/op	  144473 B/op	    7357 allocs/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/gofaster-jx_____-10         	    5226	    230626 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/valyala-fastjson-10         	    3534	    337400 ns/op	     670 B/op	       1 allocs/op
BenchmarkCalcStats/escaped_3k____________/jscan___________-10         	  842656	      1401 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/escaped_3k____________/jsoniter________-10         	  149356	      7951 ns/op	    2064 B/op	      15 allocs/op
BenchmarkCalcStats/escaped_3k____________/gofaster-jx_____-10         	  180602	      6605 ns/op	     504 B/op	       6 allocs/op
BenchmarkCalcStats/escaped_3k____________/valyala-fastjson-10         	  105752	     11354 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_int_1024_12k____/jscan___________-10         	   86401	     13895 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_int_1024_12k____/jsoniter________-10         	   31449	     38011 ns/op	   16384 B/op	    1024 allocs/op
BenchmarkCalcStats/array_int_1024_12k____/gofaster-jx_____-10         	   40638	     29928 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_int_1024_12k____/valyala-fastjson-10         	   55292	     21746 ns/op	       6 B/op	       0 allocs/op
BenchmarkCalcStats/array_dec_1024_10k____/jscan___________-10         	   81942	     12865 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_dec_1024_10k____/jsoniter________-10         	   27829	     43117 ns/op	   16384 B/op	    1024 allocs/op
BenchmarkCalcStats/array_dec_1024_10k____/gofaster-jx_____-10         	   29832	     35876 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_dec_1024_10k____/valyala-fastjson-10         	   52190	     22375 ns/op	       6 B/op	       0 allocs/op
BenchmarkCalcStats/array_nullbool_1024_5k/jscan___________-10         	  159696	      7398 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_nullbool_1024_5k/jsoniter________-10         	   56115	     21209 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_nullbool_1024_5k/gofaster-jx_____-10         	   36439	     33163 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_nullbool_1024_5k/valyala-fastjson-10         	  116670	     10292 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_str_1024_639k___/jscan___________-10         	    8140	    146574 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_str_1024_639k___/jsoniter________-10         	    1945	    599316 ns/op	  670172 B/op	    1018 allocs/op
BenchmarkCalcStats/array_str_1024_639k___/gofaster-jx_____-10         	    7257	    165041 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_str_1024_639k___/valyala-fastjson-10         	   19447	     60584 ns/op	      50 B/op	       0 allocs/op
BenchmarkValid/deeparray_____________/jscan___________-10             	73762368	        15.97 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/deeparray_____________/encoding-json___-10             	 8947839	       133.2 ns/op	     104 B/op	       5 allocs/op
BenchmarkValid/deeparray_____________/jsoniter________-10             	 3461888	       347.0 ns/op	     352 B/op	       9 allocs/op
BenchmarkValid/deeparray_____________/gofaster-jx_____-10             	 4347409	       276.1 ns/op	      80 B/op	       2 allocs/op
BenchmarkValid/deeparray_____________/tidwallgjson____-10             	296225592	         4.044 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/deeparray_____________/valyala-fastjson-10             	 1000000	      1028 ns/op	    1184 B/op	      11 allocs/op
BenchmarkValid/deeparray_____________/goccy-go-json___-10             	   15226	     78478 ns/op	   49295 B/op	    2062 allocs/op
BenchmarkValid/deeparray_____________/bytedance-sonic_-10             	 8964656	       133.7 ns/op	     104 B/op	       5 allocs/op
BenchmarkValid/unwind_stack__________/jscan___________-10             	  579032	      2071 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/unwind_stack__________/encoding-json___-10             	  232429	      5138 ns/op	      24 B/op	       1 allocs/op
BenchmarkValid/unwind_stack__________/jsoniter________-10             	   18069	     66620 ns/op	   33149 B/op	    1033 allocs/op
BenchmarkValid/unwind_stack__________/gofaster-jx_____-10             	    2990	    398754 ns/op	   65664 B/op	    1026 allocs/op
BenchmarkValid/unwind_stack__________/tidwallgjson____-10             	   85274	     14037 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/unwind_stack__________/valyala-fastjson-10             	     196	   5814109 ns/op	52454488 B/op	    4143 allocs/op
BenchmarkValid/unwind_stack__________/goccy-go-json___-10             	    7395	    154853 ns/op	  102298 B/op	    4105 allocs/op
BenchmarkValid/unwind_stack__________/bytedance-sonic_-10             	  232034	      5138 ns/op	      24 B/op	       1 allocs/op
BenchmarkValid/miniscule_1b__________/jscan___________-10             	99927835	        11.98 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/encoding-json___-10             	66537903	        18.03 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/jsoniter________-10             	21847411	        54.50 ns/op	      16 B/op	       1 allocs/op
BenchmarkValid/miniscule_1b__________/gofaster-jx_____-10             	89593746	        13.37 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/tidwallgjson____-10             	214460872	         5.622 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/valyala-fastjson-10             	138044894	         8.696 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/goccy-go-json___-10             	 5463507	       203.8 ns/op	     704 B/op	       5 allocs/op
BenchmarkValid/miniscule_1b__________/bytedance-sonic_-10             	64044830	        18.74 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/jscan___________-10             	68990955	        17.41 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/encoding-json___-10             	25837428	        45.28 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/jsoniter________-10             	26759084	        44.38 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/gofaster-jx_____-10             	40737576	        29.33 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/tidwallgjson____-10             	74113532	        16.44 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/valyala-fastjson-10             	58286025	        20.05 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/goccy-go-json___-10             	 3349239	       362.9 ns/op	    1072 B/op	       9 allocs/op
BenchmarkValid/tiny_8b_______________/bytedance-sonic_-10             	25552982	        45.47 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/jscan___________-10             	 5047078	       237.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/encoding-json___-10             	 1324250	       906.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/jsoniter________-10             	 1733530	       688.6 ns/op	      56 B/op	       7 allocs/op
BenchmarkValid/small_336b____________/gofaster-jx_____-10             	 3149773	       380.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/tidwallgjson____-10             	 3555520	       337.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/valyala-fastjson-10             	 3189453	       376.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/goccy-go-json___-10             	  475158	      2488 ns/op	    2866 B/op	      61 allocs/op
BenchmarkValid/small_336b____________/bytedance-sonic_-10             	 1322008	       907.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/jscan___________-10             	     100	  11094928 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/encoding-json___-10             	      16	  68734076 ns/op	      92 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/jsoniter________-10             	      25	  44049983 ns/op	13582697 B/op	  644361 allocs/op
BenchmarkValid/large_26m_____________/gofaster-jx_____-10             	      57	  20690905 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/tidwallgjson____-10             	      43	  27340072 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/valyala-fastjson-10             	      45	  25835529 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/goccy-go-json___-10             	       1	7223878208 ns/op	144621864 B/op	 2338098 allocs/op
BenchmarkValid/large_26m_____________/bytedance-sonic_-10             	      16	  68800419 ns/op	      80 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/jscan___________-10             	   13962	     85812 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/encoding-json___-10             	    3338	    358778 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/jsoniter________-10             	    4940	    239316 ns/op	   69236 B/op	    2121 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/gofaster-jx_____-10             	    8491	    140462 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/tidwallgjson____-10             	    9145	    130503 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/valyala-fastjson-10             	    4216	    284709 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/goccy-go-json___-10             	     416	   2917648 ns/op	  780821 B/op	   20801 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/bytedance-sonic_-10             	    3342	    357946 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/jscan___________-10             	  863632	      1393 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/encoding-json___-10             	  128832	      9295 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/jsoniter________-10             	  147468	      8023 ns/op	    2064 B/op	      15 allocs/op
BenchmarkValid/escaped_3k____________/gofaster-jx_____-10             	  203584	      5861 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/tidwallgjson____-10             	  401119	      2993 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/valyala-fastjson-10             	  131742	      9078 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/goccy-go-json___-10             	   80326	     14860 ns/op	    4480 B/op	      13 allocs/op
BenchmarkValid/escaped_3k____________/bytedance-sonic_-10             	  128679	      9298 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/jscan___________-10             	  126627	      9290 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/encoding-json___-10             	   35508	     33809 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/jsoniter________-10             	   58038	     20725 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/gofaster-jx_____-10             	   68139	     17856 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/tidwallgjson____-10             	   88935	     13486 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/valyala-fastjson-10             	   81699	     14597 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/goccy-go-json___-10             	   10000	    100410 ns/op	   73473 B/op	    2057 allocs/op
BenchmarkValid/array_int_1024_12k____/bytedance-sonic_-10             	   35664	     33710 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/jscan___________-10             	  136604	      8720 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/encoding-json___-10             	   34710	     34479 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/jsoniter________-10             	   18067	     66272 ns/op	    8755 B/op	     547 allocs/op
BenchmarkValid/array_dec_1024_10k____/gofaster-jx_____-10             	   59536	     20080 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/tidwallgjson____-10             	  101229	     11297 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/valyala-fastjson-10             	   63325	     16883 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/goccy-go-json___-10             	   10000	    104839 ns/op	   73469 B/op	    2057 allocs/op
BenchmarkValid/array_dec_1024_10k____/bytedance-sonic_-10             	   34363	     34519 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/jscan___________-10             	  333795	      3399 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/encoding-json___-10             	   58797	     20301 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/jsoniter________-10             	   72146	     16920 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/gofaster-jx_____-10             	   56263	     21241 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/tidwallgjson____-10             	  212868	      5505 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/valyala-fastjson-10             	  231019	      4857 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/goccy-go-json___-10             	   26331	     45382 ns/op	   48908 B/op	    1036 allocs/op
BenchmarkValid/array_nullbool_1024_5k/bytedance-sonic_-10             	   59068	     20323 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/jscan___________-10             	    8403	    142921 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/encoding-json___-10             	     859	   1392907 ns/op	       1 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/jsoniter________-10             	    2354	    507430 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/gofaster-jx_____-10             	    7844	    152642 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/tidwallgjson____-10             	    2380	    502848 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/valyala-fastjson-10             	    4714	    253981 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/goccy-go-json___-10             	    1384	    836698 ns/op	 2817330 B/op	    3081 allocs/op
BenchmarkValid/array_str_1024_639k___/bytedance-sonic_-10             	     859	   1392600 ns/op	       1 B/op	       0 allocs/op
PASS
ok  	github.com/romshark/jscan/v2	191.992s
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
BenchmarkCalcStats/miniscule_1b__________/jscan___________-12         	24948513	        46.74 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/miniscule_1b__________/jsoniter________-12         	 8972714	       178.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkCalcStats/miniscule_1b__________/gofaster-jx_____-12         	34174065	        32.55 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/miniscule_1b__________/valyala-fastjson-12         	33400308	        34.77 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/tiny_8b_______________/jscan___________-12         	17675479	        70.64 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/tiny_8b_______________/jsoniter________-12         	 4477546	       258.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkCalcStats/tiny_8b_______________/gofaster-jx_____-12         	14504192	        76.28 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/tiny_8b_______________/valyala-fastjson-12         	12178302	        90.75 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/small_336b____________/jscan___________-12         	 1932068	       653.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/small_336b____________/jsoniter________-12         	  519332	      2539 ns/op	      80 B/op	      11 allocs/op
BenchmarkCalcStats/small_336b____________/gofaster-jx_____-12         	 1098571	      1096 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/small_336b____________/valyala-fastjson-12         	 1059382	      1030 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/large_26m_____________/jscan___________-12         	      39	  27558632 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/large_26m_____________/jsoniter________-12         	       7	 148281244 ns/op	32851274 B/op	 1108518 allocs/op
BenchmarkCalcStats/large_26m_____________/gofaster-jx_____-12         	      24	  47501715 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/large_26m_____________/valyala-fastjson-12         	      16	  73184504 ns/op	21054339 B/op	   20684 allocs/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/jscan___________-12         	    4966	    238318 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/jsoniter________-12         	     585	   2040711 ns/op	  144473 B/op	    7357 allocs/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/gofaster-jx_____-12         	    3022	    408705 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/nasa_SxSW_2016_125k___/valyala-fastjson-12         	    2797	    435153 ns/op	     846 B/op	       1 allocs/op
BenchmarkCalcStats/escaped_3k____________/jscan___________-12         	  434394	      2892 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/escaped_3k____________/jsoniter________-12         	   31486	     37986 ns/op	    2064 B/op	      15 allocs/op
BenchmarkCalcStats/escaped_3k____________/gofaster-jx_____-12         	   55659	     20579 ns/op	     504 B/op	       6 allocs/op
BenchmarkCalcStats/escaped_3k____________/valyala-fastjson-12         	   60624	     18835 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_int_1024_12k____/jscan___________-12         	   46663	     27033 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_int_1024_12k____/jsoniter________-12         	    4386	    237690 ns/op	   16384 B/op	    1024 allocs/op
BenchmarkCalcStats/array_int_1024_12k____/gofaster-jx_____-12         	   23632	     48962 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_int_1024_12k____/valyala-fastjson-12         	   28450	     37031 ns/op	      12 B/op	       0 allocs/op
BenchmarkCalcStats/array_dec_1024_10k____/jscan___________-12         	   33384	     30885 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_dec_1024_10k____/jsoniter________-12         	   10000	    245693 ns/op	   16384 B/op	    1024 allocs/op
BenchmarkCalcStats/array_dec_1024_10k____/gofaster-jx_____-12         	   16720	     65810 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_dec_1024_10k____/valyala-fastjson-12         	   24985	     43489 ns/op	      14 B/op	       0 allocs/op
BenchmarkCalcStats/array_nullbool_1024_5k/jscan___________-12         	   85256	     14643 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_nullbool_1024_5k/jsoniter________-12         	   36159	     34770 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_nullbool_1024_5k/gofaster-jx_____-12         	   23634	     47297 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_nullbool_1024_5k/valyala-fastjson-12         	   53520	     22210 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_str_1024_639k___/jscan___________-12         	    4239	    262570 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_str_1024_639k___/jsoniter________-12         	     349	   3384101 ns/op	  670170 B/op	    1018 allocs/op
BenchmarkCalcStats/array_str_1024_639k___/gofaster-jx_____-12         	    4020	    275710 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/array_str_1024_639k___/valyala-fastjson-12         	    6916	    158612 ns/op	     142 B/op	       0 allocs/op
BenchmarkValid/deeparray_____________/jscan___________-12             	31556031	        35.70 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/deeparray_____________/encoding-json___-12             	 1353940	       904.7 ns/op	     104 B/op	       5 allocs/op
BenchmarkValid/deeparray_____________/jsoniter________-12             	  422845	      2454 ns/op	     352 B/op	       9 allocs/op
BenchmarkValid/deeparray_____________/gofaster-jx_____-12             	  731324	      1749 ns/op	      80 B/op	       2 allocs/op
BenchmarkValid/deeparray_____________/tidwallgjson____-12             	123818463	         8.951 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/deeparray_____________/valyala-fastjson-12             	  172442	      7090 ns/op	    1184 B/op	      11 allocs/op
BenchmarkValid/deeparray_____________/goccy-go-json___-12             	    2149	    528035 ns/op	   49317 B/op	    2062 allocs/op
BenchmarkValid/deeparray_____________/bytedance-sonic_-12             	37018498	        29.89 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/unwind_stack__________/jscan___________-12             	  324952	      3726 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/unwind_stack__________/encoding-json___-12             	  112858	      9809 ns/op	      24 B/op	       1 allocs/op
BenchmarkValid/unwind_stack__________/jsoniter________-12             	    3802	    344261 ns/op	   33160 B/op	    1033 allocs/op
BenchmarkValid/unwind_stack__________/gofaster-jx_____-12             	     762	   1904051 ns/op	   65664 B/op	    1026 allocs/op
BenchmarkValid/unwind_stack__________/tidwallgjson____-12             	  134110	      8508 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/unwind_stack__________/valyala-fastjson-12             	      24	  51779648 ns/op	52447407 B/op	    4142 allocs/op
BenchmarkValid/unwind_stack__________/goccy-go-json___-12             	     924	   1315229 ns/op	  102341 B/op	    4105 allocs/op
BenchmarkValid/unwind_stack__________/bytedance-sonic_-12             	  166339	      6986 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/jscan___________-12             	43311066	        27.82 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/encoding-json___-12             	25812458	        43.67 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/jsoniter________-12             	 4148998	       305.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkValid/miniscule_1b__________/gofaster-jx_____-12             	40261909	        27.52 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/tidwallgjson____-12             	85151316	        12.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/valyala-fastjson-12             	68117352	        15.81 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/miniscule_1b__________/goccy-go-json___-12             	  751773	      1414 ns/op	     704 B/op	       5 allocs/op
BenchmarkValid/miniscule_1b__________/bytedance-sonic_-12             	30546312	        37.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/jscan___________-12             	30611830	        39.16 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/encoding-json___-12             	15759042	        78.74 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/jsoniter________-12             	14405368	        76.36 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/gofaster-jx_____-12             	20923336	        59.06 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/tidwallgjson____-12             	33457417	        31.62 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/valyala-fastjson-12             	28318700	        39.34 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny_8b_______________/goccy-go-json___-12             	  417813	      2591 ns/op	    1072 B/op	       9 allocs/op
BenchmarkValid/tiny_8b_______________/bytedance-sonic_-12             	18562359	        64.30 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/jscan___________-12             	 2234793	       491.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/encoding-json___-12             	  756824	      1560 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/jsoniter________-12             	  520132	      2024 ns/op	      56 B/op	       7 allocs/op
BenchmarkValid/small_336b____________/gofaster-jx_____-12             	 1424547	       725.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/tidwallgjson____-12             	 1834096	       597.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/valyala-fastjson-12             	 1816468	       621.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small_336b____________/goccy-go-json___-12             	   76108	     19285 ns/op	    2867 B/op	      61 allocs/op
BenchmarkValid/small_336b____________/bytedance-sonic_-12             	 1027084	      1005 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/jscan___________-12             	      46	  23503030 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/encoding-json___-12             	       9	 112513145 ns/op	     171 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/jsoniter________-12             	      13	  94490440 ns/op	13582836 B/op	  644361 allocs/op
BenchmarkValid/large_26m_____________/gofaster-jx_____-12             	      31	  37538530 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/tidwallgjson____-12             	      24	  48973183 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/valyala-fastjson-12             	      24	  43518324 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/large_26m_____________/goccy-go-json___-12             	       1	29392708674 ns/op	144637240 B/op	 2338145 allocs/op
BenchmarkValid/large_26m_____________/bytedance-sonic_-12             	      33	  35061204 ns/op	    1288 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/jscan___________-12             	    6070	    180855 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/encoding-json___-12             	    1830	    569842 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/jsoniter________-12             	    1092	   1261905 ns/op	   69246 B/op	    2121 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/gofaster-jx_____-12             	    4573	    259608 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/tidwallgjson____-12             	    5019	    250844 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/valyala-fastjson-12             	    3480	    347022 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/goccy-go-json___-12             	      56	  20345949 ns/op	  779350 B/op	   20799 allocs/op
BenchmarkValid/nasa_SxSW_2016_125k___/bytedance-sonic_-12             	    4884	    226035 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/jscan___________-12             	  404980	      2900 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/encoding-json___-12             	   91657	     13109 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/jsoniter________-12             	   34366	     38919 ns/op	    2065 B/op	      15 allocs/op
BenchmarkValid/escaped_3k____________/gofaster-jx_____-12             	  112106	     10458 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/tidwallgjson____-12             	  213298	      5349 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/valyala-fastjson-12             	   93360	     12581 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/escaped_3k____________/goccy-go-json___-12             	   15820	     91174 ns/op	    4480 B/op	      13 allocs/op
BenchmarkValid/escaped_3k____________/bytedance-sonic_-12             	 2133204	       482.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/jscan___________-12             	   61797	     20063 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/encoding-json___-12             	   22114	     50884 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/jsoniter________-12             	   29547	     37295 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/gofaster-jx_____-12             	   35664	     29878 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/tidwallgjson____-12             	   44745	     24650 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/valyala-fastjson-12             	   37722	     27070 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024_12k____/goccy-go-json___-12             	    1674	    670207 ns/op	   73503 B/op	    2057 allocs/op
BenchmarkValid/array_int_1024_12k____/bytedance-sonic_-12             	   47402	     25306 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/jscan___________-12             	   56043	     22505 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/encoding-json___-12             	   16840	     66461 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/jsoniter________-12             	    5548	    264600 ns/op	    8756 B/op	     547 allocs/op
BenchmarkValid/array_dec_1024_10k____/gofaster-jx_____-12             	   23862	     46557 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/tidwallgjson____-12             	   42774	     27951 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/valyala-fastjson-12             	   30075	     35855 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024_10k____/goccy-go-json___-12             	    2166	    695147 ns/op	   73494 B/op	    2057 allocs/op
BenchmarkValid/array_dec_1024_10k____/bytedance-sonic_-12             	   34590	     30690 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/jscan___________-12             	  147784	      8039 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/encoding-json___-12             	   31860	     32317 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/jsoniter________-12             	   40851	     24840 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/gofaster-jx_____-12             	   36618	     28742 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/tidwallgjson____-12             	  108295	     11378 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/valyala-fastjson-12             	  106530	     11114 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024_5k/goccy-go-json___-12             	    3466	    324939 ns/op	   48953 B/op	    1036 allocs/op
BenchmarkValid/array_nullbool_1024_5k/bytedance-sonic_-12             	   68221	     16490 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/jscan___________-12             	    4572	    257375 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/encoding-json___-12             	     382	   2682080 ns/op	       4 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/jsoniter________-12             	    1286	    911493 ns/op	       1 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/gofaster-jx_____-12             	    4006	    277346 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/tidwallgjson____-12             	    1392	    896774 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/valyala-fastjson-12             	    2430	    466319 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024_639k___/goccy-go-json___-12             	     184	   6461298 ns/op	 2817340 B/op	    3080 allocs/op
BenchmarkValid/array_str_1024_639k___/bytedance-sonic_-12             	   11598	     89045 ns/op	       3 B/op	       0 allocs/op
PASS
ok  	github.com/romshark/jscan/v2	261.795s
```

</details>

|package|version|
|-|-|
|pkg.go.dev/encoding/json|[go1.20.5](https://pkg.go.dev/encoding/json)|
|github.com/go-faster/jx|[v1.0.0](https://github.com/go-faster/jx/releases/tag/v1.0.0)|
|github.com/go-faster/jx|[v1.0.0](https://github.com/go-faster/jx/releases/tag/v1.0.0)|
|github.com/json-iterator/go|[v1.1.12](https://github.com/json-iterator/go/releases/tag/v1.1.12)|
|github.com/tidwall/gjson|[v1.14.4](https://github.com/tidwall/gjson/releases/tag/v1.14.4)|
|github.com/valyala/fastjson|[v1.6.4](https://github.com/valyala/fastjson/releases/tag/v1.6.4)|
|github.com/goccy/go-json|[v0.10.2](https://github.com/goccy/go-json/releases/tag/v0.10.2)|
|github.com/bytedance/sonic|[v1.8.6](https://github.com/bytedance/sonic/releases/tag/v1.8.6)|
