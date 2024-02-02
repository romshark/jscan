package atoi

// U8 assumes that s is a valid unsigned 8-bit integer.
// Returns (0, true) if the value would overflow.
func U8[S ~string | ~[]byte](s S) (n uint8, overflow bool) {
	d := func(index int) uint8 { return uint8(s[index] - '0') }
	switch len(s) {
	case 1:
		return d(0), false
	case 2:
		return d(0)*1e1 + d(1), false
	case 3: // This case can overflow
		if s[0]-'0' > 2 {
			return 0, true
		}
		p := d(1)*1e1 + d(2)
		n = p + d(0)*1e2
		if n <= p {
			return 0, true
		}
		return n, false
	}
	return 0, true // Anything above 3 digits overflows uint8
}

// I8 assumes that s is a valid 8-bit signed or unsigned integer (where
// only '-' is accepted). Returns (0, true) if the value would overflow.
func I8[S ~string | ~[]byte](s S) (n int8, overflow bool) {
	d := func(index int) int8 { return int8(s[index] - '0') }
	if s[0] == '-' {
		switch len(s) {
		case 2:
			return -d(1), false
		case 3:
			return -(d(1)*1e1 + d(2)), false
		case 4: // This case can overflow
			if n = -(d(1)*1e2 + d(2)*1e1 + d(3)); n > 0 {
				return 0, true
			}
			return n, false
		}
	} else {
		switch len(s) {
		case 1:
			return d(0), false
		case 2:
			return d(0)*1e1 + d(1), false
		case 3: // This case can overflow
			if n = d(0)*1e2 + d(1)*1e1 + d(2); n < 0 {
				return 0, true
			}
			return n, false
		}
	}
	return 0, true // Anything above 3 digits overflows int8
}

// U16 assumes that s is a valid unsigned 16-bit integer.
// Returns (0, true) if the value would overflow.
func U16[S ~string | ~[]byte](s S) (n uint16, overflow bool) {
	d := func(index int) uint16 { return uint16(s[index] - '0') }
	switch len(s) {
	case 1:
		return d(0), false
	case 2:
		return d(0)*1e1 + d(1), false
	case 3:
		return d(0)*1e2 + d(1)*1e1 + d(2), false
	case 4:
		return d(0)*1e3 + d(1)*1e2 + d(2)*1e1 + d(3), false
	case 5: // This case can overflow
		if s[0]-'0' > 6 {
			return 0, true
		}
		p := d(1)*1e3 + d(2)*1e2 + d(3)*1e1 + d(4)
		n = p + d(0)*1e4
		if n <= p {
			return 0, true
		}
		return n, false
	}
	return 0, true // Anything above 5 digits overflows uint16
}

// I16 assumes that s is a valid 16-bit signed or unsigned integer (where
// only '-' is accepted). Returns (0, true) if the value would overflow.
func I16[S ~string | ~[]byte](s S) (n int16, overflow bool) {
	d := func(index int) int16 { return int16(s[index] - '0') }
	if s[0] == '-' {
		switch len(s) {
		case 2:
			return -d(1), false
		case 3:
			return -(n*1e2 + d(1)*1e1 + d(2)), false
		case 4:
			return -(d(1)*1e2 + d(2)*1e1 + d(3)), false
		case 5:
			return -(d(1)*1e3 + d(2)*1e2 + d(3)*1e1 + d(4)), false
		case 6: // This case can overflow
			n = -(d(1)*1e4 + d(2)*1e3 + d(3)*1e2 + d(4)*1e1 + d(5))
			if n > 0 {
				return 0, true
			}
			return n, false
		}
	} else {
		switch len(s) {
		case 1:
			return d(0), false
		case 2:
			return d(0)*1e1 + d(1), false
		case 3:
			return d(0)*1e2 + d(1)*1e1 + d(2), false
		case 4:
			return d(0)*1e3 + d(1)*1e2 + d(2)*1e1 + d(3), false
		case 5: // This case can overflow
			n = d(0)*1e4 + d(1)*1e3 + d(2)*1e2 + d(3)*1e1 + d(4)
			if n < 0 {
				return 0, true
			}
			return n, false
		}
	}
	return 0, true // Anything above 5 digits overflows int16
}

