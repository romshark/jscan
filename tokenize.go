package jscan

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/romshark/jscan/v2/internal/atoi"
	"github.com/romshark/jscan/v2/internal/jsonnum"
	"github.com/romshark/jscan/v2/internal/strfind"
)

// TokenType defines a token type.
type TokenType byte

const (
	_ TokenType = iota

	// TokenTypeObject is the start token of a composite object value.
	TokenTypeObject // '{'

	// TokenTypeObjectEnd is the end token of a composite object value.
	TokenTypeObjectEnd // '}'

	// TokenTypeArray is the start token of a composite array value.
	TokenTypeArray // '['

	// TokenTypeArrayEnd is the end token of a composite array value.
	TokenTypeArrayEnd // ']'

	// TokenTypeKey is an object key token.
	TokenTypeKey

	// TokenTypeTrue is the boolean true value token.
	TokenTypeTrue

	// TokenTypeFalse is the boolean false value token.
	TokenTypeFalse

	// TokenTypeNull is the null value token.
	TokenTypeNull

	// TokenTypeInteger is any (signed) integer value token.
	TokenTypeInteger

	// TokenTypeNumber is any (signed) non-integer number value token (exponents, float).
	TokenTypeNumber

	// TokenTypeString is a string value token.
	TokenTypeString
)

func (t TokenType) String() string {
	switch t {
	case TokenTypeObject:
		return "Object"
	case TokenTypeObjectEnd:
		return "ObjectEnd"
	case TokenTypeArray:
		return "Array"
	case TokenTypeArrayEnd:
		return "ArrayEnd"
	case TokenTypeKey:
		return "Key"
	case TokenTypeTrue:
		return "True"
	case TokenTypeFalse:
		return "False"
	case TokenTypeNull:
		return "Null"
	case TokenTypeInteger:
		return "Integer"
	case TokenTypeNumber:
		return "Number"
	case TokenTypeString:
		return "String"
	}
	return ""
}

// Token is any JSON token except comma, colon and space.
type Token[S ~string | ~[]byte] struct {
	// Index declares the start byte index in the source.
	Index int

	// End behaves differently for composite (TokenTypeObject, TokenTypeObjectEnd,
	// TokenTypeArray, TokenTypeArrayEnd) types and non-composite token types.
	//
	// For TokenTypeObject and TokenTypeArray End declares the index of the end token (
	// TokenTypeObjectEnd and TokenTypeArrayEnd respectively) in the token buffer.
	//
	// For TokenTypeObjectEnd and TokenTypeArrayEnd End declares the index of the
	// start token (TokenTypeObject and TokenTypeArray respectively) in the token buffer.
	//
	// For all other token types, End declares the index of the end byte of the value
	// in the source.
	//
	// End can be used to quickly skip over large sections of JSON.
	End int

	// Elements behaves differently for TokenTypeObject and TokenTypeArray.
	// For TokenTypeObject it declares the number of key-value pairs, whereas for
	// TokenTypeArray it declares the number of elements in the array.
	// Elements is meaningless for non-array and non-object tokens.
	Elements int

	Type TokenType
}

var (
	ErrOverflow  = fmt.Errorf("token value overflows integer type")
	ErrWrongType = fmt.Errorf("token value has different type")
)

const intSize = unsafe.Sizeof(int(0))

// Int returns the int value of the token.
// Expects src to be the source string provided to the tokenizer.
// Returns int(0) if the token is a null value.
// Returns ErrWrongType if the token isn't an integer value.
// Returns ErrOverflow if the value would overflow type int.
func (t Token[S]) Int(src S) (int, error) {
	if t.Type == TokenTypeNull {
		return 0, nil
	}
	if t.Type != TokenTypeInteger {
		return 0, ErrWrongType
	}
	if intSize != 8 {
		v, overflow := atoi.I32(src[t.Index:t.End])
		if overflow {
			return 0, ErrOverflow
		}
		return int(v), nil
	}
	v, overflow := atoi.I64(src[t.Index:t.End])
	if overflow {
		return 0, ErrOverflow
	}
	return int(v), nil
}

// Int8 returns the int8 value of the token.
// Expects src to be the source string provided to the tokenizer.
// Returns int8(0) if the token is a null value.
// Returns ErrWrongType if the token isn't an integer value.
// Returns ErrOverflow if the value would overflow type int8.
func (t Token[S]) Int8(src S) (int8, error) {
	if t.Type == TokenTypeNull {
		return 0, nil
	}
	if t.Type != TokenTypeInteger {
		return 0, ErrWrongType
	}
	v, overflow := atoi.I8(src[t.Index:t.End])
	if overflow {
		return 0, ErrOverflow
	}
	return v, nil
}

