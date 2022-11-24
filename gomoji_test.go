package gomoji_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/forPelevin/gomoji"
)

func TestContainsEmoji(t *testing.T) {
	tests := []struct {
		name     string
		inputStr string
		want     bool
	}{
		{
			name:     "empty input string",
			inputStr: "",
			want:     false,
		},
		{
			name:     "string without emoji",
			inputStr: "hello! This is a simple string without any emoji",
			want:     false,
		},
		{
			name:     "numbers in string",
			inputStr: "qwerty1",
			want:     false,
		},
		{
			name:     "emoji number in string",
			inputStr: "qwerty1️⃣",
			want:     true,
		},
		{
			name:     "only emoji in string",
			inputStr: `🥰`,
			want:     true,
		},
		{
			name:     "emoji in the middle of a string",
			inputStr: `hi 😀 how r u?`,
			want:     true,
		},
		{
			name:     "emoji in the end of a string",
			inputStr: `hi! how r u doing?🤔`,
			want:     true,
		},
		{
			name:     "heart emoji in string",
			inputStr: "I ❤️ you",
			want:     true,
		},
		{
			name:     "it determines the skintone emojis",
			inputStr: "I 👍🏿 you",
			want:     true,
		},
		{
			name:     "double exclamation mark in text",
			inputStr: "Hello!!",
			want:     false,
		},
		{
			name:     "double exclamation mark emoji in text",
			inputStr: "Hello‼",
			want:     true,
		},
		{
			name:     "emoji keycap # in text",
			inputStr: "Just type #⃣",
			want:     true,
		},
		{
			name:     "text keycap # in text",
			inputStr: "Just type #",
			want:     false,
		},
		{
			name:     "new emoji",
			inputStr: "🆕️ NWT H&M Corduroy Pants in 'Light Beige'",
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gomoji.ContainsEmoji(tt.inputStr); got != tt.want {
				t.Errorf("ContainsEmoji() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkContainsEmojiParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gomoji.ContainsEmoji("Hi \U0001F970")
		}
	})
}

func BenchmarkContainsEmoji(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gomoji.ContainsEmoji("Hi \U0001F970")
	}
}

func TestRemoveEmojis(t *testing.T) {
	tests := []struct {
		name     string
		inputStr string
		want     string
	}{
		{
			name:     "string without emoji",
			inputStr: "string without emoji",
			want:     "string without emoji",
		},
		{
			name:     "string with numbers",
			inputStr: "1qwerty2",
			want:     "1qwerty2",
		},
		{
			name:     "string with emoji numbers",
			inputStr: "1️⃣qwerty2",
			want:     "qwerty2",
		},
		{
			name:     "string with emojis",
			inputStr: "❤️🛶😂",
			want:     "",
		},
		{
			name:     "string with unicode 14 emoji",
			inputStr: "te\U0001FAB7st",
			want:     "test",
		},
		{
			name:     "remove rare emojis",
			inputStr: "🧖 hello 🦋world",
			want:     "hello world",
		},
		{
			name:     "new emoji",
			inputStr: "🆕️ NWT H&M Corduroy Pants in 'Light Beige'",
			want:     "NWT H&M Corduroy Pants in 'Light Beige'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gomoji.RemoveEmojis(tt.inputStr); got != tt.want {
				t.Errorf("RemoveEmojis() = \"%v\", want \"%v\"", got, tt.want)
			}
		})
	}
}

