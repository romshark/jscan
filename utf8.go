package jscan

// characters below RuneSelf are represented as themselves in a single byte.
const utf8RuneSelf = 0x80

const (
	// The default lowest and highest continuation byte.
	utf8locb = 0b10000000
	utf8hicb = 0b10111111

	// These names of these constants are chosen to give nice alignment in the
	// table below. The first nibble is an index into acceptRanges or F for
	// special one-byte cases. The second nibble is the Rune length or the
	// Status for the special one-byte case.
	utf8xx = 0xF1 // invalid: size 1
	utf8as = 0xF0 // ASCII: size 1
	utf8s1 = 0x02 // accept 0, size 2
	utf8s2 = 0x13 // accept 1, size 3
	utf8s3 = 0x03 // accept 0, size 3
	utf8s4 = 0x23 // accept 2, size 3
	utf8s5 = 0x34 // accept 3, size 4
	utf8s6 = 0x04 // accept 0, size 4
	utf8s7 = 0x44 // accept 4, size 4
)

// utf8First is information about the utf8First byte in a UTF-8 sequence.
var utf8First = [256]uint8{
	//   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
	utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as,
	utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, // 0x00-0x0F
	utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as,
	utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, // 0x10-0x1F
	utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as,
	utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, // 0x20-0x2F
	utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as,
	utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, // 0x30-0x3F
	utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as,
	utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, // 0x40-0x4F
	utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as,
	utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, // 0x50-0x5F
	utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as,
	utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, // 0x60-0x6F
	utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as,
	utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, utf8as, // 0x70-0x7F
	//   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
	utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx,
	utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, // 0x80-0x8F
	utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx,
	utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, // 0x90-0x9F
	utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx,
	utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, // 0xA0-0xAF
	utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx,
	utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, // 0xB0-0xBF
	utf8xx, utf8xx, utf8s1, utf8s1, utf8s1, utf8s1, utf8s1, utf8s1,
	utf8s1, utf8s1, utf8s1, utf8s1, utf8s1, utf8s1, utf8s1, utf8s1, // 0xC0-0xCF
	utf8s1, utf8s1, utf8s1, utf8s1, utf8s1, utf8s1, utf8s1, utf8s1,
	utf8s1, utf8s1, utf8s1, utf8s1, utf8s1, utf8s1, utf8s1, utf8s1, // 0xD0-0xDF
	utf8s2, utf8s3, utf8s3, utf8s3, utf8s3, utf8s3, utf8s3, utf8s3,
	utf8s3, utf8s3, utf8s3, utf8s3, utf8s3, utf8s4, utf8s3, utf8s3, // 0xE0-0xEF
	utf8s5, utf8s6, utf8s6, utf8s6, utf8s7, utf8xx, utf8xx, utf8xx,
	utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, utf8xx, // 0xF0-0xFF
}

// utf8AcceptRange gives the range of valid values for the second byte in a UTF-8
// sequence.
type utf8AcceptRange struct {
	lo uint8 // lowest value for second byte.
	hi uint8 // highest value for second byte.
}

// utf8AcceptRanges has size 16 to avoid bounds checks in the code that uses it.
var utf8AcceptRanges = [16]utf8AcceptRange{
	0: {utf8locb, utf8hicb},
	1: {0xA0, utf8hicb},
	2: {utf8locb, 0x9F},
	3: {0x90, utf8hicb},
	4: {utf8locb, 0x8F},
}
