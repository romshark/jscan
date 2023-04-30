package jsonnum_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/romshark/jscan/internal/jsonnum"

	"github.com/stretchr/testify/require"
)

func TestEndIndex(t *testing.T) {
	for _, tt := range []struct {
		input          string
		expectIndexEnd int
	}{
		{"0", len("0")},
		{"12", len("12")},
		{"234", len("234")},
		{"3456", len("3456")},
		{"45678", len("45678")},
		{"567890", len("567890")},
		{"6789012", len("6789012")},
		{"78901234", len("78901234")},
		{"890123456", len("890123456")},
		{"9012345678", len("9012345678")},
		{"90123456789", len("90123456789")},
		{"901234567890", len("901234567890")},
		{"9012345678901", len("9012345678901")},
		{"90123456789012", len("90123456789012")},

		{"1x", len("1")},

		{"0xxxxxxxx", len("0")},
		{"12xxxxxxxx", len("12")},
		{"234xxxxxxxx", len("234")},
		{"3456xxxxxxxx", len("3456")},
		{"45678xxxxxxxx", len("45678")},
		{"567890xxxxxxxx", len("567890")},
		{"6789012xxxxxxxx", len("6789012")},
		{"78901234xxxxxxxx", len("78901234")},
		{"890123456xxxxxxxx", len("890123456")},
		{"9012345678xxxxxxxx", len("9012345678")},
		{"90123456789xxxxxxxx", len("90123456789")},
		{"901234567890xxxxxxxx", len("901234567890")},
		{"9012345678901xxxxxxxx", len("9012345678901")},
		{"90123456789012xxxxxxxx", len("90123456789012")},

		{"0.0x", len("0.0")},

		{"0.0", len("0.0")},
		{"0.12", len("0.12")},
		{"0.345", len("0.345")},
		{"0.4567", len("0.4567")},
		{"0.56789", len("0.56789")},
		{"0.678901", len("0.678901")},
		{"0.7890123", len("0.7890123")},
		{"0.89012345", len("0.89012345")},
		{"0.901234567", len("0.901234567")},
		{"0.0123456789", len("0.0123456789")},
		{"0.12345678901", len("0.12345678901")},
		{"0.234567890123", len("0.234567890123")},
		{"0.3456789012345", len("0.3456789012345")},
		{"0.45678901234567", len("0.45678901234567")},

		{"0.0xxxxxxxx", len("0.0")},
		{"0.12xxxxxxxx", len("0.12")},
		{"0.345xxxxxxxx", len("0.345")},
		{"0.4567xxxxxxxx", len("0.4567")},
		{"0.56789xxxxxxxx", len("0.56789")},
		{"0.678901xxxxxxxx", len("0.678901")},
		{"0.7890123xxxxxxxx", len("0.7890123")},
		{"0.89012345xxxxxxxx", len("0.89012345")},
		{"0.901234567xxxxxxxx", len("0.901234567")},
		{"0.0123456789xxxxxxxx", len("0.0123456789")},
		{"0.12345678901xxxxxxxx", len("0.12345678901")},
		{"0.234567890123xxxxxxxx", len("0.234567890123")},
		{"0.3456789012345xxxxxxxx", len("0.3456789012345")},
		{"0.45678901234567xxxxxxxx", len("0.45678901234567")},

		{"0.0e2", len("0.0e2")},
		{"0.12e2", len("0.12e2")},
		{"0.345e2", len("0.345e2")},
		{"0.4567e2", len("0.4567e2")},
		{"0.56789e2", len("0.56789e2")},
		{"0.678901e2", len("0.678901e2")},
		{"0.7890123e2", len("0.7890123e2")},
		{"0.89012345e2", len("0.89012345e2")},
		{"0.901234567e2", len("0.901234567e2")},
		{"0.0123456789e2", len("0.0123456789e2")},
		{"0.12345678901e2", len("0.12345678901e2")},
		{"0.234567890123e2", len("0.234567890123e2")},
		{"0.3456789012345e2", len("0.3456789012345e2")},
		{"0.45678901234567e2", len("0.45678901234567e2")},

		{"0.0e2xxxxxxxx", len("0.0e2")},
		{"0.12e2xxxxxxxx", len("0.12e2")},
		{"0.345e2xxxxxxxx", len("0.345e2")},
		{"0.4567e2xxxxxxxx", len("0.4567e2")},
		{"0.56789e2xxxxxxxx", len("0.56789e2")},
		{"0.678901e2xxxxxxxx", len("0.678901e2")},
		{"0.7890123e2xxxxxxxx", len("0.7890123e2")},
		{"0.89012345e2xxxxxxxx", len("0.89012345e2")},
		{"0.901234567e2xxxxxxxx", len("0.901234567e2")},
		{"0.0123456789e2xxxxxxxx", len("0.0123456789e2")},
		{"0.12345678901e2xxxxxxxx", len("0.12345678901e2")},
		{"0.234567890123e2xxxxxxxx", len("0.234567890123e2")},
		{"0.3456789012345e2xxxxxxxx", len("0.3456789012345e2")},
		{"0.45678901234567e2xxxxxxxx", len("0.45678901234567e2")},

		{"1.000000000000", len("1.000000000000")},
		{"12.000000000000", len("12.000000000000")},
		{"123.000000000000", len("123.000000000000")},
		{"1234.000000000000", len("1234.000000000000")},
		{"12345.000000000000", len("12345.000000000000")},
		{"123456.000000000000", len("123456.000000000000")},
		{"1234567.000000000000", len("1234567.000000000000")},
		{"12345678.000000000000", len("12345678.000000000000")},

		{"0e0", len("0e0")},
		{"0e12", len("0e12")},
		{"0e345", len("0e345")},
		{"0e4567", len("0e4567")},
		{"0e56789", len("0e56789")},
		{"0e678901", len("0e678901")},
		{"0e7890123", len("0e7890123")},
		{"0e89012345", len("0e89012345")},
		{"0e9012345678", len("0e9012345678")},
		{"0e12345678901", len("0e12345678901")},

		{"1e2xxxxxxxx", len("1e2")},
		{"12e2xxxxxxxx", len("12e2")},
		{"123e2xxxxxxxx", len("123e2")},
		{"1234e2xxxxxxxx", len("1234e2")},
		{"12345e2xxxxxxxx", len("12345e2")},
		{"123456e2xxxxxxxx", len("123456e2")},
		{"1234567e2xxxxxxxx", len("1234567e2")},
		{"12345678e2xxxxxxxx", len("12345678e2")},
		{"123456789e2xxxxxxxx", len("123456789e2")},
		{"1234567890e2xxxxxxxx", len("1234567890e2")},

		{"1e1xxxxxxxx", len("1e1")},
		{"1e12xxxxxxxx", len("1e12")},
		{"1e123xxxxxxxx", len("1e123")},
		{"1e1234xxxxxxxx", len("1e1234")},
		{"1e12345xxxxxxxx", len("1e12345")},
		{"1e123456xxxxxxxx", len("1e123456")},
		{"1e1234567xxxxxxxx", len("1e1234567")},
		{"1e12345678xxxxxxxx", len("1e12345678")},
		{"1e123456789xxxxxxxx", len("1e123456789")},
		{"1e1234567890xxxxxxxx", len("1e1234567890")},

		{"1e1x", len("1e1")},

		{"0E1", len("0E1")},
		{"0E2", len("0E2")},
		{"0E3", len("0E3")},
		{"0E4", len("0E4")},
		{"0E5", len("0E5")},
		{"0E6", len("0E6")},
		{"0E7", len("0E7")},
		{"0E8", len("0E8")},
		{"0E9", len("0E9")},
		{"0E1234567890", len("0E1234567890")},

		{"0e+1234567890", len("0e+1234567890")},
		{"0e-1234567890", len("0e-1234567890")},

		{"0E+1234567890", len("0E+1234567890")},
		{"0E-1234567890", len("0E-1234567890")},

		{"1E-1234567890", len("1E-1234567890")},
		{"1234567890E-1234567890", len("1234567890E-1234567890")},

		{
			"1234567890.1234567890E-1234567890",
			len("1234567890.1234567890E-1234567890"),
		},

		{
			func() string {
				return strconv.FormatFloat(math.MaxFloat64, 'f', 368, 64)
			}(),
			func() int {
				return len(strconv.FormatFloat(math.MaxFloat64, 'f', 368, 64))
			}(),
		},
	} {
		t.Run(tt.input, func(t *testing.T) {
			end, err := jsonnum.EndIndex(tt.input)
			require.False(t, err)
			require.Equal(t, tt.expectIndexEnd, end)

			t.Run("negative", func(t *testing.T) {
				end, err := jsonnum.EndIndex("-" + tt.input)
				require.False(t, err)
				require.Equal(t, tt.expectIndexEnd+1, end)
			})

			t.Run("bytes", func(t *testing.T) {
				end, err := jsonnum.EndIndex([]byte(tt.input))
				require.False(t, err)
				require.Equal(t, tt.expectIndexEnd, end)

				t.Run("negative", func(t *testing.T) {
					end, err := jsonnum.EndIndex(
						[]byte("-" + tt.input),
					)
					require.False(t, err)
					require.Equal(t, tt.expectIndexEnd+1, end)
				})
			})
		})
	}
}