// Int16 returns the int16 value of the token.
// Expects src to be the source string provided to the tokenizer.
// Returns int16(0) if the token is a null value.
// Returns ErrWrongType if the token isn't an integer value.
// Returns ErrOverflow if the value would overflow type int16.
func (t Token[S]) Int16(src S) (int16, error) {
	if t.Type == TokenTypeNull {
		return 0, nil
	}
	if t.Type != TokenTypeInteger {
		return 0, ErrWrongType
	}
	v, overflow := atoi.I16(src[t.Index:t.End])
	if overflow {
		return 0, ErrOverflow
	}
	return v, nil
}

// Int32 returns the int32 value of the token.
// Expects src to be the source string provided to the tokenizer.
// Returns int32(0) if the token is a null value.
// Returns ErrWrongType if the token isn't an integer value.
// Returns ErrOverflow if the value would overflow type int32.
func (t Token[S]) Int32(src S) (int32, error) {
	if t.Type == TokenTypeNull {
		return 0, nil
	}
	if t.Type != TokenTypeInteger {
		return 0, ErrWrongType
	}
	v, overflow := atoi.I32(src[t.Index:t.End])
	if overflow {
		return 0, ErrOverflow
	}
	return v, nil
}

// Int64 returns the int64 value of the token.
// Expects src to be the source string provided to the tokenizer.
// Returns int64(0) if the token is a null value.
// Returns ErrWrongType if the token isn't an integer value.
// Returns ErrOverflow if the value would overflow type int64.
func (t Token[S]) Int64(src S) (int64, error) {
	if t.Type == TokenTypeNull {
		return 0, nil
	}
	if t.Type != TokenTypeInteger {
		return 0, ErrWrongType
	}
	v, overflow := atoi.I64(src[t.Index:t.End])
	if overflow {
		return 0, ErrOverflow
	}
	return v, nil
}

// Uint returns the uint value of the token.
// Expects src to be the source string provided to the tokenizer.
// Returns uint(0) if the token is a null value.
// Returns ErrWrongType if the token isn't an integer value.
// Returns ErrOverflow if the value would overflow type uint.
func (t Token[S]) Uint(src S) (uint, error) {
	if t.Type == TokenTypeNull {
		return 0, nil
	}
	if t.Type != TokenTypeInteger || src[t.Index] == '-' {
		return 0, ErrWrongType
	}
	if intSize != 8 {
		v, overflow := atoi.U32(src[t.Index:t.End])
		if overflow {
			return 0, ErrOverflow
		}
		return uint(v), nil
	}
	v, overflow := atoi.U64(src[t.Index:t.End])
	if overflow {
		return 0, ErrOverflow
	}
	return uint(v), nil
}

// Uint8 returns the uint8 value of the token.
// Expects src to be the source string provided to the tokenizer.
// Returns uint8(0) if the token is a null value.
// Returns ErrWrongType if the token isn't an integer value.
// Returns ErrOverflow if the value would overflow type uint8.
func (t Token[S]) Uint8(src S) (uint8, error) {
	if t.Type == TokenTypeNull {
		return 0, nil
	}
	if t.Type != TokenTypeInteger || src[t.Index] == '-' {
		return 0, ErrWrongType
	}
	v, overflow := atoi.U8(src[t.Index:t.End])
	if overflow {
		return 0, ErrOverflow
	}
	return uint8(v), nil
}

// Uint16 returns the uint16 value of the token.
// Expects src to be the source string provided to the tokenizer.
// Returns uint16(0) if the token is a null value.
// Returns ErrWrongType if the token isn't an integer value.
// Returns ErrOverflow if the value would overflow type uint16.
func (t Token[S]) Uint16(src S) (uint16, error) {
	if t.Type == TokenTypeNull {
		return 0, nil
	}
	if t.Type != TokenTypeInteger || src[t.Index] == '-' {
		return 0, ErrWrongType
	}
	v, overflow := atoi.U16(src[t.Index:t.End])
	if overflow {
		return 0, ErrOverflow
	}
	return uint16(v), nil
}

// Uint32 returns the uint32 value of the token.
// Expects src to be the source string provided to the tokenizer.
// Returns uint32(0) if the token is a null value.
// Returns ErrWrongType if the token isn't an integer value.
// Returns ErrOverflow if the value would overflow type uint32.
func (t Token[S]) Uint32(src S) (uint32, error) {
	if t.Type == TokenTypeNull {
		return 0, nil
	}
	if t.Type != TokenTypeInteger || src[t.Index] == '-' {
		return 0, ErrWrongType
	}
	v, overflow := atoi.U32(src[t.Index:t.End])
	if overflow {
		return 0, ErrOverflow
	}
	return uint32(v), nil
}

