package gomoji_test

import (
	"fmt"

	"github.com/forPelevin/gomoji"
)

func ExampleContainsEmoji() {
	fmt.Println(gomoji.ContainsEmoji("hello world"))
	fmt.Println(gomoji.ContainsEmoji("hello world 🤗"))
	// Output:
	// false
	// true
}

func ExampleRemoveEmojis() {
	fmt.Println(gomoji.RemoveEmojis("🧖 hello 🦋 world"))
	// Output:  hello  world
}

func ExampleReplaceEmojisWith() {
	fmt.Println(gomoji.ReplaceEmojisWith("🧖 hello 🦋 world", '_'))
	// Output: _ hello _ world
}

func ExampleReplaceEmojisWithSlug() {
	fmt.Println(gomoji.ReplaceEmojisWithSlug("hello 🦋 world"))
	// Output: hello butterfly world
}

func ExampleReplaceEmojisWithFunc() {
	fmt.Println(gomoji.ReplaceEmojisWithFunc("hello 🦋 world", func(em gomoji.Emoji) string {
		return ":" + em.Slug + ":"
	}))
	// Output: hello :butterfly: world
}

func ExampleFindAll() {
	for _, em := range gomoji.FindAll("hello 🦋 world 🐝") {
		fmt.Println(em.Slug)
	}
	// Output:
	// honeybee
	// butterfly
}

func ExampleGetInfo() {
	em, err := gomoji.GetInfo("🦋")
	if err != nil {
		panic(err)
	}
	fmt.Println(em.Slug, "|", em.UnicodeName, "|", em.Group, "|", em.SubGroup)
	// Output: butterfly | E3.0 butterfly | Animals & Nature | animal-bug
}
