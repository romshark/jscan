package keyescape

import (
	"bytes"
	"strings"
	"unsafe"
)

// Append appends key to dest replacing all occurrences of
// '~' with '~0' and '/' with '~1'.
func Append[S ~[]byte | ~string](dest []byte, key S) []byte {
	// variantCheckAndReplaceUnrolled performed best in benchmarks
	return variantCheckAndReplaceUnrolled[S](dest, key)
}

// variantCheckAndReplace is an implementation variant that checks
// for tilde and slash characters using standard IndexByte
// and if any is found replaces them in an append loop.
func variantCheckAndReplace[S ~[]byte | ~string](dest []byte, key S) []byte {
	var hasTilde, hasSlash bool
	switch k := any(key).(type) {
	case string:
		hasTilde = strings.IndexByte(k, '~') != -1
		hasSlash = strings.IndexByte(k, '/') != -1
	case []byte:
		hasTilde = bytes.IndexByte(k, '~') != -1
		hasSlash = bytes.IndexByte(k, '/') != -1
	}
	if hasTilde || hasSlash {
		for i := 0; i < len(key); i++ {
			switch key[i] {
			case '~':
				dest = append(dest, "~0"...)
			case '/':
				dest = append(dest, "~1"...)
			default:
				dest = append(dest, key[i])
			}
		}
		return dest
	}
	return append(dest, key...)
}

// variantCheckAndReplaceUnrolled is an implementation variant that checks
// for tilde and slash characters using standard IndexByte
// and if any is found replaces them in an unrolled append loop.
func variantCheckAndReplaceUnrolled[S ~[]byte | ~string](
	dest []byte, key S,
) []byte {
	if len(key) < 8 {
		// Naive loop is fastest on small keys
		for ; len(key) > 0; key = key[1:] {
			switch key[0] {
			case '~':
				dest = append(dest, "~0"...)
			case '/':
				dest = append(dest, "~1"...)
			default:
				dest = append(dest, key[0])
			}
		}
		return dest
	}
	if len(key) < 12 {
		// Skip slower tilde/slash search for small keys using 4-batches
		for len(key) > 0 {
			if len(key) > 3 {
				if key[0] == '~' || key[0] == '/' {
					goto CHECK_4
				}
				if key[1] == '~' || key[1] == '/' {
					key, dest = key[1:], append(dest, key[:1]...)
					goto CHECK_4
				}
				if key[2] == '~' || key[2] == '/' {
					key, dest = key[2:], append(dest, key[:2]...)
					goto CHECK_4
				}
				if key[3] == '~' || key[3] == '/' {
					key, dest = key[3:], append(dest, key[:3]...)
					goto CHECK_4
				}
				dest = append(dest, key[:4]...)
				key = key[4:]
				continue
			}
		CHECK_4:
			switch key[0] {
			case '~':
				dest = append(dest, "~0"...)
			case '/':
				dest = append(dest, "~1"...)
			default:
				dest = append(dest, key[0])
			}
			key = key[1:]
		}
		return dest
	}
	// Use accelerated search to avoid iteration if possible
	switch k := any(key).(type) {
	case string:
		if strings.IndexByte(k, '~') != -1 || strings.IndexByte(k, '/') != -1 {
			break
		}
		return append(dest, key...)
	case []byte:
		if bytes.IndexByte(k, '~') != -1 || bytes.IndexByte(k, '/') != -1 {
			break
		}
		return append(dest, key...)
	}

	for len(key) > 0 {
		if len(key) > 7 {
			if key[0] == '~' || key[0] == '/' {
				goto CHECK_8
			}
			if key[1] == '~' || key[1] == '/' {
				key, dest = key[1:], append(dest, key[:1]...)
				goto CHECK_8
			}
			if key[2] == '~' || key[2] == '/' {
				key, dest = key[2:], append(dest, key[:2]...)
				goto CHECK_8
			}
			if key[3] == '~' || key[3] == '/' {
				key, dest = key[3:], append(dest, key[:3]...)
				goto CHECK_8
			}
			if key[4] == '~' || key[4] == '/' {
				key, dest = key[4:], append(dest, key[:4]...)
				goto CHECK_8
			}
			if key[5] == '~' || key[5] == '/' {
				key, dest = key[5:], append(dest, key[:5]...)
				goto CHECK_8
			}
			if key[6] == '~' || key[6] == '/' {
				key, dest = key[6:], append(dest, key[:6]...)
				goto CHECK_8
			}
			if key[7] == '~' || key[7] == '/' {
				key, dest = key[7:], append(dest, key[:7]...)
				goto CHECK_8
			}
			dest = append(dest, key[:8]...)
			key = key[8:]
			continue
		}
	CHECK_8:
		switch key[0] {
		case '~':
			dest = append(dest, "~0"...)
		case '/':
			dest = append(dest, "~1"...)
		default:
			dest = append(dest, key[0])
		}
		key = key[1:]
	}
	return dest
}

// variantStdReplacer uses the standard library strings replacer.
func variantStdReplacer[S ~[]byte | ~string](r *strings.Replacer, dest []byte, key S) []byte {
	switch key := any(key).(type) {
	case []byte:
		return append(dest, r.Replace(unsafeB2S(key))...)
	}
	return append(dest, r.Replace(string(key))...)
}

func unsafeB2S(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
