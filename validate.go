package jscan

// Valid returns true if s is a valid JSON value, otherwise returns false.
//
// Unlike (*Validator).Valid this function will take a validator instance
// from a global pool and can therefore be less efficient.
// Consider reusing a Validator instance instead.
func Valid[S ~string | ~[]byte](s S) bool {
	return !Validate(s).IsErr()
}

// ValidateOne scans one JSON value from s and returns an error if it's invalid
// and trailing as substring of s with the scanned value cut.
// In case of an error trailing will be a substring of s cut up until the index
// where the error was encountered.
//
// Unlike (*Validator).ValidateOne this function will take a validator instance
// from a global pool and can therefore be less efficient.
// Consider reusing a Validator instance instead.
func ValidateOne[S ~string | ~[]byte](s S) (trailing S, err Error[S]) {
	var v *Validator[S]
	switch any(s).(type) {
	case string:
		x := validatorPoolString.Get()
		defer validatorPoolString.Put(x)
		v = x.(*Validator[S])
	case []byte:
		x := validatorPoolBytes.Get()
		defer validatorPoolBytes.Put(x)
		v = x.(*Validator[S])
	}
	v.stack = v.stack[:0]

	return validate(v.stack, s)
}

// Validate returns an error if s is invalid JSON.
//
// Unlike (*Validator).Validate this function will take a validator instance
// from a global pool and can therefore be less efficient.
// Consider reusing a Validator instance instead.
func Validate[S ~string | ~[]byte](s S) Error[S] {
	var v *Validator[S]
	switch any(s).(type) {
	case string:
		x := validatorPoolString.Get()
		defer validatorPoolString.Put(x)
		v = x.(*Validator[S])
	case []byte:
		x := validatorPoolBytes.Get()
		defer validatorPoolBytes.Put(x)
		v = x.(*Validator[S])
	}
	v.stack = v.stack[:0]

	t, err := validate(v.stack, s)
	if err.IsErr() {
		return err
	}
	if t = skipSpace(t); len(t) > 0 {
		if t[0] < 0x20 {
			return getError(ErrorCodeIllegalControlChar, s, t)
		}
		return getError(ErrorCodeUnexpectedToken, s, t)
	}
	return Error[S]{}
}

// NewValidator creates a new reusable validator instance.
// A higher preallocStackFrames value implies greater memory usage but also reduces
// the chance of dynamic memory allocations if the JSON depth surpasses the stack size.
// preallocStackFrames of 1024 is equivalent to ~1KiB of memory usage (1 frame = 1 byte).
// Use DefaultStackSizeValidator when not sure.
func NewValidator[S ~string | ~[]byte](preallocStackFrames int) *Validator[S] {
	return &Validator[S]{
		stack: make([]stackNodeType, 0, preallocStackFrames),
	}
}

// Validator is a reusable validator instance.
// The validator is more efficient than the parser at JSON validation.
// A validator instance can be more efficient than global Valid, Validate and ValidateOne
// function calls due to potential stack frame allocation avoidance.
type Validator[S ~string | ~[]byte] struct{ stack []stackNodeType }

// Valid returns true if s is a valid JSON value, otherwise returns false.
func (v *Validator[S]) Valid(s S) bool {
	return !v.Validate(s).IsErr()
}

// ValidateOne scans one JSON value from s and returns an error if it's invalid
// and trailing as substring of s with the scanned value cut.
// In case of an error trailing will be a substring of s cut up until the index
// where the error was encountered.
func (v *Validator[S]) ValidateOne(s S) (trailing S, err Error[S]) {
	return validate(v.stack, s)
}

// Validate returns an error if s is invalid JSON,
// otherwise returns a zero value of Error[S].
func (v *Validator[S]) Validate(s S) Error[S] {
	t, err := validate(v.stack, s)
	if err.IsErr() {
		return err
	}
	if t = skipSpace(t); len(t) > 0 {
		if t[0] < 0x20 {
			return getError(ErrorCodeIllegalControlChar, s, t)
		}
		return getError(ErrorCodeUnexpectedToken, s, t)
	}
	return Error[S]{}
}

