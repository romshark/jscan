package jscan

import (
	"strconv"
	"sync"
	"unicode/utf8"

	"github.com/romshark/jscan/internal/jsonnum"
	"github.com/romshark/jscan/internal/keyescape"
	"github.com/romshark/jscan/internal/stack"
	"github.com/romshark/jscan/internal/strfind"
)

// Iterator provides access to the recently scanned value.
type Iterator struct {
	st         *stack.Stack
	src        string
	escapePath bool
	cachedPath []byte
	expect     expectation

	ValueType                       ValueType
	Level                           int
	KeyStart, KeyEnd, KeyLenEscaped int
	ValueStart, ValueEnd            int
	ArrayIndex                      int
}

// Key returns the field key if any.
func (i *Iterator) Key() string {
	if i.KeyStart < 0 {
		return ""
	}
	return i.src[i.KeyStart:i.KeyEnd]
}

// Value returns the value if any.
func (i *Iterator) Value() string {
	if i.ValueEnd < 0 || i.ValueEnd < i.ValueStart {
		return ""
	}
	return i.src[i.ValueStart:i.ValueEnd]
}

// ScanPath calls fn for every element in the path of the value.
// If keyStart is != -1 then the element is a field value, otherwise
// arrIndex indicates the index of the item in the underlying array.
func (i *Iterator) ScanPath(fn func(keyStart, keyEnd, arrIndex int)) {
	for j := i.st.Len() - 1; j >= 0; j-- {
		t := i.st.TopOffset(j)
		if t.KeyStart > -1 {
			fn(t.KeyStart, t.KeyEnd, -1)
		}
		if t.Type == stack.NodeTypeArray {
			fn(-1, -1, t.ArrLen-1)
		}
	}
}

// Path returns the stringified path.
func (i *Iterator) Path() (s string) {
	i.ViewPath(func(p []byte) { s = string(p) })
	return
}

// ViewPath calls fn and provides the stringified path.
//
// WARNING: do not use or alias p after fn returns!
// Only viewing or copying are considered safe!
// Use (*Iterator).Path instead for a safer and more convenient API.
func (i *Iterator) ViewPath(fn func(p []byte)) {
	if len(i.cachedPath) > 0 {
		// The path is already cached
		fn(i.cachedPath[1:])
		return
	}

	var b []byte
	needDelim := false
	i.ScanPath(func(keyStart, keyEnd, arrIndex int) {
		if keyStart != -1 {
			if needDelim {
				b = append(b, '.')
			}
			k := i.src[keyStart:keyEnd]
			if i.escapePath {
				k = keyescape.Escape(k)
			}
			b = append(b, k...)
			needDelim = true
		} else {
			b = append(b, '[')
			b = strconv.AppendInt(b, int64(arrIndex), 10)
			b = append(b, ']')
			needDelim = true
		}
	})
	if i.KeyStart != -1 {
		if needDelim {
			b = append(b, '.')
		}
		k := i.src[i.KeyStart:i.KeyEnd]
		if i.escapePath {
			k = keyescape.Escape(k)
		}
		b = append(b, k...)
	}
	fn(b)
}

var itrPool = sync.Pool{
	New: func() interface{} {
		return &Iterator{
			st: stack.New(64),
		}
	},
}

func getItrFromPool(
	s string,
	escapePath bool,
	startIndex int,
) *Iterator {
	i := itrPool.Get().(*Iterator)
	i.st.Reset()
	i.escapePath = escapePath
	i.cachedPath = i.cachedPath[:0]
	i.src = s
	i.ValueType = 0
	i.Level = 0
	i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1
	i.ValueStart, i.ValueEnd = startIndex, -1
	i.ArrayIndex = 0
	i.expect = expectVal
	return i
}

// getError returns the stringified error, if any.
func (i *Iterator) getError(c ErrorCode) Error {
	return Error{
		Code:  c,
		Src:   i.src,
		Index: i.ValueStart,
	}
}

type Error struct {
	Src   string
	Index int
	Code  ErrorCode
}

// IsErr returns true if there is an error, otherwise returns false.
func (e Error) IsErr() bool { return e.Code != 0 }

func (e Error) Error() string {
	if e.Index < len(e.Src) {
		r, _ := utf8.DecodeRuneInString(e.Src[e.Index:])
		return errorMessage(e.Code, e.Index, r)
	}
	return errorMessage(e.Code, e.Index, 0)
}

