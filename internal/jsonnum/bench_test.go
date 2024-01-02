package jsonnum_test

import (
	"runtime"
	"testing"

	"github.com/romshark/jscan/v2/internal/jsonnum"
)

func BenchmarkValid(b *testing.B) {
	var rc jsonnum.ReturnCode
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
					if _, rc = jsonnum.ReadNumber(bb); rc == jsonnum.ReturnCodeErr {
						b.Fatal("unexpected error")
					}
				}
			})

			b.Run("bytes", func(b *testing.B) {
				bb := []byte(bb)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					if _, rc = jsonnum.ReadNumber(bb); rc == jsonnum.ReturnCodeErr {
						b.Fatal("unexpected error")
					}
				}
			})
		})
	}
}

func BenchmarkInvalid(b *testing.B) {
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
			// This rc will not be checked since "01" is not technically wrong
			// according to jsonnum.ReadNumber, it would return ("1", false) instead.
			// All inputs are already tested in TestReadNumberErr and TestReadNumberZero.
			var rc jsonnum.ReturnCode
			var remainderString string
			var remainderBytes []byte

			b.Run("string", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					remainderString, rc = jsonnum.ReadNumber(bb)
				}
			})

			b.Run("bytes", func(b *testing.B) {
				bb := []byte(bb)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					remainderBytes, rc = jsonnum.ReadNumber(bb)
				}
			})

			runtime.KeepAlive(rc)
			runtime.KeepAlive(remainderString)
			runtime.KeepAlive(remainderBytes)
		})
	}
}
