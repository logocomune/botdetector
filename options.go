package botdetector

import (
	"errors"
	lru "github.com/hashicorp/golang-lru/v2"
)

type Option func(*BotDetector) (*BotDetector, error)

func WithRules(r []string) Option {
	return func(b *BotDetector) (*BotDetector, error) {
		b.importRules(r)
		return b, nil
	}
}

func WithCache(size int) Option {
	return func(b *BotDetector) (*BotDetector, error) {
		if size <= 0 {
			return b, errors.New("cache size must be greater than 0")
		}
		var err error
		b.cache, err = lru.New[string, bool](size)
		return b, err
	}
}
