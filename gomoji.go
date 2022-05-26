package gomoji

import (
	"bytes"
	"errors"
	"strings"
	"unicode"

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
	for _, r := range s {
		if _, ok := emojiMap[string(r)]; ok {
			return true
		}
	}

	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		if _, ok := emojiMap[gr.Str()]; ok {
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
	cleanBuf := bytes.Buffer{}

	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		if _, ok := emojiMap[gr.Str()]; !ok {
			cleanBuf.Write(gr.Bytes())
		}
	}

	res := cleanBuf
	res.Reset()
	for _, r := range cleanBuf.String() {
		if _, ok := emojiMap[string(r)]; !ok {
			res.WriteRune(r)
		}
	}

	return strings.TrimFunc(res.String(), func(r rune) bool {
		return unicode.IsSpace(r) || !unicode.IsGraphic(r) || !unicode.IsPrint(r) || unicode.In(r, unicode.Variation_Selector)
	})
}

// GetInfo returns a gomoji.Emoji model representation of provided emoji.
// If the emoji was not found, it returns the gomoji.ErrStrNotEmoji error
func GetInfo(emoji string) (Emoji, error) {
	em, ok := emojiMap[emoji]
	if !ok {
		return Emoji{}, ErrStrNotEmoji
	}

	return em, nil
}

// CollectAll finds all emojis in given string. Unlike FindAll, this does not
// distinct repeating occurrences of emoji. If there are no emojis it returns a nil-slice.
func CollectAll(s string) []Emoji {

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
	emojis := make(map[string]Emoji)

	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		if em, ok := emojiMap[gr.Str()]; ok {
			emojis[gr.Str()] = em
		}
	}

	for _, r := range s {
		if em, ok := emojiMap[string(r)]; ok {
			emojis[string(r)] = em
		}
	}

	return emojiMapToSlice(emojis)
}

func emojiMapToSlice(em map[string]Emoji) []Emoji {
	var emojis []Emoji
	for _, emoji := range em {
		emojis = append(emojis, emoji)
	}

	return emojis
}
