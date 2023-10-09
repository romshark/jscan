// Package utf8 is copied from the Go standard package "unicode/utf8"
// https://cs.opensource.google/go/go/+/refs/tags/go1.21.2:src/unicode/utf8/utf8.go;l=8
// See LICENCES.md for more information.
package utf8

// characters below RuneSelf are represented as themselves in a single byte.
const RuneSelf = 0x80

const (
	// The default lowest and highest continuation byte.
	Locb = 0b10000000
	Hicb = 0b10111111

	// These names of these constants are chosen to give nice alignment in the
	// table below. The first nibble is an index into acceptRanges or F for
	// special one-byte cases. The second nibble is the Rune length or the
	// Status for the special one-byte case.
	XX = 0xF1 // invalid: size 1
	AS = 0xF0 // ASCII: size 1
	S1 = 0x02 // accept 0, size 2
	S2 = 0x13 // accept 1, size 3
	S3 = 0x03 // accept 0, size 3
	S4 = 0x23 // accept 2, size 3
	S5 = 0x34 // accept 3, size 4
	S6 = 0x04 // accept 0, size 4
	S7 = 0x44 // accept 4, size 4
)

// First is information about the First byte in a UTF-8 sequence.
var First = [256]uint8{
	//   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
	AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, // 0x00-0x0F
	AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, // 0x10-0x1F
	AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, // 0x20-0x2F
	AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, // 0x30-0x3F
	AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, // 0x40-0x4F
	AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, // 0x50-0x5F
	AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, // 0x60-0x6F
	AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, AS, // 0x70-0x7F
	//   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
	XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, // 0x80-0x8F
	XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, // 0x90-0x9F
	XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, // 0xA0-0xAF
	XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, // 0xB0-0xBF
	XX, XX, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, // 0xC0-0xCF
	S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, S1, // 0xD0-0xDF
	S2, S3, S3, S3, S3, S3, S3, S3, S3, S3, S3, S3, S3, S4, S3, S3, // 0xE0-0xEF
	S5, S6, S6, S6, S7, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, XX, // 0xF0-0xFF
}

// AcceptRange gives the range of valid values for the second byte in a UTF-8
// sequence.
type AcceptRange struct {
	Lo uint8 // lowest value for second byte.
	Hi uint8 // highest value for second byte.
}

// AcceptRanges has size 16 to avoid bounds checks in the code that uses it.
var AcceptRanges = [16]AcceptRange{
	0: {Locb, Hicb},
	1: {0xA0, Hicb},
	2: {Locb, 0x9F},
	3: {0x90, Hicb},
	4: {Locb, 0x8F},
}
