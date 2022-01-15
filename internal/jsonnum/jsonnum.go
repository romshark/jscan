package jsonnum

func Parse(s string) (end int, err bool) {
	var i int
	var c rune

	switch s[0] {
	case '-':
		// Signed
		s = s[1:]
		if len(s) < 1 {
			// Expected at least one digit
			return 0, true
		}
		end++
	case '0':
		if len(s) > 1 {
			// Leading zero
			switch s[1] {
			case '.':
				s = s[2:]
				end += 2
				i = -1
				goto FRACTION
			case ' ', '\t', '\r', '\n', ',', '}', ']':
				// Zero
				return end + 1, false
			case 'e', 'E':
				s = s[2:]
				end += 2
				i = -1
				goto EXPONENT_SIGN
			default:
				// Unexpected rune
				return 0, true
			}
		}
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
	default:
		// Unexpected rune
		return 0, true
	}

	// Integer
	for i, c = range s {
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case '.':
			d := i + 1
			s = s[d:]
			end += d
			i = -1
			goto FRACTION
		case ' ', '\t', '\r', '\n', ',', '}', ']':
			if i < 1 {
				// Expected at least one digit
				return 0, true
			}
			// Integer
			return end + i, false
		case 'e', 'E':
			d := i + 1
			s = s[d:]
			end += d
			i = -1
			goto EXPONENT_SIGN
		default:
			// Unexpected rune
			return 0, true
		}
	}
	s = s[i+1:]

	if s == "" {
		// Integer without exponent
		return end + i + 1, false
	}

	// Fraction
	i = 0
FRACTION:
	for i, c = range s {
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case ' ', '\t', '\r', '\n', ',', '}', ']':
			if i < 1 {
				// Expected at least one digit
				return 0, true
			}
			// Number with fraction
			return end + i, false
		case 'e', 'E':
			d := i + 1
			s = s[d:]
			end += d
			i = -1
			goto EXPONENT_SIGN
		default:
			// Unexpected rune
			return 0, true
		}
	}
	if i == -1 {
		// Unexpected end of string
		return 0, true
	}
	s = s[i+1:]

	if s == "" {
		// Number (with fraction but) without exponent
		return end + i + 1, false
	}
	// Exponent
	switch s[0] {
	case 'e', 'E':
		s = s[1:]
		end++
	default:
		// Unexpected rune
		return 0, true
	}

	// Exponent sign
	i = 0
EXPONENT_SIGN:
	if s == "" {
		// Missing exponent value
		return 0, true
	}
	switch s[0] {
	case '-', '+':
		s = s[1:]
		end++
	}

	for i, c = range s {
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case ' ', '\t', '\r', '\n', ',', '}', ']':
			if i < 1 {
				// Expected at least one digit
				return 0, true
			}
			// Number with (fraction and) exponent
			return end + i, false
		default:
			// Unexpected rune
			return 0, true
		}
	}
	if i == -1 {
		// Unexpected end of string
		return 0, true
	}

	// Number with (fraction and) exponent
	return end + i + 1, false
}

func ParseBytes(s []byte) (end int, err bool) {
	var i int
	var c byte

	switch s[0] {
	case '-':
		// Signed
		s = s[1:]
		if len(s) < 1 {
			// Expected at least one digit
			return 0, true
		}
		end++
	case '0':
		if len(s) > 1 {
			// Leading zero
			switch s[1] {
			case '.':
				s = s[2:]
				end += 2
				i = -1
				goto FRACTION
			case ' ', '\t', '\r', '\n', ',', '}', ']':
				// Zero
				return end + 1, false
			case 'e', 'E':
				s = s[2:]
				end += 2
				i = -1
				goto EXPONENT_SIGN
			default:
				// Unexpected rune
				return 0, true
			}
		}
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
	default:
		// Unexpected rune
		return 0, true
	}

	// Integer
	for i, c = range s {
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case '.':
			d := i + 1
			s = s[d:]
			end += d
			i = -1
			goto FRACTION
		case ' ', '\t', '\r', '\n', ',', '}', ']':
			if i < 1 {
				// Expected at least one digit
				return 0, true
			}
			// Integer
			return end + i, false
		case 'e', 'E':
			d := i + 1
			s = s[d:]
			end += d
			i = -1
			goto EXPONENT_SIGN
		default:
			// Unexpected rune
			return 0, true
		}
	}
	s = s[i+1:]

	if len(s) < 1 {
		// Integer without exponent
		return end + i + 1, false
	}

	// Fraction
	i = 0
FRACTION:
	for i, c = range s {
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case ' ', '\t', '\r', '\n', ',', '}', ']':
			if i < 1 {
				// Expected at least one digit
				return 0, true
			}
			// Number with fraction
			return end + i, false
		case 'e', 'E':
			d := i + 1
			s = s[d:]
			end += d
			i = -1
			goto EXPONENT_SIGN
		default:
			// Unexpected rune
			return 0, true
		}
	}
	if i == -1 {
		// Unexpected end of string
		return 0, true
	}
	s = s[i+1:]

	if len(s) < 1 {
		// Number (with fraction but) without exponent
		return end + i + 1, false
	}

	// Exponent sign
EXPONENT_SIGN:
	if len(s) < 1 {
		// Missing exponent value
		return 0, true
	}
	switch s[0] {
	case '-', '+':
		s = s[1:]
		end++
	}

	for i, c = range s {
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case ' ', '\t', '\r', '\n', ',', '}', ']':
			if i < 1 {
				// Expected at least one digit
				return 0, true
			}
			// Number with (fraction and) exponent
			return end + i, false
		default:
			// Unexpected rune
			return 0, true
		}
	}
	if i == -1 {
		// Unexpected end of string
		return 0, true
	}

	// Number with (fraction and) exponent
	return end + i + 1, false
}
