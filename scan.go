package jscan

import (
	"github.com/romshark/jscan/v2/internal/jsonnum"
	"github.com/romshark/jscan/v2/internal/strfind"
)

// ScanOne calls fn for every encountered value including objects and arrays.
// When an object or array is encountered fn will also be called for each of its
// member and element values.
//
// Unlike Scan, ScanOne doesn't return ErrorCodeUnexpectedToken when
// it encounters anything other than EOF after reading a valid JSON value.
// Returns an error if any and trailing as substring of s with the scanned value cut.
// In case of an error trailing will be a substring of s cut up until the index
// where the error was encountered.
//
// Unlike (*Parser).ScanOne this function will take an iterator instance
// from a global iterator pool and can therefore be less efficient.
// Consider reusing a Parser instance instead.
//
// TIP: Explicitly cast s to string or []byte to use the global iterator pools
// and avoid an unecessary iterator allocation such as when dealing with
// json.RawMessage and similar types derived from string or []byte.
//
//	m := json.RawMessage(`1`)
//	jscan.ScanOne([]byte(m), // Cast m to []byte to avoid allocation!
//
// WARNING: Don't use or alias *Iterator[S] after fn returns!
func ScanOne[S ~string | ~[]byte](
	s S, fn func(*Iterator[S]) (err bool),
) (trailing S, err Error[S]) {
	var i *Iterator[S]
	switch any(s).(type) {
	case string:
		x := iteratorPoolString.Get()
		defer iteratorPoolString.Put(x)
		i = x.(*Iterator[S])
	case []byte:
		x := iteratorPoolBytes.Get()
		defer iteratorPoolBytes.Put(x)
		i = x.(*Iterator[S])
	default:
		i = newIterator[S]()
	}
	i.src = s
	reset(i)
	return scan(i, fn)
}

// Scan calls fn for every encountered value including objects and arrays.
// When an object or array is encountered fn will also be called for each of its
// member and element values.
//
// Unlike (*Parser).Scan this function will take an iterator instance
// from a global iterator pool and can therefore be less efficient.
// Consider reusing a Parser instance instead.
//
// TIP: Explicitly cast s to string or []byte to use the global iterator pools
// and avoid an unecessary iterator allocation such as when dealing with
// json.RawMessage and similar types derived from string or []byte.
//
//	m := json.RawMessage(`1`)
//	jscan.Scan([]byte(m), // Cast m to []byte to avoid allocation!
//
// WARNING: Don't use or alias *Iterator[S] after fn returns!
func Scan[S ~string | ~[]byte](
	s S, fn func(*Iterator[S]) (err bool),
) (err Error[S]) {
	var i *Iterator[S]
	switch any(s).(type) {
	case string:
		x := iteratorPoolString.Get()
		defer iteratorPoolString.Put(x)
		i = x.(*Iterator[S])
	case []byte:
		x := iteratorPoolBytes.Get()
		defer iteratorPoolBytes.Put(x)
		i = x.(*Iterator[S])
	default:
		i = newIterator[S]()
	}
	i.src = s
	reset(i)
	t, err := scan(i, fn)
	if err.IsErr() {
		return err
	}
	var illegalChar bool
	t, illegalChar = strfind.EndOfWhitespaceSeq(t)
	if illegalChar {
		return getError(ErrorCodeIllegalControlChar, s, t)
	}
	if len(t) > 0 {
		return getError(ErrorCodeUnexpectedToken, s, t)
	}
	return Error[S]{}
}

// Parser wraps an iterator in a reusable instance.
// Reusing a parser instance is more efficient than global functions
// that rely on a global iterator pool.
type Parser[S ~string | ~[]byte] struct{ i *Iterator[S] }

// NewParser creates a new reusable parser instance.
// A higher preallocStackFrames value implies greater memory usage but also reduces
// the chance of dynamic memory allocations if the JSON depth surpasses the stack size.
// preallocStackFrames of 32 is equivalent to ~1KiB of memory usage on 64-bit systems
// (1 frame = ~32 bytes).
// Use DefaultStackSizeIterator when not sure.
func NewParser[S ~string | ~[]byte](preallocStackFrames int) *Parser[S] {
	i := &Iterator[S]{stack: make([]stackNode, preallocStackFrames)}
	reset(i)
	return &Parser[S]{i: i}
}