// Uint64 returns the uint64 value of the token.
// Expects src to be the source string provided to the tokenizer.
// Returns uint64(0) if the token is a null value.
// Returns ErrWrongType if the token isn't an integer value.
// Returns ErrOverflow if the value would overflow type uint64.
func (t Token[S]) Uint64(src S) (uint64, error) {
	if t.Type == TokenTypeNull {
		return 0, nil
	}
	if t.Type != TokenTypeInteger || src[t.Index] == '-' {
		return 0, ErrWrongType
	}
	v, overflow := atoi.U64(src[t.Index:t.End])
	if overflow {
		return 0, ErrOverflow
	}
	return uint64(v), nil
}

// Float32 returns the float32 value of the token.
// Expects src to be the source string provided to the tokenizer.
// Returns float32(0) if the token is a null value.
// Returns ErrWrongType if the token isn't a number value.
func (t Token[S]) Float32(src S) (float32, error) {
	if t.Type == TokenTypeNull {
		return 0, nil
	}
	if t.Type != TokenTypeNumber && t.Type != TokenTypeInteger {
		return 0, ErrWrongType
	}
	var sz S
	var s string
	if _, ok := any(sz).([]byte); ok {
		b := []byte(src[t.Index:t.End])
		s = unsafe.String(unsafe.SliceData(b), len(b))
	} else {
		s = string(src[t.Index:t.End])
	}
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(v), nil
}

