package gomoji

import (
	"errors"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/rivo/uniseg"
)

// errors
var (
	ErrStrNotEmoji = errors.New("the string is not emoji")
)

// Emoji is an entity that represents comprehensive emoji info.
type Emoji struct {
	Slug        string `json:"slug"`
	Character   string `json:"character"`
	UnicodeName string `json:"unicode_name"`
	CodePoint   string `json:"code_point"`
	Group       string `json:"group"`
	SubGroup    string `json:"sub_group"`
}

// ContainsEmoji checks whether given string contains emoji or not. It uses local emoji list as provider.
func ContainsEmoji(s string) bool {
	if isASCII(s) {
		return false
	}

	// Single-rune emojis are the common case and are much cheaper to find
	// than running grapheme segmentation.
	for _, r := range s {
		if _, ok := emojiMap[string(r)]; ok {
			return true
		}
	}

	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		if _, ok := lookupEmoji(gr.Str()); ok {
			return true
		}
	}

	return false
}

// AllEmojis gets all emojis from provider.
func AllEmojis() []Emoji {
	return emojiMapToSlice(emojiMap)
}

// RemoveEmojis removes all emojis from the s string and returns a new string.
func RemoveEmojis(s string) string {
	return ReplaceEmojisWithFunc(s, nil)
}

// ReplaceEmojisWith replaces all emojis from the s string with the specified rune and returns a new string.
func ReplaceEmojisWith(s string, c rune) string {
	replacerStr := string(c)
	return ReplaceEmojisWithFunc(s, func(_ Emoji) string {
		return replacerStr
	})
}

// ReplaceEmojisWithSlug replaces all emojis from the s string with the emoji's slug and returns a new string.
func ReplaceEmojisWithSlug(s string) string {
	return ReplaceEmojisWithFunc(s, func(em Emoji) string {
		return em.Slug
	})
}

type replacerFn func(e Emoji) string

// ReplaceEmojisWithFunc replaces all emojis from the s string with the result of the replacerFn function
// and returns a new string. Stray emoji presentation selectors (U+FE0E, U+FE0F) are removed as well.
func ReplaceEmojisWithFunc(s string, replacer replacerFn) string {
	if isASCII(s) {
		return s
	}

	var buf strings.Builder

	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		cluster := gr.Str()

		if em, ok := lookupEmoji(cluster); ok {
			if replacer != nil {
				buf.WriteString(replacer(em))
			}
			continue
		}

		buf.WriteString(removePresentationSelectors(cluster))
	}

	return buf.String()
}

// GetInfo returns a gomoji.Emoji model representation of provided emoji.
// If the emoji was not found, it returns the gomoji.ErrStrNotEmoji error
func GetInfo(emoji string) (Emoji, error) {
	em, ok := lookupEmoji(emoji)
	if !ok {
		return Emoji{}, ErrStrNotEmoji
	}

	return em, nil
}

// CollectAll finds all emojis in given string. Unlike FindAll, this does not
// distinct repeating occurrences of emoji. If there are no emojis it returns a nil-slice.
func CollectAll(s string) []Emoji {
	if isASCII(s) {
		return nil
	}

	var emojis []Emoji

	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		if em, ok := emojiMap[gr.Str()]; ok {
			emojis = append(emojis, em)
			continue
		}

		start, end := gr.Positions()
		for _, r := range s[start:end] {
			if em, ok := emojiMap[string(r)]; ok {
				emojis = append(emojis, em)
			}
		}
	}

	return emojis
}

// FindAll finds all emojis in given string. If there are no emojis it returns a nil-slice.
func FindAll(s string) []Emoji {
	if isASCII(s) {
		return nil
	}

	emojis := make(map[string]Emoji)
	gr := uniseg.NewGraphemes(s)

	for gr.Next() {
		cluster := gr.Str()
		if em, ok := emojiMap[cluster]; ok {
			emojis[cluster] = em
			continue
		}

		// Sub-cluster analysis for partial matches
		for i := len(cluster); i > 0; i-- {
			sub := cluster[:i]
			if em, ok := emojiMap[sub]; ok {
				emojis[sub] = em
				break
			}
		}
	}

	return emojiMapToSlice(emojis)
}

// lookupEmoji finds the emoji entity for s, falling back to a lookup with
// emoji presentation selectors (U+FE0E, U+FE0F) removed, so that forms like
// "emoji + extra U+FE0F" still resolve to the underlying emoji.
func lookupEmoji(s string) (Emoji, bool) {
	if em, ok := emojiMap[s]; ok {
		return em, true
	}

	if stripped := removePresentationSelectors(s); stripped != s {
		em, ok := emojiMap[stripped]
		return em, ok
	}

	return Emoji{}, false
}

// Emoji presentation selectors.
const (
	textVS  = '\uFE0E'
	emojiVS = '\uFE0F'
)

// removePresentationSelectors removes the emoji presentation selectors
// (U+FE0E text style, U+FE0F emoji style) from s. Other variation selectors
// are kept intact because they carry meaning in non-emoji text, e.g.
// ideographic and Mongolian variation sequences.
func removePresentationSelectors(s string) string {
	if !strings.ContainsRune(s, textVS) && !strings.ContainsRune(s, emojiVS) {
		return s
	}

	return strings.Map(func(r rune) rune {
		if r == textVS || r == emojiVS {
			return -1
		}
		return r
	}, s)
}

// isASCII reports whether s contains only ASCII bytes. Every emoji contains
// at least one multi-byte rune, so ASCII-only strings can skip grapheme
// segmentation entirely.
func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}

	return true
}

func emojiMapToSlice(em map[string]Emoji) []Emoji {
	if len(em) == 0 {
		return nil
	}

	emojis := make([]Emoji, 0, len(em))
	for _, emoji := range em {
		emojis = append(emojis, emoji)
	}

	sort.Slice(emojis, func(i, j int) bool {
		return emojis[i].Character < emojis[j].Character
	})

	return emojis
}
