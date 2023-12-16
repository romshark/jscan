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
	}{
		{"0", ""},
		{"12", ""},
		{"234", ""},
		{"3456", ""},
		{"45678", ""},
		{"567890", ""},
		{"6789012", ""},
		{"78901234", ""},
		{"890123456", ""},
		{"9012345678", ""},
		{"90123456789", ""},
		{"901234567890", ""},
		{"9012345678901", ""},
		{"90123456789012", ""},

		{"1x", "x"},

		{"0xxxxxxxx", "xxxxxxxx"},
		{"12xxxxxxxx", "xxxxxxxx"},
		{"234xxxxxxxx", "xxxxxxxx"},
		{"3456xxxxxxxx", "xxxxxxxx"},
		{"45678xxxxxxxx", "xxxxxxxx"},
		{"567890xxxxxxxx", "xxxxxxxx"},
		{"6789012xxxxxxxx", "xxxxxxxx"},
		{"78901234xxxxxxxx", "xxxxxxxx"},
		{"890123456xxxxxxxx", "xxxxxxxx"},
		{"9012345678xxxxxxxx", "xxxxxxxx"},
		{"90123456789xxxxxxxx", "xxxxxxxx"},
		{"901234567890xxxxxxxx", "xxxxxxxx"},
		{"9012345678901xxxxxxxx", "xxxxxxxx"},
		{"90123456789012xxxxxxxx", "xxxxxxxx"},

		{"0.0x", "x"},

		{"0.0", ""},
		{"0.12", ""},
		{"0.345", ""},
		{"0.4567", ""},
		{"0.56789", ""},
		{"0.678901", ""},
		{"0.7890123", ""},
		{"0.89012345", ""},
		{"0.901234567", ""},
		{"0.0123456789", ""},
		{"0.12345678901", ""},
		{"0.234567890123", ""},
		{"0.3456789012345", ""},
		{"0.45678901234567", ""},

		{"0.0xxxxxxxx", "xxxxxxxx"},
		{"0.12xxxxxxxx", "xxxxxxxx"},
		{"0.345xxxxxxxx", "xxxxxxxx"},
		{"0.4567xxxxxxxx", "xxxxxxxx"},
		{"0.56789xxxxxxxx", "xxxxxxxx"},
		{"0.678901xxxxxxxx", "xxxxxxxx"},
		{"0.7890123xxxxxxxx", "xxxxxxxx"},
		{"0.89012345xxxxxxxx", "xxxxxxxx"},
		{"0.901234567xxxxxxxx", "xxxxxxxx"},
		{"0.0123456789xxxxxxxx", "xxxxxxxx"},
		{"0.12345678901xxxxxxxx", "xxxxxxxx"},
		{"0.234567890123xxxxxxxx", "xxxxxxxx"},
		{"0.3456789012345xxxxxxxx", "xxxxxxxx"},
		{"0.45678901234567xxxxxxxx", "xxxxxxxx"},

		{"0.0e2", ""},
		{"0.12e2", ""},
		{"0.345e2", ""},
		{"0.4567e2", ""},
		{"0.56789e2", ""},
		{"0.678901e2", ""},
		{"0.7890123e2", ""},
		{"0.89012345e2", ""},
		{"0.901234567e2", ""},
		{"0.0123456789e2", ""},
		{"0.12345678901e2", ""},
		{"0.234567890123e2", ""},
		{"0.3456789012345e2", ""},
		{"0.45678901234567e2", ""},

		{"0.0e2xxxxxxxx", "xxxxxxxx"},
		{"0.12e2xxxxxxxx", "xxxxxxxx"},
		{"0.345e2xxxxxxxx", "xxxxxxxx"},
		{"0.4567e2xxxxxxxx", "xxxxxxxx"},
		{"0.56789e2xxxxxxxx", "xxxxxxxx"},
		{"0.678901e2xxxxxxxx", "xxxxxxxx"},
		{"0.7890123e2xxxxxxxx", "xxxxxxxx"},
		{"0.89012345e2xxxxxxxx", "xxxxxxxx"},
		{"0.901234567e2xxxxxxxx", "xxxxxxxx"},
		{"0.0123456789e2xxxxxxxx", "xxxxxxxx"},
		{"0.12345678901e2xxxxxxxx", "xxxxxxxx"},
		{"0.234567890123e2xxxxxxxx", "xxxxxxxx"},
		{"0.3456789012345e2xxxxxxxx", "xxxxxxxx"},
		{"0.45678901234567e2xxxxxxxx", "xxxxxxxx"},

		{"1.000000000000", ""},
		{"12.000000000000", ""},
		{"123.000000000000", ""},
		{"1234.000000000000", ""},
		{"12345.000000000000", ""},
		{"123456.000000000000", ""},
		{"1234567.000000000000", ""},
		{"12345678.000000000000", ""},

		{"0e0", ""},
		{"0e12", ""},
		{"0e345", ""},
		{"0e4567", ""},
		{"0e56789", ""},
		{"0e678901", ""},
		{"0e7890123", ""},
		{"0e89012345", ""},
		{"0e9012345678", ""},
		{"0e12345678901", ""},

		{"1e2xxxxxxxx", "xxxxxxxx"},
		{"12e2xxxxxxxx", "xxxxxxxx"},
		{"123e2xxxxxxxx", "xxxxxxxx"},
		{"1234e2xxxxxxxx", "xxxxxxxx"},
		{"12345e2xxxxxxxx", "xxxxxxxx"},
		{"123456e2xxxxxxxx", "xxxxxxxx"},
		{"1234567e2xxxxxxxx", "xxxxxxxx"},
		{"12345678e2xxxxxxxx", "xxxxxxxx"},
		{"123456789e2xxxxxxxx", "xxxxxxxx"},
		{"1234567890e2xxxxxxxx", "xxxxxxxx"},

		{"1e1xxxxxxxx", "xxxxxxxx"},
		{"1e12xxxxxxxx", "xxxxxxxx"},
		{"1e123xxxxxxxx", "xxxxxxxx"},
		{"1e1234xxxxxxxx", "xxxxxxxx"},
		{"1e12345xxxxxxxx", "xxxxxxxx"},
		{"1e123456xxxxxxxx", "xxxxxxxx"},
		{"1e1234567xxxxxxxx", "xxxxxxxx"},
		{"1e12345678xxxxxxxx", "xxxxxxxx"},
		{"1e123456789xxxxxxxx", "xxxxxxxx"},
		{"1e1234567890xxxxxxxx", "xxxxxxxx"},

		{"1e1x", "x"},

		{"0E1", ""},
		{"0E2", ""},
		{"0E3", ""},
		{"0E4", ""},
		{"0E5", ""},
		{"0E6", ""},
		{"0E7", ""},
		{"0E8", ""},
		{"0E9", ""},
		{"0E1234567890", ""},

		{"0e+1234567890", ""},
		{"0e-1234567890", ""},

		{"0E+1234567890", ""},
		{"0E-1234567890", ""},

		{"1E-1234567890", ""},
		{"1234567890E-1234567890", ""},

		{
			"1234567890.1234567890E-1234567890",
			"",
		},

		{
			func() string {
				return strconv.FormatFloat(math.MaxFloat64, 'f', 368, 64)
			}(),
			"",
		},
	} {
		t.Run(tt.input, func(t *testing.T) {
			trailing, err := jsonnum.ReadNumber(tt.input)
			require.False(t, err)
			require.Equal(t, tt.expectAfter, trailing)

			t.Run("negative", func(t *testing.T) {
				trailing, err := jsonnum.ReadNumber("-" + tt.input)
				require.False(t, err)
				require.Equal(t, tt.expectAfter, trailing)
			})

			t.Run("bytes", func(t *testing.T) {
				trailing, err := jsonnum.ReadNumber([]byte(tt.input))
				require.False(t, err)
				require.Equal(t, tt.expectAfter, string(trailing))

				t.Run("negative", func(t *testing.T) {
					trailing, err := jsonnum.ReadNumber(
						[]byte("-" + tt.input),
					)
					require.False(t, err)
					require.Equal(t, tt.expectAfter, string(trailing))
				})
			})
		})
	}
}

func TestReadNumberZero(t *testing.T) {
	const input = "01"

	trailing, err := jsonnum.ReadNumber(input)
	require.False(t, err)
	require.Equal(t, "1", trailing)

	t.Run("negative", func(t *testing.T) {
		trailing, err := jsonnum.ReadNumber("-" + input)
		require.False(t, err)
		require.Equal(t, "1", trailing)
	})

	t.Run("bytes", func(t *testing.T) {
		trailing, err := jsonnum.ReadNumber([]byte(input))
		require.False(t, err)
		require.Equal(t, "1", string(trailing))

		t.Run("negative", func(t *testing.T) {
			trailing, err := jsonnum.ReadNumber([]byte("-" + input))
			require.False(t, err)
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
						trailing, err := jsonnum.ReadNumber(in)
						require.True(t, err)
						require.Equal(t, tt.expectAfter+t2.term, trailing)
					})

					t.Run("bytes", func(t *testing.T) {
						in := tt.input + t2.term
						trailing, err := jsonnum.ReadNumber([]byte(in))
						require.True(t, err)
						require.Equal(t, tt.expectAfter+t2.term, string(trailing))
					})
				})
			}
		})
	}
}
