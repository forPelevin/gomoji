package gomoji

import (
	"reflect"
	"testing"
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsEmoji(tt.inputStr); got != tt.want {
				t.Errorf("ContainsEmoji() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkContainsEmojiParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ContainsEmoji("Hi \U0001F970")
		}
	})
}

func BenchmarkContainsEmoji(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ContainsEmoji("Hi \U0001F970")
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveEmojis(tt.inputStr); got != tt.want {
				t.Errorf("RemoveEmojis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkRemoveEmojisParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			RemoveEmojis("\U0001F96F Hi \U0001F970")
		}
	})
}

func BenchmarkRemoveEmojis(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RemoveEmojis("\U0001F96F Hi \U0001F970")
	}
}

func TestGetInfo(t *testing.T) {
	tests := []struct {
		name       string
		inputEmoji string
		want       Emoji
		wantErr    bool
	}{
		{
			name:       "just a number",
			inputEmoji: "1",
			want:       Emoji{},
			wantErr:    true,
		},
		{
			name:       "valid emoji number",
			inputEmoji: "1️⃣",
			want: Emoji{
				Slug:        "keycap:-1",
				Character:   "1️⃣",
				UnicodeName: "keycap: 1",
				CodePoint:   "U+0031 U+FE0F U+20E3",
				Group:       "Symbols",
				SubGroup:    "keycap",
			},
			wantErr: false,
		},
		{
			name:       "unicode 14",
			inputEmoji: "\U0001FAAC",
			want: Emoji{
				Slug:        "⊛-hamsa",
				Character:   "🪬",
				UnicodeName: "⊛ hamsa",
				CodePoint:   "U+1FAAC",
				Group:       "Activities",
				SubGroup:    "game",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInfo(tt.inputEmoji)
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
			GetInfo("\U0001F96F") // nolint:errcheck
		}
	})
}

func BenchmarkGetInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetInfo("\U0001F96F") // nolint:errcheck
	}
}

func TestFindAll(t *testing.T) {
	tests := []struct {
		name     string
		inputStr string
		want     []Emoji
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
			want: []Emoji{
				{
					Slug:        "butterfly",
					Character:   "🦋",
					UnicodeName: "butterfly",
					CodePoint:   "U+1F98B",
					Group:       "Animals & Nature",
					SubGroup:    "animal-bug",
				},
				{
					Slug:        "roll-of-paper",
					Character:   "🧻",
					UnicodeName: "roll of paper",
					CodePoint:   "U+1F9FB",
					Group:       "Objects",
					SubGroup:    "household",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindAll(tt.inputStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkFindAllParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			FindAll("\U0001F96F Hi \U0001F970")
		}
	})
}

func BenchmarkFindAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindAll("\U0001F96F Hi \U0001F970")
	}
}