// Get calls fn for the value at the given path.
// The value path is defined by keys separated by a dot and
// index access operators for arrays.
// If no value is found for the given path then fn isn't called
// and no error is returned.
// If escapePath then all dots and square brackets are expected to be escaped.
//
// WARNING: Fields exported by *Iterator in fn must not be mutated!
// Do not use or alias *Iterator after fn returns!
func Get(s, path string, escapePath bool, fn func(*Iterator)) Error {
	err := Scan(Options{
		CachePath:  true,
		EscapePath: escapePath,
	}, s, func(i *Iterator) (err bool) {
		i.ViewPath(func(p []byte) {
			if string(p) != string(path) {
				return
			}

			fn(i)
			err = true
		})
		return
	})
	if err.Code == ErrorCodeCallback {
		err.Code = 0
	}
	return err
}

// ValidateOne scans a JSON value from s and returns an error if it's invalid,
// otherwise returns s with the scanned value cut.
func ValidateOne(s string) (trailing string, err Error) {
	i := getItrFromPool(s, false, 0)
	defer itrPool.Put(i)
	return i.validate(s)
}

// Validate returns an error if s is invalid JSON.
func Validate(s string) Error {
	i := getItrFromPool(s, false, 0)
	defer itrPool.Put(i)

	trailing, err := i.validate(s)
	if err.IsErr() {
		return err
	}
	end, illegalChar := strfind.EndOfWhitespaceSeq(trailing)
	if illegalChar {
		i.ValueStart += end
		return i.getError(ErrorCodeIllegalControlChar)
	}
	if len(trailing[end:]) > 0 {
		i.ValueStart += end
		return i.getError(ErrorCodeUnexpectedToken)
	}
	return Error{}
}

