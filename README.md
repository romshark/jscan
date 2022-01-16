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
BenchmarkCalcStats/jscan/tiny-10    	17281384	        61.52 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/jsoniter/tiny-10 	12591319	        93.49 ns/op	     160 B/op	       2 allocs/op
BenchmarkCalcStats/gofaster-jx/tiny-10         	15380218	        76.69 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/valyala-fastjson/tiny-10             	14686556	        81.48 ns/op	       0 B/op	       0 allocs/op

BenchmarkCalcStats/jscan_withpath/tiny-10      	14819349	        80.68 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/jsoniter_withpath/tiny-10   	11339800	       105.5 ns/op	     160 B/op	       2 allocs/op
BenchmarkCalcStats/gofaster-jx_withpath/tiny-10         	13993488	        84.37 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/valyala-fastjson_withpath/tiny-10    	11028562	       107.5 ns/op	       8 B/op	       1 allocs/op
```

Small JSON document (335 bytes):

```
BenchmarkCalcStats/jscan/small-10   	 2055507	       572.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/jsoniter/small-10         	 1465920	       818.2 ns/op	     224 B/op	      12 allocs/op
BenchmarkCalcStats/gofaster-jx/small-10      	 1301629	       921.0 ns/op	      16 B/op	       2 allocs/op
BenchmarkCalcStats/valyala-fastjson/small-10          	 1667216	       719.2 ns/op	       0 B/op	       0 allocs/op

BenchmarkCalcStats/jscan_withpath/small-10   	 1593405	       747.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/jsoniter_withpath/small-10         	  974680	      1211 ns/op	     288 B/op	      21 allocs/op
BenchmarkCalcStats/gofaster-jx_withpath/small-10      	  909572	      1296 ns/op	      80 B/op	      13 allocs/op
BenchmarkCalcStats/valyala-fastjson_withpath/small-10 	 1222809	       981.1 ns/op	      88 B/op	      10 allocs/op
```

Large JSON document (26.1 MB):

```
BenchmarkCalcStats
BenchmarkCalcStats/jscan/large-10   	      30	  39817379 ns/op	     271 B/op	       2 allocs/op
BenchmarkCalcStats/jsoniter/large-10         	      19	  59779695 ns/op	33060224 B/op	 1108611 allocs/op
BenchmarkCalcStats/gofaster-jx/large-10      	      22	  50123860 ns/op	      58 B/op	       0 allocs/op
BenchmarkCalcStats/valyala-fastjson/large-10          	      30	  40858881 ns/op	11381580 B/op	   11033 allocs/op

BenchmarkCalcStats/jscan_withpath/large-10   	      26	  44913019 ns/op	     322 B/op	       2 allocs/op
BenchmarkCalcStats/jsoniter_withpath/large-10         	      13	  85893750 ns/op	55805724 B/op	 1757457 allocs/op
BenchmarkCalcStats/gofaster-jx_withpath/large-10      	      12	  95080243 ns/op	52296336 B/op	 1544965 allocs/op
BenchmarkCalcStats/valyala-fastjson_withpath/large-10 	      21	  52270306 ns/op	29394618 B/op	  340766 allocs/op
```

Array of 1024 integers:

```
BenchmarkCalcStats/jscan/array_int_1024-10         	   49778	     23505 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/jsoniter/array_int_1024-10      	   29674	     40478 ns/op	   16528 B/op	    1025 allocs/op
BenchmarkCalcStats/gofaster-jx/array_int_1024-10   	   36080	     33218 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/valyala-fastjson/array_int_1024-10       	   46993	     25360 ns/op	       5 B/op	       0 allocs/op

