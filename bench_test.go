package jscan_test

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/romshark/jscan"

	gofasterjx "github.com/go-faster/jx"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

type Stats struct {
	TotalStrings  int
	TotalNulls    int
	TotalBooleans int
	TotalNumbers  int
	TotalObjects  int
	TotalArrays   int
	TotalKeys     int
	MaxKeyLen     int
	MaxDepth      int
	MaxArrayLen   int
	MaxPathLen    int
}

func CalcStatsJscan(str string) (s Stats) {
	if err := jscan.Scan(
		jscan.Options{},
		str,
		func(i *jscan.Iterator) (err bool) {
			if l := i.KeyEnd - i.KeyStart; l > 0 {
				s.TotalKeys++
				if l > s.MaxKeyLen {
					s.MaxKeyLen = l
				}
			}
			switch i.ValueType {
			case jscan.ValueTypeObject:
				s.TotalObjects++
			case jscan.ValueTypeArray:
				s.TotalArrays++
			case jscan.ValueTypeNull:
				s.TotalNulls++
			case jscan.ValueTypeFalse, jscan.ValueTypeTrue:
				s.TotalBooleans++
			case jscan.ValueTypeNumber:
				s.TotalNumbers++
			case jscan.ValueTypeString:
				s.TotalStrings++
			}
			if i.Level > s.MaxDepth {
				s.MaxDepth = i.Level
			}
			if l := i.ArrayIndex + 1; l > s.MaxArrayLen {
				s.MaxArrayLen = l
			}
			return false
		},
	); err.IsErr() {
		panic(fmt.Errorf("unexpected error: %s", err))
	}
	return
}

func CalcStatsJscanWithPath(str string) (s Stats) {
	if err := jscan.Scan(jscan.Options{
		CachePath:  true,
		EscapePath: false,
	}, str, func(i *jscan.Iterator) (err bool) {
		if l := i.KeyEnd - i.KeyStart; l > 0 {
			s.TotalKeys++
			if l > s.MaxKeyLen {
				s.MaxKeyLen = l
			}
		}
		switch i.ValueType {
		case jscan.ValueTypeObject:
			s.TotalObjects++
		case jscan.ValueTypeArray:
			s.TotalArrays++
		case jscan.ValueTypeNull:
			s.TotalNulls++
		case jscan.ValueTypeFalse, jscan.ValueTypeTrue:
			s.TotalBooleans++
		case jscan.ValueTypeNumber:
			s.TotalNumbers++
		case jscan.ValueTypeString:
			s.TotalStrings++
		}
		if i.Level > s.MaxDepth {
			s.MaxDepth = i.Level
		}
		if l := i.ArrayIndex + 1; l > s.MaxArrayLen {
			s.MaxArrayLen = l
		}
		i.ViewPath(func(p []byte) {
			if l := len(p); l > s.MaxPathLen {
				s.MaxPathLen = l
			}
		})
		return false
	}); err.IsErr() {
		panic(fmt.Errorf("unexpected error: %s", err))
	}
	return
}

func CalcStatsJsoniter(str string) (s Stats) {
	i := jsoniter.ParseString(jsoniter.ConfigDefault, str)
	var readValue func(lv int, k string, ai int, i *jsoniter.Iterator)
	readValue = func(
		level int,
		key string,
		arrIndex int,
		i *jsoniter.Iterator,
	) {
		if level > s.MaxDepth {
			s.MaxDepth = level
		}
		if l := len(key); l > 0 {
			s.TotalKeys++
			if l > s.MaxKeyLen {
				s.MaxKeyLen = l
			}
		}
		switch i.WhatIsNext() {
		case jsoniter.StringValue:
			i.ReadString()
			s.TotalStrings++
		case jsoniter.NumberValue:
			i.ReadNumber()
			s.TotalNumbers++
		case jsoniter.NilValue:
			i.ReadNil()
			s.TotalNulls++
		case jsoniter.BoolValue:
			i.ReadBool()
			s.TotalBooleans++
		case jsoniter.ArrayValue:
			s.TotalArrays++
			l := level + 1
			index := 0
			for e := i.ReadArray(); e; e = i.ReadArray() {
				readValue(l, "", index, i)
				index++
				if index > s.MaxArrayLen {
					s.MaxArrayLen = index
				}
			}
		case jsoniter.ObjectValue:
			s.TotalObjects++
			l := level + 1
			for f := i.ReadObject(); f != ""; f = i.ReadObject() {
				readValue(l, f, -1, i)
			}
		}
	}
	readValue(0, "", -1, i)
	return
}

