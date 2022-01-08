package jsonnum_test

import (
	"testing"

	"github.com/romshark/jscan/internal/jsonnum"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	for _, tt := range []struct {
		in        string
		expectEnd int
	}{
		// All numbers will be tested with
		// valid terminators and prefixes
		{"0", 1},
		{"1", 1},
		{"2", 1},
		{"3", 1},
		{"4", 1},
		{"5", 1},
		{"6", 1},
		{"7", 1},
		{"8", 1},
		{"9", 1},
		{"42", 2},
		{"1234567890", 10},

		{"0.0", 3},
		{"0.0123456789", 12},

		{"0e1", 3},
		{"0e1234567890", 12},

		{"0E1", 3},
		{"0E1234567890", 12},

		{"0e+1234567890", 13},
		{"0e-1234567890", 13},

		{"0E+1234567890", 13},
		{"0E-1234567890", 13},

		{"1E-1234567890", 13},
		{"1234567890E-1234567890", 22},

		{"1234567890.1234567890E-1234567890", 33},
	} {
		t.Run(tt.in, func(t *testing.T) {
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
					end, err := jsonnum.Parse(tt.in + t2.term)
					require.False(t, err)
					require.Equal(t, tt.expectEnd, end)

					t.Run("negative", func(t *testing.T) {
						end, err := jsonnum.Parse("-" + tt.in + t2.term)
						require.False(t, err)
						require.Equal(t, tt.expectEnd+1, end)
					})
				})
			}
		})
	}
}

func TestParseErr(t *testing.T) {
	for _, tt := range []string{
		"a",
		"-",
		"0.",
		"1234567890.",
		"01",
		"0e",
		"1e",
		"e1",
		"1234567890e",
		"1e-",
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
					require.True(t, err)
					require.Zero(t, end)
				})
			}

		})
	}
}
