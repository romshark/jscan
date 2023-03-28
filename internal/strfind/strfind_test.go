package strfind_test

import (
	_ "embed"
	"testing"

	"github.com/romshark/jscan/internal/strfind"

	"github.com/stretchr/testify/require"
)

//go:embed test_longstr.txt
var longStrTXT string

var testsIndexTerm = []struct {
	name            string
	input           string
	i               int
	expectIndexEnd  int
	expectErrorCode strfind.ErrCode
}{
	{
		name:           "ok_empty_string",
		input:          `"`,
		i:              0,
		expectIndexEnd: 0,
	},
	{
		name:           "ok_escaped_quotes",
		input:          `\""`,
		i:              0,
		expectIndexEnd: 2,
	},
	{
		name:           "ok_escaped_backslash",
		input:          `\\"`,
		i:              0,
		expectIndexEnd: 2,
	},
	{
		name:           "ok_escaped_bashslash_and_escaped_quotes",
		input:          `\\\""`,
		i:              0,
		expectIndexEnd: 4,
	},
	{
		name:           "ok_text_followed_by_escape_sequences",
		input:          `abcd\\\""`,
		i:              3,
		expectIndexEnd: 8,
	},
	{
		name:           "ok_escaped_slash",
		input:          `\/"`,
		i:              0,
		expectIndexEnd: 2,
	},
	{
		name:           "ok_escaped_backspace",
		input:          `\b"`,
		i:              0,
		expectIndexEnd: 2,
	},
	{
		name:           "ok_escaped_formfeed",
		input:          `\f"`,
		i:              0,
		expectIndexEnd: 2,
	},
	{
		name:           "ok_escaped_newline",
		input:          `\n"`,
		i:              0,
		expectIndexEnd: 2,
	},
	{
		name:           "ok_escaped_carriage_return",
		input:          `\r"`,
		i:              0,
		expectIndexEnd: 2,
	},
	{
		name:           "ok_escaped_tab",
		input:          `\t"`,
		i:              0,
		expectIndexEnd: 2,
	},
	{
		name:           "ok_escaped_hex",
		input:          `\uffff"`,
		i:              0,
		expectIndexEnd: 6,
	},
	{
		name:           "ok_longstr",
		input:          longStrTXT,
		i:              1,
		expectIndexEnd: len(longStrTXT) - 1,
	},
	{
		name:           "ok_escaped_at_0",
		input:          `\"tailtext"`,
		i:              0,
		expectIndexEnd: len(`\"tailtext`),
	},
	{
		name:           "ok_escaped_at_1",
		input:          `0\"tailtext"`,
		i:              0,
		expectIndexEnd: len(`0\"tailtext`),
	},
	{
		name:           "ok_escaped_at_2",
		input:          `01\"tailtext"`,
		i:              0,
		expectIndexEnd: len(`01\"tailtext`),
	},
	{
		name:           "ok_escaped_at_3",
		input:          `012\"tailtext"`,
		i:              0,
		expectIndexEnd: len(`012\"tailtext`),
	},
	{
		name:           "ok_escaped_at_4",
		input:          `0123\"tailtext"`,
		i:              0,
		expectIndexEnd: len(`0123\"tailtext`),
	},
	{
		name:           "ok_escaped_at_5",
		input:          `01234\"tailtext"`,
		i:              0,
		expectIndexEnd: len(`01234\"tailtext`),
	},
	{
		name:           "ok_escaped_at_6",
		input:          `012345\"tailtext"`,
		i:              0,
		expectIndexEnd: len(`012345\"tailtext`),
	},
	{
		name:           "ok_escaped_at_7",
		input:          `0123456\"tailtext"`,
		i:              0,
		expectIndexEnd: len(`0123456\"tailtext`),
	},
	{
		name:           "ok_escaped_at_8",
		input:          `01234567\"tailtext"`,
		i:              0,
		expectIndexEnd: len(`01234567\"tailtext`),
	},

	// Errors
	{
		name:            "err_unexpeof_no_terminator",
		input:           ``,
		i:               0,
		expectIndexEnd:  0,
		expectErrorCode: strfind.ErrCodeUnexpectedEOF,
	},
	{
		name:            "err_unexpeof_text_followed_by_no_terminator",
		input:           `value`,
		i:               3,
		expectIndexEnd:  len("value"),
		expectErrorCode: strfind.ErrCodeUnexpectedEOF,
	},
	{
		name:            "err_unexpeof_after_escape",
		input:           `value\`,
		i:               0,
		expectIndexEnd:  len(`value\`),
		expectErrorCode: strfind.ErrCodeUnexpectedEOF,
	},
	{
		name:            "err_controlchar",
		input:           `ab` + string(byte(0x1F)) + `c"`,
		i:               0,
		expectIndexEnd:  2,
		expectErrorCode: strfind.ErrCodeIllegalControlChar,
	},
	{
		name:            "err_escapechar",
		input:           `\0"`,
		i:               0,
		expectIndexEnd:  1,
		expectErrorCode: strfind.ErrCodeInvalidEscapeSeq,
	},
	{
		name:            "err_illegal_escape_sequence",
		input:           `escaped: \u000k"`,
		i:               0,
		expectIndexEnd:  len(`escaped: \`),
		expectErrorCode: strfind.ErrCodeInvalidEscapeSeq,
	},
}

func TestIndexTerm(t *testing.T) {
	for _, tt := range testsIndexTerm {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("string", func(t *testing.T) {
				a, errCode := strfind.IndexTerm(tt.input, tt.i)
				require.Equal(t, tt.expectErrorCode, errCode, "error code")
				require.Equal(t, tt.expectIndexEnd, a)
			})

			t.Run("bytes", func(t *testing.T) {
				a, errCode := strfind.IndexTerm([]byte(tt.input), tt.i)
				require.Equal(t, tt.expectErrorCode, errCode, "error code")
				require.Equal(t, tt.expectIndexEnd, a)
			})
		})
	}
}

func TestLastIndexUnescaped(t *testing.T) {
	for _, tt := range []struct {
		input  string
		expect int
	}{
		{``, -1},
		{`x`, 0},
		{`\x`, -1},
		{`\\x`, 2},
		{`\\\x`, -1},
		{`x\\\x`, 0},
		{`xxxxx`, 4},
	} {
		t.Run("", func(t *testing.T) {
			a := strfind.LastIndexUnescaped([]byte(tt.input), 'x')
			require.Equal(t, tt.expect, a)
		})
	}
}

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

		{"\u0000", 0, true},
		{"   \u0000a", 3, true},
	} {
		t.Run("", func(t *testing.T) {
			t.Run("string", func(t *testing.T) {
				a, ilc := strfind.EndOfWhitespaceSeq(tt.input)
				require.Equal(t, tt.expect, a)
				require.Equal(t, tt.expectIllegalChars, ilc)
			})

			t.Run("bytes", func(t *testing.T) {
				a, ilc := strfind.EndOfWhitespaceSeqBytes([]byte(tt.input))
				require.Equal(t, tt.expect, a)
				require.Equal(t, tt.expectIllegalChars, ilc)
			})
		})
	}
}
