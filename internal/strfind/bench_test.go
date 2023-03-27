package strfind_test

import (
	"testing"

	"github.com/romshark/jscan/internal/strfind"
)

var (
	GI int
	GE strfind.ErrCode
)

func BenchmarkIndexTerm(b *testing.B) {
	for _, bb := range testsIndexTerm {
		b.Run(bb.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				GI, GE = strfind.IndexTerm(bb.input, bb.i)
			}
		})
	}
}