// validate returns an error if s is invalid JSON.
func (i *Iterator) validate(s string) (string, Error) {
	startIndex, illegal := strfind.EndOfWhitespaceSeq(s)
	if illegal {
		return "", Error{
			Src:   s,
			Index: startIndex,
			Code:  ErrorCodeIllegalControlChar,
		}
	}

	if startIndex >= len(s) {
		return "", Error{
			Src:   s,
			Index: startIndex,
			Code:  ErrorCodeUnexpectedEOF,
		}
	}

	for i.ValueStart < len(s) {
		switch s[i.ValueStart] {
		case ' ', '\t', '\r', '\n':
			e, illegal := strfind.EndOfWhitespaceSeq(s[i.ValueStart:])
			i.ValueStart += e
			if illegal {
				return "", i.getError(ErrorCodeIllegalControlChar)
			}

		case ',':
			switch i.expect {
			case expectCommaOrArrTerm:
				i.expect = expectVal
			case expectCommaOrObjTerm:
				i.expect = expectKey
			default:
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			i.ValueStart++

		case '}':
			if i.expect != expectCommaOrObjTerm &&
				i.expect != expectKeyOrObjTerm {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			i.ValueStart++

			i.st.Pop()
			t := i.st.Top()
			if t == nil {
				return s[i.ValueStart:], Error{}
			}
			switch t.Type {
			case stack.NodeTypeArray:
				i.expect = expectCommaOrArrTerm
			case stack.NodeTypeObject:
				i.expect = expectCommaOrObjTerm
			}

			i.KeyStart, i.KeyEnd = -1, -1

		case ']':
			if i.expect != expectCommaOrArrTerm &&
				i.expect != expectValOrArrTerm {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			i.ValueStart++

			i.st.Pop()
			t := i.st.Top()
			if t == nil {
				return s[i.ValueStart:], Error{}
			}
			switch t.Type {
			case stack.NodeTypeArray:
				i.expect = expectCommaOrArrTerm
			case stack.NodeTypeObject:
				i.expect = expectCommaOrObjTerm
			}

			i.KeyStart, i.KeyEnd = -1, -1

		case '{': // Object
			if i.expect != expectVal && i.expect != expectValOrArrTerm {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			i.expect = expectKeyOrObjTerm

			i.KeyStart, i.KeyEnd = -1, -1

			i.st.Push(stack.NodeTypeObject, 0, 0, 0)
			i.ValueStart++

		case '[': // Array
			if i.expect != expectVal && i.expect != expectValOrArrTerm {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			i.expect = expectValOrArrTerm

			i.KeyStart, i.KeyEnd = -1, -1

			i.st.Push(stack.NodeTypeArray, 0, 0, 0)
			i.ValueStart++

		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if i.expect != expectVal && i.expect != expectValOrArrTerm {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}

			var err bool
			i.ValueEnd, err = jsonnum.EndIndex(i.src[i.ValueStart:])
			if err {
				return "", i.getError(ErrorCodeMalformedNumber)
			}
			if _, err := strconv.ParseFloat(
				i.src[i.ValueStart:i.ValueStart+i.ValueEnd], 64,
			); err != nil {
				return "", i.getError(ErrorCodeMalformedNumber)
			}
			i.ValueStart += i.ValueEnd

			t := i.st.Top()
			if t == nil {
				return s[i.ValueStart:], Error{}
			}
			switch t.Type {
			case stack.NodeTypeArray:
				i.expect = expectCommaOrArrTerm
			case stack.NodeTypeObject:
				i.expect = expectCommaOrObjTerm
			}

			i.KeyStart, i.KeyEnd = -1, -1

		case '"': // String
			i.ValueStart++
			var errCode strfind.ErrCode
			i.ValueEnd, errCode = strfind.IndexTerm(s, i.ValueStart)
			if errCode > strfind.ErrCodeOK {
				i.ValueStart = i.ValueEnd
				return "", i.getError(ErrorCode(errCode))
			}
			i.ValueStart = i.ValueEnd + 1
			t := i.st.Top()
			if t != nil && t.Type == stack.NodeTypeArray {
				// Array item string value
				if i.expect != expectVal && i.expect != expectValOrArrTerm {
					return "", i.getError(ErrorCodeUnexpectedToken)
				}
				i.expect = expectCommaOrArrTerm
				i.KeyStart, i.KeyEnd = -1, -1
			} else if i.KeyStart != -1 || i.st.Len() == 0 {
				// String value
				if i.expect != expectVal {
					return "", i.getError(ErrorCodeUnexpectedToken)
				}
				if t == nil {
					return s[i.ValueStart:], Error{}
				}
				i.expect = expectCommaOrObjTerm
				i.KeyStart, i.KeyEnd = -1, -1
			} else {
				// Key
				if i.expect != expectKey && i.expect != expectKeyOrObjTerm {
					return "", i.getError(ErrorCodeUnexpectedToken)
				}
				i.expect = expectVal

				i.KeyStart, i.KeyEnd = i.ValueStart, i.ValueEnd
				if i.ValueEnd > len(s) {
					return "", i.getError(ErrorCodeUnexpectedEOF)
				}
				if x, err := i.parseColumn(s[i.ValueStart:]); err {
					return "", i.getError(ErrorCodeUnexpectedToken)
				} else {
					i.ValueStart += x + 1
				}
			}

		case 'n': // Null
			if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
				i.expect3(i.ValueStart+1, 'u', 'l', 'l') {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			i.ValueStart += len("null")

			t := i.st.Top()
			if t == nil {
				return s[i.ValueStart:], Error{}
			}
			switch t.Type {
			case stack.NodeTypeArray:
				i.expect = expectCommaOrArrTerm
			case stack.NodeTypeObject:
				i.expect = expectCommaOrObjTerm
			}

			i.KeyStart, i.KeyEnd = -1, -1

		case 'f': // False
			if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
				i.read4(i.ValueStart+1, 'a', 'l', 's', 'e') {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			i.ValueStart += len("false")

			t := i.st.Top()
			if t == nil {
				return s[i.ValueStart:], Error{}
			}
			switch t.Type {
			case stack.NodeTypeArray:
				i.expect = expectCommaOrArrTerm
			case stack.NodeTypeObject:
				i.expect = expectCommaOrObjTerm
			}

			i.KeyStart, i.KeyEnd = -1, -1

		case 't': // True
			if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
				i.expect3(i.ValueStart+1, 'r', 'u', 'e') {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			i.ValueStart += len("true")

			t := i.st.Top()
			if t == nil {
				return s[i.ValueStart:], Error{}
			}
			switch t.Type {
			case stack.NodeTypeArray:
				i.expect = expectCommaOrArrTerm
			case stack.NodeTypeObject:
				i.expect = expectCommaOrObjTerm
			}

			i.KeyStart, i.KeyEnd = -1, -1

		default:
			if s[i.ValueStart] < 0x20 {
				return "", i.getError(ErrorCodeIllegalControlChar)
			}
			return "", i.getError(ErrorCodeUnexpectedToken)
		}
	}

	if i.st.Len() > 0 {
		return "", i.getError(ErrorCodeUnexpectedEOF)
	}

	return s[i.ValueStart:], Error{}
}

// Valid returns true if s is valid JSON, otherwise returns false.
func Valid(s string) bool {
	return !Validate(s).IsErr()
}

// ScanOne calls fn for every scanned value including objects and arrays.
// When an object or array is encountered fn will also be called for each
// field value and array item respectively.
// Unlike Scan, ScanOne doesn't return ErrorCodeUnexpectedToken when
// it encounters anything other than EOF after reading a valid JSON value.
// Returns s with the scanned value cut.
//
// WARNING: Fields exported by *Iterator in fn must not be mutated!
// Do not use or alias *Iterator after fn returns!
func ScanOne(
	o Options,
	s string,
	fn func(*Iterator) (err bool),
) (trailing string, err Error) {
	startIndex, illegal := strfind.EndOfWhitespaceSeq(s)
	if illegal {
		return "", Error{
			Src:   s,
			Index: startIndex,
			Code:  ErrorCodeIllegalControlChar,
		}
	}

	if startIndex >= len(s) {
		return "", Error{
			Src:   s,
			Index: startIndex,
			Code:  ErrorCodeUnexpectedEOF,
		}
	}

	i := getItrFromPool(s, o.EscapePath, startIndex)
	defer itrPool.Put(i)

	if o.CachePath {
		return i.scanWithCachedPath(s, fn)
	}
	return i.scan(s, fn)
}

// Scan calls fn for every scanned value including objects and arrays.
// When an object or array is encountered fn will also be called for each
// field value and array item respectively.
//
// WARNING: Fields exported by *Iterator in fn must not be mutated!
// Do not use or alias *Iterator after fn returns!
func Scan(
	o Options,
	s string,
	fn func(*Iterator) (err bool),
) (err Error) {
	startIndex, illegal := strfind.EndOfWhitespaceSeq(s)
	if illegal {
		return Error{
			Src:   s,
			Index: startIndex,
			Code:  ErrorCodeIllegalControlChar,
		}
	}

	if startIndex >= len(s) {
		return Error{
			Src:   s,
			Index: startIndex,
			Code:  ErrorCodeUnexpectedEOF,
		}
	}

	i := getItrFromPool(s, o.EscapePath, startIndex)
	defer itrPool.Put(i)

	var trailing string
	if o.CachePath {
		trailing, err = i.scanWithCachedPath(s, fn)
	} else {
		trailing, err = i.scan(s, fn)
	}
	if err.IsErr() {
		return err
	}
	end, illegalChar := strfind.EndOfWhitespaceSeq(trailing)
	if illegalChar {
		i.ValueStart += end
		return i.getError(ErrorCodeIllegalControlChar)
	}
	if len(trailing[end:]) > 0 {
		i.ValueStart += end
		return i.getError(ErrorCodeUnexpectedToken)
	}
	return Error{}
}

func (i *Iterator) scan(
	s string,
	fn func(*Iterator) (err bool),
) (trailing string, err Error) {
	for i.ValueStart < len(s) {
		switch s[i.ValueStart] {
		case ' ', '\t', '\r', '\n':
			e, illegal := strfind.EndOfWhitespaceSeq(s[i.ValueStart:])
			i.ValueStart += e
			if illegal {
				return "", i.getError(ErrorCodeIllegalControlChar)
			}

		case ',':
			switch i.expect {
			case expectCommaOrArrTerm:
				i.expect = expectVal
			case expectCommaOrObjTerm:
				i.expect = expectKey
			default:
				return "", i.getError(ErrorCodeUnexpectedToken)
			}

			i.ValueStart++

		case '}':
			if i.expect != expectCommaOrObjTerm &&
				i.expect != expectKeyOrObjTerm {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			i.ValueStart++

			i.st.Pop()
			t := i.st.Top()
			if t == nil {
				return s[i.ValueStart:], Error{}
			}
			switch t.Type {
			case stack.NodeTypeArray:
				i.expect = expectCommaOrArrTerm
			case stack.NodeTypeObject:
				i.expect = expectCommaOrObjTerm
			}

			i.Level--
			i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1

		case ']':
			if i.expect != expectCommaOrArrTerm &&
				i.expect != expectValOrArrTerm {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			i.ValueStart++

			i.st.Pop()
			t := i.st.Top()
			if t == nil {
				return s[i.ValueStart:], Error{}
			}
			switch t.Type {
			case stack.NodeTypeArray:
				i.expect = expectCommaOrArrTerm
			case stack.NodeTypeObject:
				i.expect = expectCommaOrObjTerm
			}

			i.Level--
			i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1

		case '{': // Object
			if i.expect != expectVal && i.expect != expectValOrArrTerm {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			i.expect = expectKeyOrObjTerm

			i.ValueType = ValueTypeObject
			i.ValueEnd = -1
			ks, ke := i.KeyStart, i.KeyEnd
			if i.callFn(i.st.Top(), fn) {
				return "", i.getError(ErrorCodeCallback)
			}
			i.st.Push(stack.NodeTypeObject, 0, ks, ke)
			i.Level++
			i.ValueStart++

		case '[': // Array
			if i.expect != expectVal && i.expect != expectValOrArrTerm {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			i.expect = expectValOrArrTerm

			i.ValueType = ValueTypeArray
			i.ValueEnd = -1
			ks, ke := i.KeyStart, i.KeyEnd
			if i.callFn(i.st.Top(), fn) {
				return "", i.getError(ErrorCodeCallback)
			}
			i.st.Push(stack.NodeTypeArray, 0, ks, ke)
			i.Level++
			i.ValueStart++

		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if i.expect != expectVal && i.expect != expectValOrArrTerm {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}

			t := i.st.Top()

			var err bool
			i.ValueEnd, err = jsonnum.EndIndex(i.src[i.ValueStart:])
			if err {
				return "", i.getError(ErrorCodeMalformedNumber)
			}
			i.ValueEnd += +i.ValueStart

			i.ValueType = ValueTypeNumber
			if i.callFn(t, fn) {
				return "", i.getError(ErrorCodeCallback)
			}
			i.ValueStart = i.ValueEnd

			if t == nil {
				return s[i.ValueStart:], Error{}
			}
			switch t.Type {
			case stack.NodeTypeArray:
				i.expect = expectCommaOrArrTerm
			case stack.NodeTypeObject:
				i.expect = expectCommaOrObjTerm
			}

		case '"': // String
			i.ValueStart++
			i.ArrayIndex = -1
			var errCode strfind.ErrCode
			i.ValueEnd, errCode = strfind.IndexTerm(s, i.ValueStart)
			if errCode > strfind.ErrCodeOK {
				i.ValueStart = i.ValueEnd
				return "", i.getError(ErrorCode(errCode))
			}
			t := i.st.Top()
			if t != nil && t.Type == stack.NodeTypeArray {
				// Array item string value
				if i.expect != expectVal && i.expect != expectValOrArrTerm {
					i.ValueStart--
					return "", i.getError(ErrorCodeUnexpectedToken)
				}
				i.expect = expectCommaOrArrTerm

				i.ArrayIndex = t.ArrLen
				t.ArrLen++
				i.ValueType = ValueTypeString
				i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1
				if fn(i) {
					i.ValueStart--
					return "", i.getError(ErrorCodeCallback)
				}
				i.ValueStart = i.ValueEnd + 1
			} else if i.KeyStart != -1 || i.Level == 0 {
				// String value
				if i.expect != expectVal {
					i.ValueStart--
					return "", i.getError(ErrorCodeUnexpectedToken)
				}

				i.ValueType = ValueTypeString
				if i.callFn(nil, fn) {
					i.ValueStart--
					return "", i.getError(ErrorCodeCallback)
				}
				i.ValueStart = i.ValueEnd + 1

				if t == nil {
					return s[i.ValueStart:], Error{}
				}
				i.expect = expectCommaOrObjTerm
			} else {
				// Key
				if i.expect != expectKey && i.expect != expectKeyOrObjTerm {
					i.ValueStart--
					return "", i.getError(ErrorCodeUnexpectedToken)
				}
				i.expect = expectVal

				i.KeyStart, i.KeyEnd = i.ValueStart, i.ValueEnd
				i.ValueStart = i.ValueEnd + 1
				if i.ValueEnd > len(s) {
					return "", i.getError(ErrorCodeUnexpectedEOF)
				}
				if x, err := i.parseColumn(s[i.ValueStart:]); err {
					return "", i.getError(ErrorCodeUnexpectedToken)
				} else {
					i.ValueStart += x + 1
				}
			}

		case 'n': // Null
			if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
				i.expect3(i.ValueStart+1, 'u', 'l', 'l') {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			t := i.st.Top()

			i.ValueType = ValueTypeNull
			i.ValueEnd = i.ValueStart + len("null")
			if i.callFn(t, fn) {
				return "", i.getError(ErrorCodeCallback)
			}
			i.ValueStart += len("null")

			if t == nil {
				return s[i.ValueStart:], Error{}
			}
			switch t.Type {
			case stack.NodeTypeArray:
				i.expect = expectCommaOrArrTerm
			case stack.NodeTypeObject:
				i.expect = expectCommaOrObjTerm
			}

		case 'f': // False
			if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
				i.read4(i.ValueStart+1, 'a', 'l', 's', 'e') {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			t := i.st.Top()

			i.ValueType = ValueTypeFalse
			i.ValueEnd = i.ValueStart + len("false")
			if i.callFn(t, fn) {
				return "", i.getError(ErrorCodeCallback)
			}
			i.ValueStart += len("false")

			if t == nil {
				return s[i.ValueStart:], Error{}
			}
			switch t.Type {
			case stack.NodeTypeArray:
				i.expect = expectCommaOrArrTerm
			case stack.NodeTypeObject:
				i.expect = expectCommaOrObjTerm
			}

		case 't': // True
			if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
				i.expect3(i.ValueStart+1, 'r', 'u', 'e') {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			t := i.st.Top()

			i.ValueType = ValueTypeTrue
			i.ValueEnd = i.ValueStart + len("true")
			if i.callFn(t, fn) {
				return "", i.getError(ErrorCodeCallback)
			}
			i.ValueStart += len("true")

			if t == nil {
				return s[i.ValueStart:], Error{}
			}
			switch t.Type {
			case stack.NodeTypeArray:
				i.expect = expectCommaOrArrTerm
			case stack.NodeTypeObject:
				i.expect = expectCommaOrObjTerm
			}

		default:
			if s[i.ValueStart] < 0x20 {
				return "", i.getError(ErrorCodeIllegalControlChar)
			}
			return "", i.getError(ErrorCodeUnexpectedToken)
		}
	}

	if i.st.Len() > 0 {
		return "", i.getError(ErrorCodeUnexpectedEOF)
	}

	return s[i.ValueStart:], Error{}
}

func (i *Iterator) expect3(at int, a, b, c byte) (err bool) {
	return len(i.src)-at < 3 ||
		i.src[at] != a ||
		i.src[at+1] != b ||
		i.src[at+2] != c
}

func (i *Iterator) read4(at int, a, b, c, d byte) (err bool) {
	return len(i.src)-at < 4 ||
		i.src[at] != a ||
		i.src[at+1] != b ||
		i.src[at+2] != c ||
		i.src[at+3] != d
}

func (i *Iterator) scanWithCachedPath(
	s string,
	fn func(*Iterator) bool,
) (trailing string, err Error) {
	i.cachedPath = append(i.cachedPath, '.')

	for i.ValueStart < len(s) {
		switch s[i.ValueStart] {
		case ' ', '\t', '\r', '\n':
			e, illegal := strfind.EndOfWhitespaceSeq(s[i.ValueStart:])
			i.ValueStart += e
			if illegal {
				return "", i.getError(ErrorCodeIllegalControlChar)
			}

		case ',':
			switch i.expect {
			case expectCommaOrArrTerm:
				i.expect = expectVal
			case expectCommaOrObjTerm:
				i.expect = expectKey
			default:
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			i.ValueStart++

		case '}':
			if i.expect != expectCommaOrObjTerm &&
				i.expect != expectKeyOrObjTerm {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			i.ValueStart++

			i.st.Pop()
			t := i.st.Top()
			if t == nil {
				return s[i.ValueStart:], Error{}
			}
			switch t.Type {
			case stack.NodeTypeArray:
				i.expect = expectCommaOrArrTerm
			case stack.NodeTypeObject:
				i.expect = expectCommaOrObjTerm
			}

			i.Level--
			i.cacheCleanupAfterContainer()
			i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1

		case ']':
			if i.expect != expectCommaOrArrTerm &&
				i.expect != expectValOrArrTerm {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			i.ValueStart++

			i.st.Pop()
			t := i.st.Top()
			if t == nil {
				return s[i.ValueStart:], Error{}
			}
			switch t.Type {
			case stack.NodeTypeArray:
				i.expect = expectCommaOrArrTerm
			case stack.NodeTypeObject:
				i.expect = expectCommaOrObjTerm
			}

			i.Level--
			if len(i.cachedPath) > 1 {
				i.cachedPath = i.cachedPath[:len(i.cachedPath)-1]
				i.cacheCleanupAfterContainer()
			}
			i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1

		case '{': // Object
			if i.expect != expectVal && i.expect != expectValOrArrTerm {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			i.expect = expectKeyOrObjTerm

			i.ValueType = ValueTypeObject
			i.ValueEnd = -1
			ks, ke := i.KeyStart, i.KeyEnd
			if i.callFnWithPathCache(i.st.Top(), false, fn) {
				return "", i.getError(ErrorCodeCallback)
			}
			i.st.Push(stack.NodeTypeObject, 0, ks, ke)
			i.Level++
			i.ValueStart++

		case '[': // Array
			if i.expect != expectVal && i.expect != expectValOrArrTerm {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			i.expect = expectValOrArrTerm

			i.ValueType = ValueTypeArray
			i.ValueEnd = -1
			ks, ke := i.KeyStart, i.KeyEnd
			if i.callFnWithPathCache(i.st.Top(), false, fn) {
				return "", i.getError(ErrorCodeCallback)
			}
			i.st.Push(stack.NodeTypeArray, 0, ks, ke)
			i.Level++
			i.ValueStart++
			i.cachedPath = append(i.cachedPath, '[')

		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if i.expect != expectVal && i.expect != expectValOrArrTerm {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}

			t := i.st.Top()

			var err bool
			i.ValueEnd, err = jsonnum.EndIndex(i.src[i.ValueStart:])
			if err {
				return "", i.getError(ErrorCodeMalformedNumber)
			}
			i.ValueEnd += +i.ValueStart

			i.ValueType = ValueTypeNumber
			if i.callFnWithPathCache(t, true, fn) {
				return "", i.getError(ErrorCodeCallback)
			}
			i.ValueStart = i.ValueEnd

			if t == nil {
				return s[i.ValueStart:], Error{}
			}
			switch t.Type {
			case stack.NodeTypeArray:
				i.expect = expectCommaOrArrTerm
			case stack.NodeTypeObject:
				i.expect = expectCommaOrObjTerm
			}

		case '"': // String
			i.ValueStart++
			i.ArrayIndex = -1
			var errCode strfind.ErrCode
			i.ValueEnd, errCode = strfind.IndexTerm(s, i.ValueStart)
			if errCode > strfind.ErrCodeOK {
				i.ValueStart = i.ValueEnd
				return "", i.getError(ErrorCode(errCode))
			}
			t := i.st.Top()
			if t != nil && t.Type == stack.NodeTypeArray {
				// Array item string value
				if i.expect != expectVal && i.expect != expectValOrArrTerm {
					i.ValueStart--
					return "", i.getError(ErrorCodeUnexpectedToken)
				}
				i.expect = expectCommaOrArrTerm

				i.ArrayIndex = t.ArrLen
				t.ArrLen++
				i.ValueType = ValueTypeString
				i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1

				initArrIndex := len(i.cachedPath)
				i.cachedPath = strconv.AppendInt(
					i.cachedPath, int64(i.ArrayIndex), 10,
				)
				i.cachedPath = append(i.cachedPath, ']')

				if fn(i) {
					i.ValueStart--
					return "", i.getError(ErrorCodeCallback)
				}
				// Remove index
				i.cachedPath = i.cachedPath[:initArrIndex]

				i.ValueStart = i.ValueEnd + 1
			} else if i.KeyStart != -1 || i.Level == 0 {
				// String value
				if i.expect != expectVal {
					i.ValueStart--
					return "", i.getError(ErrorCodeUnexpectedToken)
				}

				i.ValueType = ValueTypeString
				if i.callFnWithPathCache(nil, true, fn) {
					i.ValueStart--
					return "", i.getError(ErrorCodeCallback)
				}
				i.ValueStart = i.ValueEnd + 1

				if t == nil {
					return s[i.ValueStart:], Error{}
				}
				i.expect = expectCommaOrObjTerm
			} else {
				// Key
				if i.expect != expectKey && i.expect != expectKeyOrObjTerm {
					i.ValueStart--
					return "", i.getError(ErrorCodeUnexpectedToken)
				}
				i.expect = expectVal

				i.KeyStart, i.KeyEnd = i.ValueStart, i.ValueEnd
				i.ValueStart = i.ValueEnd + 1
				if i.ValueEnd > len(s) {
					return "", i.getError(ErrorCodeUnexpectedEOF)
				}
				if x, err := i.parseColumn(s[i.ValueStart:]); err {
					return "", i.getError(ErrorCodeUnexpectedToken)
				} else {
					i.ValueStart += x + 1
				}

				if len(i.cachedPath) > 1 {
					i.cachedPath = append(i.cachedPath, '.')
				}
				var e string
				if i.escapePath {
					e = keyescape.Escape(s[i.KeyStart:i.KeyEnd])
				} else {
					e = s[i.KeyStart:i.KeyEnd]
				}
				i.KeyLenEscaped = len(e)
				i.cachedPath = append(i.cachedPath, e...)
			}

		case 'n': // Null
			if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
				i.expect3(i.ValueStart+1, 'u', 'l', 'l') {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			t := i.st.Top()

			i.ValueType = ValueTypeNull
			i.ValueEnd = i.ValueStart + len("null")
			if i.callFnWithPathCache(t, true, fn) {
				return "", i.getError(ErrorCodeCallback)
			}
			i.ValueStart += len("null")

			if t == nil {
				return s[i.ValueStart:], Error{}
			}
			switch t.Type {
			case stack.NodeTypeArray:
				i.expect = expectCommaOrArrTerm
			case stack.NodeTypeObject:
				i.expect = expectCommaOrObjTerm
			}

		case 'f': // False
			if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
				i.read4(i.ValueStart+1, 'a', 'l', 's', 'e') {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			t := i.st.Top()

			i.ValueType = ValueTypeFalse
			i.ValueEnd = i.ValueStart + len("false")
			if i.callFnWithPathCache(t, true, fn) {
				return "", i.getError(ErrorCodeCallback)
			}
			i.ValueStart += len("false")

			if t == nil {
				return s[i.ValueStart:], Error{}
			}
			switch t.Type {
			case stack.NodeTypeArray:
				i.expect = expectCommaOrArrTerm
			case stack.NodeTypeObject:
				i.expect = expectCommaOrObjTerm
			}

		case 't': // True
			if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
				i.expect3(i.ValueStart+1, 'r', 'u', 'e') {
				return "", i.getError(ErrorCodeUnexpectedToken)
			}
			t := i.st.Top()

			i.ValueType = ValueTypeTrue
			i.ValueEnd = i.ValueStart + len("true")
			if i.callFnWithPathCache(t, true, fn) {
				return "", i.getError(ErrorCodeCallback)
			}
			i.ValueStart += len("true")

			if t == nil {
				return s[i.ValueStart:], Error{}
			}
			switch t.Type {
			case stack.NodeTypeArray:
				i.expect = expectCommaOrArrTerm
			case stack.NodeTypeObject:
				i.expect = expectCommaOrObjTerm
			}

		default:
			if s[i.ValueStart] < 0x20 {
				return "", i.getError(ErrorCodeIllegalControlChar)
			}
			return "", i.getError(ErrorCodeUnexpectedToken)
		}
	}

	if i.st.Len() > 0 {
		return "", i.getError(ErrorCodeUnexpectedEOF)
	}

	return s[i.ValueStart:], Error{}
}

func (i *Iterator) parseColumn(s string) (index int, err bool) {
	if len(s) > 0 && s[0] == ':' {
		return 0, false
	}
	for j, c := range s {
		if c == ':' {
			return j, false
		} else if !isSpace(c) {
			break
		}
	}
	return -1, true
}

func isSpace(b rune) bool {
	switch b {
	case ' ', '\r', '\n', '\t':
		return true
	}
	return false
}

func (i *Iterator) cacheCleanupAfterContainer() {
	if l := len(i.cachedPath) - 1; l > -1 {
		switch i.cachedPath[l] {
		case ']':
			if x := strfind.LastIndexUnescaped(i.cachedPath, '['); x > -1 {
				i.cachedPath = i.cachedPath[:x+1]
			}
		default:
			if x := strfind.LastIndexUnescaped(i.cachedPath, '.'); x > -1 {
				i.cachedPath = i.cachedPath[:x]
			}
		}
		if len(i.cachedPath) < 1 {
			i.cachedPath = append(i.cachedPath, '.')
		}
	}
}

func (i *Iterator) callFn(
	t *stack.Node,
	fn func(i *Iterator) (err bool),
) (err bool) {
	i.ArrayIndex = -1
	if t != nil && t.Type == stack.NodeTypeArray {
		i.ArrayIndex = t.ArrLen
		t.ArrLen++
	}

	if fn(i) {
		return true
	}

	i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1
	return false
}

func (i *Iterator) callFnWithPathCache(
	t *stack.Node,
	nonComposite bool,
	fn func(i *Iterator) bool,
) (err bool) {
	i.ArrayIndex = -1
	initArrIndex := -1
	if t != nil && t.Type == stack.NodeTypeArray {
		i.ArrayIndex = t.ArrLen
		t.ArrLen++
	}

	if l := len(i.cachedPath) - 1; l > 0 && i.cachedPath[l] == '[' {
		i.cachedPath = strconv.AppendInt(
			i.cachedPath, int64(i.ArrayIndex), 10,
		)
		i.cachedPath = append(i.cachedPath, ']')
		initArrIndex = l
	}

	if fn(i) {
		return true
	}

	if nonComposite {
		if initArrIndex > -1 {
			// Remove index
			i.cachedPath = i.cachedPath[:initArrIndex+1]
		} else if i.KeyLenEscaped > 0 {
			// Remove field key
			if i.KeyLenEscaped+1 < len(i.cachedPath) &&
				i.cachedPath[len(i.cachedPath)-(i.KeyLenEscaped+1)] == '.' {
				i.KeyLenEscaped += 1
			}
			i.cachedPath = i.cachedPath[:len(i.cachedPath)-i.KeyLenEscaped]
		}
	}
	i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1
	return false
}