func TestReplaceEmojisWith(t *testing.T) {
	replacementChar := '_'

	tests := []struct {
		name     string
		inputStr string
		want     string
	}{
		{
			name:     "string without emoji",
			inputStr: "string without emoji",
			want:     "string without emoji",
		},
		{
			name:     "string with numbers",
			inputStr: "1qwerty2",
			want:     "1qwerty2",
		},
		{
			name:     "string with emoji number",
			inputStr: "1️⃣qwerty2",
			want:     fmt.Sprintf("%cqwerty2", replacementChar),
		},
		{
			name:     "string with emojis",
			inputStr: "❤️🛶😂",
			want:     fmt.Sprintf("%c%c%c", replacementChar, replacementChar, replacementChar),
		},
		{
			name:     "string with unicode 14 emoji",
			inputStr: "te\U0001FAB7st",
			want:     fmt.Sprintf("te%cst", replacementChar),
		},
		{
			name:     "replace rare emojis",
			inputStr: "🧖 hello 🦋world",
			want:     fmt.Sprintf("%c hello %cworld", replacementChar, replacementChar),
		},
		{
			name:     "new emoji",
			inputStr: "🆕 NWT H&M Corduroy Pants in 'Light Beige'",
			want:     fmt.Sprintf("%c NWT H&M Corduroy Pants in 'Light Beige'", replacementChar),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gomoji.ReplaceEmojisWith(tt.inputStr, replacementChar); got != tt.want {
				t.Errorf("ReplaceEmojisWith() = \"%v\", want \"%v\"", got, tt.want)
			}
		})
	}
}

func TestReplaceEmojisWithSlug(t *testing.T) {
	tests := []struct {
		name     string
		inputStr string
		want     string
	}{
		{
			name:     "string without emoji",
			inputStr: "string without emoji",
			want:     "string without emoji",
		},
		{
			name:     "string with numbers",
			inputStr: "1qwerty2",
			want:     "1qwerty2",
		},
		{
			name:     "string with emoji number",
			inputStr: "1️⃣qwerty2",
			want:     "keycap-1qwerty2",
		},
		{
			name:     "string with emojis",
			inputStr: "❤️🛶😂",
			want:     "red-heartcanoeface-with-tears-of-joy",
		},
		{
			name:     "string with unicode 14 emoji",
			inputStr: "te\U0001FAB7st",
			want:     "telotusst",
		},
		{
			name:     "replace rare emojis",
			inputStr: "🧖 hello 🦋world",
			want:     "person-in-steamy-room hello butterflyworld",
		},
		{
			name:     "new emoji",
			inputStr: "🆕 NWT H&M Corduroy Pants in 'Light Beige'",
			want:     "new-button NWT H&M Corduroy Pants in 'Light Beige'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gomoji.ReplaceEmojisWithSlug(tt.inputStr); got != tt.want {
				t.Errorf("ReplaceEmojisWithSlug() = \"%v\", want \"%v\"", got, tt.want)
			}
		})
	}
}

func BenchmarkReplaceEmojisWithSlugParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gomoji.ReplaceEmojisWithSlug("🧖 hello 🦋world")
		}
	})
}

func BenchmarkReplaceEmojisWithSlug(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gomoji.ReplaceEmojisWithSlug("🧖 hello 🦋world")
	}
}

func TestReplaceEmojisWithFunc(t *testing.T) {
	tests := []struct {
		name     string
		inputStr string
		replacer func(e gomoji.Emoji) string
		want     string
	}{
		{
			name:     "string without emoji",
			inputStr: "string without emoji",
			replacer: func(e gomoji.Emoji) string {
				return ":" + e.Slug
			},
			want: "string without emoji",
		},
		{
			name:     "string with numbers",
			inputStr: "1qwerty2",
			replacer: func(e gomoji.Emoji) string {
				return ":" + e.Slug
			},
			want: "1qwerty2",
		},
		{
			name:     "string with emoji number",
			inputStr: "1️⃣qwerty2",
			replacer: func(e gomoji.Emoji) string {
				return ":" + e.Slug
			},
			want: ":keycap-1qwerty2",
		},
		{
			name:     "string with emojis",
			inputStr: "❤️🛶😂",
			replacer: func(e gomoji.Emoji) string {
				return "emo"
			},
			want: "emoemoemo",
		},
		{
			name:     "string with unicode 14 emoji",
			inputStr: "te\U0001FAB7st",
			replacer: func(e gomoji.Emoji) string {
				return "_"
			},
			want: "te_st",
		},
		{
			name:     "replace rare emojis",
			inputStr: "🧖 hello 🦋world",
			replacer: func(e gomoji.Emoji) string {
				return e.SubGroup
			},
			want: "person-activity hello animal-bugworld",
		},
		{
			name:     "new emoji",
			inputStr: "🆕 NWT H&M Corduroy Pants in 'Light Beige'",
			replacer: func(e gomoji.Emoji) string {
				return e.Slug
			},
			want: "new-button NWT H&M Corduroy Pants in 'Light Beige'",
		},
		{
			name:     "replacer is nil, so all emojis are simply removed",
			inputStr: "🧖 hello 🦋world",
			want:     "hello world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gomoji.ReplaceEmojisWithFunc(tt.inputStr, tt.replacer); got != tt.want {
				t.Errorf("ReplaceEmojisWithFunc() = \"%v\", want \"%v\"", got, tt.want)
			}
		})
	}
}

func BenchmarkRemoveEmojisParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gomoji.RemoveEmojis("\U0001F96F Hi \U0001F970")
		}
	})
}

func BenchmarkRemoveEmojis(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gomoji.RemoveEmojis("\U0001F96F Hi \U0001F970")
	}
}

func TestGetInfo(t *testing.T) {
	tests := []struct {
		name       string
		inputEmoji string
		want       gomoji.Emoji
		wantErr    bool
	}{
		{
			name:       "just a number",
			inputEmoji: "1",
			want:       gomoji.Emoji{},
			wantErr:    true,
		},
		{
			name:       "valid emoji number",
			inputEmoji: "1️⃣",
			want: gomoji.Emoji{
				Slug:        "keycap-1",
				Character:   "1️⃣",
				UnicodeName: "E0.6 keycap: 1",
				CodePoint:   "0031 FE0F 20E3",
				Group:       "Symbols",
				SubGroup:    "keycap",
			},
			wantErr: false,
		},
		{
			name:       "unicode 14",
			inputEmoji: "\U0001FAAC",
			want: gomoji.Emoji{
				Slug:        "hamsa",
				Character:   "🪬",
				UnicodeName: "E14.0 hamsa",
				CodePoint:   "1FAAC",
				Group:       "Objects",
				SubGroup:    "other-object",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gomoji.GetInfo(tt.inputEmoji)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkGetInfoParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gomoji.GetInfo("\U0001F96F") // nolint:errcheck
		}
	})
}

func BenchmarkGetInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gomoji.GetInfo("\U0001F96F") // nolint:errcheck
	}
}

func TestFindAll(t *testing.T) {
	tests := []struct {
		name     string
		inputStr string
		want     []gomoji.Emoji
	}{
		{
			name:     "empty string",
			inputStr: "",
			want:     nil,
		},
		{
			name:     "string without emoji",
			inputStr: "hello world",
			want:     nil,
		},
		{
			name:     "string with 2 emoji",
			inputStr: "hello 🦋 world \U0001F9FB",
			want: []gomoji.Emoji{
				{
					Slug:        "butterfly",
					Character:   "🦋",
					UnicodeName: "E3.0 butterfly",
					CodePoint:   "1F98B",
					Group:       "Animals & Nature",
					SubGroup:    "animal-bug",
				},
				{
					Slug:        "roll-of-paper",
					Character:   "🧻",
					UnicodeName: "E11.0 roll of paper",
					CodePoint:   "1F9FB",
					Group:       "Objects",
					SubGroup:    "household",
				},
			},
		},
		{
			name:     "string with 1 emoji",
			inputStr: "🆕️ NWT H&M Corduroy Pants in 'Light Beige'",
			want: []gomoji.Emoji{
				{
					Slug:        "new-button",
					Character:   "🆕",
					UnicodeName: "E0.6 NEW button",
					CodePoint:   "1F195",
					Group:       "Symbols",
					SubGroup:    "alphanum",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gomoji.FindAll(tt.inputStr)

			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i].Character < tt.want[j].Character
			})
			sort.Slice(got, func(i, j int) bool {
				return got[i].Character < got[j].Character
			})

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkFindAllParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gomoji.FindAll("\U0001F96F Hi \U0001F970")
		}
	})
}

func BenchmarkFindAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gomoji.FindAll("\U0001F96F Hi \U0001F970")
	}
}

