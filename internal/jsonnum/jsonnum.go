package jsonnum

type ReturnCode byte

const (
	ReturnCodeErr     ReturnCode = 0
	ReturnCodeInteger ReturnCode = 9  // Translates to jscan.TokenTypeInteger
	ReturnCodeNumber  ReturnCode = 10 // Translates to jscan.TokenTypeNumber
)

// ReadNumber returns s with the read number value cut off
// and err=ReturnCodeErr if a syntax error was encountered.
func ReadNumber[S ~string | ~[]byte](s S) (trailing S, err ReturnCode) {
	var i int

	if s[0] == '-' {
		// Signed
		s = s[1:]
		if len(s) < 1 {
			// Expected at least one digit
			return s, ReturnCodeErr
		}
	}

	if s[0] == '0' {
		s = s[1:]
		if len(s) < 1 {
			// Zero
			return s, ReturnCodeInteger
		}
		// Leading zero
		switch s[0] {
		case '.':
			s = s[1:]
			goto FRACTION
		case 'e', 'E':
			s = s[1:]
			goto EXPONENT_SIGN
		default:
			// Zero
			return s, ReturnCodeInteger
		}
	}

	// Integer
	if len(s) < 1 || (s[0] < '1' || s[0] > '9') {
		// Expected at least one digit
		return s, ReturnCodeErr
	}
	s = s[1:]
	for len(s) >= 8 {
		if s[0] < '0' || s[0] > '9' {
			if s[0] == 'e' || s[0] == 'E' {
				s = s[1:]
				goto EXPONENT_SIGN
			} else if s[0] == '.' {
				s = s[1:]
				goto FRACTION
			}
			// Integer
			return s, ReturnCodeInteger
		}
		if s[1] < '0' || s[1] > '9' {
			if s[1] == 'e' || s[1] == 'E' {
				s = s[2:]
				goto EXPONENT_SIGN
			} else if s[1] == '.' {
				s = s[2:]
				goto FRACTION
			}
			// Integer
			s = s[1:]
			return s, ReturnCodeInteger
		}
		if s[2] < '0' || s[2] > '9' {
			if s[2] == 'e' || s[2] == 'E' {
				s = s[3:]
				goto EXPONENT_SIGN
			} else if s[2] == '.' {
				s = s[3:]
				goto FRACTION
			}
			// Integer
			s = s[2:]
			return s, ReturnCodeInteger
		}
		if s[3] < '0' || s[3] > '9' {
			if s[3] == 'e' || s[3] == 'E' {
				s = s[4:]
				goto EXPONENT_SIGN
			} else if s[3] == '.' {
				s = s[4:]
				goto FRACTION
			}
			// Integer
			s = s[3:]
			return s, ReturnCodeInteger
		}
		if s[4] < '0' || s[4] > '9' {
			if s[4] == 'e' || s[4] == 'E' {
				s = s[5:]
				goto EXPONENT_SIGN
			} else if s[4] == '.' {
				s = s[5:]
				goto FRACTION
			}
			// Integer
			s = s[4:]
			return s, ReturnCodeInteger
		}
		if s[5] < '0' || s[5] > '9' {
			if s[5] == 'e' || s[5] == 'E' {
				s = s[6:]
				goto EXPONENT_SIGN
			} else if s[5] == '.' {
				s = s[6:]
				goto FRACTION
			}
			// Integer
			s = s[5:]
			return s, ReturnCodeInteger
		}
		if s[6] < '0' || s[6] > '9' {
			if s[6] == 'e' || s[6] == 'E' {
				s = s[7:]
				goto EXPONENT_SIGN
			} else if s[6] == '.' {
				s = s[7:]
				goto FRACTION
			}
			// Integer
			s = s[6:]
			return s, ReturnCodeInteger
		}
		if s[7] < '0' || s[7] > '9' {
			if s[7] == 'e' || s[7] == 'E' {
				s = s[8:]
				goto EXPONENT_SIGN
			} else if s[7] == '.' {
				s = s[8:]
				goto FRACTION
			}
			// Integer
			s = s[7:]
			return s, ReturnCodeInteger
		}
		s = s[8:]
	}
	for i = 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			if s[i] == 'e' || s[i] == 'E' {
				s = s[i+1:]
				goto EXPONENT_SIGN
			} else if s[i] == '.' {
				s = s[i+1:]
				goto FRACTION
			}
			// Integer
			s = s[i:]
			return s, ReturnCodeInteger
		}
	}
	s = s[i:]

	if len(s) < 1 {
		// Integer without exponent
		return s, ReturnCodeInteger
	}

