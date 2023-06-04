package jscan_test

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/romshark/jscan"

	encodingjson "encoding/json"

	bytedancesonic "github.com/bytedance/sonic"
	gofasterjx "github.com/go-faster/jx"
	goccygojson "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
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

func MustCalcStatsJscan(p *jscan.Parser[[]byte], str []byte) (s Stats) {
	if err := p.Scan(
		str,
		func(i *jscan.Iterator[[]byte]) (err bool) {
			if i.KeyIndex() != -1 {
				// Calculate key length excluding the quotes
				l := i.KeyIndexEnd() - i.KeyIndex() - 2
				s.TotalKeys++
				if l > s.MaxKeyLen {
					s.MaxKeyLen = l
				}
			}
			switch i.ValueType() {
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
			if i.Level() > s.MaxDepth {
				s.MaxDepth = i.Level()
			}
			if l := i.ArrayIndex() + 1; l > s.MaxArrayLen {
				s.MaxArrayLen = l
			}
			return false
		},
	); err.IsErr() {
		panic(fmt.Errorf("unexpected error: %s", err))
	}
	return
}

func MustCalcStatsJsoniter(p *jsoniter.Iterator, str []byte) (s Stats) {
	p.ResetBytes(str)
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
	readValue(0, "", -1, p)
	return
}

func MustCalcStatsGofasterJx(p *gofasterjx.Decoder, str []byte) (s Stats) {
	p.ResetBytes(str)
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

		switch p.Next() {
		case gofasterjx.String:
			s.TotalStrings++
			if err := p.Skip(); err != nil {
				return err
			}
		case gofasterjx.Null:
			s.TotalNulls++
			if err := p.Skip(); err != nil {
				return err
			}
		case gofasterjx.Bool:
			s.TotalBooleans++
			if err := p.Skip(); err != nil {
				return err
			}
		case gofasterjx.Number:
			s.TotalNumbers++
			if err := p.Skip(); err != nil {
				return err
			}
		case gofasterjx.Array:
			s.TotalArrays++
			i := 0
			if err := p.Arr(func(d *gofasterjx.Decoder) error {
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
			if err := p.ObjBytes(func(d *gofasterjx.Decoder, key []byte) error {
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

func MustCalcStatsValyalaFastjson(p *valyalafastjson.Parser, str []byte) (s Stats) {
	v, err := p.ParseBytes(str)
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
	const input = `{
		"s":"value",
		"t":true,
		"f":false,
		"0":null,
		"n":-9.123e3,
		"o0":{},
		"a0":[],
		"o":{
			"k":"\"v\"",
			"a":[
				true,
				null,
				"item",
				-67.02e9,
				["foo"]
			]
		},
		"[abc]":[0]
	}`
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
		p := jscan.NewParser[[]byte](64)
		require.Equal(t, expect, MustCalcStatsJscan(p, []byte(input)))
	})
	t.Run("jsoniter", func(t *testing.T) {
		p := jsoniter.NewIterator(jsoniter.ConfigFastest)
		require.Equal(t, expect, MustCalcStatsJsoniter(p, []byte(input)))
	})
	t.Run("gofaster-jx", func(t *testing.T) {
		p := new(gofasterjx.Decoder)
		require.Equal(t, expect, MustCalcStatsGofasterJx(p, []byte(input)))
	})

	t.Run("valyala-fastjson", func(t *testing.T) {
		p := new(valyalafastjson.Parser)
		require.Equal(t, expect, MustCalcStatsValyalaFastjson(p, []byte(input)))
	})
}

var gs Stats

func BenchmarkCalcStats(b *testing.B) {
	for _, bd := range []struct {
		name  string
		input SourceProvider
	}{
		{"miniscule_1b__________", SrcFile("miniscule_1b.json")},
		{"tiny_8b_______________", SrcFile("tiny_8b.json")},
		{"small_336b____________", SrcFile("small_336b.json")},
		{"large_26m_____________", SrcFile("large_26m.json.gz")},
		{"nasa_SxSW_2016_125k___", SrcFile("nasa_SxSW_2016_125k.json.gz")},
		{"escaped_3k____________", SrcFile("escaped_3k.json")},
		{"array_int_1024_12k____", SrcFile("array_int_1024_12k.json")},
		{"array_dec_1024_10k____", SrcFile("array_dec_1024_10k.json")},
		{"array_nullbool_1024_5k", SrcFile("array_nullbool_1024_5k.json")},
		{"array_str_1024_639k___", SrcFile("array_str_1024_639k.json")},
	} {
		b.Run(bd.name, func(b *testing.B) {
			src, err := bd.input.GetJSON()
			require.NoError(b, err)

			b.Run("jscan___________", func(b *testing.B) {
				p := jscan.NewParser[[]byte](64)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gs = MustCalcStatsJscan(p, src)
				}
			})

			b.Run("jsoniter________", func(b *testing.B) {
				p := jsoniter.NewIterator(jsoniter.ConfigFastest)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gs = MustCalcStatsJsoniter(p, src)
				}
			})

			b.Run("gofaster-jx_____", func(b *testing.B) {
				p := new(gofasterjx.Decoder)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gs = MustCalcStatsGofasterJx(p, src)
				}
			})

			b.Run("valyala-fastjson", func(b *testing.B) {
				p := new(valyalafastjson.Parser)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gs = MustCalcStatsValyalaFastjson(p, src)
				}
			})
		})
	}
}

func TestValid(t *testing.T) {
	j := `[false,[[2, {"[foo]":[{"bar-baz":"fuz"}]}]]]`
	require.True(t, encodingjson.Valid([]byte(j)))

	t.Run("jscan", func(t *testing.T) {
		require.True(t, jscan.Valid(j))
	})

	t.Run("jscan_verifier", func(t *testing.T) {
		require.True(t, jscan.Valid(j))
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

func repeat(s string, n int) string {
	b := make([]byte, n*len(s))
	for i := 0; i < n; i++ {
		b = append(b, s...)
	}
	return string(b)
}

func BenchmarkValid(b *testing.B) {
	for _, bd := range []struct {
		name  string
		input SourceProvider
	}{
		{"deeparray_____________", SrcMake(func() []byte {
			return []byte(repeat("[", 1024) + repeat("]", 1024))
		})},
		{"unwind_stack__________", SrcMake(func() []byte {
			return MakeRepeated("[", 1024)
		})},
		{"miniscule_1b__________", SrcFile("miniscule_1b.json")},
		{"tiny_8b_______________", SrcFile("tiny_8b.json")},
		{"small_336b____________", SrcFile("small_336b.json")},
		{"large_26m_____________", SrcFile("large_26m.json.gz")},
		{"nasa_SxSW_2016_125k___", SrcFile("nasa_SxSW_2016_125k.json.gz")},
		{"escaped_3k____________", SrcFile("escaped_3k.json")},
		{"array_int_1024_12k____", SrcFile("array_int_1024_12k.json")},
		{"array_dec_1024_10k____", SrcFile("array_dec_1024_10k.json")},
		{"array_nullbool_1024_5k", SrcFile("array_nullbool_1024_5k.json")},
		{"array_str_1024_639k___", SrcFile("array_str_1024_639k.json")},
	} {
		b.Run(bd.name, func(b *testing.B) {
			src, err := bd.input.GetJSON()
			require.NoError(b, err)

			b.Run("jscan___________", func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gbool = jscan.Valid(src)
				}
			})

			b.Run("encoding-json___", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					gbool = json.Valid(src)
				}
			})

			b.Run("jsoniter________", func(b *testing.B) {
				jb := []byte(src)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gbool = jsoniter.Valid(jb)
				}
			})

			b.Run("gofaster-jx_____", func(b *testing.B) {
				jb := []byte(src)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gbool = gofasterjx.Valid(jb)
				}
			})

			b.Run("tidwallgjson____", func(b *testing.B) {
				j := string(src)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gbool = tidwallgjson.Valid(j)
				}
			})

			b.Run("valyala-fastjson", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					gbool = (valyalafastjson.ValidateBytes(src) != nil)
				}
			})

			b.Run("goccy-go-json___", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					gbool = goccygojson.Valid(src)
				}
			})

			b.Run("bytedance-sonic_", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					gbool = bytedancesonic.ConfigFastest.Valid(src)
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

type SourceProvider interface{ GetJSON() ([]byte, error) }

type SrcMake func() []byte

func (s SrcMake) GetJSON() ([]byte, error) { return s(), nil }

type SrcFile string

func (s SrcFile) GetJSON() ([]byte, error) {
	p := filepath.Join("testdata", string(s))
	switch {
	case strings.HasSuffix(string(s), ".json"):
		return os.ReadFile(p)
	case strings.HasSuffix(string(s), ".gz"):
		f, err := os.Open(p)
		if err != nil {
			return nil, fmt.Errorf("opening archive file: %w", err)
		}
		r, err := gzip.NewReader(f)
		if err != nil {
			return nil, fmt.Errorf("initializing gzip reader: %w", err)
		}
		b, err := io.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("reading from gzip reader: %w", err)
		}
		return b, nil
	}
	return nil, fmt.Errorf("unsupported file: %q", s)
}