func TestCollectAll(t *testing.T) {
	tests := []struct {
		name     string
		inputStr string
		want     []gomoji.Emoji
	}{
		{
			name:     "empty string",
			inputStr: "",
			want:     nil,
		},
		{
			name:     "string without emoji",
			inputStr: "hello world",
			want:     nil,
		},
		{
			name:     "string with 2 emoji",
			inputStr: "hello 🦋 world \U0001F9FB",
			want: []gomoji.Emoji{
				{
					Slug:        "butterfly",
					Character:   "🦋",
					UnicodeName: "E3.0 butterfly",
					CodePoint:   "1F98B",
					Group:       "Animals & Nature",
					SubGroup:    "animal-bug",
				},
				{
					Slug:        "roll-of-paper",
					Character:   "🧻",
					UnicodeName: "E11.0 roll of paper",
					CodePoint:   "1F9FB",
					Group:       "Objects",
					SubGroup:    "household",
				},
			},
		},
		{
			name:     "string with 1 emoji",
			inputStr: "🆕️ NWT H&M Corduroy Pants in 'Light Beige'",
			want: []gomoji.Emoji{
				{
					Slug:        "new-button",
					Character:   "🆕",
					UnicodeName: "E0.6 NEW button",
					CodePoint:   "1F195",
					Group:       "Symbols",
					SubGroup:    "alphanum",
				},
			},
		},
		{
			name:     "string with 6 emoji, mixed and repeating",
			inputStr: "🆕️ NWT H&M Corduroy 🧻🦋🧻 Pants in 'Light Beige'🦋🆕",
			want: []gomoji.Emoji{
				{
					Slug:        "new-button",
					Character:   "🆕",
					UnicodeName: "E0.6 NEW button",
					CodePoint:   "1F195",
					Group:       "Symbols",
					SubGroup:    "alphanum",
				},
				{
					Slug:        "roll-of-paper",
					Character:   "🧻",
					UnicodeName: "E11.0 roll of paper",
					CodePoint:   "1F9FB",
					Group:       "Objects",
					SubGroup:    "household",
				},
				{
					Slug:        "butterfly",
					Character:   "🦋",
					UnicodeName: "E3.0 butterfly",
					CodePoint:   "1F98B",
					Group:       "Animals & Nature",
					SubGroup:    "animal-bug",
				},
				{
					Slug:        "roll-of-paper",
					Character:   "🧻",
					UnicodeName: "E11.0 roll of paper",
					CodePoint:   "1F9FB",
					Group:       "Objects",
					SubGroup:    "household",
				},
				{
					Slug:        "butterfly",
					Character:   "🦋",
					UnicodeName: "E3.0 butterfly",
					CodePoint:   "1F98B",
					Group:       "Animals & Nature",
					SubGroup:    "animal-bug",
				},
				{
					Slug:        "new-button",
					Character:   "🆕",
					UnicodeName: "E0.6 NEW button",
					CodePoint:   "1F195",
					Group:       "Symbols",
					SubGroup:    "alphanum",
				},
			},
		},
		{
			name:     "regional indicators",
			inputStr: "🇦 🇧 🇨",
			want: []gomoji.Emoji{
				{
					Slug:        "regional-indicator-symbol-letter-A",
					Character:   "🇦",
					UnicodeName: "regional indicator symbol letter A",
					CodePoint:   "1F1E6",
					Group:       "symbols",
					SubGroup:    "symbols",
				},
				{
					Slug:        "regional-indicator-symbol-letter-B",
					Character:   "🇧",
					UnicodeName: "regional indicator symbol letter B",
					CodePoint:   "1F1E7",
					Group:       "symbols",
					SubGroup:    "symbols",
				},
				{
					Slug:        "regional-indicator-symbol-letter-C",
					Character:   "🇨",
					UnicodeName: "regional indicator symbol letter C",
					CodePoint:   "1F1E8",
					Group:       "symbols",
					SubGroup:    "symbols",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gomoji.CollectAll(tt.inputStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollectAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
