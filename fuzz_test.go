package gomoji

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/rivo/uniseg"
)

// fuzzSeeds covers the shapes that make grapheme handling hard. Invisible
// characters are written as escapes so the corpus stays readable: \u200d is a
// zero-width joiner, \ufe0e and \ufe0f are the emoji presentation selectors,
// and \ufe00 and \U000e0100 are variation selectors that carry meaning outside
// of emoji.
var fuzzSeeds = []string{
	// Plain text and simple emojis.
	"",
	"hello world",
	"\U0001f381",
	"a\U0001f381b",

	// Multi-code-point emojis: ZWJ sequences, flags and skin tones.
	"\U0001f469\u200d\U0001f469\u200d\U0001f467\u200d\U0001f466",
	"\U0001f3f3\ufe0f\u200d\U0001f308",
	"\U0001f1fa\U0001f1e6",
	"\U0001f1fa\U0001f1e6\U0001f1fa",
	"\U0001f44b\U0001f3fd",
	"#\ufe0f\u20e3",

	// Presentation selectors: stray, doubled, and attached to non-emoji text.
	"☺\ufe0f\ufe0f",
	"\ufe0f",
	"\ufe0e",
	"a\ufe0fb",
	"\U0001f3f3\ufe0fX\u200d\U0001f308",

	// Variation selectors that are not emoji presentation selectors and must
	// survive processing untouched.
	"漢\ufe00字",
	"\U000e0100",

	// A dangling ZWJ makes the emoji part of a larger grapheme cluster that is
	// not itself a known emoji, so the Replace* family leaves it in place even
	// though ContainsEmoji and FindAll still report it.
	"\U0001f469\u200d",
	"a\U0001f469\u200db",

	// Invalid UTF-8. Deleting an emoji here can splice the surrounding stray
	// bytes into a rune that was never in the input, and that rune can itself
	// be an emoji, so a single pass is not a fixed point.
	"\xff\xfe",
	"\U0001f381\xff",
	"000\xf0\U0001f618\x98\x94\x94",
	"\xf0\U0001f618\x9f\x98\x94",
}

func addFuzzSeeds(f *testing.F) {
	f.Helper()

	for _, s := range fuzzSeeds {
		f.Add(s)
	}
}

// FuzzRemoveEmojis asserts that removal only ever deletes text: for valid
// UTF-8 the result is a rune subsequence of the input, never grows, stays
// valid UTF-8, and reaches a fixed point after one pass.
func FuzzRemoveEmojis(f *testing.F) {
	addFuzzSeeds(f)

	f.Fuzz(func(t *testing.T, s string) {
		got := RemoveEmojis(s)

		if len(got) > len(s) {
			t.Fatalf("RemoveEmojis(%q) = %q grew the input: %d > %d bytes", s, got, len(got), len(s))
		}

		// Every guarantee below holds for valid input only: deleting bytes from
		// invalid UTF-8 can splice neighbouring stray bytes into a rune that
		// was never in the input, and that rune can itself be an emoji.
		if !utf8.ValidString(s) {
			return
		}

		if !utf8.ValidString(got) {
			t.Fatalf("RemoveEmojis(%q) = %q produced invalid UTF-8 from valid input", s, got)
		}

		if !isRuneSubsequence(got, s) {
			t.Fatalf("RemoveEmojis(%q) = %q is not a subsequence of the input", s, got)
		}

		if again := RemoveEmojis(got); again != got {
			t.Fatalf("RemoveEmojis is not idempotent: %q -> %q -> %q", s, got, again)
		}
	})
}

// FuzzFindAll asserts that every reported emoji is real: it is present in the
// input verbatim, resolves through the public GetInfo lookup, and is reported
// at most once.
func FuzzFindAll(f *testing.F) {
	addFuzzSeeds(f)

	f.Fuzz(func(t *testing.T, s string) {
		found := FindAll(s)

		seen := make(map[string]struct{}, len(found))
		for _, em := range found {
			if !strings.Contains(s, em.Character) {
				t.Fatalf("FindAll(%q) reported %q, which is not present in the input", s, em.Character)
			}

			if _, err := GetInfo(em.Character); err != nil {
				t.Fatalf("FindAll(%q) reported %q, which GetInfo rejects: %v", s, em.Character, err)
			}

			if _, dup := seen[em.Character]; dup {
				t.Fatalf("FindAll(%q) reported %q more than once", s, em.Character)
			}
			seen[em.Character] = struct{}{}
		}

		if len(found) > 0 && !ContainsEmoji(s) {
			t.Fatalf("FindAll(%q) found %d emojis but ContainsEmoji reported false", s, len(found))
		}
	})
}

// FuzzReplaceEmojisWithSlug asserts that slug replacement stays well-formed
// and that every emoji it is contracted to replace contributes its slug to
// the output.
func FuzzReplaceEmojisWithSlug(f *testing.F) {
	addFuzzSeeds(f)

	f.Fuzz(func(t *testing.T, s string) {
		got := ReplaceEmojisWithSlug(s)

		if utf8.ValidString(s) && !utf8.ValidString(got) {
			t.Fatalf("ReplaceEmojisWithSlug(%q) = %q produced invalid UTF-8 from valid input", s, got)
		}

		// Replacement operates on whole grapheme clusters, so the guarantee is
		// scoped to clusters that are themselves a known emoji. An emoji that
		// merely starts a larger cluster, as with a dangling ZWJ, is left in
		// place even though FindAll reports it.
		gr := uniseg.NewGraphemes(s)
		for gr.Next() {
			em, ok := lookupEmoji(gr.Str())
			if !ok {
				continue
			}

			if strings.Contains(s, em.Slug) {
				// The slug already occurs in the input, so its presence in the
				// output proves nothing.
				continue
			}

			if !strings.Contains(got, em.Slug) {
				t.Fatalf("ReplaceEmojisWithSlug(%q) = %q is missing the slug %q for cluster %q", s, got, em.Slug, em.Character)
			}
		}
	})
}

// FuzzCollectAll asserts that CollectAll never reports an emoji the public
// lookup does not recognise.
func FuzzCollectAll(f *testing.F) {
	addFuzzSeeds(f)

	f.Fuzz(func(t *testing.T, s string) {
		for _, em := range CollectAll(s) {
			if _, err := GetInfo(em.Character); err != nil {
				t.Fatalf("CollectAll(%q) reported %q, which GetInfo rejects: %v", s, em.Character, err)
			}
		}
	})
}

// isRuneSubsequence reports whether sub can be obtained from s by deleting
// runes, i.e. whether every rune of sub appears in s in the same order.
func isRuneSubsequence(sub, s string) bool {
	rest := s
	for _, r := range sub {
		idx := strings.IndexRune(rest, r)
		if idx < 0 {
			return false
		}

		_, size := utf8.DecodeRuneInString(rest[idx:])
		rest = rest[idx+size:]
	}

	return true
}
