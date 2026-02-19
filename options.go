package botdetector

import (
	"errors"
	lru "github.com/hashicorp/golang-lru/v2"
)

// Option is a functional option that configures a BotDetector instance.
// Options are applied in the order they are passed to New.
type Option func(*BotDetector) (*BotDetector, error)

// WithRules returns an Option that replaces the default rule set with the provided rules.
// The custom rules are processed according to the same prefix/suffix convention:
// "^pattern" for prefix match, "pattern$" for suffix match, "^pattern$" for exact match,
// and plain "pattern" for substring match.
func WithRules(r []string) Option {
	return func(b *BotDetector) (*BotDetector, error) {
		b.importRules(r)
		return b, nil
	}
}

// WithCache returns an Option that enables an LRU cache of the given size for IsBot results.
// Once a user-agent is evaluated its result is cached so subsequent calls with the same
// normalised user-agent string are resolved without re-scanning all rules.
// size must be greater than 0; otherwise New will return an error.
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