func CalcStatsJsoniterWithPath(str string) (s Stats) {
	i := jsoniter.ParseString(jsoniter.ConfigDefault, str)
	var readValue func(lv int, k, p string, ai int, i *jsoniter.Iterator)
	readValue = func(
		level int,
		key, path string,
		arrIndex int,
		i *jsoniter.Iterator,
	) {
		if level > s.MaxDepth {
			s.MaxDepth = level
		}
		if l := len(key); l > 0 {
			s.TotalKeys++
			if l > s.MaxKeyLen {
				s.MaxKeyLen = l
			}
		}
		if key != "" {
			if path != "" {
				path += "." + key
			} else {
				path += key
			}
		} else if arrIndex > -1 {
			path += "[" + strconv.Itoa(arrIndex) + "]"
		}
		if l := len(path); l > s.MaxPathLen {
			s.MaxPathLen = l
		}
		switch i.WhatIsNext() {
		case jsoniter.StringValue:
			i.ReadString()
			s.TotalStrings++
		case jsoniter.NumberValue:
			i.ReadNumber()
			s.TotalNumbers++
		case jsoniter.NilValue:
			i.ReadNil()
			s.TotalNulls++
		case jsoniter.BoolValue:
			i.ReadBool()
			s.TotalBooleans++
		case jsoniter.ArrayValue:
			s.TotalArrays++
			l := level + 1
			index := 0
			for e := i.ReadArray(); e; e = i.ReadArray() {
				readValue(l, "", path, index, i)
				index++
				if index > s.MaxArrayLen {
					s.MaxArrayLen = index
				}
			}
		case jsoniter.ObjectValue:
			s.TotalObjects++
			l := level + 1
			for f := i.ReadObject(); f != ""; f = i.ReadObject() {
				readValue(l, f, path, -1, i)
			}
		}
	}
	readValue(0, "", "", -1, i)
	return
}

func CalcStatsGofasterJx(str string) (s Stats) {
	d := gofasterjx.GetDecoder()
	defer gofasterjx.PutDecoder(d)
	d.Reset(strings.NewReader(str))

	var jxParseValue func(lv int, k string, ai int) error
	jxParseValue = func(
		level int,
		key string,
		arrayIndex int,
	) error {
		if level > s.MaxDepth {
			s.MaxDepth = level
		}
		if l := len(key); l > 0 {
			s.TotalKeys++
			if l > s.MaxKeyLen {
				s.MaxKeyLen = l
			}
		}

		switch d.Next() {
		case gofasterjx.String:
			s.TotalStrings++
			_, err := d.Str()
			if err != nil {
				return err
			}
		case gofasterjx.Null:
			s.TotalNulls++
			if err := d.Null(); err != nil {
				return err
			}
		case gofasterjx.Bool:
			s.TotalBooleans++
			_, err := d.Bool()
			if err != nil {
				return err
			}
		case gofasterjx.Number:
			s.TotalNumbers++
			_, err := d.Num()
			if err != nil {
				return err
			}
		case gofasterjx.Array:
			s.TotalArrays++
			i := 0
			if err := d.Arr(func(d *gofasterjx.Decoder) error {
				if err := jxParseValue(level+1, "", i); err != nil {
					return err
				}
				i++
				if i > s.MaxArrayLen {
					s.MaxArrayLen = i
				}
				return nil
			}); err != nil {
				return err
			}
		case gofasterjx.Object:
			s.TotalObjects++
			if err := d.Obj(func(d *gofasterjx.Decoder, key string) error {
				return jxParseValue(level+1, key, -1)
			}); err != nil {
				return err
			}
		}
		return nil
	}
	if err := jxParseValue(0, "", -1); err != nil {
		panic(err)
	}
	return
}

func CalcStatsGofasterJxWithPath(str string) (s Stats) {
	d := gofasterjx.GetDecoder()
	defer gofasterjx.PutDecoder(d)
	d.Reset(strings.NewReader(str))

	var jxParseValue func(lv int, k, path string, ai int) error
	jxParseValue = func(
		level int,
		key, path string,
		arrayIndex int,
	) error {
		if level > s.MaxDepth {
			s.MaxDepth = level
		}
		if l := len(key); l > 0 {
			s.TotalKeys++
			if l > s.MaxKeyLen {
				s.MaxKeyLen = l
			}
		}
		if key != "" {
			if path != "" {
				path += "." + key
			} else {
				path += key
			}
		} else if arrayIndex > -1 {
			path += "[" + strconv.Itoa(arrayIndex) + "]"
		}
		if l := len(path); l > s.MaxPathLen {
			s.MaxPathLen = l
		}

		switch d.Next() {
		case gofasterjx.String:
			s.TotalStrings++
			_, err := d.Str()
			if err != nil {
				return err
			}
		case gofasterjx.Null:
			s.TotalNulls++
			if err := d.Null(); err != nil {
				return err
			}
		case gofasterjx.Bool:
			s.TotalBooleans++
			_, err := d.Bool()
			if err != nil {
				return err
			}
		case gofasterjx.Number:
			s.TotalNumbers++
			_, err := d.Num()
			if err != nil {
				return err
			}
		case gofasterjx.Array:
			s.TotalArrays++
			i := 0
			if err := d.Arr(func(d *gofasterjx.Decoder) error {
				if err := jxParseValue(level+1, "", path, i); err != nil {
					return err
				}
				i++
				if i > s.MaxArrayLen {
					s.MaxArrayLen = i
				}
				return nil
			}); err != nil {
				return err
			}
		case gofasterjx.Object:
			s.TotalObjects++
			if err := d.Obj(func(d *gofasterjx.Decoder, key string) error {
				return jxParseValue(level+1, key, path, -1)
			}); err != nil {
				return err
			}
		}
		return nil
	}
	if err := jxParseValue(0, "", "", -1); err != nil {
		panic(err)
	}
	return
}