// validate returns the remainder of i.src and an error if any is encountered.
func validate[S ~string | ~[]byte](st []stackNodeType, s S) (S, Error[S]) {
	src := s
	var top stackNodeType

	stPop := func() { st = st[:len(st)-1] }
	stTop := func() {
		if len(st) < 1 {
			top = 0
			return
		}
		top = st[len(st)-1]
	}
	stPush := func(t stackNodeType) { st = append(st, t) }

VALUE:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if lutSX[s[0]] == 1 {
		s = skipSpace(s)
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
	}
	switch s[0] {
	case '{':
		goto VALUE_OBJECT
	case '[':
		goto VALUE_ARRAY
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		goto VALUE_NUMBER
	case '"':
		goto VALUE_STRING
	case 'n':
		goto VALUE_NULL
	case 'f':
		goto VALUE_FALSE
	case 't':
		goto VALUE_TRUE
	}
	if s[0] < 0x20 {
		return s, getError(ErrorCodeIllegalControlChar, src, s)
	}
	return s, getError(ErrorCodeUnexpectedToken, src, s)

VALUE_OBJECT:
	s = s[1:]
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if lutSX[s[0]] == 1 {
		s = skipSpace(s)
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
	}
	if s[0] == '}' {
		s = s[1:]
		goto AFTER_VALUE
	}
	stPush(stackNodeTypeObject)
	goto OBJ_KEY

VALUE_ARRAY:
	stPush(stackNodeTypeArray)
	s = s[1:]
	goto VALUE_OR_ARR_TERM

VALUE_NUMBER:
	if s[0] == '-' {
		// Signed
		s = s[1:]
		if len(s) < 1 {
			// Expected at least one digit
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
	}

	if s[0] == '0' {
		s = s[1:]
		if len(s) < 1 {
			// Number zero
			goto AFTER_VALUE
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
			goto AFTER_VALUE
		}
	}

	// Integer
	if s[0] < '1' || s[0] > '9' {
		// Expected at least one digit
		return s, getError(ErrorCodeMalformedNumber, src, s)
	}
	s = s[1:]
	for len(s) > 7 {
		if lutED[s[0]] != 2 {
			goto CHECK_INT
		}
		if lutED[s[1]] != 2 {
			s = s[1:]
			goto CHECK_INT
		}
		if lutED[s[2]] != 2 {
			s = s[2:]
			goto CHECK_INT
		}
		if lutED[s[3]] != 2 {
			s = s[3:]
			goto CHECK_INT
		}
		if lutED[s[4]] != 2 {
			s = s[4:]
			goto CHECK_INT
		}
		if lutED[s[5]] != 2 {
			s = s[5:]
			goto CHECK_INT
		}
		if lutED[s[6]] != 2 {
			s = s[6:]
			goto CHECK_INT
		}
		if lutED[s[7]] != 2 {
			s = s[7:]
			goto CHECK_INT
		}
		s = s[8:]
	}
	for ; len(s) > 0; s = s[1:] {
		if s[0] < '0' || s[0] > '9' {
			if s[0] == 'e' || s[0] == 'E' {
				s = s[1:]
				goto EXPONENT_SIGN
			} else if s[0] == '.' {
				s = s[1:]
				goto FRACTION
			}
			// Integer
			goto AFTER_VALUE
		}
	}

	if len(s) < 1 {
		// Integer without exponent
		goto AFTER_VALUE
	}

CHECK_INT:
	_ = 0 // Fix test coverage misreport
	switch s[0] {
	case 'e', 'E':
		s = s[1:]
		goto EXPONENT_SIGN
	case '.':
		s = s[1:]
		goto FRACTION
	default:
		// Integer
		goto AFTER_VALUE
	}

FRACTION:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] < '0' || s[0] > '9' {
		// Expected at least one digit
		return s, getError(ErrorCodeMalformedNumber, src, s)
	}
	s = s[1:]
	for len(s) > 7 {
		if lutED[s[0]] != 2 {
			goto CHECK_FRAC
		}
		if lutED[s[1]] != 2 {
			s = s[1:]
			goto CHECK_FRAC
		}
		if lutED[s[2]] != 2 {
			s = s[2:]
			goto CHECK_FRAC
		}
		if lutED[s[3]] != 2 {
			s = s[3:]
			goto CHECK_FRAC
		}
		if lutED[s[4]] != 2 {
			s = s[4:]
			goto CHECK_FRAC
		}
		if lutED[s[5]] != 2 {
			s = s[5:]
			goto CHECK_FRAC
		}
		if lutED[s[6]] != 2 {
			s = s[6:]
			goto CHECK_FRAC
		}
		if lutED[s[7]] != 2 {
			s = s[7:]
			goto CHECK_FRAC
		}
		s = s[8:]
	}
	for ; len(s) > 0; s = s[1:] {
		if s[0] < '0' || s[0] > '9' {
			if s[0] == 'e' || s[0] == 'E' {
				s = s[1:]
				goto EXPONENT_SIGN
			}
			goto AFTER_VALUE
		}
	}

	if len(s) < 1 {
		// Number (with fraction but) without exponent
		goto AFTER_VALUE
	}

CHECK_FRAC:
	if s[0] == 'e' || s[0] == 'E' {
		s = s[1:]
		goto EXPONENT_SIGN
	}
	goto AFTER_VALUE

EXPONENT_SIGN:
	if len(s) < 1 {
		// Missing exponent value
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] == '-' || s[0] == '+' {
		s = s[1:]
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
	}
	if s[0] < '0' || s[0] > '9' {
		// Expected at least one digit
		return s, getError(ErrorCodeMalformedNumber, src, s)
	}
	s = s[1:]
	for len(s) > 7 {
		if lutED[s[0]] != 2 {
			// Number with (fraction and) exponent
			goto AFTER_VALUE
		}
		if lutED[s[1]] != 2 {
			// Number with (fraction and) exponent
			s = s[1:]
			goto AFTER_VALUE
		}
		if lutED[s[2]] != 2 {
			// Number with (fraction and) exponent
			s = s[2:]
			goto AFTER_VALUE
		}
		if lutED[s[3]] != 2 {
			// Number with (fraction and) exponent
			s = s[3:]
			goto AFTER_VALUE
		}
		if lutED[s[4]] != 2 {
			// Number with (fraction and) exponent
			s = s[4:]
			goto AFTER_VALUE
		}
		if lutED[s[5]] != 2 {
			// Number with (fraction and) exponent
			s = s[5:]
			goto AFTER_VALUE
		}
		if lutED[s[6]] != 2 {
			// Number with (fraction and) exponent
			s = s[6:]
			goto AFTER_VALUE
		}
		if lutED[s[7]] != 2 {
			// Number with (fraction and) exponent
			s = s[7:]
			goto AFTER_VALUE
		}
		s = s[8:]
	}
	for ; len(s) > 0; s = s[1:] {
		if s[0] < '0' || s[0] > '9' {
			// Number with (fraction and) exponent
			goto AFTER_VALUE
		}
	}
	goto AFTER_VALUE

VALUE_STRING:
	s = s[1:]
	for {
		for ; len(s) > 15; s = s[16:] {
			if lutStr[s[0]] != 0 {
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[1]] != 0 {
				s = s[1:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[2]] != 0 {
				s = s[2:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[3]] != 0 {
				s = s[3:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[4]] != 0 {
				s = s[4:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[5]] != 0 {
				s = s[5:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[6]] != 0 {
				s = s[6:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[7]] != 0 {
				s = s[7:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[8]] != 0 {
				s = s[8:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[9]] != 0 {
				s = s[9:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[10]] != 0 {
				s = s[10:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[11]] != 0 {
				s = s[11:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[12]] != 0 {
				s = s[12:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[13]] != 0 {
				s = s[13:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[14]] != 0 {
				s = s[14:]
				goto CHECK_STRING_CHARACTER
			}
			if lutStr[s[15]] != 0 {
				s = s[15:]
				goto CHECK_STRING_CHARACTER
			}
			continue
		}

	CHECK_STRING_CHARACTER:
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
		switch s[0] {
		case '\\':
			if len(s) < 2 {
				s = s[1:]
				return s, getError(ErrorCodeUnexpectedEOF, src, s)
			}
			if lutED[s[1]] == 1 {
				s = s[2:]
				continue
			}
			if s[1] != 'u' {
				return s, getError(ErrorCodeInvalidEscape, src, s)
			}
			if len(s) < 6 ||
				lutSX[s[5]] != 2 ||
				lutSX[s[4]] != 2 ||
				lutSX[s[3]] != 2 ||
				lutSX[s[2]] != 2 {
				return s, getError(ErrorCodeInvalidEscape, src, s)
			}
			s = s[5:]
		case '"':
			s = s[1:]
			goto AFTER_VALUE
		default:
			if s[0] < 0x20 {
				return s, getError(ErrorCodeIllegalControlChar, src, s)
			}
			s = s[1:]
		}
	}

VALUE_NULL:
	if len(s) < 4 || string(s[:4]) != "null" {
		return s, getError(ErrorCodeUnexpectedToken, src, s)
	}
	s = s[len("null"):]
	goto AFTER_VALUE

VALUE_FALSE:
	if len(s) < 5 || string(s[:5]) != "false" {
		return s, getError(ErrorCodeUnexpectedToken, src, s)
	}
	s = s[len("false"):]
	goto AFTER_VALUE

VALUE_TRUE:
	if s := s; len(s) < 4 || string(s[:4]) != "true" {
		return s, getError(ErrorCodeUnexpectedToken, src, s)
	}
	s = s[len("true"):]
	goto AFTER_VALUE

OBJ_KEY:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if lutSX[s[0]] == 1 {
		s = skipSpace(s)
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
	}
	if s[0] != '"' {
		if s[0] < 0x20 {
			return s, getError(ErrorCodeIllegalControlChar, src, s)
		}
		return s, getError(ErrorCodeUnexpectedToken, src, s)
	}

	s = s[1:]
	for {
		for ; len(s) > 15; s = s[16:] {
			if lutStr[s[0]] != 0 {
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[1]] != 0 {
				s = s[1:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[2]] != 0 {
				s = s[2:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[3]] != 0 {
				s = s[3:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[4]] != 0 {
				s = s[4:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[5]] != 0 {
				s = s[5:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[6]] != 0 {
				s = s[6:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[7]] != 0 {
				s = s[7:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[8]] != 0 {
				s = s[8:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[9]] != 0 {
				s = s[9:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[10]] != 0 {
				s = s[10:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[11]] != 0 {
				s = s[11:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[12]] != 0 {
				s = s[12:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[13]] != 0 {
				s = s[13:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[14]] != 0 {
				s = s[14:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			if lutStr[s[15]] != 0 {
				s = s[15:]
				goto CHECK_FIELDNAME_STRING_CHARACTER
			}
			continue
		}

	CHECK_FIELDNAME_STRING_CHARACTER:
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
		switch s[0] {
		case '\\':
			if len(s) < 2 {
				s = s[1:]
				return s, getError(ErrorCodeUnexpectedEOF, src, s)
			}
			if lutED[s[1]] == 1 {
				s = s[2:]
				continue
			}
			if s[1] != 'u' {
				return s, getError(ErrorCodeInvalidEscape, src, s)
			}
			if len(s) < 6 ||
				lutSX[s[5]] != 2 ||
				lutSX[s[4]] != 2 ||
				lutSX[s[3]] != 2 ||
				lutSX[s[2]] != 2 {
				return s, getError(ErrorCodeInvalidEscape, src, s)
			}
			s = s[5:]
		case '"':
			s = s[1:]
			goto AFTER_OBJ_KEY_STRING
		default:
			if s[0] < 0x20 {
				return s, getError(ErrorCodeIllegalControlChar, src, s)
			}
			s = s[1:]
		}
	}
AFTER_OBJ_KEY_STRING:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if lutSX[s[0]] == 1 {
		s = skipSpace(s)
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
	}
	if s[0] != ':' {
		if s[0] < 0x20 {
			return s, getError(ErrorCodeIllegalControlChar, src, s)
		}
		return s, getError(ErrorCodeUnexpectedToken, src, s)
	}
	s = s[1:]
	goto VALUE

VALUE_OR_ARR_TERM:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if lutSX[s[0]] == 1 {
		s = skipSpace(s)
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
	}
	switch s[0] {
	case ']':
		s = s[1:]
		stPop()
		goto AFTER_VALUE
	case '{':
		goto VALUE_OBJECT
	case '[':
		goto VALUE_ARRAY
	case '"':
		goto VALUE_STRING
	case 't':
		goto VALUE_TRUE
	case 'f':
		goto VALUE_FALSE
	case 'n':
		goto VALUE_NULL
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		goto VALUE_NUMBER
	}
	if s[0] < 0x20 {
		return s, getError(ErrorCodeIllegalControlChar, src, s)
	}
	return s, getError(ErrorCodeUnexpectedToken, src, s)

AFTER_VALUE:
	stTop()
	if top == 0 {
		return s, Error[S]{}
	}
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if lutSX[s[0]] == 1 {
		s = skipSpace(s)
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
	}
	switch s[0] {
	case ',':
		s = s[1:]
		if top == stackNodeTypeArray {
			goto VALUE
		}
		goto OBJ_KEY
	case '}':
		if top != stackNodeTypeObject {
			return s, getError(ErrorCodeUnexpectedToken, src, s)
		}
		s = s[1:]
		stPop()
		goto AFTER_VALUE
	case ']':
		if top != stackNodeTypeArray {
			return s, getError(ErrorCodeUnexpectedToken, src, s)
		}
		s = s[1:]
		stPop()
		goto AFTER_VALUE
	}
	if s[0] < 0x20 {
		return s, getError(ErrorCodeIllegalControlChar, src, s)
	}
	return s, getError(ErrorCodeUnexpectedToken, src, s)
}