// U32 assumes that s is a valid unsigned 32-bit integer.
// Returns (0, true) if the value would overflow.
func U32[S ~string | ~[]byte](s S) (n uint32, overflow bool) {
	d := func(index int) uint32 { return uint32(s[index] - '0') }
	switch len(s) {
	case 1:
		return d(0), false
	case 2:
		return d(0)*1e1 + d(1), false
	case 3:
		return d(0)*1e2 + d(1)*1e1 + d(2), false
	case 4:
		return d(0)*1e3 + d(1)*1e2 + d(2)*1e1 + d(3), false
	case 5:
		return d(0)*1e4 + d(1)*1e3 + d(2)*1e2 + d(3)*1e1 + d(4), false
	case 6:
		return d(0)*1e5 + d(1)*1e4 +
			d(2)*1e3 + d(3)*1e2 +
			d(4)*1e1 + d(5), false
	case 7:
		return d(0)*1e6 + d(1)*1e5 +
			d(2)*1e4 + d(3)*1e3 +
			d(4)*1e2 + d(5)*1e1 +
			d(6), false
	case 8:
		return d(0)*1e7 + d(1)*1e6 +
			d(2)*1e5 + d(3)*1e4 +
			d(4)*1e3 + d(5)*1e2 +
			d(6)*1e1 + d(7), false
	case 9:
		return d(0)*1e8 + d(1)*1e7 +
			d(2)*1e6 + d(3)*1e5 +
			d(4)*1e4 + d(5)*1e3 +
			d(6)*1e2 + d(7)*1e1 +
			d(8), false
	case 10: // This case can overflow
		if s[0]-'0' > 4 {
			return 0, true
		}
		p := d(1)*1e8 + d(2)*1e7 +
			d(3)*1e6 + d(4)*1e5 +
			d(5)*1e4 + d(6)*1e3 +
			d(7)*1e2 + d(8)*1e1 +
			d(9)
		n = p + d(0)*1e9
		if n <= p {
			return 0, true
		}
		return n, false
	}
	return 0, true // Anything above 10 digits overflows uint32
}

