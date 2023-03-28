package jsonnum_test

import (
	"testing"

	"github.com/romshark/jscan/internal/jsonnum"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	// All numbers will be tested with valid terminators and prefixes
	for _, tt := range []string{
		"0",
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"42",
		"1234567890",

		"0.0",
		"0.0123456789",

		"0e1",
		"0e1234567890",

		"0E1",
		"0E1234567890",

		"0e+1234567890",
		"0e-1234567890",

		"0E+1234567890",
		"0E-1234567890",

		"1E-1234567890",
		"1234567890E-1234567890",

		"1234567890.1234567890E-1234567890",
	} {
		t.Run(tt, func(t *testing.T) {
			for _, t2 := range []struct {
				term string
				name string
			}{
				{"", "void"},
				{",", "comma"},
				{"}", "object term"},
				{"]", "array term"},
				{" ", "space"},
				{"\t", "horizontal tab"},
				{"\r", "carriage return"},
			} {
				t.Run("terminated by "+t2.name, func(t *testing.T) {
					end, err := jsonnum.Parse(tt + t2.term)
					require.False(t, err)
					require.Equal(t, len(tt), end)

					t.Run("negative", func(t *testing.T) {
						end, err := jsonnum.Parse("-" + tt + t2.term)
						require.False(t, err)
						require.Equal(t, len(tt)+1, end)
					})

					t.Run("bytes", func(t *testing.T) {
						end, err := jsonnum.ParseBytes([]byte(tt + t2.term))
						require.False(t, err)
						require.Equal(t, len(tt), end)

						t.Run("negative", func(t *testing.T) {
							end, err := jsonnum.ParseBytes(
								[]byte("-" + tt + t2.term),
							)
							require.False(t, err)
							require.Equal(t, len(tt)+1, end)
						})
					})
				})
			}
		})
	}
}

func TestParseZero(t *testing.T) {
	const input = "01"

	end, err := jsonnum.Parse(input)
	require.False(t, err)
	require.Equal(t, 1, end)

	t.Run("negative", func(t *testing.T) {
		end, err := jsonnum.Parse("-" + input)
		require.False(t, err)
		require.Equal(t, 2, end)
	})

	t.Run("bytes", func(t *testing.T) {
		end, err := jsonnum.ParseBytes([]byte(input))
		require.False(t, err)
		require.Equal(t, 1, end)

		t.Run("negative", func(t *testing.T) {
			end, err := jsonnum.ParseBytes(
				[]byte("-" + input),
			)
			require.False(t, err)
			require.Equal(t, 2, end)
		})
	})
}

func TestParseErr(t *testing.T) {
	for _, tt := range []string{
		"a",
		"-",
		"0.",
		"1234567890.",
		"0e",
		"1e",
		"e1",
		"1234567890e",
		"1e-",
		"-1x",
		"-1.1x",
		"1ex",
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
				{"}", "object term"},
				{"]", "array term"},
				{" ", "space"},
				{"\t", "horizontal tab"},
				{"\r", "carriage return"},
			} {
				t.Run("terminated by "+t2.name, func(t *testing.T) {
					t.Run("string", func(t *testing.T) {
						end, err := jsonnum.Parse(tt + t2.term)
						require.True(t, err)
						require.Zero(t, end)
					})

					t.Run("bytes", func(t *testing.T) {
						end, err := jsonnum.ParseBytes([]byte(tt + t2.term))
						require.True(t, err)
						require.Zero(t, end)
					})
				})
			}
		})
	}
}
