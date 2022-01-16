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
			t.Run("string", func(t *testing.T) {
				a := strfind.IndexTerm(tt.in, tt.i)
				require.Equal(t, tt.exp, a)
			})

			t.Run("bytes", func(t *testing.T) {
				a := strfind.IndexTermBytes([]byte(tt.in), tt.i)
				require.Equal(t, tt.exp, a)
			})
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

func TestEndOfWhitespaceSeq(t *testing.T) {
	for _, tt := range []struct {
		in  string
		exp int
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
	} {
		t.Run("", func(t *testing.T) {
			t.Run("string", func(t *testing.T) {
				a := strfind.EndOfWhitespaceSeq(tt.in)
				require.Equal(t, tt.exp, a)
			})

			t.Run("bytes", func(t *testing.T) {
				a := strfind.EndOfWhitespaceSeqBytes([]byte(tt.in))
				require.Equal(t, tt.exp, a)
			})
		})
	}
}
