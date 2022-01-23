package jscan

import (
	"fmt"
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

// ErrorCode defines the error type.
type ErrorCode int8

// All error codes
const (
	_ ErrorCode = iota
	ErrorCodeUnexpectedToken
	ErrorCodeMalformedNumber
	ErrorCodeUnexpectedEOF
	ErrorCodeIllegalControlCharacter
	ErrorCallback
)

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

func getItrFromPool(s string, escapePath bool) *Iterator {
	i := itrPool.Get().(*Iterator)
	i.st.Reset()
	i.escapePath = escapePath
	i.cachedPath = i.cachedPath[:0]
	i.src = s
	i.ValueType = 0
	i.Level = 0
	i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1
	i.ValueStart, i.ValueEnd = 0, -1
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
// WARNING: Fields exported by *Iterator provided in fn must not be mutated!
// Do not use or alias *Iterator after fn returns!
func Get(s, path string, escapePath bool, fn func(*Iterator)) Error {
	err := Scan(Options{
		CachePath:  true,
		EscapePath: escapePath,
	}, s, func(i *Iterator) (err bool) {
		i.ViewPath(func(p []byte) {
			if len(p) != len(path) {
				return
			}
			for i := range p {
				if p[i] != path[i] {
					return
				}
			}

			fn(i)
			err = true
		})
		return
	})
	if err.Code == ErrorCallback {
		err.Code = 0
	}
	return err
}

// Validate returns an error if s is invalid JSON.
func Validate(s string) Error {
	if s == "" {
		return Error{Code: ErrorCodeUnexpectedEOF}
	}

	i := getItrFromPool(s, false)
	defer itrPool.Put(i)

	for i.ValueStart < len(s) {
		switch s[i.ValueStart] {
		case ' ', '\t', '\r', '\n':
			e, illegal := strfind.EndOfWhitespaceSeq(s[i.ValueStart:])
			i.ValueStart += e
			if illegal {
				return i.getError(ErrorCodeIllegalControlCharacter)
			}

		case ',':
			switch i.expect {
			case expectCommaOrArrTerm:
				i.expect = expectVal
			case expectCommaOrObjTerm:
				i.expect = expectKey
			default:
				return i.getError(ErrorCodeUnexpectedToken)
			}
			i.ValueStart++

		case '}':
			if i.expect != expectCommaOrObjTerm &&
				i.expect != expectKeyOrObjTerm {
				return i.getError(ErrorCodeUnexpectedToken)
			}

			i.st.Pop()
			if t := i.st.Top(); t != nil {
				switch t.Type {
				case stack.NodeTypeArray:
					i.expect = expectCommaOrArrTerm
				case stack.NodeTypeObject:
					i.expect = expectCommaOrObjTerm
				}
			}

			i.ValueStart++
			i.KeyStart, i.KeyEnd = -1, -1

		case ']':
			if i.expect != expectCommaOrArrTerm &&
				i.expect != expectValOrArrTerm {
				return i.getError(ErrorCodeUnexpectedToken)
			}

			i.st.Pop()

			if t := i.st.Top(); t != nil {
				switch t.Type {
				case stack.NodeTypeArray:
					i.expect = expectCommaOrArrTerm
				case stack.NodeTypeObject:
					i.expect = expectCommaOrObjTerm
				}
			}

			i.ValueStart++
			i.KeyStart, i.KeyEnd = -1, -1

		case '{': // Object
			if i.expect != expectVal && i.expect != expectValOrArrTerm {
				return i.getError(ErrorCodeUnexpectedToken)
			}
			i.expect = expectKeyOrObjTerm

			i.KeyStart, i.KeyEnd = -1, -1

			i.st.Push(stack.NodeTypeObject, 0, 0, 0)
			i.ValueStart++

		case '[': // Array
			if i.expect != expectVal && i.expect != expectValOrArrTerm {
				return i.getError(ErrorCodeUnexpectedToken)
			}
			i.expect = expectValOrArrTerm

			i.KeyStart, i.KeyEnd = -1, -1

			i.st.Push(stack.NodeTypeArray, 0, 0, 0)
			i.ValueStart++

		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if i.expect != expectVal && i.expect != expectValOrArrTerm {
				return i.getError(ErrorCodeUnexpectedToken)
			}

			if t := i.st.Top(); t != nil {
				switch t.Type {
				case stack.NodeTypeArray:
					i.expect = expectCommaOrArrTerm
				case stack.NodeTypeObject:
					i.expect = expectCommaOrObjTerm
				}
			}

			var err bool
			i.ValueEnd, err = jsonnum.Parse(i.src[i.ValueStart:])
			if err {
				return i.getError(ErrorCodeMalformedNumber)
			}
			i.ValueStart += i.ValueEnd
			i.KeyStart, i.KeyEnd = -1, -1

		case '"': // String
			i.ValueStart++
			i.ValueEnd = strfind.IndexTerm(s, i.ValueStart)
			if i.ValueEnd < 0 {
				return i.getError(ErrorCodeUnexpectedEOF)
			}
			for _, c := range s[i.ValueStart:i.ValueEnd] {
				if c < 0x20 {
					return i.getError(ErrorCodeIllegalControlCharacter)
				}
			}
			i.ValueStart = i.ValueEnd + 1
			t := i.st.Top()
			if t != nil && t.Type == stack.NodeTypeArray {
				// Array item string value
				if i.expect != expectVal && i.expect != expectValOrArrTerm {
					return i.getError(ErrorCodeUnexpectedToken)
				}
				i.expect = expectCommaOrArrTerm
				i.KeyStart, i.KeyEnd = -1, -1
			} else if i.KeyStart != -1 || i.st.Len() == 0 {
				// String field value
				if i.expect != expectVal {
					return i.getError(ErrorCodeUnexpectedToken)
				}
				if t != nil {
					i.expect = expectCommaOrObjTerm
				}

				i.KeyStart, i.KeyEnd = -1, -1
			} else {
				// Key
				if i.expect != expectKey && i.expect != expectKeyOrObjTerm {
					return i.getError(ErrorCodeUnexpectedToken)
				}
				i.expect = expectVal

				i.KeyStart, i.KeyEnd = i.ValueStart, i.ValueEnd
				if i.ValueEnd > len(s) {
					return i.getError(ErrorCodeUnexpectedEOF)
				}
				if x, err := i.parseColumn(s[i.ValueStart:]); err {
					return i.getError(ErrorCodeUnexpectedToken)
				} else {
					i.ValueStart += x + 1
				}
			}

		case 'n': // Null
			if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
				i.expect3(i.ValueStart+1, 'u', 'l', 'l') {
				return i.getError(ErrorCodeUnexpectedToken)
			}
			if t := i.st.Top(); t != nil {
				switch t.Type {
				case stack.NodeTypeArray:
					i.expect = expectCommaOrArrTerm
				case stack.NodeTypeObject:
					i.expect = expectCommaOrObjTerm
				}
			}

			i.KeyStart, i.KeyEnd = -1, -1
			i.ValueStart += len("null")

		case 'f': // False
			if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
				i.read4(i.ValueStart+1, 'a', 'l', 's', 'e') {
				return i.getError(ErrorCodeUnexpectedToken)
			}
			if t := i.st.Top(); t != nil {
				switch t.Type {
				case stack.NodeTypeArray:
					i.expect = expectCommaOrArrTerm
				case stack.NodeTypeObject:
					i.expect = expectCommaOrObjTerm
				}
			}

			i.KeyStart, i.KeyEnd = -1, -1
			i.ValueStart += len("false")

		case 't': // True
			if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
				i.expect3(i.ValueStart+1, 'r', 'u', 'e') {
				return i.getError(ErrorCodeUnexpectedToken)
			}
			if t := i.st.Top(); t != nil {
				switch t.Type {
				case stack.NodeTypeArray:
					i.expect = expectCommaOrArrTerm
				case stack.NodeTypeObject:
					i.expect = expectCommaOrObjTerm
				}
			}

			i.KeyStart, i.KeyEnd = -1, -1
			i.ValueStart += len("true")

		default:
			if s[i.ValueStart] < 0x20 {
				return i.getError(ErrorCodeIllegalControlCharacter)
			}
			return i.getError(ErrorCodeUnexpectedToken)
		}
	}

	if i.st.Len() > 0 {
		return i.getError(ErrorCodeUnexpectedEOF)
	}

	return Error{}
}

// Valid returns true if s is valid JSON, otherwise returns false.
func Valid(s string) bool {
	return !Validate(s).IsErr()
}

type Options struct {
	CachePath  bool
	EscapePath bool
}

// Scan calls fn for every scanned value including objects and arrays.
// Scan returns true if there was an error or if fn returned true,
// otherwise it returns false.
// If cachePath == true then paths are generated and cached on the fly
// reducing their performance penalty.
//
// WARNING: Fields exported by *Iterator provided in fn must not be mutated!
// Do not use or alias *Iterator after fn returns!
func Scan(
	o Options,
	s string,
	fn func(*Iterator) (err bool),
) Error {
	if s == "" {
		return Error{Code: ErrorCodeUnexpectedEOF}
	}

	i := getItrFromPool(s, o.EscapePath)
	defer itrPool.Put(i)

	if o.CachePath {
		return i.scanWithCachedPath(s, fn)
	}
	return i.scan(s, fn)
}

func (i *Iterator) scan(
	s string,
	fn func(*Iterator) (err bool),
) Error {
	for i.ValueStart < len(s) {
		switch s[i.ValueStart] {
		case ' ', '\t', '\r', '\n':
			e, illegal := strfind.EndOfWhitespaceSeq(s[i.ValueStart:])
			i.ValueStart += e
			if illegal {
				return i.getError(ErrorCodeIllegalControlCharacter)
			}

		case ',':
			switch i.expect {
			case expectCommaOrArrTerm:
				i.expect = expectVal
			case expectCommaOrObjTerm:
				i.expect = expectKey
			default:
				return i.getError(ErrorCodeUnexpectedToken)
			}

			i.ValueStart++

		case '}':
			if i.expect != expectCommaOrObjTerm &&
				i.expect != expectKeyOrObjTerm {
				return i.getError(ErrorCodeUnexpectedToken)
			}

			i.st.Pop()

			if t := i.st.Top(); t != nil {
				switch t.Type {
				case stack.NodeTypeArray:
					i.expect = expectCommaOrArrTerm
				case stack.NodeTypeObject:
					i.expect = expectCommaOrObjTerm
				}
			}

			i.Level--
			i.ValueStart++
			i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1

		case ']':
			if i.expect != expectCommaOrArrTerm &&
				i.expect != expectValOrArrTerm {
				return i.getError(ErrorCodeUnexpectedToken)
			}

			i.st.Pop()

			if t := i.st.Top(); t != nil {
				switch t.Type {
				case stack.NodeTypeArray:
					i.expect = expectCommaOrArrTerm
				case stack.NodeTypeObject:
					i.expect = expectCommaOrObjTerm
				}
			}

			i.Level--
			i.ValueStart++
			i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1

		case '{': // Object
			if i.expect != expectVal && i.expect != expectValOrArrTerm {
				return i.getError(ErrorCodeUnexpectedToken)
			}
			i.expect = expectKeyOrObjTerm

			i.ValueType = ValueTypeObject
			i.ValueEnd = -1
			ks, ke := i.KeyStart, i.KeyEnd
			if i.callFn(i.st.Top(), fn) {
				return i.getError(ErrorCallback)
			}
			i.st.Push(stack.NodeTypeObject, 0, ks, ke)
			i.Level++
			i.ValueStart++

		case '[': // Array
			if i.expect != expectVal && i.expect != expectValOrArrTerm {
				return i.getError(ErrorCodeUnexpectedToken)
			}
			i.expect = expectValOrArrTerm

			i.ValueType = ValueTypeArray
			i.ValueEnd = -1
			ks, ke := i.KeyStart, i.KeyEnd
			if i.callFn(i.st.Top(), fn) {
				return i.getError(ErrorCallback)
			}
			i.st.Push(stack.NodeTypeArray, 0, ks, ke)
			i.Level++
			i.ValueStart++

		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if i.expect != expectVal && i.expect != expectValOrArrTerm {
				return i.getError(ErrorCodeUnexpectedToken)
			}

			t := i.st.Top()
			if t != nil {
				switch t.Type {
				case stack.NodeTypeArray:
					i.expect = expectCommaOrArrTerm
				case stack.NodeTypeObject:
					i.expect = expectCommaOrObjTerm
				}
			}

			var err bool
			i.ValueEnd, err = jsonnum.Parse(i.src[i.ValueStart:])
			if err {
				return i.getError(ErrorCodeMalformedNumber)
			}
			i.ValueEnd += +i.ValueStart

			i.ValueType = ValueTypeNumber
			if i.callFn(t, fn) {
				return i.getError(ErrorCallback)
			}
			i.ValueStart = i.ValueEnd

		case '"': // String
			i.ValueStart++
			i.ArrayIndex = -1
			i.ValueEnd = strfind.IndexTerm(s, i.ValueStart)
			if i.ValueEnd < 0 {
				i.ValueStart--
				return i.getError(ErrorCodeUnexpectedEOF)
			}
			for _, c := range s[i.ValueStart:i.ValueEnd] {
				if c < 0x20 {
					return i.getError(ErrorCodeIllegalControlCharacter)
				}
			}
			t := i.st.Top()
			if t != nil && t.Type == stack.NodeTypeArray {
				// Array item string value
				if i.expect != expectVal && i.expect != expectValOrArrTerm {
					i.ValueStart--
					return i.getError(ErrorCodeUnexpectedToken)
				}
				i.expect = expectCommaOrArrTerm

				i.ArrayIndex = t.ArrLen
				t.ArrLen++
				i.ValueType = ValueTypeString
				i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1
				if fn(i) {
					i.ValueStart--
					return i.getError(ErrorCallback)
				}
				i.ValueStart = i.ValueEnd + 1
			} else if i.KeyStart != -1 || i.Level == 0 {
				// String field value
				if i.expect != expectVal {
					i.ValueStart--
					return i.getError(ErrorCodeUnexpectedToken)
				}
				if t != nil {
					i.expect = expectCommaOrObjTerm
				}

				i.ValueType = ValueTypeString
				if i.callFn(nil, fn) {
					i.ValueStart--
					return i.getError(ErrorCallback)
				}
				i.ValueStart = i.ValueEnd + 1
			} else {
				// Key
				if i.expect != expectKey && i.expect != expectKeyOrObjTerm {
					i.ValueStart--
					return i.getError((ErrorCodeUnexpectedToken))
				}
				i.expect = expectVal

				i.KeyStart, i.KeyEnd = i.ValueStart, i.ValueEnd
				i.ValueStart = i.ValueEnd + 1
				if i.ValueEnd > len(s) {
					return i.getError(ErrorCodeUnexpectedEOF)
				}
				if x, err := i.parseColumn(s[i.ValueStart:]); err {
					return i.getError(ErrorCodeUnexpectedToken)
				} else {
					i.ValueStart += x + 1
				}
			}

		case 'n': // Null
			if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
				i.expect3(i.ValueStart+1, 'u', 'l', 'l') {
				return i.getError(ErrorCodeUnexpectedToken)
			}
			t := i.st.Top()
			if t != nil {
				switch t.Type {
				case stack.NodeTypeArray:
					i.expect = expectCommaOrArrTerm
				case stack.NodeTypeObject:
					i.expect = expectCommaOrObjTerm
				}
			}

			i.ValueType = ValueTypeNull
			i.ValueEnd = i.ValueStart + len("null")
			if i.callFn(t, fn) {
				return i.getError(ErrorCallback)
			}
			i.ValueStart += len("null")

		case 'f': // False
			if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
				i.read4(i.ValueStart+1, 'a', 'l', 's', 'e') {
				return i.getError(ErrorCodeUnexpectedToken)
			}
			t := i.st.Top()
			if t != nil {
				switch t.Type {
				case stack.NodeTypeArray:
					i.expect = expectCommaOrArrTerm
				case stack.NodeTypeObject:
					i.expect = expectCommaOrObjTerm
				}
			}

			i.ValueType = ValueTypeFalse
			i.ValueEnd = i.ValueStart + len("false")
			if i.callFn(t, fn) {
				return i.getError(ErrorCallback)
			}
			i.ValueStart += len("false")

		case 't': // True
			if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
				i.expect3(i.ValueStart+1, 'r', 'u', 'e') {
				return i.getError(ErrorCodeUnexpectedToken)
			}
			t := i.st.Top()
			if t != nil {
				switch t.Type {
				case stack.NodeTypeArray:
					i.expect = expectCommaOrArrTerm
				case stack.NodeTypeObject:
					i.expect = expectCommaOrObjTerm
				}
			}

			i.ValueType = ValueTypeTrue
			i.ValueEnd = i.ValueStart + len("true")
			if i.callFn(t, fn) {
				return i.getError(ErrorCallback)
			}
			i.ValueStart += len("true")

		default:
			if s[i.ValueStart] < 0x20 {
				return i.getError(ErrorCodeIllegalControlCharacter)
			}
			return i.getError(ErrorCodeUnexpectedToken)
		}
	}

	if i.st.Len() > 0 {
		return i.getError(ErrorCodeUnexpectedEOF)
	}

	return Error{}
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

type expectation int8

const (
	_ expectation = iota
	expectVal
	expectCommaOrObjTerm
	expectCommaOrArrTerm
	expectKeyOrObjTerm
	expectValOrArrTerm
	expectKey
)

func (i *Iterator) scanWithCachedPath(
	s string,
	fn func(*Iterator) bool,
) Error {
	i.cachedPath = append(i.cachedPath, '.')

	for i.ValueStart < len(s) {
		switch s[i.ValueStart] {
		case ' ', '\t', '\r', '\n':
			e, illegal := strfind.EndOfWhitespaceSeq(s[i.ValueStart:])
			i.ValueStart += e
			if illegal {
				return i.getError(ErrorCodeIllegalControlCharacter)
			}

		case ',':
			switch i.expect {
			case expectCommaOrArrTerm:
				i.expect = expectVal
			case expectCommaOrObjTerm:
				i.expect = expectKey
			default:
				return i.getError(ErrorCodeUnexpectedToken)
			}
			i.ValueStart++

		case '}':
			if i.expect != expectCommaOrObjTerm &&
				i.expect != expectKeyOrObjTerm {
				return i.getError(ErrorCodeUnexpectedToken)
			}

			i.st.Pop()

			if t := i.st.Top(); t != nil {
				switch t.Type {
				case stack.NodeTypeArray:
					i.expect = expectCommaOrArrTerm
				case stack.NodeTypeObject:
					i.expect = expectCommaOrObjTerm
				}
			}

			i.Level--
			i.ValueStart++
			i.cacheCleanupAfterContainer()
			i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1

		case ']':
			if i.expect != expectCommaOrArrTerm &&
				i.expect != expectValOrArrTerm {
				return i.getError(ErrorCodeUnexpectedToken)
			}

			i.st.Pop()

			if t := i.st.Top(); t != nil {
				switch t.Type {
				case stack.NodeTypeArray:
					i.expect = expectCommaOrArrTerm
				case stack.NodeTypeObject:
					i.expect = expectCommaOrObjTerm
				}
			}

			i.Level--
			i.ValueStart++
			if len(i.cachedPath) > 1 {
				i.cachedPath = i.cachedPath[:len(i.cachedPath)-1]
				i.cacheCleanupAfterContainer()
			}
			i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1

		case '{': // Object
			if i.expect != expectVal && i.expect != expectValOrArrTerm {
				return i.getError(ErrorCodeUnexpectedToken)
			}
			i.expect = expectKeyOrObjTerm

			i.ValueType = ValueTypeObject
			i.ValueEnd = -1
			ks, ke := i.KeyStart, i.KeyEnd
			if i.callFnWithPathCache(i.st.Top(), false, fn) {
				return i.getError(ErrorCallback)
			}
			i.st.Push(stack.NodeTypeObject, 0, ks, ke)
			i.Level++
			i.ValueStart++

		case '[': // Array
			if i.expect != expectVal && i.expect != expectValOrArrTerm {
				return i.getError(ErrorCodeUnexpectedToken)
			}
			i.expect = expectValOrArrTerm

			i.ValueType = ValueTypeArray
			i.ValueEnd = -1
			ks, ke := i.KeyStart, i.KeyEnd
			if i.callFnWithPathCache(i.st.Top(), false, fn) {
				return i.getError(ErrorCallback)
			}
			i.st.Push(stack.NodeTypeArray, 0, ks, ke)
			i.Level++
			i.ValueStart++
			i.cachedPath = append(i.cachedPath, '[')

		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if i.expect != expectVal && i.expect != expectValOrArrTerm {
				return i.getError(ErrorCodeUnexpectedToken)
			}

			t := i.st.Top()
			if t != nil {
				switch t.Type {
				case stack.NodeTypeArray:
					i.expect = expectCommaOrArrTerm
				case stack.NodeTypeObject:
					i.expect = expectCommaOrObjTerm
				}
			}

			var err bool
			i.ValueEnd, err = jsonnum.Parse(i.src[i.ValueStart:])
			if err {
				return i.getError(ErrorCodeMalformedNumber)
			}
			i.ValueEnd += +i.ValueStart

			i.ValueType = ValueTypeNumber
			if i.callFnWithPathCache(t, true, fn) {
				return i.getError(ErrorCallback)
			}
			i.ValueStart = i.ValueEnd

		case '"': // String
			i.ValueStart++
			i.ArrayIndex = -1
			i.ValueEnd = strfind.IndexTerm(s, i.ValueStart)
			if i.ValueEnd < 0 {
				i.ValueStart--
				return i.getError(ErrorCodeUnexpectedEOF)
			}
			for _, c := range s[i.ValueStart:i.ValueEnd] {
				if c < 0x20 {
					return i.getError(ErrorCodeIllegalControlCharacter)
				}
			}
			t := i.st.Top()
			if t != nil && t.Type == stack.NodeTypeArray {
				// Array item string value
				if i.expect != expectVal && i.expect != expectValOrArrTerm {
					i.ValueStart--
					return i.getError(ErrorCodeUnexpectedToken)
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
					return i.getError(ErrorCallback)
				}
				// Remove index
				i.cachedPath = i.cachedPath[:initArrIndex]

				i.ValueStart = i.ValueEnd + 1
			} else if i.KeyStart != -1 || i.Level == 0 {
				// String field value
				if i.expect != expectVal {
					i.ValueStart--
					return i.getError(ErrorCodeUnexpectedToken)
				}
				if t != nil {
					i.expect = expectCommaOrObjTerm
				}

				i.ValueType = ValueTypeString
				if i.callFnWithPathCache(nil, true, fn) {
					i.ValueStart--
					return i.getError(ErrorCallback)
				}
				i.ValueStart = i.ValueEnd + 1
			} else {
				// Key
				if i.expect != expectKey && i.expect != expectKeyOrObjTerm {
					i.ValueStart--
					return i.getError((ErrorCodeUnexpectedToken))
				}
				i.expect = expectVal

				i.KeyStart, i.KeyEnd = i.ValueStart, i.ValueEnd
				i.ValueStart = i.ValueEnd + 1
				if i.ValueEnd > len(s) {
					return i.getError(ErrorCodeUnexpectedEOF)
				}
				if x, err := i.parseColumn(s[i.ValueStart:]); err {
					return i.getError(ErrorCodeUnexpectedToken)
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
				return i.getError(ErrorCodeUnexpectedToken)
			}
			t := i.st.Top()
			if t != nil {
				switch t.Type {
				case stack.NodeTypeArray:
					i.expect = expectCommaOrArrTerm
				case stack.NodeTypeObject:
					i.expect = expectCommaOrObjTerm
				}
			}

			i.ValueType = ValueTypeNull
			i.ValueEnd = i.ValueStart + len("null")
			if i.callFnWithPathCache(t, true, fn) {
				return i.getError(ErrorCallback)
			}
			i.ValueStart += len("null")

		case 'f': // False
			if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
				i.read4(i.ValueStart+1, 'a', 'l', 's', 'e') {
				return i.getError(ErrorCodeUnexpectedToken)
			}
			t := i.st.Top()
			if t != nil {
				switch t.Type {
				case stack.NodeTypeArray:
					i.expect = expectCommaOrArrTerm
				case stack.NodeTypeObject:
					i.expect = expectCommaOrObjTerm
				}
			}

			i.ValueType = ValueTypeFalse
			i.ValueEnd = i.ValueStart + len("false")
			if i.callFnWithPathCache(t, true, fn) {
				return i.getError(ErrorCallback)
			}
			i.ValueStart += len("false")

		case 't': // True
			if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
				i.expect3(i.ValueStart+1, 'r', 'u', 'e') {
				return i.getError(ErrorCodeUnexpectedToken)
			}
			t := i.st.Top()
			if t != nil {
				switch t.Type {
				case stack.NodeTypeArray:
					i.expect = expectCommaOrArrTerm
				case stack.NodeTypeObject:
					i.expect = expectCommaOrObjTerm
				}
			}

			i.ValueType = ValueTypeTrue
			i.ValueEnd = i.ValueStart + len("true")
			if i.callFnWithPathCache(t, true, fn) {
				return i.getError(ErrorCallback)
			}
			i.ValueStart += len("true")

		default:
			if s[i.ValueStart] < 0x20 {
				return i.getError(ErrorCodeIllegalControlCharacter)
			}
			return i.getError(ErrorCodeUnexpectedToken)
		}
	}

	if i.st.Len() > 0 {
		return i.getError(ErrorCodeUnexpectedEOF)
	}

	return Error{}
}

// ValueType defines a JSON value type
type ValueType int

// JSON value types
const (
	_ ValueType = iota
	ValueTypeObject
	ValueTypeArray
	ValueTypeNull
	ValueTypeFalse
	ValueTypeTrue
	ValueTypeString
	ValueTypeNumber
)

func (t ValueType) String() string {
	switch t {
	case ValueTypeObject:
		return "object"
	case ValueTypeArray:
		return "array"
	case ValueTypeNull:
		return "null"
	case ValueTypeFalse:
		return "false"
	case ValueTypeTrue:
		return "true"
	case ValueTypeString:
		return "string"
	case ValueTypeNumber:
		return "number"
	}
	return ""
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

func errorMessage(c ErrorCode, index int, atIndex rune) string {
	errMsg := ""
	switch c {
	case ErrorCodeUnexpectedToken:
		errMsg = "unexpected token"
	case ErrorCodeMalformedNumber:
		errMsg = "malformed number"
	case ErrorCodeUnexpectedEOF:
		if atIndex == 0 {
			return fmt.Sprintf(
				"error at index %d: unexpected EOF",
				index,
			)
		}
		errMsg = "unexpected EOF"
	case ErrorCodeIllegalControlCharacter:
		errMsg = "illegal control character"
	case ErrorCallback:
		errMsg = "callback error"
	default:
		return ""
	}
	if index >= 0 {
		if atIndex < 32 {
			return fmt.Sprintf(
				"error at index %d (0x%x): %s",
				index, atIndex, errMsg,
			)
		}
		return fmt.Sprintf(
			"error at index %d ('%s'): %s",
			index, string(atIndex), errMsg,
		)
	}
	return fmt.Sprintf(
		"error at index %d: %s",
		index, errMsg,
	)
}
