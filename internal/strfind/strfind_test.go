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
		in              string
		exp             int
		expIllegalChars bool
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

		{"\u0000", 0, true},
		{"   \u0000a", 3, true},
	} {
		t.Run("", func(t *testing.T) {
			t.Run("string", func(t *testing.T) {
				a, ilc := strfind.EndOfWhitespaceSeq(tt.in)
				require.Equal(t, tt.exp, a)
				require.Equal(t, tt.expIllegalChars, ilc)
			})

			t.Run("bytes", func(t *testing.T) {
				a, ilc := strfind.EndOfWhitespaceSeqBytes([]byte(tt.in))
				require.Equal(t, tt.exp, a)
				require.Equal(t, tt.expIllegalChars, ilc)
			})
		})
	}
}

func TestContainsCtrlCharFalse(t *testing.T) {
	for _, tt := range []string{
		"",
		" ",
		"abc ",
		`\u0000\t\r\n\x00`,
	} {
		t.Run("", func(t *testing.T) {
			t.Run("string", func(t *testing.T) {
				require.False(t, strfind.ContainsCtrlChar(tt))
			})

			t.Run("bytes", func(t *testing.T) {
				require.False(t, strfind.ContainsCtrlCharBytes([]byte(tt)))
			})
		})
	}
}

func TestContainsCtrlChar(t *testing.T) {
	forStringIn2049 := func(withIllegalPostfix bool, fn func(s []byte)) {
		for i := 1; i <= 2049; i++ {
			s := make([]byte, i)
			for i := range s {
				s[i] = 'x'
			}
			if withIllegalPostfix {
				s[len(s)-1] = 0x00
			}
			fn(s)
		}
	}

	t.Run("contains", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			forStringIn2049(true, func(s []byte) {
				require.True(
					t, strfind.ContainsCtrlChar(string(s)),
					"undetected control character at %q", string(s),
				)
			})
		})

		t.Run("bytes", func(t *testing.T) {
			forStringIn2049(true, func(s []byte) {
				require.True(
					t, strfind.ContainsCtrlCharBytes(s),
					"undetected control character at %q", string(s),
				)
			})
		})
	})

	t.Run("notcontains", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			forStringIn2049(false, func(s []byte) {
				require.False(
					t, strfind.ContainsCtrlChar(string(s)),
					"false positive at %q", string(s),
				)
			})
		})

		t.Run("bytes", func(t *testing.T) {
			forStringIn2049(false, func(s []byte) {
				require.False(
					t, strfind.ContainsCtrlCharBytes(s),
					"false positive at %q", string(s),
				)
			})
		})
	})
}
