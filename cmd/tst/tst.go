package main

import (
	"fmt"

	"github.com/romshark/jscan"
)

func main() {
	in := `{
		"s":"value",
		"t":true
	}`
	totalKeys := 0
	err := jscan.Scan(in, func(i *jscan.Iterator[string]) (err bool) {
		if i.KeyIndex() != -1 {
			totalKeys++
			fmt.Printf(
				"key %d %q (len: %d)\n",
				totalKeys,
				i.Key()[1:len(i.Key())-1],
				i.KeyIndexEnd()-i.KeyIndex(),
			)
		}
		return false
	})
	if err.IsErr() {
		panic(err.Error())
	}
	fmt.Println("")
	fmt.Println("KEYS: ", totalKeys)
}

// func main() {
// 	// in := `123"foo""bar"`
// 	in := `-"okay"`

// 	for x := in; x != ""; {
// 		trailing, err := jscan.ValidateOne(x)
// 		if err.IsErr() {
// 			panic(err.Error())
// 		}
// 		fmt.Printf("trailing: %q\n\n", trailing)
// 		x = trailing
// 	}
// }

// func main() {
// 	// in := `123"foo""bar"`
// 	in := `0e0"bar"0.12`

// 	for x := in; x != ""; {
// 		trailing, err := jscan.ScanOne(
// 			jscan.Options{},
// 			x,
// 			func(
// 				i *jscan.Iterator,
// 			) (err bool) {
// 				fmt.Println(
// 					i.ValueType.String(),
// 					":", i.Value(),
// 				)
// 				return false
// 			},
// 		)
// 		if err.IsErr() {
// 			panic(err.Error())
// 		}
// 		fmt.Printf("trailing: %q\n\n", trailing)
// 		x = trailing
// 	}
// }
