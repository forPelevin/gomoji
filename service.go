package gomoji

import (
	"context"
	"fmt"
	"sync"
)

type (
	// Provider is an emoji provider.
	Provider interface {
		AllEmojis(ctx context.Context) ([]Emoji, error)
	}

	// Service is a service to work with emojis using a provider. The service methods are concurrent safe.
	Service struct {
		p        Provider
		emojiMap map[rune]Emoji
		mu       sync.RWMutex
	}
)

// NewService creates a new instance of service.
func NewService(p Provider) *Service {
	return &Service{
		p: p,
	}
}

func (s *Service) initEmojiMap(ctx context.Context) error {
	if s.emojiMap != nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	emojis, err := s.p.AllEmojis(ctx)
	if err != nil {
		return fmt.Errorf("get all emojis: %w", err)
	}

	s.emojiMap = make(map[rune]Emoji, len(emojis))
	for _, e := range emojis {
		s.emojiMap[[]rune(e.Character)[0]] = e
	}

	return nil
}

// ContainsEmoji checks whether given string contains emoji or not. It takes emoji data from provider.
func (s *Service) ContainsEmoji(ctx context.Context, str string) (bool, error) {
	err := s.initEmojiMap(ctx)
	if err != nil {
		return false, fmt.Errorf("fill emoji map: %w", err)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, r := range str {
		if _, ok := s.emojiMap[r]; ok {
			return true, nil
		}
	}

	return false, nil
}

// AllEmojis gets all emojis from provider.
func (s *Service) AllEmojis(ctx context.Context) ([]Emoji, error) {
	err := s.initEmojiMap(ctx)
	if err != nil {
		return nil, fmt.Errorf("fill emoji map: %w", err)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return emojiMapToSlice(s.emojiMap), nil
}