// ScanOne calls fn for every encountered value including objects and arrays.
// When an object or array is encountered fn will also be called for each of its
// member and element values.
//
// Unlike Scan, ScanOne doesn't return ErrorCodeUnexpectedToken when
// it encounters anything other than EOF after reading a valid JSON value.
// Returns an error if any and trailing as substring of s with the scanned value cut.
// In case of an error trailing will be a substring of s cut up until the index
// where the error was encountered.
//
// WARNING: Don't use or alias *Iterator[S] after fn returns!
func (p *Parser[S]) ScanOne(
	s S, fn func(*Iterator[S]) (err bool),
) (trailing S, err Error[S]) {
	reset(p.i)
	p.i.src = s
	return scan(p.i, fn)
}

// Scan calls fn for every encountered value including objects and arrays.
// When an object or array is encountered fn will also be called for each of its
// member and element values.
//
// WARNING: Don't use or alias *Iterator[S] after fn returns!
func (p *Parser[S]) Scan(
	s S, fn func(*Iterator[S]) (err bool),
) Error[S] {
	reset(p.i)
	p.i.src = s

	t, err := scan(p.i, fn)
	if err.IsErr() {
		return err
	}
	var illegalChar bool
	t, illegalChar = strfind.EndOfWhitespaceSeq(t)
	if illegalChar {
		return getError(ErrorCodeIllegalControlChar, s, t)
	}
	if len(t) > 0 {
		return getError(ErrorCodeUnexpectedToken, s, t)
	}
	return Error[S]{}
}

// scan calls fn for every value encountered.
// Returns the remainder of i.src and an error if any is encountered.
func scan[S ~string | ~[]byte](
	i *Iterator[S], fn func(*Iterator[S]) (err bool),
) (S, Error[S]) {
	var (
		rollback S // Used as fallback for error report
		s        = i.src
		b        bool
		ks, ke   int
	)

VALUE:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, i.src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
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
		return s, getError(ErrorCodeIllegalControlChar, i.src, s)
	}
	return s, getError(ErrorCodeUnexpectedToken, i.src, s)

VALUE_OBJECT:
	i.valueType = ValueTypeObject
	i.valueIndex, i.valueIndexEnd = len(i.src)-len(s), -1
	s = s[1:]
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, i.src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
		}
	}
	ks, ke = i.keyIndex, i.keyIndexEnd

	{ // Invoke callback
		i.arrayIndex = -1
		if len(i.stack) != 0 &&
			i.stack[len(i.stack)-1].Type == stackNodeTypeArray {
			i.arrayIndex = i.stack[len(i.stack)-1].ArrLen
			i.stack[len(i.stack)-1].ArrLen++
		}
		if fn(i) {
			return s, i.getError(ErrorCodeCallback)
		}
		i.keyIndex = -1
	}

	if s[0] == '}' {
		s = s[1:]
		goto AFTER_VALUE
	}
	i.stack = append(i.stack, stackNode{
		Type:        stackNodeTypeObject,
		KeyIndex:    ks,
		KeyIndexEnd: ke,
	})
	goto OBJ_KEY

VALUE_ARRAY:
	i.valueType = ValueTypeArray
	i.valueIndex, i.valueIndexEnd = len(i.src)-len(s), -1
	s = s[1:]
	ks, ke = i.keyIndex, i.keyIndexEnd

	{ // Invoke callback
		i.arrayIndex = -1
		if len(i.stack) != 0 &&
			i.stack[len(i.stack)-1].Type == stackNodeTypeArray {
			i.arrayIndex = i.stack[len(i.stack)-1].ArrLen
			i.stack[len(i.stack)-1].ArrLen++
		}
		if fn(i) {
			return s, i.getError(ErrorCodeCallback)
		}
		i.keyIndex = -1
	}

	i.stack = append(i.stack, stackNode{
		Type:        stackNodeTypeArray,
		KeyIndex:    ks,
		KeyIndexEnd: ke,
	})
	goto VALUE_OR_ARR_TERM

