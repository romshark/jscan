package jsonnum_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/romshark/jscan/v2/internal/jsonnum"

	"github.com/stretchr/testify/require"
)

func TestReadNumber(t *testing.T) {
	for _, tt := range []struct {
		input       string
		expectAfter string
		expect      jsonnum.ReturnCode
	}{
		{"0", "", jsonnum.ReturnCodeInteger},
		{"12", "", jsonnum.ReturnCodeInteger},
		{"234", "", jsonnum.ReturnCodeInteger},
		{"3456", "", jsonnum.ReturnCodeInteger},
		{"45678", "", jsonnum.ReturnCodeInteger},
		{"567890", "", jsonnum.ReturnCodeInteger},
		{"6789012", "", jsonnum.ReturnCodeInteger},
		{"78901234", "", jsonnum.ReturnCodeInteger},
		{"890123456", "", jsonnum.ReturnCodeInteger},
		{"9012345678", "", jsonnum.ReturnCodeInteger},
		{"90123456789", "", jsonnum.ReturnCodeInteger},
		{"901234567890", "", jsonnum.ReturnCodeInteger},
		{"9012345678901", "", jsonnum.ReturnCodeInteger},
		{"90123456789012", "", jsonnum.ReturnCodeInteger},

		{"1x", "x", jsonnum.ReturnCodeInteger},

		{"0xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeInteger},
		{"12xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeInteger},
		{"234xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeInteger},
		{"3456xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeInteger},
		{"45678xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeInteger},
		{"567890xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeInteger},
		{"6789012xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeInteger},
		{"78901234xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeInteger},
		{"890123456xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeInteger},
		{"9012345678xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeInteger},
		{"90123456789xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeInteger},
		{"901234567890xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeInteger},
		{"9012345678901xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeInteger},
		{"90123456789012xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeInteger},

		{"0.0x", "x", jsonnum.ReturnCodeNumber},

		{"0.0", "", jsonnum.ReturnCodeNumber},
		{"0.12", "", jsonnum.ReturnCodeNumber},
		{"0.345", "", jsonnum.ReturnCodeNumber},
		{"0.4567", "", jsonnum.ReturnCodeNumber},
		{"0.56789", "", jsonnum.ReturnCodeNumber},
		{"0.678901", "", jsonnum.ReturnCodeNumber},
		{"0.7890123", "", jsonnum.ReturnCodeNumber},
		{"0.89012345", "", jsonnum.ReturnCodeNumber},
		{"0.901234567", "", jsonnum.ReturnCodeNumber},
		{"0.0123456789", "", jsonnum.ReturnCodeNumber},
		{"0.12345678901", "", jsonnum.ReturnCodeNumber},
		{"0.234567890123", "", jsonnum.ReturnCodeNumber},
		{"0.3456789012345", "", jsonnum.ReturnCodeNumber},
		{"0.45678901234567", "", jsonnum.ReturnCodeNumber},

		{"0.0xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.12xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.345xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.4567xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.56789xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.678901xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.7890123xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.89012345xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.901234567xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.0123456789xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.12345678901xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.234567890123xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.3456789012345xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.45678901234567xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},

		{"0.0e2", "", jsonnum.ReturnCodeNumber},
		{"0.12e2", "", jsonnum.ReturnCodeNumber},
		{"0.345e2", "", jsonnum.ReturnCodeNumber},
		{"0.4567e2", "", jsonnum.ReturnCodeNumber},
		{"0.56789e2", "", jsonnum.ReturnCodeNumber},
		{"0.678901e2", "", jsonnum.ReturnCodeNumber},
		{"0.7890123e2", "", jsonnum.ReturnCodeNumber},
		{"0.89012345e2", "", jsonnum.ReturnCodeNumber},
		{"0.901234567e2", "", jsonnum.ReturnCodeNumber},
		{"0.0123456789e2", "", jsonnum.ReturnCodeNumber},
		{"0.12345678901e2", "", jsonnum.ReturnCodeNumber},
		{"0.234567890123e2", "", jsonnum.ReturnCodeNumber},
		{"0.3456789012345e2", "", jsonnum.ReturnCodeNumber},
		{"0.45678901234567e2", "", jsonnum.ReturnCodeNumber},

		{"0.0e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.12e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.345e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.4567e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.56789e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.678901e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.7890123e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.89012345e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.901234567e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.0123456789e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.12345678901e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.234567890123e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.3456789012345e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"0.45678901234567e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},

		{"1.000000000000", "", jsonnum.ReturnCodeNumber},
		{"12.000000000000", "", jsonnum.ReturnCodeNumber},
		{"123.000000000000", "", jsonnum.ReturnCodeNumber},
		{"1234.000000000000", "", jsonnum.ReturnCodeNumber},
		{"12345.000000000000", "", jsonnum.ReturnCodeNumber},
		{"123456.000000000000", "", jsonnum.ReturnCodeNumber},
		{"1234567.000000000000", "", jsonnum.ReturnCodeNumber},
		{"12345678.000000000000", "", jsonnum.ReturnCodeNumber},

		{"0e0", "", jsonnum.ReturnCodeNumber},
		{"0e12", "", jsonnum.ReturnCodeNumber},
		{"0e345", "", jsonnum.ReturnCodeNumber},
		{"0e4567", "", jsonnum.ReturnCodeNumber},
		{"0e56789", "", jsonnum.ReturnCodeNumber},
		{"0e678901", "", jsonnum.ReturnCodeNumber},
		{"0e7890123", "", jsonnum.ReturnCodeNumber},
		{"0e89012345", "", jsonnum.ReturnCodeNumber},
		{"0e9012345678", "", jsonnum.ReturnCodeNumber},
		{"0e12345678901", "", jsonnum.ReturnCodeNumber},

		{"1e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"12e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"123e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"1234e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"12345e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"123456e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"1234567e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"12345678e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"123456789e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"1234567890e2xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},

		{"1e1xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"1e12xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"1e123xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"1e1234xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"1e12345xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"1e123456xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"1e1234567xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"1e12345678xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"1e123456789xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},
		{"1e1234567890xxxxxxxx", "xxxxxxxx", jsonnum.ReturnCodeNumber},

		{"1e1x", "x", jsonnum.ReturnCodeNumber},

		{"0E1", "", jsonnum.ReturnCodeNumber},
		{"0E2", "", jsonnum.ReturnCodeNumber},
		{"0E3", "", jsonnum.ReturnCodeNumber},
		{"0E4", "", jsonnum.ReturnCodeNumber},
		{"0E5", "", jsonnum.ReturnCodeNumber},
		{"0E6", "", jsonnum.ReturnCodeNumber},
		{"0E7", "", jsonnum.ReturnCodeNumber},
		{"0E8", "", jsonnum.ReturnCodeNumber},
		{"0E9", "", jsonnum.ReturnCodeNumber},
		{"0E1234567890", "", jsonnum.ReturnCodeNumber},

		{"0e+1234567890", "", jsonnum.ReturnCodeNumber},
		{"0e-1234567890", "", jsonnum.ReturnCodeNumber},

		{"0E+1234567890", "", jsonnum.ReturnCodeNumber},
		{"0E-1234567890", "", jsonnum.ReturnCodeNumber},

		{"1E-1234567890", "", jsonnum.ReturnCodeNumber},
		{"1234567890E-1234567890", "", jsonnum.ReturnCodeNumber},

		{
			"1234567890.1234567890E-1234567890",
			"",
			jsonnum.ReturnCodeNumber,
		},

		{
			func() string { return strconv.FormatFloat(math.MaxFloat64, 'f', 368, 64) }(),
			"",
			jsonnum.ReturnCodeNumber,
		},
	} {
		t.Run(tt.input, func(t *testing.T) {
			trailing, rc := jsonnum.ReadNumber(tt.input)
			require.Equal(t, tt.expect, rc)
			require.Equal(t, tt.expectAfter, trailing)

			t.Run("negative", func(t *testing.T) {
				trailing, rc := jsonnum.ReadNumber("-" + tt.input)
				require.Equal(t, tt.expect, rc)
				require.Equal(t, tt.expectAfter, trailing)
			})

			t.Run("bytes", func(t *testing.T) {
				trailing, rc := jsonnum.ReadNumber([]byte(tt.input))
				require.Equal(t, tt.expect, rc)
				require.Equal(t, tt.expectAfter, string(trailing))

				t.Run("negative", func(t *testing.T) {
					trailing, rc := jsonnum.ReadNumber([]byte("-" + tt.input))
					require.Equal(t, tt.expect, rc)
					require.Equal(t, tt.expectAfter, string(trailing))
				})
			})
		})
	}
}

func TestReadNumberZero(t *testing.T) {
	const input = "01"

	trailing, rc := jsonnum.ReadNumber(input)
	require.Equal(t, jsonnum.ReturnCodeInteger, rc)
	require.Equal(t, "1", trailing)

	t.Run("negative", func(t *testing.T) {
		trailing, rc := jsonnum.ReadNumber("-" + input)
		require.Equal(t, jsonnum.ReturnCodeInteger, rc)
		require.Equal(t, "1", trailing)
	})

	t.Run("bytes", func(t *testing.T) {
		trailing, rc := jsonnum.ReadNumber([]byte(input))
		require.Equal(t, jsonnum.ReturnCodeInteger, rc)
		require.Equal(t, "1", string(trailing))

		t.Run("negative", func(t *testing.T) {
			trailing, rc := jsonnum.ReadNumber([]byte("-" + input))
			require.Equal(t, jsonnum.ReturnCodeInteger, rc)
			require.Equal(t, "1", string(trailing))
		})
	})
}

func TestReadNumberErr(t *testing.T) {
	for _, tt := range []struct {
		input       string
		expectAfter string
	}{
		{"a", "a"},
		{"-", ""},
		{"0.", ""},
		{"-0.", ""},
		{"1234567890.", ""},
		{"0e", ""},
		{"-0e", ""},
		{"1e", ""},
		{"e1", "e1"},
		{"1234567890e", ""},
		{"1e-", ""},
		{"0.E0", "E0"},
		{"0.e0", "e0"},
		{"0E.0", ".0"},
		{"42.E0", "E0"},
		{"42.e0", "e0"},
		{"42E.0", ".0"},
		{"0.1234567890e", ""},
		{"1234567890.1234567890e", ""},
	} {
		t.Run(tt.input, func(t *testing.T) {
			for _, t2 := range []struct {
				term string
				name string
			}{
				{"", "void"},
				{",", "comma"},
				{"}", "object_term"},
				{"]", "array_term"},
				{" ", "space"},
				{"\t", "horizontal_tab"},
				{"\r", "carriage_return"},
				{`{"x":null}`, "valid_object"},
			} {
				t.Run("terminated by "+t2.name, func(t *testing.T) {
					t.Run("string", func(t *testing.T) {
						in := tt.input + t2.term
						trailing, rc := jsonnum.ReadNumber(in)
						require.Equal(t, jsonnum.ReturnCodeErr, rc)
						require.Equal(t, tt.expectAfter+t2.term, trailing)
					})

					t.Run("bytes", func(t *testing.T) {
						in := tt.input + t2.term
						trailing, rc := jsonnum.ReadNumber([]byte(in))
						require.Equal(t, jsonnum.ReturnCodeErr, rc)
						require.Equal(t, tt.expectAfter+t2.term, string(trailing))
					})
				})
			}
		})
	}
}
