package jscan_test

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/romshark/jscan/v2"

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

func MustCalcStatsJscanTokenizer[S []byte | string](p *jscan.Tokenizer[S], str S) (s Stats) {
	if err := p.Tokenize(
		str,
		func(tokens []jscan.Token[S]) (err bool) {
			depth := 0
			for i := range tokens {
				switch tokens[i].Type {
				case jscan.TokenTypeKey:
					l := tokens[i].End - tokens[i].Index - 2
					s.TotalKeys++
					if l > s.MaxKeyLen {
						s.MaxKeyLen = l
					}
				case jscan.TokenTypeObject:
					depth++
					s.TotalObjects++
				case jscan.TokenTypeArray:
					depth++
					if depth > s.MaxDepth {
						s.MaxDepth = depth
					}
					s.TotalArrays++
					if tokens[i].Elements > s.MaxArrayLen {
						s.MaxArrayLen = tokens[i].Elements
					}
				case jscan.TokenTypeNull:
					s.TotalNulls++
				case jscan.TokenTypeFalse, jscan.TokenTypeTrue:
					s.TotalBooleans++
				case jscan.TokenTypeNumber, jscan.TokenTypeInteger:
					s.TotalNumbers++
				case jscan.TokenTypeString:
					s.TotalStrings++
				case jscan.TokenTypeObjectEnd, jscan.TokenTypeArrayEnd:
					depth--
				}
			}
			return false
		},
	); err.IsErr() {
		panic(fmt.Errorf("unexpected error: %s", err))
	}
	return
}

func TestCalcStats(t *testing.T) {
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

	p := jscan.NewParser[[]byte](64)
	require.Equal(t, expect, MustCalcStatsJscan(p, []byte(input)))

	k := jscan.NewTokenizer[[]byte](128, 10)
	require.Equal(t, expect, MustCalcStatsJscanTokenizer(k, []byte(input)))
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
			b.Run("parser", func(b *testing.B) {
				src, err := bd.input.GetJSON()
				require.NoError(b, err)

				p := jscan.NewParser[[]byte](1024)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gs = MustCalcStatsJscan(p, src)
				}
			})
			b.Run("tokenizer", func(b *testing.B) {
				src, err := bd.input.GetJSON()
				require.NoError(b, err)

				k := jscan.NewTokenizer[[]byte](64, 1024)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gs = MustCalcStatsJscanTokenizer(k, src)
				}
			})
		})
	}
}

func TestValid(t *testing.T) {
	j := `[false,[[2, {"[foo]":[{"bar-baz":"fuz"}]}]]]`
	require.True(t, json.Valid([]byte(j)))
	require.True(t, jscan.Valid(j))
}

var GB bool

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
			return []byte(repeat("[", 1024))
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

			v := jscan.NewValidator[[]byte](1024)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				GB = v.Valid(src)
			}
		})
	}
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