VALUE_NUMBER:
	{
		i.valueIndex = len(i.src) - len(s)
		{
			rollback = s
			var rc jsonnum.ReturnCode
			if s, rc = jsonnum.ReadNumber(s); rc == jsonnum.ReturnCodeErr {
				return s, getError(ErrorCodeMalformedNumber, i.src, rollback)
			}
		}
		i.valueIndexEnd = len(i.src) - len(s)
		i.valueType = ValueTypeNumber

		{ // Invoke callback
			i.arrayIndex = -1
			if len(i.stack) != 0 &&
				i.stack[len(i.stack)-1].Type == stackNodeTypeArray {
				i.arrayIndex = i.stack[len(i.stack)-1].ArrLen
				i.stack[len(i.stack)-1].ArrLen++
			}
			if fn(i) {
				return s, i.getError(ErrorCodeCallback)
			}
			i.keyIndex = -1
		}
	}
	goto AFTER_VALUE

VALUE_STRING:
	s = s[1:]
	i.valueIndex = len(i.src) - len(s) - 1
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
			return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
		}
		switch s[0] {
		case '\\':
			if len(s) < 2 {
				s = s[1:]
				return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
			}
			if lutEscape[s[1]] == 1 {
				s = s[2:]
				continue
			}
			if s[1] != 'u' {
				return s, getError(ErrorCodeInvalidEscape, i.src, s)
			}
			if len(s) < 6 ||
				lutSX[s[5]] != 2 ||
				lutSX[s[4]] != 2 ||
				lutSX[s[3]] != 2 ||
				lutSX[s[2]] != 2 {
				return s, getError(ErrorCodeInvalidEscape, i.src, s)
			}
			s = s[5:]
		case '"':
			s = s[1:]
			i.valueIndexEnd = len(i.src) - len(s)
			i.valueType = ValueTypeString

			{ // Invoke callback
				i.arrayIndex = -1
				if len(i.stack) != 0 &&
					i.stack[len(i.stack)-1].Type == stackNodeTypeArray {
					i.arrayIndex = i.stack[len(i.stack)-1].ArrLen
					i.stack[len(i.stack)-1].ArrLen++
				}
				if fn(i) {
					return s, i.getError(ErrorCodeCallback)
				}
				i.keyIndex = -1
			}

			goto AFTER_VALUE
		default:
			if s[0] < 0x20 {
				return s, getError(ErrorCodeIllegalControlChar, i.src, s)
			}
			s = s[1:]
		}
	}

VALUE_NULL:
	if len(s) < 4 || string(s[:4]) != "null" {
		return s, getError(ErrorCodeUnexpectedToken, i.src, s)
	}
	i.valueType = ValueTypeNull
	i.valueIndex = len(i.src) - len(s)
	i.valueIndexEnd = len(i.src) - len(s) + len("null")
	s = s[len("null"):]

	{ // Invoke callback
		i.arrayIndex = -1
		if len(i.stack) != 0 &&
			i.stack[len(i.stack)-1].Type == stackNodeTypeArray {
			i.arrayIndex = i.stack[len(i.stack)-1].ArrLen
			i.stack[len(i.stack)-1].ArrLen++
		}
		if fn(i) {
			return s, i.getError(ErrorCodeCallback)
		}
		i.keyIndex = -1
	}

	goto AFTER_VALUE

VALUE_FALSE:
	if len(s) < 5 || string(s[:5]) != "false" {
		return s, getError(ErrorCodeUnexpectedToken, i.src, s)
	}
	i.valueType = ValueTypeFalse
	i.valueIndex = len(i.src) - len(s)
	i.valueIndexEnd = len(i.src) - len(s) + len("false")
	s = s[len("false"):]

	{ // Invoke callback
		i.arrayIndex = -1
		if len(i.stack) != 0 &&
			i.stack[len(i.stack)-1].Type == stackNodeTypeArray {
			i.arrayIndex = i.stack[len(i.stack)-1].ArrLen
			i.stack[len(i.stack)-1].ArrLen++
		}
		if fn(i) {
			return s, i.getError(ErrorCodeCallback)
		}
		i.keyIndex = -1
	}

	goto AFTER_VALUE