// I32 assumes that s is a valid 32-bit signed or unsigned integer (where
// only '-' is accepted). Returns (0, true) if the value would overflow.
func I32[S ~string | ~[]byte](s S) (n int32, overflow bool) {
	d := func(index int) int32 { return int32(s[index] - '0') }
	if s[0] == '-' {
		switch len(s) {
		case 2:
			return -d(1), false
		case 3:
			return -(n*1e2 + d(1)*1e1 + d(2)), false
		case 4:
			return -(d(1)*1e2 + d(2)*1e1 + d(3)), false
		case 5:
			return -(d(1)*1e3 + d(2)*1e2 + d(3)*1e1 + d(4)), false
		case 6:
			return -(d(1)*1e4 + d(2)*1e3 + d(3)*1e2 + d(4)*1e1 + d(5)), false
		case 7:
			return -(d(1)*1e5 + d(2)*1e4 +
				d(3)*1e3 + d(4)*1e2 +
				d(5)*1e1 + d(6)), false
		case 8:
			return -(d(1)*1e6 + d(2)*1e5 +
				d(3)*1e4 + d(4)*1e3 +
				d(5)*1e2 + d(6)*1e1 +
				d(7)), false
		case 9:
			return -(d(1)*1e7 + d(2)*1e6 +
				d(3)*1e5 + d(4)*1e4 +
				d(5)*1e3 + d(6)*1e2 +
				d(7)*1e1 + d(8)), false
		case 10:
			return -(d(1)*1e8 + d(2)*1e7 +
				d(3)*1e6 + d(4)*1e5 +
				d(5)*1e4 + d(6)*1e3 +
				d(7)*1e2 + d(8)*1e1 +
				d(9)), false
		case 11: // This case can overflow
			n = d(1)*1e9 + d(2)*1e8 +
				d(3)*1e7 + d(4)*1e6 +
				d(5)*1e5 + d(6)*1e4 +
				d(7)*1e3 + d(8)*1e2 +
				d(9)*1e1 + d(10)
			n = -n
			if n > 0 {
				return 0, true
			}
			return n, false
		}
	} else {
		switch len(s) {
		case 1:
			return d(0), false
		case 2:
			return d(0)*1e1 + d(1), false
		case 3:
			return d(0)*1e2 + d(1)*1e1 + d(2), false
		case 4:
			return d(0)*1e3 + d(1)*1e2 + d(2)*1e1 + d(3), false
		case 5:
			return d(0)*1e4 + d(1)*1e3 + d(2)*1e2 + d(3)*1e1 + d(4), false
		case 6:
			return d(0)*1e5 + d(1)*1e4 +
				d(2)*1e3 + d(3)*1e2 +
				d(4)*1e1 + d(5), false
		case 7:
			return d(0)*1e6 + d(1)*1e5 +
				d(2)*1e4 + d(3)*1e3 +
				d(4)*1e2 + d(5)*1e1 +
				d(6), false
		case 8:
			return d(0)*1e7 + d(1)*1e6 +
				d(2)*1e5 + d(3)*1e4 +
				d(4)*1e3 + d(5)*1e2 +
				d(6)*1e1 + d(7), false
		case 9:
			return d(0)*1e8 + d(1)*1e7 +
				d(2)*1e6 + d(3)*1e5 +
				d(4)*1e4 + d(5)*1e3 +
				d(6)*1e2 + d(7)*1e1 +
				d(8), false
		case 10: // This case can overflow
			n = d(0)*1e9 + d(1)*1e8 +
				d(2)*1e7 + d(3)*1e6 +
				d(4)*1e5 + d(5)*1e4 +
				d(6)*1e3 + d(7)*1e2 +
				d(8)*1e1 + d(9)
			if n < 0 {
				return 0, true
			}
			return n, false
		}
	}
	return 0, true // Anything above 10 digits overflows int32
}

