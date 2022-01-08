package jscan

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/romshark/jscan/internal/jsonnum"
	"github.com/romshark/jscan/internal/stack"
	"github.com/romshark/jscan/internal/strfind"
)

var keyEscapeReplacer = strings.NewReplacer(
	".", `\.`,
	"[", `\[`,
	"]", `\]`,
)

// Iterator provides access to the recently scanned value.
type Iterator struct {
	st         *stack.Stack
	src        string
	escapePath bool
	cachedPath []byte
	expect     expectation
	errCode    ErrorCode

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
// If key is != "" then the element is a field, otherwise
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
// If escape == true then all '.', '[' and ']' in keys will be escaped.
func (i *Iterator) Path() (s string) {
	i.ViewPath(func(p []byte) { s = string(p) })
	return
}

// ViewPath calls fn and provides the stringified path.
// If escape == true then all '.', '[' and ']' in keys will be escaped.
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
				k = keyEscapeReplacer.Replace(k)
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
			k = keyEscapeReplacer.Replace(k)
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

// getError returns the stringified error, if any
func (i *Iterator) getError() Error {
	return Error{
		Src:   i.src,
		Index: i.ValueStart,
		Code:  i.errCode,
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
	errMsg := ""
	switch e.Code {
	case ErrorCodeUnexpectedToken:
		errMsg = "unexpected token"
	case ErrorCodeMalformedNumber:
		errMsg = "malformed number"
	case ErrorCodeUnexpectedEOF:
		errMsg = "unexpected EOF"
	case ErrorCallback:
		errMsg = "callback error"
	default:
		return ""
	}
	r := ""
	if e.Index < len(e.Src) {
		r = " ('" + string(e.Src[e.Index]) + "')"
	}
	return fmt.Sprintf(
		"error at index %d%s: %s",
		e.Index, r, errMsg,
	)
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
func Scan(
	o Options,
	s string,
	fn func(*Iterator) (err bool),
) Error {
	i := itrPool.Get().(*Iterator)
	defer itrPool.Put(i)
	i.st.Reset()
	i.escapePath = o.EscapePath
	i.cachedPath = i.cachedPath[:0]
	i.src = s
	i.ValueType = 0
	i.Level = 0
	i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1
	i.ValueStart, i.ValueEnd = 0, -1
	i.ArrayIndex = 0
	i.expect = expectVal

	if o.CachePath {
		if i.scanWithCachedPath(s, fn) {
			return i.getError()
		}
	}
	if i.scan(s, fn) {
		return i.getError()
	}
	return Error{}
}

func (i *Iterator) onComma() (err bool) {
	switch i.expect {
	case expectCommaOrArrTerm:
		i.expect = expectValOrArrTerm
		return false
	case expectCommaOrObjTerm:
		i.expect = expectKeyOrObjTerm
		return false
	}
	i.errCode = ErrorCodeUnexpectedToken
	return true
}

func (i *Iterator) onObjectTerm() (err bool) {
	if i.expect != expectCommaOrObjTerm && i.expect != expectKeyOrObjTerm {
		i.errCode = ErrorCodeUnexpectedToken
		return true
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
	return false
}

func (i *Iterator) onArrayTerm() (err bool) {
	if i.expect != expectCommaOrArrTerm && i.expect != expectValOrArrTerm {
		i.errCode = ErrorCodeUnexpectedToken
		return true
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
	return false
}

func (i *Iterator) onObjectBegin() (err bool) {
	if i.expect != expectVal && i.expect != expectValOrArrTerm {
		i.errCode = ErrorCodeUnexpectedToken
		return true
	}
	i.expect = expectKeyOrObjTerm
	return false
}

func (i *Iterator) onArrayBegin() (err bool) {
	if i.expect != expectVal && i.expect != expectValOrArrTerm {
		i.errCode = ErrorCodeUnexpectedToken
		return true
	}
	i.expect = expectValOrArrTerm
	return false
}

func (i *Iterator) onNumber() (top *stack.Node, err bool) {
	if i.expect != expectVal && i.expect != expectValOrArrTerm {
		i.errCode = ErrorCodeUnexpectedToken
		return nil, true
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

	i.ValueEnd, err = jsonnum.Parse(i.src[i.ValueStart:])
	if err {
		i.errCode = ErrorCodeMalformedNumber
		return nil, true
	}
	i.ValueEnd += +i.ValueStart

	return t, false
}

func (i *Iterator) onStringArrayItem() (err bool) {
	if i.expect != expectValOrArrTerm {
		i.errCode = ErrorCodeUnexpectedToken
		return true
	}
	i.expect = expectCommaOrArrTerm
	return false
}

func (i *Iterator) onStringFieldValue(t *stack.Node) (err bool) {
	if i.expect != expectVal {
		i.errCode = ErrorCodeUnexpectedToken
		return true
	}
	if t != nil {
		switch t.Type {
		case stack.NodeTypeArray:
			i.expect = expectCommaOrArrTerm
		case stack.NodeTypeObject:
			i.expect = expectCommaOrObjTerm
		}
	}
	return false
}

func (i *Iterator) onKey() (err bool) {
	if i.expect != expectKeyOrObjTerm {
		i.errCode = ErrorCodeUnexpectedToken
		return true
	}
	i.expect = expectVal
	return false
}

func (i *Iterator) onNull() (t *stack.Node, err bool) {
	if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
		i.expect3(i.ValueStart+1, 'u', 'l', 'l') {
		i.errCode = ErrorCodeUnexpectedToken
		return nil, true
	}
	t = i.st.Top()
	if t != nil {
		switch t.Type {
		case stack.NodeTypeArray:
			i.expect = expectCommaOrArrTerm
		case stack.NodeTypeObject:
			i.expect = expectCommaOrObjTerm
		}
	}
	return t, false
}

func (i *Iterator) onFalse() (t *stack.Node, err bool) {
	if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
		i.read4(i.ValueStart+1, 'a', 'l', 's', 'e') {
		i.errCode = ErrorCodeUnexpectedToken
		return nil, true
	}
	t = i.st.Top()
	if t != nil {
		switch t.Type {
		case stack.NodeTypeArray:
			i.expect = expectCommaOrArrTerm
		case stack.NodeTypeObject:
			i.expect = expectCommaOrObjTerm
		}
	}
	return t, false
}

func (i *Iterator) onTrue() (t *stack.Node, err bool) {
	if (i.expect != expectVal && i.expect != expectValOrArrTerm) ||
		i.expect3(i.ValueStart+1, 'r', 'u', 'e') {
		i.errCode = ErrorCodeUnexpectedToken
		return nil, true
	}
	t = i.st.Top()
	if t != nil {
		switch t.Type {
		case stack.NodeTypeArray:
			i.expect = expectCommaOrArrTerm
		case stack.NodeTypeObject:
			i.expect = expectCommaOrObjTerm
		}
	}
	return t, false
}

func (i *Iterator) scan(
	s string,
	fn func(*Iterator) (err bool),
) (err bool) {
	for i.ValueStart < len(s) {
		switch s[i.ValueStart] {
		case ' ', '\t', '\r', '\n':
			i.ValueStart++

		case ',':
			if i.onComma() {
				return true
			}
			i.ValueStart++

		case '}':
			if i.onObjectTerm() {
				return true
			}

			i.Level--
			i.ValueStart++
			i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1

		case ']':
			if i.onArrayTerm() {
				return true
			}

			i.Level--
			i.ValueStart++
			i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1

		case '{': // Object
			if i.onObjectBegin() {
				return true
			}

			i.ValueType = ValueTypeObject
			i.ValueEnd = -1
			ks, ke := i.KeyStart, i.KeyEnd
			if i.callFn(i.st.Top(), fn) {
				return true
			}
			i.st.Push(stack.NodeTypeObject, 0, ks, ke)
			i.Level++
			i.ValueStart++

		case '[': // Array
			if i.onArrayBegin() {
				return true
			}

			i.ValueType = ValueTypeArray
			i.ValueEnd = -1
			ks, ke := i.KeyStart, i.KeyEnd
			if i.callFn(i.st.Top(), fn) {
				return true
			}
			i.st.Push(stack.NodeTypeArray, 0, ks, ke)
			i.Level++
			i.ValueStart++

		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			t, err := i.onNumber()
			if err {
				return true
			}

			i.ValueType = ValueTypeNumber
			if i.callFn(t, fn) {
				return true
			}
			i.ValueStart = i.ValueEnd

		case '"': // String
			i.ValueStart++
			i.ArrayIndex = -1
			i.ValueEnd = strfind.IndexTerm(s, i.ValueStart)
			if i.ValueEnd < 0 {
				i.ValueStart--
				i.errCode = ErrorCodeUnexpectedEOF
				return true
			}
			t := i.st.Top()
			if t != nil && t.Type == stack.NodeTypeArray {
				// Array item string value
				if i.onStringArrayItem() {
					i.ValueStart--
					return true
				}

				i.ArrayIndex = t.ArrLen
				t.ArrLen++
				i.ValueType = ValueTypeString
				i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1
				fn(i)
				i.ValueStart = i.ValueEnd + 1
			} else if i.KeyStart != -1 || i.Level == 0 {
				// String field value
				if i.onStringFieldValue(t) {
					i.ValueStart--
					return true
				}

				i.ValueType = ValueTypeString
				if i.callFn(nil, fn) {
					i.ValueStart--
					return true
				}
				i.ValueStart = i.ValueEnd + 1
			} else {
				// Key
				if i.onKey() {
					i.ValueStart--
					return true
				}
				i.KeyStart, i.KeyEnd = i.ValueStart, i.ValueEnd
				i.ValueStart = i.ValueEnd + 1
				if i.ValueEnd > len(s) {
					i.errCode = ErrorCodeUnexpectedEOF
					return true
				}
				if x, err := i.parseColumn(s[i.ValueStart:]); err {
					return true
				} else {
					i.ValueStart += x + 1
				}
			}

		case 'n': // Null
			t, err := i.onNull()
			if err {
				return true
			}

			i.ValueType = ValueTypeNull
			i.ValueEnd = i.ValueStart + len("null")
			if i.callFn(t, fn) {
				return true
			}
			i.ValueStart += len("null")

		case 'f': // False
			t, err := i.onFalse()
			if err {
				return true
			}

			i.ValueType = ValueTypeFalse
			i.ValueEnd = i.ValueStart + len("false")
			if i.callFn(t, fn) {
				return true
			}
			i.ValueStart += len("false")

		case 't': // True
			t, err := i.onTrue()
			if err {
				return true
			}

			i.ValueType = ValueTypeTrue
			i.ValueEnd = i.ValueStart + len("true")
			if i.callFn(t, fn) {
				return true
			}
			i.ValueStart += len("true")

		default:
			i.errCode = ErrorCodeUnexpectedToken
			return true
		}
	}

	if i.st.Len() > 0 {
		i.errCode = ErrorCodeUnexpectedEOF
		return true
	}

	return false
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
)

func (i *Iterator) scanWithCachedPath(
	s string,
	fn func(*Iterator) bool,
) (err bool) {
	i.cachedPath = append(i.cachedPath, '.')

	for i.ValueStart < len(s) {
		switch s[i.ValueStart] {
		case ' ', '\t', '\r', '\n':
			i.ValueStart++

		case ',':
			if i.onComma() {
				return true
			}
			i.ValueStart++

		case '}':
			if i.onObjectTerm() {
				return true
			}

			i.Level--
			i.ValueStart++
			i.cacheCleanupAfterContainer()
			i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1

		case ']':
			if i.onArrayTerm() {
				return true
			}

			i.Level--
			i.ValueStart++
			if len(i.cachedPath) > 1 {
				i.cachedPath = i.cachedPath[:len(i.cachedPath)-1]
				i.cacheCleanupAfterContainer()
			}
			i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1

		case '{': // Object
			if i.onObjectBegin() {
				return true
			}

			i.ValueType = ValueTypeObject
			i.ValueEnd = -1
			ks, ke := i.KeyStart, i.KeyEnd
			if i.callFnWithPathCache(i.st.Top(), false, fn) {
				return true
			}
			i.st.Push(stack.NodeTypeObject, 0, ks, ke)
			i.Level++
			i.ValueStart++

		case '[': // Array
			if i.onArrayBegin() {
				return true
			}

			i.ValueType = ValueTypeArray
			i.ValueEnd = -1
			ks, ke := i.KeyStart, i.KeyEnd
			if i.callFnWithPathCache(i.st.Top(), false, fn) {
				return true
			}
			i.st.Push(stack.NodeTypeArray, 0, ks, ke)
			i.Level++
			i.ValueStart++
			i.cachedPath = append(i.cachedPath, '[')

		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			t, err := i.onNumber()
			if err {
				return true
			}

			i.ValueType = ValueTypeNumber
			if i.callFnWithPathCache(t, true, fn) {
				return true
			}
			i.ValueStart = i.ValueEnd

		case '"': // String
			i.ValueStart++
			i.ArrayIndex = -1
			i.ValueEnd = strfind.IndexTerm(s, i.ValueStart)
			if i.ValueEnd < 0 {
				i.ValueStart--
				i.errCode = ErrorCodeUnexpectedEOF
				return true
			}
			t := i.st.Top()
			if t != nil && t.Type == stack.NodeTypeArray {
				// Array item string value
				if i.onStringArrayItem() {
					return true
				}

				i.ArrayIndex = t.ArrLen
				t.ArrLen++
				i.ValueType = ValueTypeString
				i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1

				initArrIndex := len(i.cachedPath)
				i.cachedPath = strconv.AppendInt(
					i.cachedPath, int64(i.ArrayIndex), 10,
				)
				i.cachedPath = append(i.cachedPath, ']')

				fn(i)
				// Remove index
				i.cachedPath = i.cachedPath[:initArrIndex]

				i.ValueStart = i.ValueEnd + 1
			} else if i.KeyStart != -1 || i.Level == 0 {
				// String field value
				if i.onStringFieldValue(t) {
					return true
				}

				i.ValueType = ValueTypeString
				if i.callFnWithPathCache(nil, true, fn) {
					return true
				}
				i.ValueStart = i.ValueEnd + 1
			} else {
				// Key
				if i.onKey() {
					i.ValueStart--
					return true
				}
				i.KeyStart, i.KeyEnd = i.ValueStart, i.ValueEnd
				i.ValueStart = i.ValueEnd + 1
				if i.ValueEnd > len(s) {
					i.errCode = ErrorCodeUnexpectedEOF
					return true
				}
				if x, err := i.parseColumn(s[i.ValueStart:]); err {
					return true
				} else {
					i.ValueStart += x + 1
				}

				if len(i.cachedPath) > 1 {
					i.cachedPath = append(i.cachedPath, '.')
				}
				var e string
				if i.escapePath {
					e = keyEscapeReplacer.Replace(s[i.KeyStart:i.KeyEnd])
				} else {
					e = s[i.KeyStart:i.KeyEnd]
				}
				i.KeyLenEscaped = len(e)
				i.cachedPath = append(i.cachedPath, e...)
			}

		case 'n': // Null
			t, err := i.onNull()
			if err {
				return true
			}

			i.ValueType = ValueTypeNull
			i.ValueEnd = i.ValueStart + len("null")
			if i.callFnWithPathCache(t, true, fn) {
				return true
			}
			i.ValueStart += len("null")

		case 'f': // False
			t, err := i.onFalse()
			if err {
				return true
			}

			i.ValueType = ValueTypeFalse
			i.ValueEnd = i.ValueStart + len("false")
			if i.callFnWithPathCache(t, true, fn) {
				return true
			}
			i.ValueStart += len("false")

		case 't': // True
			t, err := i.onTrue()
			if err {
				return true
			}

			i.ValueType = ValueTypeTrue
			i.ValueEnd = i.ValueStart + len("true")
			if i.callFnWithPathCache(t, true, fn) {
				return true
			}
			i.ValueStart += len("true")

		default:
			i.errCode = ErrorCodeUnexpectedToken
			return true
		}
	}

	if i.st.Len() > 0 {
		i.errCode = ErrorCodeUnexpectedEOF
		return true
	}

	return false
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
	i.errCode = ErrorCodeUnexpectedToken
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

func (i *Iterator) callFn(t *stack.Node, fn func(i *Iterator) (err bool)) bool {
	i.ArrayIndex = -1
	if t != nil && t.Type == stack.NodeTypeArray {
		i.ArrayIndex = t.ArrLen
		t.ArrLen++
	}

	if fn(i) {
		i.errCode = ErrorCallback
		return true
	}

	i.KeyStart, i.KeyEnd, i.KeyLenEscaped = -1, -1, -1
	return false
}

func (i *Iterator) callFnWithPathCache(
	t *stack.Node,
	nonComposite bool,
	fn func(i *Iterator) bool,
) bool {
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
		i.errCode = ErrorCallback
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