FRACTION:
	if len(s) < 1 || (s[0] < '0' || s[0] > '9') {
		// Expected at least one digit
		return s, ReturnCodeErr
	}
	s = s[1:]

	for len(s) >= 8 {
		if s[0] < '0' || s[0] > '9' {
			if s[0] == 'e' || s[0] == 'E' {
				s = s[1:]
				goto EXPONENT_SIGN
			}
			return s, ReturnCodeNumber
		}
		if s[1] < '0' || s[1] > '9' {
			if s[1] == 'e' || s[1] == 'E' {
				s = s[2:]
				goto EXPONENT_SIGN
			}
			s = s[1:]
			return s, ReturnCodeNumber
		}
		if s[2] < '0' || s[2] > '9' {
			if s[2] == 'e' || s[2] == 'E' {
				s = s[3:]
				goto EXPONENT_SIGN
			}
			s = s[2:]
			return s, ReturnCodeNumber
		}
		if s[3] < '0' || s[3] > '9' {
			if s[3] == 'e' || s[3] == 'E' {
				s = s[4:]
				goto EXPONENT_SIGN
			}
			s = s[3:]
			return s, ReturnCodeNumber
		}
		if s[4] < '0' || s[4] > '9' {
			if s[4] == 'e' || s[4] == 'E' {
				s = s[5:]
				goto EXPONENT_SIGN
			}
			s = s[4:]
			return s, ReturnCodeNumber
		}
		if s[5] < '0' || s[5] > '9' {
			if s[5] == 'e' || s[5] == 'E' {
				s = s[6:]
				goto EXPONENT_SIGN
			}
			s = s[5:]
			return s, ReturnCodeNumber
		}
		if s[6] < '0' || s[6] > '9' {
			if s[6] == 'e' || s[6] == 'E' {
				s = s[7:]
				goto EXPONENT_SIGN
			}
			s = s[6:]
			return s, ReturnCodeNumber
		}
		if s[7] < '0' || s[7] > '9' {
			if s[7] == 'e' || s[7] == 'E' {
				s = s[8:]
				goto EXPONENT_SIGN
			}
			s = s[7:]
			return s, ReturnCodeNumber
		}
		s = s[8:]
	}
	for i = 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			if s[i] == 'e' || s[i] == 'E' {
				s = s[i+1:]
				goto EXPONENT_SIGN
			}
			s = s[i:]
			return s, ReturnCodeNumber
		}
	}
	s = s[i:]

	if len(s) < 1 {
		// Number (with fraction but) without exponent
		return s, ReturnCodeNumber
	}

EXPONENT_SIGN:
	if len(s) < 1 {
		// Missing exponent value
		return s, ReturnCodeErr
	}
	if s[0] == '-' || s[0] == '+' {
		s = s[1:]
	}
	if len(s) < 1 || (s[0] < '0' || s[0] > '9') {
		// Expected at least one digit
		return s, ReturnCodeErr
	}
	s = s[1:]

	for len(s) >= 8 {
		if s[0] < '0' || s[0] > '9' {
			// Number with (fraction and) exponent
			return s, ReturnCodeNumber
		}
		if s[1] < '0' || s[1] > '9' {
			// Number with (fraction and) exponent
			s = s[1:]
			return s, ReturnCodeNumber
		}
		if s[2] < '0' || s[2] > '9' {
			// Number with (fraction and) exponent
			s = s[2:]
			return s, ReturnCodeNumber
		}
		if s[3] < '0' || s[3] > '9' {
			// Number with (fraction and) exponent
			s = s[3:]
			return s, ReturnCodeNumber
		}
		if s[4] < '0' || s[4] > '9' {
			// Number with (fraction and) exponent
			s = s[4:]
			return s, ReturnCodeNumber
		}
		if s[5] < '0' || s[5] > '9' {
			// Number with (fraction and) exponent
			s = s[5:]
			return s, ReturnCodeNumber
		}
		if s[6] < '0' || s[6] > '9' {
			// Number with (fraction and) exponent
			s = s[6:]
			return s, ReturnCodeNumber
		}
		if s[7] < '0' || s[7] > '9' {
			// Number with (fraction and) exponent
			s = s[7:]
			return s, ReturnCodeNumber
		}
		s = s[8:]
	}
	for i = 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			// Number with (fraction and) exponent
			s = s[i:]
			return s, ReturnCodeNumber
		}
	}
	s = s[i:]
	return s, ReturnCodeNumber
}