// U64 assumes that s is a valid unsigned 64-bit integer.
// Returns (0, true) if the value would overflow.
func U64[S ~string | ~[]byte](s S) (n uint64, overflow bool) {
	d := func(index int) uint64 { return uint64(s[index] - '0') }
	switch len(s) {
	case 1:
		return d(0), false
	case 2:
		return d(0)*1e1 + d(1), false
	case 3:
		return d(0)*1e2 + d(1)*1e1 +
			d(2), false
	case 4:
		return d(0)*1e3 + d(1)*1e2 +
			d(2)*1e1 + d(3), false
	case 5:
		return d(0)*1e4 + d(1)*1e3 +
			d(2)*1e2 + d(3)*1e1 +
			d(4), false
	case 6:
		return d(0)*1e5 + d(1)*1e4 +
			d(2)*1e3 + d(3)*1e2 +
			d(4)*1e1 + d(5), false
	case 7:
		return d(0)*1e6 + d(1)*1e5 +
			d(2)*1e4 + d(3)*1e3 +
			d(4)*1e2 + d(5)*1e1 +
			d(6), false
	case 8:
		return d(0)*1e7 + d(1)*1e6 +
			d(2)*1e5 + d(3)*1e4 +
			d(4)*1e3 + d(5)*1e2 +
			d(6)*1e1 + d(7), false
	case 9:
		return d(0)*1e8 + d(1)*1e7 +
			d(2)*1e6 + d(3)*1e5 +
			d(4)*1e4 + d(5)*1e3 +
			d(6)*1e2 + d(7)*1e1 +
			d(8), false
	case 10:
		return d(0)*1e9 + d(1)*1e8 +
			d(2)*1e7 + d(3)*1e6 +
			d(4)*1e5 + d(5)*1e4 +
			d(6)*1e3 + d(7)*1e2 +
			d(8)*1e1 + d(9), false
	case 11:
		return d(0)*1e10 + d(1)*1e9 +
			d(2)*1e8 + d(3)*1e7 +
			d(4)*1e6 + d(5)*1e5 +
			d(6)*1e4 + d(7)*1e3 +
			d(8)*1e2 + d(9)*1e1 +
			d(10), false
	case 12:
		return d(0)*1e11 + d(1)*1e10 +
			d(2)*1e9 + d(3)*1e8 +
			d(4)*1e7 + d(5)*1e6 +
			d(6)*1e5 + d(7)*1e4 +
			d(8)*1e3 + d(9)*1e2 +
			d(10)*1e1 + d(11), false
	case 13:
		return d(0)*1e12 + d(1)*1e11 +
			d(2)*1e10 + d(3)*1e9 +
			d(4)*1e8 + d(5)*1e7 +
			d(6)*1e6 + d(7)*1e5 +
			d(8)*1e4 + d(9)*1e3 +
			d(10)*1e2 + d(11)*1e1 +
			d(12), false
	case 14:
		return d(0)*1e13 + d(1)*1e12 +
			d(2)*1e11 + d(3)*1e10 +
			d(4)*1e9 + d(5)*1e8 +
			d(6)*1e7 + d(7)*1e6 +
			d(8)*1e5 + d(9)*1e4 +
			d(10)*1e3 + d(11)*1e2 +
			d(12)*1e1 + d(13), false
	case 15:
		return d(0)*1e14 + d(1)*1e13 +
			d(2)*1e12 + d(3)*1e11 +
			d(4)*1e10 + d(5)*1e9 +
			d(6)*1e8 + d(7)*1e7 +
			d(8)*1e6 + d(9)*1e5 +
			d(10)*1e4 + d(11)*1e3 +
			d(12)*1e2 + d(13)*1e1 +
			d(14), false
	case 16:
		return d(0)*1e15 + d(1)*1e14 +
			d(2)*1e13 + d(3)*1e12 +
			d(4)*1e11 + d(5)*1e10 +
			d(6)*1e9 + d(7)*1e8 +
			d(8)*1e7 + d(9)*1e6 +
			d(10)*1e5 + d(11)*1e4 +
			d(12)*1e3 + d(13)*1e2 +
			d(14)*1e1 + d(15), false
	case 17:
		return d(0)*1e16 + d(1)*1e15 +
			d(2)*1e14 + d(3)*1e13 +
			d(4)*1e12 + d(5)*1e11 +
			d(6)*1e10 + d(7)*1e9 +
			d(8)*1e8 + d(9)*1e7 +
			d(10)*1e6 + d(11)*1e5 +
			d(12)*1e4 + d(13)*1e3 +
			d(14)*1e2 + d(15)*1e1 +
			d(16), false
	case 18:
		return d(0)*1e17 + d(1)*1e16 +
			d(2)*1e15 + d(3)*1e14 +
			d(4)*1e13 + d(5)*1e12 +
			d(6)*1e11 + d(7)*1e10 +
			d(8)*1e9 + d(9)*1e8 +
			d(10)*1e7 + d(11)*1e6 +
			d(12)*1e5 + d(13)*1e4 +
			d(14)*1e3 + d(15)*1e2 +
			d(16)*1e1 + d(17), false
	case 19:
		return d(0)*1e18 + d(1)*1e17 +
			d(2)*1e16 + d(3)*1e15 +
			d(4)*1e14 + d(5)*1e13 +
			d(6)*1e12 + d(7)*1e11 +
			d(8)*1e10 + d(9)*1e9 +
			d(10)*1e8 + d(11)*1e7 +
			d(12)*1e6 + d(13)*1e5 +
			d(14)*1e4 + d(15)*1e3 +
			d(16)*1e2 + d(17)*1e1 +
			d(18), false
	case 20: // This case can overflow
		if s[0] != '1' {
			return 0, true
		}
		n = d(1)*1e18 + d(2)*1e17 +
			d(3)*1e16 + d(4)*1e15 +
			d(5)*1e14 + d(6)*1e13 +
			d(7)*1e12 + d(8)*1e11 +
			d(9)*1e10 + d(10)*1e9 +
			d(11)*1e8 + d(12)*1e7 +
			d(13)*1e6 + d(14)*1e5 +
			d(15)*1e4 + d(16)*1e3 +
			d(17)*1e2 + d(18)*1e1 +
			d(19)
		if n > 8446744073709551615 {
			return 0, true
		}
		return 1e19 + n, false
	}
	return 0, true // Anything above 20 digits overflows uint64
}

