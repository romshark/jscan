package jscan

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSkipSpace(t *testing.T) {
	for _, tt := range []struct {
		input  string
		expect int
	}{
		{"", 0},
		{"e", 0},
		{" ", 1},
		{" \r\n\t", 4},

		{"\n", 1},
		{"\t", 1},
		{"\r", 1},

		{" e", 1},
		{"\ne", 1},
		{"\te", 1},
		{"\re", 1},

		{"   abc", 3},
		{"  \nabc", 3},
		{"  \tabc", 3},
		{"  \rabc", 3},

		{"0123456789", 0},
		{" 0123456789", 1},
		{" \r\n\t0123456789", 4},

		{"\n0123456789", 1},
		{"\t0123456789", 1},
		{"\r0123456789", 1},

		{" e0123456789", 1},
		{"\ne0123456789", 1},
		{"\te0123456789", 1},
		{"\re0123456789", 1},

		{"   abc0123456789", 3},
		{"  \nabc0123456789", 3},
		{"  \tabc0123456789", 3},
		{"  \rabc0123456789", 3},

		{repeat(" ", 1) + repeat("x", 64), 1},
		{repeat(" ", 2) + repeat("x", 64), 2},
		{repeat(" ", 3) + repeat("x", 64), 3},
		{repeat(" ", 4) + repeat("x", 64), 4},
		{repeat(" ", 5) + repeat("x", 64), 5},
		{repeat(" ", 6) + repeat("x", 64), 6},
		{repeat(" ", 7) + repeat("x", 64), 7},
		{repeat(" ", 8) + repeat("x", 64), 8},
		{repeat(" ", 9) + repeat("x", 64), 9},
		{repeat(" ", 10) + repeat("x", 64), 10},
		{repeat(" ", 11) + repeat("x", 64), 11},
		{repeat(" ", 12) + repeat("x", 64), 12},
		{repeat(" ", 13) + repeat("x", 64), 13},
		{repeat(" ", 14) + repeat("x", 64), 14},
		{repeat(" ", 15) + repeat("x", 64), 15},
		{repeat(" ", 16) + repeat("x", 64), 16},

		{string(byte(0x1F)), 0},
		{"\00123456789", 0},
		{repeat(" ", 1) + "\001" + repeat("x", 64), 1},
		{repeat(" ", 2) + "\001" + repeat("x", 64), 2},
		{repeat(" ", 3) + "\001" + repeat("x", 64), 3},
		{repeat(" ", 4) + "\001" + repeat("x", 64), 4},
		{repeat(" ", 5) + "\001" + repeat("x", 64), 5},
		{repeat(" ", 6) + "\001" + repeat("x", 64), 6},
		{repeat(" ", 7) + "\001" + repeat("x", 64), 7},
		{repeat(" ", 8) + "\001" + repeat("x", 64), 8},
		{repeat(" ", 9) + "\001" + repeat("x", 64), 9},
		{repeat(" ", 10) + "\001" + repeat("x", 64), 10},
		{repeat(" ", 11) + "\001" + repeat("x", 64), 11},
		{repeat(" ", 12) + "\001" + repeat("x", 64), 12},
		{repeat(" ", 13) + "\001" + repeat("x", 64), 13},
		{repeat(" ", 14) + "\001" + repeat("x", 64), 14},
		{repeat(" ", 15) + "\001" + repeat("x", 64), 15},
		{repeat(" ", 16) + "\001" + repeat("x", 64), 16},

		{"\x000123456789", 0},
		{"   \x00a0123456789", 3},
	} {
		t.Run("", func(t *testing.T) {
			trailing := skipSpace(tt.input)
			require.Equal(t, tt.expect, len(tt.input)-len(trailing))
		})
	}
}

func repeat(x string, n int) string {
	s := make([]byte, 0, n*len(x))
	for i := 0; i < n; i++ {
		s = append(s, x...)
	}
	return string(s)
}