VALUE_TRUE:
	if s := s; len(s) < 4 || string(s[:4]) != "true" {
		return s, getError(ErrorCodeUnexpectedToken, i.src, s)
	}
	i.valueType = ValueTypeTrue
	i.valueIndex = len(i.src) - len(s)
	i.valueIndexEnd = len(i.src) - len(s) + len("true")
	s = s[len("true"):]

	{ // Invoke callback
		i.arrayIndex = -1
		if len(i.stack) != 0 &&
			i.stack[len(i.stack)-1].Type == stackNodeTypeArray {
			i.arrayIndex = i.stack[len(i.stack)-1].ArrLen
			i.stack[len(i.stack)-1].ArrLen++
		}
		if fn(i) {
			return s, i.getError(ErrorCodeCallback)
		}
		i.keyIndex = -1
	}

	goto AFTER_VALUE

OBJ_KEY:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, i.src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
		}
	}
	if s[0] != '"' {
		if s[0] < 0x20 {
			return s, getError(ErrorCodeIllegalControlChar, i.src, s)
		}
		return s, getError(ErrorCodeUnexpectedToken, i.src, s)
	}

	s = s[1:]

	i.valueIndex = len(i.src) - len(s) - 1
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
			return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
		}
		switch s[0] {
		case '\\':
			if len(s) < 2 {
				s = s[1:]
				return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
			}
			if lutEscape[s[1]] == 1 {
				s = s[2:]
				continue
			}
			if s[1] != 'u' {
				return s, getError(ErrorCodeInvalidEscape, i.src, s)
			}
			if len(s) < 6 ||
				lutSX[s[5]] != 2 ||
				lutSX[s[4]] != 2 ||
				lutSX[s[3]] != 2 ||
				lutSX[s[2]] != 2 {
				return s, getError(ErrorCodeInvalidEscape, i.src, s)
			}
			s = s[5:]
		case '"':
			s = s[1:]
			i.keyIndex, i.keyIndexEnd = i.valueIndex, len(i.src)-len(s)
			goto AFTER_OBJ_KEY_STRING
		default:
			if s[0] < 0x20 {
				return s, getError(ErrorCodeIllegalControlChar, i.src, s)
			}
			s = s[1:]
		}
	}
AFTER_OBJ_KEY_STRING:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, i.src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
		}
	}
	if s[0] != ':' {
		if s[0] < 0x20 {
			return s, getError(ErrorCodeIllegalControlChar, i.src, s)
		}
		return s, getError(ErrorCodeUnexpectedToken, i.src, s)
	}
	s = s[1:]
	goto VALUE

VALUE_OR_ARR_TERM:
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, i.src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
		}
	}
	switch s[0] {
	case ']':
		s = s[1:]
		i.stack = i.stack[:len(i.stack)-1]
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
		return s, getError(ErrorCodeIllegalControlChar, i.src, s)
	}
	return s, getError(ErrorCodeUnexpectedToken, i.src, s)

AFTER_VALUE:
	if len(i.stack) == 0 {
		return s, Error[S]{}
	}
	if len(s) < 1 {
		return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
	}
	if s[0] <= ' ' {
		switch s[0] {
		case ' ', '\t', '\r', '\n':
			s, b = strfind.EndOfWhitespaceSeq(s)
			if b {
				return s, getError(ErrorCodeIllegalControlChar, i.src, s)
			}
		}
		if len(s) < 1 {
			return s, getError(ErrorCodeUnexpectedEOF, i.src, s)
		}
	}
	switch s[0] {
	case ',':
		s = s[1:]
		if i.stack[len(i.stack)-1].Type == stackNodeTypeArray {
			goto VALUE
		}
		goto OBJ_KEY
	case '}':
		if i.stack[len(i.stack)-1].Type != stackNodeTypeObject {
			return s, getError(ErrorCodeUnexpectedToken, i.src, s)
		}
		s = s[1:]
		i.stack = i.stack[:len(i.stack)-1]
		i.keyIndex, i.keyIndexEnd = -1, -1
		goto AFTER_VALUE
	case ']':
		if i.stack[len(i.stack)-1].Type != stackNodeTypeArray {
			return s, getError(ErrorCodeUnexpectedToken, i.src, s)
		}
		s = s[1:]
		i.stack = i.stack[:len(i.stack)-1]
		i.keyIndex, i.keyIndexEnd = -1, -1
		goto AFTER_VALUE
	}
	if s[0] < 0x20 {
		return s, getError(ErrorCodeIllegalControlChar, i.src, s)
	}
	return s, getError(ErrorCodeUnexpectedToken, i.src, s)
}