BenchmarkCalcStats/jscan_withpath/array_int_1024-10         	   33222	     36090 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/jsoniter_withpath/array_int_1024-10      	   13416	     89293 ns/op	   24496 B/op	    2973 allocs/op
BenchmarkCalcStats/gofaster-jx_withpath/array_int_1024-10   	   14683	     81523 ns/op	    7970 B/op	    1948 allocs/op
BenchmarkCalcStats/valyala-fastjson_withpath/array_int_1024-10         	   24570	     48887 ns/op	    8194 B/op	    1024 allocs/op
```

Array of 1024 floating point numbers:

```
BenchmarkCalcStats/jscan/array_dec_1024-10         	   43645	     25321 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/jsoniter/array_dec_1024-10      	   27168	     44132 ns/op	   16528 B/op	    1025 allocs/op
BenchmarkCalcStats/gofaster-jx/array_dec_1024-10   	   13220	     90298 ns/op	    6498 B/op	     547 allocs/op
BenchmarkCalcStats/valyala-fastjson/array_dec_1024-10       	   42273	     27560 ns/op	       0 B/op	       0 allocs/op

BenchmarkCalcStats/jscan_withpath/array_dec_1024-10         	   28917	     40965 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/jsoniter_withpath/array_dec_1024-10      	   12783	     94376 ns/op	   24496 B/op	    2973 allocs/op
BenchmarkCalcStats/gofaster-jx_withpath/array_dec_1024-10   	   13364	     90193 ns/op	    7970 B/op	    1948 allocs/op
BenchmarkCalcStats/valyala-fastjson_withpath/array_dec_1024-10         	   23299	     50773 ns/op	    8204 B/op	    1024 allocs/op
```

Array of 1024 strings:

```
BenchmarkCalcStats/jscan/array_str_1024-10         	   36171	     32719 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/jsoniter/array_str_1024-10      	    1909	    611724 ns/op	  670313 B/op	    1019 allocs/op
BenchmarkCalcStats/gofaster-jx/array_str_1024-10   	    2246	    520334 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/valyala-fastjson/array_str_1024-10       	   15608	     77085 ns/op	      55 B/op	       0 allocs/op

BenchmarkCalcStats/jscan_withpath/array_str_1024-10         	   26394	     45593 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/jsoniter_withpath/array_str_1024-10      	    1767	    667022 ns/op	  678284 B/op	    2967 allocs/op
BenchmarkCalcStats/gofaster-jx_withpath/array_str_1024-10   	    1773	    684708 ns/op	  677621 B/op	    2927 allocs/op
BenchmarkCalcStats/valyala-fastjson_withpath/array_str_1024-10         	   12117	     98868 ns/op	    8195 B/op	    1024 allocs/op
```

Array of 1024 nullable booleans:

```
BenchmarkCalcStats/jscan/array_nullbool_1024-10         	   78883	     14034 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/jsoniter/array_nullbool_1024-10      	   44224	     26998 ns/op	     144 B/op	       1 allocs/op
BenchmarkCalcStats/gofaster-jx/array_nullbool_1024-10   	   35641	     33411 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/valyala-fastjson/array_nullbool_1024-10       	   73539	     16104 ns/op	       0 B/op	       0 allocs/op

BenchmarkCalcStats/jscan_withpath/array_nullbool_1024-10         	   43722	     27344 ns/op	       0 B/op	       0 allocs/op
BenchmarkCalcStats/jsoniter_withpath/array_nullbool_1024-10      	   15586	     76804 ns/op	    8112 B/op	    1949 allocs/op
BenchmarkCalcStats/gofaster-jx_withpath/array_nullbool_1024-10   	   14474	     83515 ns/op	    7970 B/op	    1948 allocs/op
BenchmarkCalcStats/valyala-fastjson_withpath/array_nullbool_1024-10         	   30286	     39500 ns/op	    8194 B/op	    1024 allocs/op
```

Get by path:

```
BenchmarkGet/jscan-10         	 3957463	       288.8 ns/op	      16 B/op	       2 allocs/op
BenchmarkGet/jsoniter-10      	 1272576	       941.7 ns/op	     496 B/op	      19 allocs/op
BenchmarkGet/tidwallgjson-10  	 6308622	       190.5 ns/op	      16 B/op	       2 allocs/op
BenchmarkGet/valyalafastjson-10         	 6188893	       193.9 ns/op	       0 B/op	       0 allocs/op
```

Validation:

```
BenchmarkValid/tiny/jscan-10   	19409760	        61.44 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny/jsoniter-10         	20232646	        59.52 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny/gofaster-jx-10      	15790867	        75.55 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny/encoding-json-10    	27169221	        43.90 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny/tidwallgjson-10     	72519966	        16.46 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/tiny/valyala-fastjson-10 	41914624	        28.70 ns/op	       0 B/op	       0 allocs/op