func TestEndIndexZero(t *testing.T) {
	const input = "01"

	end, err := jsonnum.EndIndex(input)
	require.False(t, err)
	require.Equal(t, 1, end)

	t.Run("negative", func(t *testing.T) {
		end, err := jsonnum.EndIndex("-" + input)
		require.False(t, err)
		require.Equal(t, 2, end)
	})

	t.Run("bytes", func(t *testing.T) {
		end, err := jsonnum.EndIndex([]byte(input))
		require.False(t, err)
		require.Equal(t, 1, end)

		t.Run("negative", func(t *testing.T) {
			end, err := jsonnum.EndIndex(
				[]byte("-" + input),
			)
			require.False(t, err)
			require.Equal(t, 2, end)
		})
	})
}

func TestEndIndexErr(t *testing.T) {
	for _, tt := range []string{
		"a",
		"-",
		"0.",
		"-0.",
		"1234567890.",
		"0e",
		"-0e",
		"1e",
		"e1",
		"1234567890e",
		"1e-",
		"0.E0",
		"0.e0",
		"0E.0",
		"42.E0",
		"42.e0",
		"42E.0",
	} {
		t.Run(tt, func(t *testing.T) {
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
						end, err := jsonnum.EndIndex(tt + t2.term)
						require.True(t, err)
						require.Zero(t, end)
					})

					t.Run("bytes", func(t *testing.T) {
						end, err := jsonnum.EndIndex([]byte(tt + t2.term))
						require.True(t, err)
						require.Zero(t, end)
					})
				})
			}
		})
	}
}
