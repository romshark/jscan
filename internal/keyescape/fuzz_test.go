package keyescape

import (
	"strings"
	"testing"
)

// ref is the obvious RFC 6901 escaper used as the reference.
func ref(key string) string {
	var dest []byte
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
	return string(dest)
}

var refReplacer = strings.NewReplacer("~", "~0", "/", "~1")

// checkAll runs every implementation variant against key in both
// instantiations of the type parameter and compares against ref.
func checkAll(t *testing.T, key string) {
	t.Helper()
	expect := ref(key)

	check := func(name, actual string) {
		t.Helper()
		if actual != expect {
			t.Fatalf("%s(%q): received %q, expected %q", name, key, actual, expect)
		}
	}

	check("Append[string]", string(Append(nil, key)))
	check("Append[[]byte]", string(Append(nil, []byte(key))))

	check("variantCheckAndReplace[string]",
		string(variantCheckAndReplace(nil, key)))
	check("variantCheckAndReplace[[]byte]",
		string(variantCheckAndReplace(nil, []byte(key))))

	check("variantCheckAndReplaceUnrolled[string]",
		string(variantCheckAndReplaceUnrolled(nil, key)))
	check("variantCheckAndReplaceUnrolled[[]byte]",
		string(variantCheckAndReplaceUnrolled(nil, []byte(key))))

	check("variantStdReplacer[string]",
		string(variantStdReplacer(refReplacer, nil, key)))
	check("variantStdReplacer[[]byte]",
		string(variantStdReplacer(refReplacer, nil, []byte(key))))

	// Appending must leave whatever dest already holds untouched.
	if actual := string(Append([]byte("PRE"), key)); actual != "PRE"+expect {
		t.Fatalf("Append(dest, %q): received %q, expected %q",
			key, actual, "PRE"+expect)
	}
}

// enumerate calls checkAll for every string up to maxLen over alphabet.
func enumerate(t *testing.T, alphabet string, maxLen int) {
	t.Helper()
	buf := make([]byte, 0, maxLen)
	var rec func(depth int)
	rec = func(depth int) {
		checkAll(t, string(buf))
		if depth == 0 {
			return
		}
		for i := 0; i < len(alphabet); i++ {
			buf = append(buf, alphabet[i])
			rec(depth - 1)
			buf = buf[:len(buf)-1]
		}
	}
	rec(maxLen)
}

// TestExhaustive enumerates every string up to length 10 over an alphabet
// of both special characters and a plain byte, which fully covers the
// sub-8-byte and 8-to-11-byte strategies of Append.
func TestExhaustive(t *testing.T) {
	enumerate(t, "~/a", 10)
}

// TestExhaustiveAtBoundary enumerates every string up to length 15 over a
// two-character alphabet, covering the 12-byte boundary at which Append
// switches to the accelerated search and the unrolled 8-byte loop.
func TestExhaustiveAtBoundary(t *testing.T) {
	enumerate(t, "~a", 15)
	enumerate(t, "/a", 15)
	enumerate(t, "~/", 15)
}

// TestSpecialAtEveryOffset covers keys longer than the table tests reach,
// with a special character at every position, which exercises each branch
// of the unrolled 4- and 8-byte loops individually.
func TestSpecialAtEveryOffset(t *testing.T) {
	for _, special := range []string{"~", "/"} {
		for n := 0; n <= 40; n++ {
			base := strings.Repeat("a", n)
			for pos := 0; pos <= n; pos++ {
				checkAll(t, base[:pos]+special+base[pos:])
			}
		}
	}
	// Two special characters at every pair of positions.
	const n = 24
	base := strings.Repeat("a", n)
	for i := 0; i <= n; i++ {
		k := base[:i] + "~" + base[i:]
		for j := 0; j <= len(k); j++ {
			checkAll(t, k[:j]+"/"+k[j:])
		}
	}
}

// FuzzAppend cross-checks every variant against ref.
func FuzzAppend(f *testing.F) {
	for _, s := range []string{
		"", "~", "/", "~0", "~1", "a", "a/b", "a~b",
		"abcdefgh", "abcdefgh/", "abcdefghijkl", "abcdefghijkl~",
		"~~~~~~~~~~~~", "////////////", "0123456789012345678901234567~",
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, key string) {
		if len(key) > 512 {
			return
		}
		checkAll(t, key)
	})
}
