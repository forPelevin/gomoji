package gomoji

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
		if _, ok := emojiMap[r]; ok {
			return true
		}
	}

	return false
}

// AllEmojis gets all emojis from provider.
func AllEmojis() []Emoji {
	return emojiMapToSlice(emojiMap)
}

// FindAll finds all emojis in given string. If there are no emojis it returns a nil-slice.
func FindAll(s string) []Emoji {
	var emojis []Emoji

	for _, r := range s {
		if em, ok := emojiMap[r]; ok {
			emojis = append(emojis, em)
		}
	}

	return emojis
}

func emojiMapToSlice(em map[rune]Emoji) []Emoji {
	emojis := make([]Emoji, 0, len(em))
	for _, emoji := range em {
		emojis = append(emojis, emoji)
	}

	return emojis
}
