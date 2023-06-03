package strfind_test

import (
	"testing"

	"github.com/romshark/jscan/internal/strfind"
)

var (
	GS string
	GE strfind.ErrCode
)

func BenchmarkReadString(b *testing.B) {
	for _, bb := range testsReadString {
		b.Run(bb.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				GS, GE = strfind.ReadString(bb.input)
			}
		})
	}
}
