package keyescape_test

import (
	"testing"

	"github.com/romshark/jscan/internal/keyescape"
	"github.com/stretchr/testify/require"
)

var testdata = []struct {
	in  string
	exp string
}{
	{"", ""},
	{"abc", "abc"},
	{".[]", `\.\[\]`},
	{"foo.bar[12].baz[1]", `foo\.bar\[12\]\.baz\[1\]`},
}

func TestEscape(t *testing.T) {
	for _, tt := range testdata {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.exp, keyescape.Escape(tt.in))
		})
	}
}

func TestEscapeAppend(t *testing.T) {
	for _, tt := range testdata {
		t.Run("", func(t *testing.T) {
			var buf []byte
			buf = keyescape.EscapeAppend(buf, []byte(tt.in))
			require.Equal(t, tt.exp, string(buf))
		})
	}
}