// I64 assumes that s is a valid 64-bit signed or unsigned integer (where
// only '-' is accepted). Returns (0, true) if the value would overflow.
func I64[S ~string | ~[]byte](s S) (n int64, overflow bool) {
	d := func(index int) int64 { return int64(s[index] - '0') }
	if s[0] == '-' {
		switch len(s) {
		case 2:
			return -d(1), false
		case 3:
			return -(n*1e2 + d(1)*1e1 + d(2)), false
		case 4:
			return -(d(1)*1e2 + d(2)*1e1 +
				d(3)), false
		case 5:
			return -(d(1)*1e3 + d(2)*1e2 +
				d(3)*1e1 + d(4)), false
		case 6:
			return -(d(1)*1e4 + d(2)*1e3 +
				d(3)*1e2 + d(4)*1e1 +
				d(5)), false
		case 7:
			return -(d(1)*1e5 + d(2)*1e4 +
				d(3)*1e3 + d(4)*1e2 +
				d(5)*1e1 + d(6)), false
		case 8:
			return -(d(1)*1e6 + d(2)*1e5 +
				d(3)*1e4 + d(4)*1e3 +
				d(5)*1e2 + d(6)*1e1 +
				d(7)), false
		case 9:
			return -(d(1)*1e7 + d(2)*1e6 +
				d(3)*1e5 + d(4)*1e4 +
				d(5)*1e3 + d(6)*1e2 +
				d(7)*1e1 + d(8)), false
		case 10:
			return -(d(1)*1e8 + d(2)*1e7 +
				d(3)*1e6 + d(4)*1e5 +
				d(5)*1e4 + d(6)*1e3 +
				d(7)*1e2 + d(8)*1e1 +
				d(9)), false
		case 11:
			return -(d(1)*1e9 + d(2)*1e8 +
				d(3)*1e7 + d(4)*1e6 +
				d(5)*1e5 + d(6)*1e4 +
				d(7)*1e3 + d(8)*1e2 +
				d(9)*1e1 + d(10)), false
		case 12:
			return -(d(1)*1e10 + d(2)*1e9 +
				d(3)*1e8 + d(4)*1e7 +
				d(5)*1e6 + d(6)*1e5 +
				d(7)*1e4 + d(8)*1e3 +
				d(9)*1e2 + d(10)*1e1 +
				d(11)), false
		case 13:
			return -(d(1)*1e11 + d(2)*1e10 +
				d(3)*1e9 + d(4)*1e8 +
				d(5)*1e7 + d(6)*1e6 +
				d(7)*1e5 + d(8)*1e4 +
				d(9)*1e3 + d(10)*1e2 +
				d(11)*1e1 + d(12)), false
		case 14:
			return -(d(1)*1e12 + d(2)*1e11 +
				d(3)*1e10 + d(4)*1e9 +
				d(5)*1e8 + d(6)*1e7 +
				d(7)*1e6 + d(8)*1e5 +
				d(9)*1e4 + d(10)*1e3 +
				d(11)*1e2 + d(12)*1e1 +
				d(13)), false
		case 15:
			return -(d(1)*1e13 + d(2)*1e12 +
				d(3)*1e11 + d(4)*1e10 +
				d(5)*1e9 + d(6)*1e8 +
				d(7)*1e7 + d(8)*1e6 +
				d(9)*1e5 + d(10)*1e4 +
				d(11)*1e3 + d(12)*1e2 +
				d(13)*1e1 + d(14)), false
		case 16:
			return -(d(1)*1e14 + d(2)*1e13 +
				d(3)*1e12 + d(4)*1e11 +
				d(5)*1e10 + d(6)*1e9 +
				d(7)*1e8 + d(8)*1e7 +
				d(9)*1e6 + d(10)*1e5 +
				d(11)*1e4 + d(12)*1e3 +
				d(13)*1e2 + d(14)*1e1 +
				d(15)), false
		case 17:
			return -(d(1)*1e15 + d(2)*1e14 +
				d(3)*1e13 + d(4)*1e12 +
				d(5)*1e11 + d(6)*1e10 +
				d(7)*1e9 + d(8)*1e8 +
				d(9)*1e7 + d(10)*1e6 +
				d(11)*1e5 + d(12)*1e4 +
				d(13)*1e3 + d(14)*1e2 +
				d(15)*1e1 + d(16)), false
		case 18:
			return -(d(1)*1e16 + d(2)*1e15 +
				d(3)*1e14 + d(4)*1e13 +
				d(5)*1e12 + d(6)*1e11 +
				d(7)*1e10 + d(8)*1e9 +
				d(9)*1e8 + d(10)*1e7 +
				d(11)*1e6 + d(12)*1e5 +
				d(13)*1e4 + d(14)*1e3 +
				d(15)*1e2 + d(16)*1e1 +
				d(17)), false
		case 19:
			return -(d(1)*1e17 + d(2)*1e16 +
				d(3)*1e15 + d(4)*1e14 +
				d(5)*1e13 + d(6)*1e12 +
				d(7)*1e11 + d(8)*1e10 +
				d(9)*1e9 + d(10)*1e8 +
				d(11)*1e7 + d(12)*1e6 +
				d(13)*1e5 + d(14)*1e4 +
				d(15)*1e3 + d(16)*1e2 +
				d(17)*1e1 + d(18)), false
		case 20: // This case can overflow
			n = d(1)*1e18 + d(2)*1e17 + d(3)*1e16 + d(4)*1e15 +
				d(5)*1e14 + d(6)*1e13 + d(7)*1e12 + d(8)*1e11 +
				d(9)*1e10 + d(10)*1e9 + d(11)*1e8 + d(12)*1e7 +
				d(13)*1e6 + d(14)*1e5 + d(15)*1e4 + d(16)*1e3 +
				d(17)*1e2 + d(18)*1e1 + d(19)
			n = -n
			if n > 0 {
				return 0, true
			}
			return n, false
		}
	} else {
		switch len(s) {
		case 1:
			return d(0), false
		case 2:
			return d(0)*1e1 + d(1), false
		case 3:
			return d(0)*1e2 + d(1)*1e1 +
				d(2), false
		case 4:
			return d(0)*1e3 + d(1)*1e2 +
				d(2)*1e1 + d(3), false
		case 5:
			return d(0)*1e4 + d(1)*1e3 +
				d(2)*1e2 + d(3)*1e1 +
				d(4), false
		case 6:
			return d(0)*1e5 + d(1)*1e4 +
				d(2)*1e3 + d(3)*1e2 +
				d(4)*1e1 + d(5), false
		case 7:
			return d(0)*1e6 + d(1)*1e5 +
				d(2)*1e4 + d(3)*1e3 +
				d(4)*1e2 + d(5)*1e1 +
				d(6), false
		case 8:
			return d(0)*1e7 + d(1)*1e6 +
				d(2)*1e5 + d(3)*1e4 +
				d(4)*1e3 + d(5)*1e2 +
				d(6)*1e1 + d(7), false
		case 9:
			return d(0)*1e8 + d(1)*1e7 +
				d(2)*1e6 + d(3)*1e5 +
				d(4)*1e4 + d(5)*1e3 +
				d(6)*1e2 + d(7)*1e1 +
				d(8), false
		case 10:
			return d(0)*1e9 + d(1)*1e8 +
				d(2)*1e7 + d(3)*1e6 +
				d(4)*1e5 + d(5)*1e4 +
				d(6)*1e3 + d(7)*1e2 +
				d(8)*1e1 + d(9), false
		case 11:
			return d(0)*1e10 + d(1)*1e9 +
				d(2)*1e8 + d(3)*1e7 +
				d(4)*1e6 + d(5)*1e5 +
				d(6)*1e4 + d(7)*1e3 +
				d(8)*1e2 + d(9)*1e1 +
				d(10), false
		case 12:
			return d(0)*1e11 + d(1)*1e10 +
				d(2)*1e9 + d(3)*1e8 +
				d(4)*1e7 + d(5)*1e6 +
				d(6)*1e5 + d(7)*1e4 +
				d(8)*1e3 + d(9)*1e2 +
				d(10)*1e1 + d(11), false
		case 13:
			return d(0)*1e12 + d(1)*1e11 +
				d(2)*1e10 + d(3)*1e9 +
				d(4)*1e8 + d(5)*1e7 +
				d(6)*1e6 + d(7)*1e5 +
				d(8)*1e4 + d(9)*1e3 +
				d(10)*1e2 + d(11)*1e1 +
				d(12), false
		case 14:
			return d(0)*1e13 + d(1)*1e12 +
				d(2)*1e11 + d(3)*1e10 +
				d(4)*1e9 + d(5)*1e8 +
				d(6)*1e7 + d(7)*1e6 +
				d(8)*1e5 + d(9)*1e4 +
				d(10)*1e3 + d(11)*1e2 +
				d(12)*1e1 + d(13), false
		case 15:
			return d(0)*1e14 + d(1)*1e13 +
				d(2)*1e12 + d(3)*1e11 +
				d(4)*1e10 + d(5)*1e9 +
				d(6)*1e8 + d(7)*1e7 +
				d(8)*1e6 + d(9)*1e5 +
				d(10)*1e4 + d(11)*1e3 +
				d(12)*1e2 + d(13)*1e1 +
				d(14), false
		case 16:
			return d(0)*1e15 + d(1)*1e14 +
				d(2)*1e13 + d(3)*1e12 +
				d(4)*1e11 + d(5)*1e10 +
				d(6)*1e9 + d(7)*1e8 +
				d(8)*1e7 + d(9)*1e6 +
				d(10)*1e5 + d(11)*1e4 +
				d(12)*1e3 + d(13)*1e2 +
				d(14)*1e1 + d(15), false
		case 17:
			return d(0)*1e16 + d(1)*1e15 +
				d(2)*1e14 + d(3)*1e13 +
				d(4)*1e12 + d(5)*1e11 +
				d(6)*1e10 + d(7)*1e9 +
				d(8)*1e8 + d(9)*1e7 +
				d(10)*1e6 + d(11)*1e5 +
				d(12)*1e4 + d(13)*1e3 +
				d(14)*1e2 + d(15)*1e1 +
				d(16), false
		case 18:
			return d(0)*1e17 + d(1)*1e16 +
				d(2)*1e15 + d(3)*1e14 +
				d(4)*1e13 + d(5)*1e12 +
				d(6)*1e11 + d(7)*1e10 +
				d(8)*1e9 + d(9)*1e8 +
				d(10)*1e7 + d(11)*1e6 +
				d(12)*1e5 + d(13)*1e4 +
				d(14)*1e3 + d(15)*1e2 +
				d(16)*1e1 + d(17), false
		case 19: // This case can overflow
			n = d(0)*1e18 + d(1)*1e17 +
				d(2)*1e16 + d(3)*1e15 +
				d(4)*1e14 + d(5)*1e13 +
				d(6)*1e12 + d(7)*1e11 +
				d(8)*1e10 + d(9)*1e9 +
				d(10)*1e8 + d(11)*1e7 +
				d(12)*1e6 + d(13)*1e5 +
				d(14)*1e4 + d(15)*1e3 +
				d(16)*1e2 + d(17)*1e1 +
				d(18)
			if n < 0 {
				return 0, true
			}
			return n, false
		}
	}
	return 0, true // Anything above 19 digits overflows int64
}
