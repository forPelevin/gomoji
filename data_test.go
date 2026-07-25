package gomoji

import (
	"strconv"
	"strings"
	"testing"
	"unicode"
)

// canonicalGroups are the group names defined by the Unicode emoji-test.txt
// data file, plus "Component" used for Emoji_Component entries (skin tones,
// hair styles, regional indicators).
var canonicalGroups = map[string]struct{}{
	"Smileys & Emotion": {},
	"People & Body":     {},
	"Animals & Nature":  {},
	"Food & Drink":      {},
	"Travel & Places":   {},
	"Activities":        {},
	"Objects":           {},
	"Symbols":           {},
	"Flags":             {},
	"Component":         {},
}

// TestEmojiDataIntegrity guards the generated data.go against malformed
// updates from the automated emoji data pipeline.
func TestEmojiDataIntegrity(t *testing.T) {
	if len(emojiMap) < 5000 {
		t.Fatalf("emojiMap has only %d entries, expected at least 5000; the data set looks truncated", len(emojiMap))
	}

	for key, em := range emojiMap {
		if key != em.Character {
			t.Errorf("map key %q does not match Character %q", key, em.Character)
		}
		if em.Slug == "" || em.Character == "" || em.UnicodeName == "" || em.CodePoint == "" || em.Group == "" || em.SubGroup == "" {
			t.Errorf("entry %q has empty fields: %+v", key, em)
		}

		if _, ok := canonicalGroups[em.Group]; !ok {
			t.Errorf("entry %q has non-canonical group %q", key, em.Group)
		}

		if strings.ToLower(em.Slug) != em.Slug {
			t.Errorf("entry %q has non-lowercase slug %q", key, em.Slug)
		}
		if strings.ContainsAny(em.Slug, " :,") {
			t.Errorf("entry %q has malformed slug %q", key, em.Slug)
		}
		if containsCJK(em.Slug) || containsCJK(em.UnicodeName) {
			t.Errorf("entry %q has non-English name data: %+v", key, em)
		}

		if got := codePointsToString(t, em.CodePoint); got != em.Character {
			t.Errorf("entry %q: CodePoint %q decodes to %q, which does not match Character %q", key, em.CodePoint, got, em.Character)
		}
	}
}

func codePointsToString(t *testing.T, codePoints string) string {
	t.Helper()

	var b strings.Builder
	for _, cp := range strings.Fields(codePoints) {
		n, err := strconv.ParseUint(cp, 16, 32)
		if err != nil {
			t.Errorf("invalid code point %q: %v", codePoints, err)
			return ""
		}
		b.WriteRune(rune(n))
	}

	return b.String()
}

func containsCJK(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana, unicode.Hangul) {
			return true
		}
	}

	return false
}
