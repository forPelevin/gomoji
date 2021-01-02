package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/forPelevin/gomoji"
)

var (
	accessKey = flag.String("accessKey", ":8080", "Access key to interact with Open Emoji API. You can get it from https://emoji-api.com.")
)

func main() {
	flag.Parse()

	s := gomoji.NewService(gomoji.NewOpenEmojiProvider(*accessKey))

	containsEmoji(s)
	allEmojis(s)
}

func containsEmoji(s *gomoji.Service) {
	stringWithoutEmoji := "Hello world!"
	res, err := s.ContainsEmoji(context.Background(), stringWithoutEmoji)
	if err != nil {
		log.Fatalf("Contains emoji: %s", err)
	}

	fmt.Printf("String: %s, contains emoji: %t\n", stringWithoutEmoji, res) // false

	stringWithEmoji := "Hello world 🤗"
	res, err = s.ContainsEmoji(context.Background(), stringWithEmoji)
	if err != nil {
		log.Fatalf("Contains emoji: %s", err)
	}

	fmt.Printf("String: %s, contains emoji: %t\n", stringWithEmoji, res) // true
}

func allEmojis(s *gomoji.Service) {
	emojis, err := s.AllEmojis(context.Background())
	if err != nil {
		log.Fatalf("All emojis: %s", err)
	}

	fmt.Printf("Emojis count: %d\n", len(emojis))
}
