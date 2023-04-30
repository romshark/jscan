package jsonnum_test

import (
	"testing"

	"github.com/romshark/jscan/internal/jsonnum"
)

var gend int

func BenchmarkValid(b *testing.B) {
	var err bool
	for _, bb := range []string{
		"0,",
		"1e10,",
		"1234567890,",
		"-1234567890,",
		"-1234567890e1,",
		"-1234567890e-123456789,",
		"-1234567890.123123424234234e-123456789,",
	} {
		b.Run("", func(b *testing.B) {
			b.Run("string", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					gend, err = jsonnum.EndIndex(bb)
					if err {
						b.Fatal("unexpected error")
					}
				}
			})

			b.Run("bytes", func(b *testing.B) {
				bb := []byte(bb)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gend, err = jsonnum.EndIndex(bb)
					if err {
						b.Fatal("unexpected error")
					}
				}
			})
		})
	}
}

func BenchmarkInvalid(b *testing.B) {
	var err bool
	for _, bb := range []string{
		"a",
		"-",
		"0.",
		"1234567890.",
		"01",
		"0e",
		"1e",
		"1234567890e",
		"1e-",
		"1234567890.1234567890e",
		"0.1234567890e",
	} {
		b.Run("", func(b *testing.B) {
			b.Run("string", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					gend, err = jsonnum.EndIndex(bb)
					if !err {
						b.Fatal("error expected")
					}
				}
			})

			b.Run("bytes", func(b *testing.B) {
				bb := []byte(bb)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					gend, err = jsonnum.EndIndex(bb)
					if !err {
						b.Fatal("error expected")
					}
				}
			})
		})
	}
}
