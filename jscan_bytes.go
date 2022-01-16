package jscan

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"

	"github.com/romshark/jscan/internal/jsonnum"
	"github.com/romshark/jscan/internal/keyescape"
	"github.com/romshark/jscan/internal/stack"
	"github.com/romshark/jscan/internal/strfind"
)

// IteratorBytes provides access to the recently scanned value.
type IteratorBytes struct {
	st         *stack.Stack
	src        []byte
	escapePath bool
	cachedPath []byte
	expect     expectation
	errCode    ErrorCode
	keyBuf     []byte

	ValueType                       ValueType
	Level                           int
	KeyStart, KeyEnd, KeyLenEscaped int
	ValueStart, ValueEnd            int
	ArrayIndex                      int
}

// Key returns the field key if any.
func (i *IteratorBytes) Key() []byte {
	if i.KeyStart < 0 {
		return nil
	}
	return i.src[i.KeyStart:i.KeyEnd]
}

// Value returns the value if any.
func (i *IteratorBytes) Value() []byte {
	if i.ValueEnd < 0 || i.ValueEnd < i.ValueStart {
		return nil
	}
	return i.src[i.ValueStart:i.ValueEnd]
}

// ScanPath calls fn for every element in the path of the value.
// If keyStart is != -1 then the element is a field value, otherwise
// arrIndex indicates the index of the item in the underlying array.
func (i *IteratorBytes) ScanPath(fn func(keyStart, keyEnd, arrIndex int)) {
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
func (i *IteratorBytes) Path() (s []byte) {
	i.ViewPath(func(p []byte) {
		s = make([]byte, len(p))
		copy(s, p)
	})
	return
}

// ViewPath calls fn and provides the stringified path.
//
// WARNING: do not use or alias p after fn returns!
// Only viewing or copying are considered safe!
// Use (*IteratorBytes).Path instead for a safer and more convenient API.
func (i *IteratorBytes) ViewPath(fn func(p []byte)) {
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
				k = keyescape.EscapeAppend(nil, k)
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
			k = keyescape.EscapeAppend(nil, k)
		}
		b = append(b, k...)
	}
	fn(b)
}

var itrPoolBytes = sync.Pool{
	New: func() interface{} {
		return &IteratorBytes{
			st:     stack.New(64),
			keyBuf: make([]byte, 0, 4096),
		}
	},
}

// getError returns the stringified error, if any.
func (i *IteratorBytes) getError() ErrorBytes {
	return ErrorBytes{
		Src:   i.src,
		Index: i.ValueStart,
		Code:  i.errCode,
	}
}

type ErrorBytes struct {
	Src   []byte
	Index int
	Code  ErrorCode
}

// IsErr returns true if there is an error, otherwise returns false.
func (e ErrorBytes) IsErr() bool { return e.Code != 0 }

