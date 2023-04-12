package jscan_test

import (
	"bytes"
	_ "embed"
	"encoding/json"
	encodingjson "encoding/json"
	"fmt"
	"testing"

	"github.com/romshark/jscan"

	bytedancesonic "github.com/bytedance/sonic"
	gofasterjx "github.com/go-faster/jx"
	goccygojson "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
	sinhashubham95jsonic "github.com/sinhashubham95/jsonic"
	"github.com/stretchr/testify/require"
	tidwallgjson "github.com/tidwall/gjson"
	valyalafastjson "github.com/valyala/fastjson"
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
}

func CalcStatsJscan(parser *jscan.ParserBytes, str []byte) (s Stats) {
	if err := parser.Scan(
		jscan.Options{},
		str,
		func(i *jscan.IteratorBytes) (err bool) {
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

func CalcStatsJsoniter(i *jsoniter.Iterator, str []byte) (s Stats) {
	i = i.ResetBytes(str)
	var readValue func(lv int, k string)
	readValue = func(
		level int,
		key string,
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
				readValue(l, "")
				index++
				if index > s.MaxArrayLen {
					s.MaxArrayLen = index
				}
			}
		case jsoniter.ObjectValue:
			s.TotalObjects++
			l := level + 1
			for f := i.ReadObject(); f != ""; f = i.ReadObject() {
				readValue(l, f)
			}
		}
	}
	readValue(0, "")
	return
}

func CalcStatsGofasterJx(parser *gofasterjx.Decoder, str []byte) (s Stats) {
	parser.ResetBytes(str)
	var jxParseValue func(lv int, k []byte, ai int) error
	jxParseValue = func(
		level int,
		key []byte,
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

		switch parser.Next() {
		case gofasterjx.String:
			s.TotalStrings++
			if err := parser.Skip(); err != nil {
				return err
			}
		case gofasterjx.Null:
			s.TotalNulls++
			if err := parser.Skip(); err != nil {
				return err
			}
		case gofasterjx.Bool:
			s.TotalBooleans++
			if err := parser.Skip(); err != nil {
				return err
			}
		case gofasterjx.Number:
			s.TotalNumbers++
			if err := parser.Skip(); err != nil {
				return err
			}
		case gofasterjx.Array:
			s.TotalArrays++
			i := 0
			if err := parser.Arr(func(d *gofasterjx.Decoder) error {
				if err := jxParseValue(level+1, nil, i); err != nil {
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
			if err := parser.ObjBytes(func(d *gofasterjx.Decoder, key []byte) error {
				return jxParseValue(level+1, key, -1)
			}); err != nil {
				return err
			}
		}
		return nil
	}
	if err := jxParseValue(0, nil, -1); err != nil {
		panic(err)
	}
	return
}

func CalcStatsValyalaFastjson(parser *valyalafastjson.Parser, str []byte) (s Stats) {
	v, err := parser.ParseBytes(str)
	if err != nil {
		panic(err)
	}

	var parseValue func(v *valyalafastjson.Value, lv int, k []byte, a int) error
	parseValue = func(
		v *valyalafastjson.Value,
		level int,
		key []byte,
		arrayIndex int,
	) error {
		if level > s.MaxDepth {
			s.MaxDepth = level
		}
		if l := len(key); l > 0 && arrayIndex == -1 {
			s.TotalKeys++
			if l > s.MaxKeyLen {
				s.MaxKeyLen = l
			}
		}

		switch v.Type() {
		case valyalafastjson.TypeString:
			s.TotalStrings++

		case valyalafastjson.TypeNull:
			s.TotalNulls++

		case valyalafastjson.TypeTrue, valyalafastjson.TypeFalse:
			s.TotalBooleans++

		case valyalafastjson.TypeNumber:
			s.TotalNumbers++

		case valyalafastjson.TypeArray:
			s.TotalArrays++
			values, err := v.Array()
			if err != nil {
				return err
			}
			lv := level + 1
			for i, v := range values {
				if err := parseValue(v, lv, key, i); err != nil {
					return err
				}
				if i := i + 1; i > s.MaxArrayLen {
					s.MaxArrayLen = i
				}
			}
			return nil

		case valyalafastjson.TypeObject:
			s.TotalObjects++
			o, err := v.Object()
			if err != nil {
				return err
			}
			lv := level + 1
			o.Visit(func(key []byte, v *valyalafastjson.Value) {
				if err = parseValue(v, lv, key, -1); err != nil {
					return
				}
			})
			return nil
		}
		return nil
	}
	if err := parseValue(v, 0, nil, -1); err != nil {
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

	t.Run("jscan", func(t *testing.T) {
		p := jscan.NewBytes(64)
		a := CalcStatsJscan(p, []byte(input))
		require.Equal(t, expect, a)
	})
	t.Run("jsoniter", func(t *testing.T) {
		p := jsoniter.NewIterator(jsoniter.ConfigFastest)
		a := CalcStatsJsoniter(p, []byte(input))
		require.Equal(t, expect, a)
	})
	t.Run("gofaster-jx", func(t *testing.T) {
		p := new(gofasterjx.Decoder)
		a := CalcStatsGofasterJx(p, []byte(input))
		require.Equal(t, expect, a)
	})
	t.Run("valyala-fastjson", func(t *testing.T) {
		p := &valyalafastjson.Parser{}
		a := CalcStatsValyalaFastjson(p, []byte(input))
		require.Equal(t, expect, a)
	})
}

var gs Stats

//go:embed miniscule.json
var jsonMiniscule []byte

//go:embed tiny.json
var jsonTiny []byte

//go:embed small.json
var jsonSmall []byte

//go:embed large.json
var jsonLarge []byte

//go:embed escaped.json
var jsonEscaped []byte

//go:embed array_int_1024.json
var jsonArrayInt1024 []byte

//go:embed array_dec_1024.json
var jsonArrayDec1024 []byte

//go:embed array_nullbool_1024.json
var jsonArrayNullBool1024 []byte

//go:embed array_str_1024.json
var jsonArrayStr1024 []byte

func BenchmarkCalcStats(b *testing.B) {
	for _, bb := range []struct {
		name string
		json []byte
	}{
		{"miniscule", jsonMiniscule},
		{"tiny", jsonTiny},
		{"small", jsonSmall},
		{"large", jsonLarge},
		{"escaped", jsonEscaped},
		{"array_int_1024", jsonArrayInt1024},
		{"array_dec_1024", jsonArrayDec1024},
		{"array_nullbool_1024", jsonArrayNullBool1024},
		{"array_str_1024", jsonArrayStr1024},
	} {
		b.Run(bb.name, func(b *testing.B) {
			b.Run("jscan", func(b *testing.B) {
				p := jscan.NewBytes(64)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gs = CalcStatsJscan(p, bb.json)
				}
			})
			b.Run("jsoniter", func(b *testing.B) {
				p := jsoniter.NewIterator(jsoniter.ConfigFastest)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gs = CalcStatsJsoniter(p, bb.json)
				}
			})
			b.Run("gofaster-jx", func(b *testing.B) {
				p := new(gofasterjx.Decoder)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gs = CalcStatsGofasterJx(p, bb.json)
				}
			})
			b.Run("valyala-fastjson", func(b *testing.B) {
				p := new(valyalafastjson.Parser)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gs = CalcStatsValyalaFastjson(p, bb.json)
				}
			})
		})
	}
}

func TestBenchGet(t *testing.T) {
	j := `[false,[[2, {"[foo]":[{"bar-baz":true}]}]]]`

	t.Run("jscan", func(t *testing.T) {
		path := `[1][0][1].\[foo\][0].bar-baz`
		err := jscan.GetBytes(
			[]byte(j), []byte(path), true, func(i *jscan.IteratorBytes) {
				require.Equal(t, jscan.ValueTypeTrue, i.ValueType)
				require.Equal(t, "true", string(i.Value()))
			},
		)
		require.False(t, err.IsErr())
	})

	t.Run("jsoniter", func(t *testing.T) {
		r := jsoniter.Get([]byte(j), 1, 0, 1, "[foo]", 0, "bar-baz")
		r.MustBeValid()
		require.NoError(t, r.LastError())
		require.True(t, r.ToBool())
	})

	t.Run("tidwallgjson", func(t *testing.T) {
		r := tidwallgjson.GetBytes([]byte(j), `1.0.1.\[foo\].0.bar-baz`)
		require.True(t, r.Exists())
		require.True(t, r.Bool())
	})

	t.Run("valyalafastjson", func(t *testing.T) {
		path := []string{"1", "0", "1", "[foo]", "0", "bar-baz"}
		require.True(t, valyalafastjson.Exists([]byte(j), path...))
		v := valyalafastjson.GetBool([]byte(j), path...)
		require.True(t, v)
	})

	t.Run("sinhashubham95jsonic", func(t *testing.T) {
		p, err := sinhashubham95jsonic.New([]byte(j))
		require.NoError(t, err)
		v, err := p.GetBool(`[1].[0].[1].[foo].[0].bar-baz`)
		require.NoError(t, err)
		require.True(t, v)
	})
}

func BenchmarkGet(b *testing.B) {
	json := `[false,[[2, {"[foo]":[{"bar-baz":true}]}]]]`

	b.Run("jscan", func(b *testing.B) {
		bytesTrue := []byte("true")
		path := []byte(`[1][0][1].\[foo\][0].bar-baz`)
		json := []byte(json)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = jscan.GetBytes(
				json, path, true, func(i *jscan.IteratorBytes) {
					gbool = bytes.Equal(i.Value(), bytesTrue)
				},
			)
		}
	})

	b.Run("jsoniter", func(b *testing.B) {
		json := []byte(json)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			gbool = jsoniter.Get(
				json, 1, 0, 1, "[foo]", 0, "bar-baz",
			).ToBool()
		}
	})

	b.Run("tidwallgjson", func(b *testing.B) {
		json := []byte(json)
		path := `1.0.1.\[foo\].0.bar-baz`
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			gbool = tidwallgjson.GetBytes(json, path).Bool()
		}
	})

	b.Run("valyalafastjson", func(b *testing.B) {
		jb := []byte(json)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			gbool = valyalafastjson.GetBool(
				jb, "1", "0", "1", `\[foo\]`, "0", `bar-baz`,
			)
		}
	})

	b.Run("sinhashubham95jsonic", func(b *testing.B) {
		p, err := sinhashubham95jsonic.New([]byte(json))
		if err != nil {
			panic(err)
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			gbool, _ = p.GetBool(`[1].[0].[1].[foo].[0].bar-baz`)
		}
	})
}