func TestImplementations(t *testing.T) {
	const input = `{"s":"value","t":true,"f":false,"0":null,"n":-9.123e3,` +
		`"o0":{},"a0":[],"o":{"k":"\"v\"",` +
		`"a":[true,null,"item",-67.02e9,["foo"]]},"[abc]":[0]}`
	expect := Stats{
		TotalStrings:  4,
		TotalNulls:    2,
		TotalBooleans: 3,
		TotalNumbers:  3,
		TotalObjects:  3,
		TotalArrays:   4,
		MaxDepth:      4,
		TotalKeys:     11,
		MaxKeyLen:     5,
		MaxArrayLen:   5,
	}
	for _, tt := range []struct {
		name string
		fn   func(string) Stats
	}{
		{name: "jscan", fn: CalcStatsJscan},
		{name: "jsoniter", fn: CalcStatsJsoniter},
		{name: "gofaster-jx", fn: CalcStatsGofasterJx},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, expect, tt.fn(input))
		})
	}
}

func TestImplementationsWithPath(t *testing.T) {
	const input = `{"s":"value","t":true,"f":false,"0":null,"n":-9.123e3,` +
		`"o0":{},"a0":[],"o":{"k":"\"v\"",` +
		`"a":[true,null,"item",-67.02e9,["foo"]]},"[abc]":[0]}`
	expect := Stats{
		TotalStrings:  4,
		TotalNulls:    2,
		TotalBooleans: 3,
		TotalNumbers:  3,
		TotalObjects:  3,
		TotalArrays:   4,
		MaxDepth:      4,
		TotalKeys:     11,
		MaxKeyLen:     5,
		MaxArrayLen:   5,
		MaxPathLen:    9,
	}
	for _, tt := range []struct {
		name string
		fn   func(string) Stats
	}{
		{name: "jscan", fn: CalcStatsJscanWithPath},
		{name: "jsoniter", fn: CalcStatsJsoniterWithPath},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, expect, tt.fn(input))
		})
	}
}

var gs Stats

//go:embed tiny.json
var jsonTiny string

//go:embed small.json
var jsonSmall string

//go:embed large.json
var jsonLarge string

func BenchmarkCalcStats(b *testing.B) {
	for _, bb := range []struct {
		name string
		fn   func(string) Stats
	}{
		{name: "jscan", fn: CalcStatsJscan},
		{name: "jsoniter", fn: CalcStatsJsoniter},
		{name: "gofaster-jx", fn: CalcStatsGofasterJx},
		{name: "jscan_withpath", fn: CalcStatsJscanWithPath},
		{name: "jsoniter_withpath", fn: CalcStatsJsoniterWithPath},
		{name: "gofaster-jx_withpath", fn: CalcStatsGofasterJxWithPath},
	} {
		b.Run(bb.name, func(b *testing.B) {
			for _, b2 := range []struct {
				name string
				json string
			}{
				{"tiny", jsonTiny},
				{"small", jsonSmall},
				{"large", jsonLarge},
			} {
				b.Run(b2.name, func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						gs = bb.fn(b2.json)
					}
				})
			}
		})
	}
}

func TestBenchGet(t *testing.T) {
	j := `[false,[[2, {"[foo]":[{"bar-baz":"fuz"}]}]]]`

	t.Run("jsoniter", func(t *testing.T) {
		r := jsoniter.Get([]byte(j), 1, 0, 1, "[foo]", 0, "bar-baz")
		r.MustBeValid()
		require.NoError(t, r.LastError())
		require.Equal(t, "fuz", r.ToString())
	})

	t.Run("jscan", func(t *testing.T) {
		path := `[1][0][1].\[foo\][0].bar-baz`
		err := jscan.Get(j, path, true, func(i *jscan.Iterator) {
			require.Equal(t, jscan.ValueTypeString, i.ValueType)
			require.Equal(t, "fuz", i.Value())
		})
		require.False(t, err.IsErr())
	})
}

var gstr string

func BenchmarkGet(b *testing.B) {
	j := `[false,[[2, {"[foo]":[{"bar-baz":"fuz"}]}]]]`

	b.Run("jsoniter", func(b *testing.B) {
		jb := []byte(j)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			gstr = jsoniter.Get(jb, 1, 0, 1, "[foo]", 0, "bar-baz").ToString()
		}
	})

	b.Run("jscan", func(b *testing.B) {
		path := `[1][0][1].\[foo\][0].bar-baz`
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = jscan.Get(j, path, true, func(i *jscan.Iterator) {
				gstr = i.Value()
			})
		}
	})
}
