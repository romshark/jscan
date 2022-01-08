package strfind_test

import (
	"testing"

	"github.com/romshark/jscan/internal/strfind"

	"github.com/stretchr/testify/require"
)

func TestIndexTerm(t *testing.T) {
	for _, tt := range []struct {
		in  string
		i   int
		exp int
	}{
		{`value`, 1, -1},
		{`"`, 0, 0},
		{`\""`, 0, 2},
		{`\\"`, 0, 2},
		{`\\\""`, 0, 4},
		{`abcd\\\""`, 3, 8},
	} {
		t.Run("", func(t *testing.T) {
			a := strfind.IndexTerm(tt.in, tt.i)
			require.Equal(t, tt.exp, a)
		})
	}
}

func TestLastIndexUnescaped(t *testing.T) {
	for _, tt := range []struct {
		in  string
		exp int
	}{
		{``, -1},
		{`x`, 0},
		{`\x`, -1},
		{`\\x`, 2},
		{`\\\x`, -1},
		{`x\\\x`, 0},
	} {
		t.Run("", func(t *testing.T) {
			a := strfind.LastIndexUnescaped([]byte(tt.in), 'x')
			require.Equal(t, tt.exp, a)
		})
	}
}