func (e ErrorBytes) Error() string {
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

// Get calls fn for the value at the given path.
// The value path is defined by keys separated by a dot and
// index access operators for arrays.
// If no value is found for the given path then fn isn't called
// and no error is returned.
// If escapePath then all dots and square brackets are expected to be escaped.
//
// WARNING: Fields exported by *IteratorBytes in fn must not be mutated!
// Do not use or alias *IteratorBytes after fn returns!
func GetBytes(
	s, path []byte,
	escapePath bool,
	fn func(*IteratorBytes),
) ErrorBytes {
	err := ScanBytes(Options{
		CachePath:  true,
		EscapePath: escapePath,
	}, s, func(i *IteratorBytes) (err bool) {
		i.ViewPath(func(p []byte) {
			if !bytes.Equal(p, path) {
				return
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

// ValidBytes returns true if s is valid JSON, otherwise returns false.
func ValidBytes(s []byte) bool {
	err := ScanBytes(Options{
		CachePath:  false,
		EscapePath: false,
	}, s, func(*IteratorBytes) (err bool) { return false })
	return !err.IsErr()
}

// ScanBytes calls fn for every scanned value including objects and arrays.
// Scan returns true if there was an error or if fn returned true,
// otherwise it returns false.
// If cachePath == true then paths are generated and cached on the fly
// reducing their performance penalty.
//
// WARNING: Fields exported by *IteratorBytes in fn must not be mutated!
// Do not use or alias *IteratorBytes after fn returns!
func ScanBytes(
	o Options,
	s []byte,
	fn func(*IteratorBytes) (err bool),
) ErrorBytes {
	i := itrPoolBytes.Get().(*IteratorBytes)
	defer itrPoolBytes.Put(i)
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
	i.keyBuf = i.keyBuf[:0]

	if o.CachePath {
		if i.scanWithCachedPath(s, fn) {
			return i.getError()
		}
	}
	if i.scan(s, fn) {
		return i.getError()
	}
	return ErrorBytes{}
}

func (i *IteratorBytes) onComma() (err bool) {
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

func (i *IteratorBytes) onObjectTerm() (err bool) {
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

func (i *IteratorBytes) onArrayTerm() (err bool) {
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

func (i *IteratorBytes) onObjectBegin() (err bool) {
	if i.expect != expectVal && i.expect != expectValOrArrTerm {
		i.errCode = ErrorCodeUnexpectedToken
		return true
	}
	i.expect = expectKeyOrObjTerm
	return false
}

func (i *IteratorBytes) onArrayBegin() (err bool) {
	if i.expect != expectVal && i.expect != expectValOrArrTerm {
		i.errCode = ErrorCodeUnexpectedToken
		return true
	}
	i.expect = expectValOrArrTerm
	return false
}

func (i *IteratorBytes) onNumber() (top *stack.Node, err bool) {
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

	i.ValueEnd, err = jsonnum.ParseBytes(i.src[i.ValueStart:])
	if err {
		i.errCode = ErrorCodeMalformedNumber
		return nil, true
	}
	i.ValueEnd += +i.ValueStart

	return t, false
}

func (i *IteratorBytes) onStringArrayItem() (err bool) {
	if i.expect != expectValOrArrTerm {
		i.errCode = ErrorCodeUnexpectedToken
		return true
	}
	i.expect = expectCommaOrArrTerm
	return false
}

func (i *IteratorBytes) onStringFieldValue(t *stack.Node) (err bool) {
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

func (i *IteratorBytes) onKey() (err bool) {
	if i.expect != expectKeyOrObjTerm {
		i.errCode = ErrorCodeUnexpectedToken
		return true
	}
	i.expect = expectVal
	return false
}

func (i *IteratorBytes) onNull() (t *stack.Node, err bool) {
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

func (i *IteratorBytes) onFalse() (t *stack.Node, err bool) {
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

func (i *IteratorBytes) onTrue() (t *stack.Node, err bool) {
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

func (i *IteratorBytes) scan(
	s []byte,
	fn func(*IteratorBytes) (err bool),
) (err bool) {
	for i.ValueStart < len(s) {
		switch s[i.ValueStart] {
		case ' ', '\t', '\r', '\n':
			i.ValueStart += strfind.EndOfWhitespaceSeqBytes(s[i.ValueStart:])

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
			i.ValueEnd = strfind.IndexTermBytes(s, i.ValueStart)
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

func (i *IteratorBytes) expect3(at int, a, b, c byte) (err bool) {
	return len(i.src)-at < 3 ||
		i.src[at] != a ||
		i.src[at+1] != b ||
		i.src[at+2] != c
}

func (i *IteratorBytes) read4(at int, a, b, c, d byte) (err bool) {
	return len(i.src)-at < 4 ||
		i.src[at] != a ||
		i.src[at+1] != b ||
		i.src[at+2] != c ||
		i.src[at+3] != d
}

func (i *IteratorBytes) scanWithCachedPath(
	s []byte,
	fn func(*IteratorBytes) bool,
) (err bool) {
	i.cachedPath = append(i.cachedPath, '.')

	for i.ValueStart < len(s) {
		switch s[i.ValueStart] {
		case ' ', '\t', '\r', '\n':
			i.ValueStart += strfind.EndOfWhitespaceSeqBytes(s[i.ValueStart:])

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
			i.ValueEnd = strfind.IndexTermBytes(s, i.ValueStart)
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
				var e []byte
				if i.escapePath {
					e = keyescape.EscapeAppend(nil, s[i.KeyStart:i.KeyEnd])
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

func (i *IteratorBytes) parseColumn(s []byte) (index int, err bool) {
	if len(s) > 0 && s[0] == ':' {
		return 0, false
	}
	for j, c := range s {
		if c == ':' {
			return j, false
		} else if !isSpaceByte(c) {
			break
		}
	}
	i.errCode = ErrorCodeUnexpectedToken
	return -1, true
}

func isSpaceByte(b byte) bool {
	switch b {
	case ' ', '\r', '\n', '\t':
		return true
	}
	return false
}

func (i *IteratorBytes) cacheCleanupAfterContainer() {
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

func (i *IteratorBytes) callFn(
	t *stack.Node,
	fn func(i *IteratorBytes) (err bool),
) bool {
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

func (i *IteratorBytes) callFnWithPathCache(
	t *stack.Node,
	nonComposite bool,
	fn func(i *IteratorBytes) bool,
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
