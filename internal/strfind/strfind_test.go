package strfind_test

import (
	_ "embed"
	"testing"

	"github.com/romshark/jscan/v2/internal/strfind"

	"github.com/stretchr/testify/require"
)

func TestEndOfWhitespaceSeq(t *testing.T) {
	for _, tt := range []struct {
		input              string
		expect             int
		expectIllegalChars bool
	}{
		{"", 0, false},
		{"e", 0, false},
		{" ", 1, false},
		{" \r\n\t", 4, false},

		{"\n", 1, false},
		{"\t", 1, false},
		{"\r", 1, false},

		{" e", 1, false},
		{"\ne", 1, false},
		{"\te", 1, false},
		{"\re", 1, false},

		{"   abc", 3, false},
		{"  \nabc", 3, false},
		{"  \tabc", 3, false},
		{"  \rabc", 3, false},

		{"0123456789", 0, false},
		{" 0123456789", 1, false},
		{" \r\n\t0123456789", 4, false},

		{"\n0123456789", 1, false},
		{"\t0123456789", 1, false},
		{"\r0123456789", 1, false},

		{" e0123456789", 1, false},
		{"\ne0123456789", 1, false},
		{"\te0123456789", 1, false},
		{"\re0123456789", 1, false},

		{"   abc0123456789", 3, false},
		{"  \nabc0123456789", 3, false},
		{"  \tabc0123456789", 3, false},
		{"  \rabc0123456789", 3, false},

		{repeat(" ", 1) + repeat("x", 64), 1, false},
		{repeat(" ", 2) + repeat("x", 64), 2, false},
		{repeat(" ", 3) + repeat("x", 64), 3, false},
		{repeat(" ", 4) + repeat("x", 64), 4, false},
		{repeat(" ", 5) + repeat("x", 64), 5, false},
		{repeat(" ", 6) + repeat("x", 64), 6, false},
		{repeat(" ", 7) + repeat("x", 64), 7, false},
		{repeat(" ", 8) + repeat("x", 64), 8, false},
		{repeat(" ", 9) + repeat("x", 64), 9, false},
		{repeat(" ", 10) + repeat("x", 64), 10, false},
		{repeat(" ", 11) + repeat("x", 64), 11, false},
		{repeat(" ", 12) + repeat("x", 64), 12, false},
		{repeat(" ", 13) + repeat("x", 64), 13, false},
		{repeat(" ", 14) + repeat("x", 64), 14, false},
		{repeat(" ", 15) + repeat("x", 64), 15, false},
		{repeat(" ", 16) + repeat("x", 64), 16, false},

		{string(byte(0x1F)), 0, true},
		{"\00123456789", 0, true},
		{repeat(" ", 1) + "\001" + repeat("x", 64), 1, true},
		{repeat(" ", 2) + "\001" + repeat("x", 64), 2, true},
		{repeat(" ", 3) + "\001" + repeat("x", 64), 3, true},
		{repeat(" ", 4) + "\001" + repeat("x", 64), 4, true},
		{repeat(" ", 5) + "\001" + repeat("x", 64), 5, true},
		{repeat(" ", 6) + "\001" + repeat("x", 64), 6, true},
		{repeat(" ", 7) + "\001" + repeat("x", 64), 7, true},
		{repeat(" ", 8) + "\001" + repeat("x", 64), 8, true},
		{repeat(" ", 9) + "\001" + repeat("x", 64), 9, true},
		{repeat(" ", 10) + "\001" + repeat("x", 64), 10, true},
		{repeat(" ", 11) + "\001" + repeat("x", 64), 11, true},
		{repeat(" ", 12) + "\001" + repeat("x", 64), 12, true},
		{repeat(" ", 13) + "\001" + repeat("x", 64), 13, true},
		{repeat(" ", 14) + "\001" + repeat("x", 64), 14, true},
		{repeat(" ", 15) + "\001" + repeat("x", 64), 15, true},
		{repeat(" ", 16) + "\001" + repeat("x", 64), 16, true},

		{"\x000123456789", 0, true},
		{"   \x00a0123456789", 3, true},
	} {
		t.Run("", func(t *testing.T) {
			trailing, ilc := strfind.EndOfWhitespaceSeq(tt.input)
			require.Equal(t, tt.expect, len(tt.input)-len(trailing))
			require.Equal(t, tt.expectIllegalChars, ilc)
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
