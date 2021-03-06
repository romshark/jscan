package jscan

import (
	"bytes"
	"strconv"
	"sync"
	"unicode/utf8"

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

func getItrBytesFromPool(
	s []byte,
	escapePath bool,
	startIndex int,
) *IteratorBytes {
	i := itrPoolBytes.Get().(*IteratorBytes)
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
	i.keyBuf = i.keyBuf[:0]
	return i
}

// getError returns the stringified error, if any.
func (i *IteratorBytes) getError(c ErrorCode) ErrorBytes {
	return ErrorBytes{
		Code:  c,
		Src:   i.src,
		Index: i.ValueStart,
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
	if e.Index < len(e.Src) {
		r, _ := utf8.DecodeRune(e.Src[e.Index:])
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

// ValidateBytes returns an error if s is invalid JSON.
func ValidateBytes(s []byte) ErrorBytes {
	startIndex, illegal := strfind.EndOfWhitespaceSeqBytes(s)
	if illegal {
		return ErrorBytes{
			Src:   s,
			Index: startIndex,
			Code:  ErrorCodeIllegalControlCharacter,
		}
	}

	if startIndex >= len(s) {
		return ErrorBytes{
			Src:   s,
			Index: startIndex,
			Code:  ErrorCodeUnexpectedEOF,
		}
	}

	i := getItrBytesFromPool(s, false, startIndex)
	defer itrPoolBytes.Put(i)

	for i.ValueStart < len(s) {
		switch s[i.ValueStart] {
		case ' ', '\t', '\r', '\n':
			e, illegal := strfind.EndOfWhitespaceSeqBytes(s[i.ValueStart:])
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
			i.ValueEnd, err = jsonnum.ParseBytes(i.src[i.ValueStart:])
			if err {
				return i.getError(ErrorCodeMalformedNumber)
			}
			i.ValueStart += i.ValueEnd
			i.KeyStart, i.KeyEnd = -1, -1

		case '"': // String
			i.ValueStart++
			i.ValueEnd = strfind.IndexTermBytes(s, i.ValueStart)
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

	return ErrorBytes{}
}

// ValidBytes returns true if s is valid JSON, otherwise returns false.
func ValidBytes(s []byte) bool {
	return !ValidateBytes(s).IsErr()
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
	startIndex, illegal := strfind.EndOfWhitespaceSeqBytes(s)
	if illegal {
		return ErrorBytes{
			Src:   s,
			Index: startIndex,
			Code:  ErrorCodeIllegalControlCharacter,
		}
	}

	if startIndex >= len(s) {
		return ErrorBytes{
			Src:   s,
			Index: startIndex,
			Code:  ErrorCodeUnexpectedEOF,
		}
	}

	i := getItrBytesFromPool(s, o.EscapePath, startIndex)
	defer itrPoolBytes.Put(i)

	if o.CachePath {
		return i.scanWithCachedPath(s, fn)
	}
	return i.scan(s, fn)
}

func (i *IteratorBytes) scan(
	s []byte,
	fn func(*IteratorBytes) (err bool),
) ErrorBytes {
	for i.ValueStart < len(s) {
		switch s[i.ValueStart] {
		case ' ', '\t', '\r', '\n':
			e, illegal := strfind.EndOfWhitespaceSeqBytes(s[i.ValueStart:])
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
			i.ValueEnd, err = jsonnum.ParseBytes(i.src[i.ValueStart:])
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
			i.ValueEnd = strfind.IndexTermBytes(s, i.ValueStart)
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

	return ErrorBytes{}
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
) ErrorBytes {
	i.cachedPath = append(i.cachedPath, '.')

	for i.ValueStart < len(s) {
		switch s[i.ValueStart] {
		case ' ', '\t', '\r', '\n':
			e, illegal := strfind.EndOfWhitespaceSeqBytes(s[i.ValueStart:])
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
			i.ValueEnd, err = jsonnum.ParseBytes(i.src[i.ValueStart:])
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
			i.ValueEnd = strfind.IndexTermBytes(s, i.ValueStart)
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

	return ErrorBytes{}
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

func (i *IteratorBytes) callFnWithPathCache(
	t *stack.Node,
	nonComposite bool,
	fn func(i *IteratorBytes) bool,
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