BenchmarkValid/small/jscan-10           	 2263227	       530.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small/jsoniter-10        	 1482616	       809.6 ns/op	      56 B/op	       7 allocs/op
BenchmarkValid/small/gofaster-jx-10     	 1414510	       846.3 ns/op	      16 B/op	       2 allocs/op
BenchmarkValid/small/encoding-json-10   	 1206177	       993.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small/tidwallgjson-10    	 3561247	       337.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/small/valyala-fastjson-10         	 2633577	       452.2 ns/op	       0 B/op	       0 allocs/op

BenchmarkValid/large/jscan-10                    	      30	  39348897 ns/op	     271 B/op	       2 allocs/op
BenchmarkValid/large/jsoniter-10                 	      24	  47381846 ns/op	13792181 B/op	  644454 allocs/op
BenchmarkValid/large/gofaster-jx-10              	      26	  44456761 ns/op	      52 B/op	       0 allocs/op
BenchmarkValid/large/encoding-json-10            	      16	  70171828 ns/op	      80 B/op	       0 allocs/op
BenchmarkValid/large/tidwallgjson-10             	      40	  28859103 ns/op	       2 B/op	       0 allocs/op
BenchmarkValid/large/valyala-fastjson-10         	      40	  28703367 ns/op	       0 B/op	       0 allocs/op

BenchmarkValid/unwind_stack/jscan-10             	  457160	      2624 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/unwind_stack/jsoniter-10          	   15862	     75624 ns/op	   33145 B/op	    1033 allocs/op
BenchmarkValid/unwind_stack/gofaster-jx-10       	    1560	    763111 ns/op	  131117 B/op	    2048 allocs/op
BenchmarkValid/unwind_stack/encoding-json-10     	  212856	      5536 ns/op	      24 B/op	       1 allocs/op
BenchmarkValid/unwind_stack/tidwallgjson-10      	   80488	     14748 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/unwind_stack/valyala-fastjson-10  	     222	   5544191 ns/op	52472494 B/op	    4133 allocs/op

BenchmarkValid/array_int_1024/jscan-10           	   66847	     17774 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024/jsoniter-10        	   51276	     23477 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024/gofaster-jx-10     	   44049	     26980 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024/encoding-json-10   	   34149	     35169 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024/tidwallgjson-10    	   77713	     15338 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_int_1024/valyala-fastjson-10         	   65236	     18438 ns/op	       0 B/op	       0 allocs/op

BenchmarkValid/array_dec_1024/jscan-10                    	   70789	     15665 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024/jsoniter-10                 	   15679	     76847 ns/op	    8754 B/op	     547 allocs/op
BenchmarkValid/array_dec_1024/gofaster-jx-10              	   14700	     81322 ns/op	    6498 B/op	     547 allocs/op
BenchmarkValid/array_dec_1024/encoding-json-10            	   32197	     36578 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024/tidwallgjson-10             	   84492	     12129 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_dec_1024/valyala-fastjson-10         	   51250	     19518 ns/op	       0 B/op	       0 allocs/op

BenchmarkValid/array_nullbool_1024/jscan-10               	  226735	      5242 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024/jsoniter-10            	   59954	     20019 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024/gofaster-jx-10         	   44670	     26863 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024/encoding-json-10       	   56584	     21165 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024/tidwallgjson-10        	  212652	      5307 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_nullbool_1024/valyala-fastjson-10    	  130365	      9181 ns/op	       0 B/op	       0 allocs/op

BenchmarkValid/array_str_1024/jscan-10                    	   35758	     33207 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024/jsoniter-10                 	    2346	    509377 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024/gofaster-jx-10              	    2325	    513011 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024/encoding-json-10            	     858	   1392125 ns/op	       1 B/op	       0 allocs/op
BenchmarkValid/array_str_1024/tidwallgjson-10             	    2370	    503992 ns/op	       0 B/op	       0 allocs/op
BenchmarkValid/array_str_1024/valyala-fastjson-10         	    4519	    260450 ns/op	       0 B/op	       0 allocs/op
```
