package strfind_test

import (
	_ "embed"
	"testing"

	"github.com/romshark/jscan/internal/strfind"

	"github.com/stretchr/testify/require"
)

//go:embed test_longstr.txt
var longStrTXT string

var testsReadString = []struct {
	name            string
	input           string
	expectAfter     *string // `"` by default if not set
	expectErrorCode strfind.ErrCode
}{
	{
		name:  "ok_empty_string",
		input: `"`,
	},
	{
		name:  "ok_escaped_quotes",
		input: `\""`,
	},
	{
		name:  "ok_escaped_backslash",
		input: `\\"`,
	},
	{
		name:  "ok_escaped_bashslash_and_escaped_quotes",
		input: `\\\""`,
	},
	{
		name:  "ok_text_followed_by_escape_sequences",
		input: `abcd\\\""`,
	},
	{
		name:  "ok_escaped_slash",
		input: `\/"`,
	},
	{
		name:  "ok_escaped_backspace",
		input: `\b"`,
	},
	{
		name:  "ok_escaped_formfeed",
		input: `\f"`,
	},
	{
		name:  "ok_escaped_newline",
		input: `\n"`,
	},
	{
		name:  "ok_escaped_carriage_return",
		input: `\r"`,
	},
	{
		name:  "ok_escaped_tab",
		input: `\t"`,
	},
	{
		name:  "ok_escaped_hex",
		input: `\uffff"`,
	},
	{
		name:  "ok_longstr",
		input: longStrTXT,
	},
	{
		name:  "ok_escaped_at_0",
		input: `\""`,
	},
	{
		name:  "ok_escaped_at_1",
		input: `0\""`,
	},
	{
		name:  "ok_escaped_at_2",
		input: `01\""`,
	},
	{
		name:  "ok_escaped_at_3",
		input: `012\""`,
	},
	{
		name:  "ok_escaped_at_4",
		input: `0123\""`,
	},
	{
		name:  "ok_escaped_at_5",
		input: `01234\""`,
	},
	{
		name:  "ok_escaped_at_6",
		input: `012345\""`,
	},
	{
		name:  "ok_escaped_at_7",
		input: `0123456\""`,
	},
	{
		name:  "ok_escaped_at_8",
		input: `01234567\""`,
	},
	{
		name:  "ok_escaped_at_9",
		input: `012345678\""`,
	},
	{
		name:  "ok_escaped_at_10",
		input: `0123456789\""`,
	},
	{
		name:  "ok_escaped_at_0",
		input: `\"` + repeat("x", 64) + `"`,
	},
	{
		name:  "ok_escaped_at_1",
		input: `0\"` + repeat("x", 64) + `"`,
	},
	{
		name:  "ok_escaped_at_2",
		input: `01\"` + repeat("x", 64) + `"`,
	},
	{
		name:  "ok_escaped_at_3",
		input: `012\"` + repeat("x", 64) + `"`,
	},
	{
		name:  "ok_escaped_at_4",
		input: `0123\"` + repeat("x", 64) + `"`,
	},
	{
		name:  "ok_escaped_at_5",
		input: `01234\"` + repeat("x", 64) + `"`,
	},
	{
		name:  "ok_escaped_at_6",
		input: `012345\"` + repeat("x", 64) + `"`,
	},
	{
		name:  "ok_escaped_at_7",
		input: `0123456\"` + repeat("x", 64) + `"`,
	},
	{
		name:  "ok_escaped_at_8",
		input: `01234567\"` + repeat("x", 64) + `"`,
	},
	{
		name:  "ok_escaped_at_9",
		input: `012345678\"` + repeat("x", 64) + `"`,
	},
	{
		name:  "ok_escaped_at_10",
		input: `0123456789\"` + repeat("x", 64) + `"`,
	},
	{
		name:  "ok_escaped_at_11_p64",
		input: repeat("x", 11) + `\"` + repeat("x", 64) + `"`,
	},
	{
		name:  "ok_escaped_at_12_p64",
		input: repeat("x", 12) + `\"` + repeat("x", 64) + `"`,
	},
	{
		name:  "ok_escaped_at_13_p64",
		input: repeat("x", 13) + `\"` + repeat("x", 64) + `"`,
	},
	{
		name:  "ok_escaped_at_14_p64",
		input: repeat("x", 14) + `\"` + repeat("x", 64) + `"`,
	},
	{
		name:  "ok_escaped_at_15_p64",
		input: repeat("x", 15) + `\"` + repeat("x", 64) + `"`,
	},
	{
		name:  "ok_escaped_at_16_p64",
		input: repeat("x", 16) + `\"` + repeat("x", 64) + `"`,
	},

	// Errors
	{
		name:            "err_unexpeof_no_terminator",
		input:           ``,
		expectAfter:     str(""),
		expectErrorCode: strfind.ErrCodeUnexpectedEOF,
	},
	{
		name:            "err_unexpeof_text_followed_by_no_terminator",
		input:           `value`,
		expectAfter:     str(""),
		expectErrorCode: strfind.ErrCodeUnexpectedEOF,
	},
	{
		name:            "err_unexpeof_after_escape",
		input:           `value\`,
		expectAfter:     str(`\`),
		expectErrorCode: strfind.ErrCodeUnexpectedEOF,
	},
	{
		name:            "err_escapechar",
		input:           `\0"`,
		expectAfter:     str(`\0"`),
		expectErrorCode: strfind.ErrCodeInvalidEscapeSeq,
	},
	{
		name:            "err_illegal_escape_sequence",
		input:           `escaped: \u000k"`,
		expectAfter:     str(`\u000k"`),
		expectErrorCode: strfind.ErrCodeInvalidEscapeSeq,
	},
	{
		name:            "err_controlchar_at_0",
		input:           "\x00",
		expectAfter:     str("\x00"),
		expectErrorCode: strfind.ErrCodeIllegalControlChar,
	},
	{
		name:            "err_controlchar_at_1",
		input:           "0\x00",
		expectAfter:     str("\x00"),
		expectErrorCode: strfind.ErrCodeIllegalControlChar,
	},
	{
		name:            "err_controlchar_at_2",
		input:           "01\x00",
		expectAfter:     str("\x00"),
		expectErrorCode: strfind.ErrCodeIllegalControlChar,
	},
	{
		name:            "err_controlchar_at_3",
		input:           "012\x00",
		expectAfter:     str("\x00"),
		expectErrorCode: strfind.ErrCodeIllegalControlChar,
	},
	{
		name:            "err_controlchar_at_4",
		input:           "0123\x00",
		expectAfter:     str("\x00"),
		expectErrorCode: strfind.ErrCodeIllegalControlChar,
	},
	{
		name:            "err_controlchar_at_5",
		input:           "01234\x00",
		expectAfter:     str("\x00"),
		expectErrorCode: strfind.ErrCodeIllegalControlChar,
	},
	{
		name:            "err_controlchar_at_6",
		input:           "012345\x00",
		expectAfter:     str("\x00"),
		expectErrorCode: strfind.ErrCodeIllegalControlChar,
	},
	{
		name:            "err_controlchar_at_7",
		input:           "0123456\x00",
		expectAfter:     str("\x00"),
		expectErrorCode: strfind.ErrCodeIllegalControlChar,
	},
	{
		name:            "err_controlchar_at_8",
		input:           "01234567\x00",
		expectAfter:     str("\x00"),
		expectErrorCode: strfind.ErrCodeIllegalControlChar,
	},
	{
		name:            "err_controlchar_at_9",
		input:           "012345678\x00",
		expectAfter:     str("\x00"),
		expectErrorCode: strfind.ErrCodeIllegalControlChar,
	},
	{
		name:            "err_controlchar_at_10",
		input:           "0123456789\x00",
		expectAfter:     str("\x00"),
		expectErrorCode: strfind.ErrCodeIllegalControlChar,
	},
}

func TestReadString(t *testing.T) {
	for _, tt := range testsReadString {
		t.Run(tt.name, func(t *testing.T) {
			expectAfter := `"`
			if tt.expectAfter != nil {
				expectAfter = *tt.expectAfter
			}

			t.Run("string", func(t *testing.T) {
				a, errCode := strfind.ReadString(tt.input)
				require.Equal(t, tt.expectErrorCode, errCode, "error code")
				require.Equal(t, expectAfter, a)
			})

			t.Run("bytes", func(t *testing.T) {
				a, errCode := strfind.ReadString([]byte(tt.input))
				require.Equal(t, tt.expectErrorCode, errCode, "error code")
				require.Equal(t, expectAfter, string(a))
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

func str(s string) *string { return &s }