// Float64 returns the float64 value of the token.
// Expects src to be the source string provided to the tokenizer.
// Returns float64(0) if the token is a null value.
// Returns ErrWrongType if the token isn't a number value.
func (t Token[S]) Float64(src S) (float64, error) {
	if t.Type == TokenTypeNull {
		return 0, nil
	}
	if t.Type != TokenTypeNumber && t.Type != TokenTypeInteger {
		return 0, ErrWrongType
	}
	var sz S
	var s string
	if _, ok := any(sz).([]byte); ok {
		b := []byte(src[t.Index:t.End])
		s = unsafe.String(unsafe.SliceData(b), len(b))
	} else {
		s = string(src[t.Index:t.End])
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

// Bool returns the bool value of the token.
// Expects src to be the source string provided to the tokenizer.
// Returns bool(false) if the token is a null value.
// Returns ErrWrongType if the token isn't a boolean value.
func (t Token[S]) Bool(src S) (bool, error) {
	switch t.Type {
	case TokenTypeTrue:
		return true, nil
	case TokenTypeFalse, TokenTypeNull:
		return false, nil
	}
	return false, ErrWrongType
}

// String returns the string value of the token.
// Expects src to be the source string provided to the tokenizer.
// Returns string("") if the token is a null value.
// Returns ErrWrongType if the token isn't a string value.
func (t Token[S]) String(src S) (string, error) {
	switch t.Type {
	case TokenTypeNull:
		return "", nil
	case TokenTypeString:
		return string(src[t.Index+1 : t.End-1]), nil
	}
	return "", ErrWrongType
}

// Tokenizer is a reusable tokenizer instance holding a stack and a token buffer
// which are reused across method calls.
type Tokenizer[S ~string | ~[]byte] struct {
	buffer []Token[S]
	stack  []int // Buffer index
}

// NewTokenizer creates a new reusable tokenizer instance.
//
// A higher preallocStackFrames value implies greater memory usage but also reduces
// the chance of dynamic memory allocations if the JSON depth surpasses the stack size.
// Use DefaultStackSizeTokenizer when not sure, which is equivalent to ~1KiB bytes of
// memory usage on 64-bit systems (1 frame = 8 bytes).

// A higher preallocTokenBuffer value also implies greater memory usage and also reduces
// the chance of dynamic memory allocations if the number of JSON tokens encountered
// surpasses the buffer size. Use DefaultTokenBufferSize when not sure, which is
// equivalent to ~32KiB of memory usage on 64-bit systems (1 token = 32 bytes).
func NewTokenizer[S ~string | ~[]byte](
	preallocStackFrames, preallocTokenBuffer int,
) *Tokenizer[S] {
	t := &Tokenizer[S]{
		buffer: make([]Token[S], preallocTokenBuffer),
		stack:  make([]int, preallocStackFrames),
	}
	return t
}

// TokenizeOne calls fn with the tokens of the first value from s.
//
// Unlike Tokenize, TokenizeOne doesn't return ErrorCodeUnexpectedToken when
// it encounters anything other than EOF after reading a valid JSON value.
// Returns an error if any and trailing as substring of s with the tokenized value cut.
// In case of an error trailing will be a substring of s cut up until the index
// where the error was encountered.
//
// WARNING: Don't use or alias tokens after fn returns!
func (t *Tokenizer[S]) TokenizeOne(
	s S, fn func(tokens []Token[S]) (err bool),
) (trailing S, err Error[S]) {
	return t.tokenize(s, fn)
}

func (t *Tokenizer[S]) topStackType() TokenType {
	return t.buffer[t.stack[len(t.stack)-1]].Type
}

// Tokenize calls fn with the tokens of the value from s.
//
// WARNING: Don't use or alias tokens after fn returns!
func (t *Tokenizer[S]) Tokenize(
	s S, fn func(tokens []Token[S]) (err bool),
) Error[S] {
	tail, err := t.tokenize(s, fn)
	if err.IsErr() {
		return err
	}
	var illegalChar bool
	tail, illegalChar = strfind.EndOfWhitespaceSeq(tail)
	if illegalChar {
		return getError(ErrorCodeIllegalControlChar, s, tail)
	}
	if len(tail) > 0 {
		return getError(ErrorCodeUnexpectedToken, s, tail)
	}
	return Error[S]{}
}

// tokenize calls fn once all tokens are parsed to the buffer.
// Returns the remainder of src and an error if any is encountered.
func (t *Tokenizer[S]) tokenize(src S, fn func(tokens []Token[S]) (err bool)) (S, Error[S]) {
	// Reset tokenizer
	t.buffer = t.buffer[:0]
	t.stack = t.stack[:0]

	var (
		index    int
		rollback S // Used as fallback for error report
		s        = src
		err      bool
	)

VALUE:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, err = strfind.EndOfWhitespaceSeq(s)
			if err {
				return s, getError(ErrorCodeIllegalControlChar, src, s)
			}
		}
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
	index = len(src) - len(s)
	s = s[1:]
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, err = strfind.EndOfWhitespaceSeq(s)
			if err {
				return s, getError(ErrorCodeIllegalControlChar, src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
	}

	if s[0] == '}' {
		t.buffer = append(
			t.buffer,
			Token[S]{
				Index:    index,
				End:      len(t.buffer) + 1,
				Type:     TokenTypeObject,
				Elements: 0,
			},
			Token[S]{
				Index:    len(src) - len(s),
				End:      len(t.buffer),
				Type:     TokenTypeObjectEnd,
				Elements: 0,
			},
		)
		s = s[1:]
		goto AFTER_VALUE
	}

	t.stack = append(t.stack, len(t.buffer))
	t.buffer = append(t.buffer, Token[S]{
		Index:    index,
		End:      0, // To be set once the end is discovered.
		Type:     TokenTypeObject,
		Elements: 0, // Will be updated on every key encountered.
	})

	goto OBJ_KEY

VALUE_ARRAY:
	index = len(src) - len(s)
	s = s[1:]

	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, err = strfind.EndOfWhitespaceSeq(s)
			if err {
				return s, getError(ErrorCodeIllegalControlChar, src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
	}

	if s[0] == ']' {
		t.buffer = append(
			t.buffer,
			Token[S]{
				Index:    index,
				End:      len(t.buffer) + 1,
				Type:     TokenTypeArray,
				Elements: 0,
			},
			Token[S]{
				Index:    len(src) - len(s),
				End:      len(t.buffer),
				Type:     TokenTypeArrayEnd,
				Elements: 0,
			},
		)
		s = s[1:]
		goto AFTER_VALUE
	}

	t.stack = append(t.stack, len(t.buffer))
	t.buffer = append(
		t.buffer,
		Token[S]{
			Index: index,
			End:   0, // To be set once the end is discovered.
			Type:  TokenTypeArray,
			// There must be at least one element if
			// end of array wasn't encountered immediately.
			Elements: 1,
		},
	)

	switch s[0] {
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

VALUE_NUMBER:
	{
		index = len(src) - len(s)
		var rc jsonnum.ReturnCode
		rollback = s
		if s, rc = jsonnum.ReadNumber(s); rc == jsonnum.ReturnCodeErr {
			return s, getError(ErrorCodeMalformedNumber, src, rollback)
		}
		t.buffer = append(t.buffer, Token[S]{
			Index:    index,
			End:      len(src) - len(s),
			Type:     TokenType(rc),
			Elements: 0,
		})
	}
	goto AFTER_VALUE

VALUE_STRING:
	index = len(src) - len(s)
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
			if lutEscape[s[1]] == 1 {
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

			t.buffer = append(t.buffer, Token[S]{
				Type:     TokenTypeString,
				Index:    index,
				Elements: 0,
				End:      len(src) - len(s),
			})

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
	index = len(src) - len(s)
	t.buffer = append(t.buffer, Token[S]{
		Type:     TokenTypeNull,
		Index:    index,
		End:      index + len("null"),
		Elements: 0,
	})
	s = s[len("null"):]

	goto AFTER_VALUE

VALUE_FALSE:
	if len(s) < 5 || string(s[:5]) != "false" {
		return s, getError(ErrorCodeUnexpectedToken, src, s)
	}
	index = len(src) - len(s)
	t.buffer = append(t.buffer, Token[S]{
		Type:     TokenTypeFalse,
		Index:    index,
		End:      index + len("false"),
		Elements: 0,
	})
	s = s[len("false"):]

	goto AFTER_VALUE

VALUE_TRUE:
	if s := s; len(s) < 4 || string(s[:4]) != "true" {
		return s, getError(ErrorCodeUnexpectedToken, src, s)
	}
	index = len(src) - len(s)
	t.buffer = append(t.buffer, Token[S]{
		Type:     TokenTypeTrue,
		Index:    index,
		End:      index + len("true"),
		Elements: 0,
	})
	s = s[len("true"):]

	goto AFTER_VALUE

OBJ_KEY:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, err = strfind.EndOfWhitespaceSeq(s)
			if err {
				return s, getError(ErrorCodeIllegalControlChar, src, s)
			}
		}
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

	index = len(src) - len(s)
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
			if lutEscape[s[1]] == 1 {
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
			t.buffer = append(t.buffer, Token[S]{
				Type:     TokenTypeKey,
				Index:    index,
				End:      len(src) - len(s),
				Elements: 0,
			})
			t.buffer[t.stack[len(t.stack)-1]].Elements++ // Update object keys count
			if len(s) < 1 {
				return s, getError(ErrorCodeUnexpectedEOF, src, s)
			}
			if s[0] <= ' ' {
				switch s[0] {
				case ' ', '\t', '\r', '\n':
					s, err = strfind.EndOfWhitespaceSeq(s)
					if err {
						return s, getError(ErrorCodeIllegalControlChar, src, s)
					}
				}
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
		default:
			if s[0] < 0x20 {
				return s, getError(ErrorCodeIllegalControlChar, src, s)
			}
			s = s[1:]
		}
	}

AFTER_VALUE:
	if len(t.stack) == 0 {
		if fn(t.buffer) {
			return s, getError(ErrorCodeCallback, src, s)
		}
		return s, Error[S]{}
	}
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, err = strfind.EndOfWhitespaceSeq(s)
			if err {
				return s, getError(ErrorCodeIllegalControlChar, src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, src, s)
		}
	}
	switch s[0] {
	case ',':
		s = s[1:]
		if t.topStackType() == TokenTypeArray {
			t.buffer[t.stack[len(t.stack)-1]].Elements++ // Update array elements count
			goto VALUE
		}
		goto OBJ_KEY
	case '}':
		if t.topStackType() != stackNodeTypeObject {
			return s, getError(ErrorCodeUnexpectedToken, src, s)
		}

		t.buffer[t.stack[len(t.stack)-1]].End = len(t.buffer) // Link start token
		t.buffer = append(t.buffer, Token[S]{
			Index: len(src) - len(s),
			End:   t.stack[len(t.stack)-1],
			Type:  TokenTypeObjectEnd,
		})

		s = s[1:]
		t.stack = t.stack[:len(t.stack)-1]
		goto AFTER_VALUE
	case ']':
		if t.topStackType() != TokenTypeArray {
			return s, getError(ErrorCodeUnexpectedToken, src, s)
		}

		t.buffer[t.stack[len(t.stack)-1]].End = len(t.buffer) // Link start token
		t.buffer = append(t.buffer, Token[S]{
			Index: len(src) - len(s),
			End:   t.stack[len(t.stack)-1],
			Type:  TokenTypeArrayEnd,
		})

		s = s[1:]
		t.stack = t.stack[:len(t.stack)-1]
		goto AFTER_VALUE
	}
	if s[0] < 0x20 {
		return s, getError(ErrorCodeIllegalControlChar, src, s)
	}
	return s, getError(ErrorCodeUnexpectedToken, src, s)
}