func TestValid(t *testing.T) {
	j := `[false,[[2, {"[foo]":[{"bar-baz":"fuz"}]}]]]`

	t.Run("jscan", func(t *testing.T) {
		require.True(t, jscan.Valid(j))
	})

	t.Run("encoding-json", func(t *testing.T) {
		require.True(t, encodingjson.Valid([]byte(j)))
	})

	t.Run("jsoniter", func(t *testing.T) {
		require.True(t, jsoniter.Valid([]byte(j)))
	})

	t.Run("tidwallgjson", func(t *testing.T) {
		require.True(t, tidwallgjson.Valid(j))
	})

	t.Run("fast-json", func(t *testing.T) {
		require.NoError(t, valyalafastjson.Validate(j))
	})

	t.Run("goccy-go-json", func(t *testing.T) {
		require.True(t, goccygojson.Valid([]byte(j)))
	})

	t.Run("bytedance-sonic", func(t *testing.T) {
		require.True(t, bytedancesonic.ConfigFastest.Valid([]byte(j)))
	})
}

var gbool bool

func BenchmarkValid(b *testing.B) {
	for _, bb := range []struct {
		name  string
		input []byte
	}{
		{"miniscule", jsonMiniscule},
		{"tiny", jsonTiny},
		{"small", jsonSmall},
		{"large", jsonLarge},
		{"escaped", jsonEscaped},
		{"unwind_stack", MakeRepeated("[", 1024)},
		{"array_int_1024", jsonArrayInt1024},
		{"array_dec_1024", jsonArrayDec1024},
		{"array_nullbool_1024", jsonArrayNullBool1024},
		{"array_str_1024", jsonArrayStr1024},
	} {
		b.Run(bb.name, func(b *testing.B) {
			b.Run("jscan", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					gbool = jscan.ValidBytes(bb.input)
				}
			})

			b.Run("encoding-json", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					gbool = json.Valid(bb.input)
				}
			})

			b.Run("jsoniter", func(b *testing.B) {
				jb := []byte(bb.input)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gbool = jsoniter.Valid(jb)
				}
			})

			b.Run("gofaster-jx", func(b *testing.B) {
				jb := []byte(bb.input)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gbool = gofasterjx.Valid(jb)
				}
			})

			b.Run("tidwallgjson", func(b *testing.B) {
				j := string(bb.input)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gbool = tidwallgjson.Valid(j)
				}
			})

			b.Run("valyala-fastjson", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					gbool = (valyalafastjson.ValidateBytes(bb.input) != nil)
				}
			})

			b.Run("goccy-go-json", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					gbool = goccygojson.Valid(bb.input)
				}
			})

			b.Run("bytedance-sonic", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					gbool = bytedancesonic.ConfigFastest.Valid(bb.input)
				}
			})
		})
	}
}

func MakeRepeated(s string, n int) []byte {
	var b bytes.Buffer
	b.Grow(len(s) * n)
	for i := 0; i < n; i++ {
		b.WriteString(s)
	}
	return b.Bytes()
}
