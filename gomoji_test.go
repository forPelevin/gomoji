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
					Slug:        "e3-0-butterfly",
					Character:   "🦋",
					UnicodeName: "E3.0 butterfly",
					CodePoint:   "1F98B",
					Group:       "animals-nature",
					SubGroup:    "animal-bug",
				},
				{
					Slug:        "roll-of-paper",
					Character:   "🧻",
					UnicodeName: "roll of paper",
					CodePoint:   "1F9FB",
					Group:       "objects",
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
