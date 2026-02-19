package botdetector

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"strings"
)

// BotDetector holds the compiled detection rules and an optional LRU result cache.
//
// Rules are stored in four typed collections for cache-friendly sequential scans
// and O(1) exact-match lookups:
//   - strictRules    – map for O(1) exact match
//   - startWithRules – slice for prefix matching
//   - endWithRules   – slice for suffix matching
//   - containsRules  – slice for substring matching
type BotDetector struct {
	strictRules    map[string]struct{}
	startWithRules []string
	endWithRules   []string
	containsRules  []string
	cache          *lru.Cache[string, bool]
}

// New creates a new instance of BotDetector using predefined rules.
func New(opt ...Option) (*BotDetector, error) {
	b := &BotDetector{}
	b.importRules(rules)

	var err error
	for i := range opt {
		b, err = opt[i](b)
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

// importRules replaces the current rule set with the provided list of rule strings.
// If r is nil or empty the rule sets are reset to an empty state.
func (b *BotDetector) importRules(r []string) {
	b.strictRules = make(map[string]struct{})
	b.startWithRules = nil
	b.endWithRules = nil
	b.containsRules = nil

	if len(r) == 0 {
		return
	}

	// Temporary sets used to deduplicate slice-based rule categories.
	seenStart := make(map[string]struct{}, len(r))
	seenEnd := make(map[string]struct{}, len(r))
	seenContains := make(map[string]struct{}, len(r))

	for _, s := range r {
		b.addExpression(s, seenStart, seenEnd, seenContains)
	}
}

// NewWithRules initializes a new instance of BotDetector with provided rules.
func NewWithRules(rules []string) *BotDetector {
	uBot := &BotDetector{}
	uBot.importRules(rules)
	return uBot
}

// addExpression parses a single rule string and registers it in the appropriate rule set.
// Rules prefixed with "^" are treated as prefix matches; rules suffixed with "$" are
// treated as suffix matches; rules with both are exact matches; all others are substring matches.
// Duplicate patterns within a category are deduplicated via the provided seen maps.
func (b *BotDetector) addExpression(original string, seenStart, seenEnd, seenContains map[string]struct{}) {
	lowered := strings.ToLower(original)
	isStrict := strings.HasPrefix(lowered, "^") && strings.HasSuffix(lowered, "$")
	isStartWith := strings.HasPrefix(lowered, "^")
	isEndWith := strings.HasSuffix(lowered, "$")

	switch {
	case isStrict:
		b.strictRules[lowered[1:len(lowered)-1]] = struct{}{}
	case isStartWith:
		detector := lowered[1:]
		if _, ok := seenStart[detector]; !ok {
			seenStart[detector] = struct{}{}
			b.startWithRules = append(b.startWithRules, detector)
		}
	case isEndWith:
		detector := lowered[:len(lowered)-1]
		if _, ok := seenEnd[detector]; !ok {
			seenEnd[detector] = struct{}{}
			b.endWithRules = append(b.endWithRules, detector)
		}
	default:
		if _, ok := seenContains[lowered]; !ok {
			seenContains[lowered] = struct{}{}
			b.containsRules = append(b.containsRules, lowered)
		}
	}
}

// IsBot tests whether the useragent is a bot, crawler or a spider.
//
// Performance notes:
//   - O3: the LRU cache is checked against the raw UA before normalize() is called,
//     so cache hits pay zero normalization cost.
//   - O1: matching exits immediately on the first rule that fires.
//   - O2: exact rules use a map for O(1) lookup; prefix/suffix/contains rules are
//     stored as plain slices for CPU cache-friendly sequential scans.
func (b *BotDetector) IsBot(ua string) bool {
	// O3: check cache with raw UA before paying normalize() cost.
	if b.cache != nil {
		if ret, ok := b.cache.Get(ua); ok {
			return ret
		}
	}

	uaNormalized := normalize(ua)

	// O1+O2: strict (exact) match — O(1) map lookup.
	if _, ok := b.strictRules[uaNormalized]; ok {
		if b.cache != nil {
			b.cache.Add(ua, true)
		}
		return true
	}

	// O1+O2: prefix match — sequential slice scan with early return.
	for _, pattern := range b.startWithRules {
		if strings.HasPrefix(uaNormalized, pattern) {
			if b.cache != nil {
				b.cache.Add(ua, true)
			}
			return true
		}
	}

	// O1+O2: suffix match — sequential slice scan with early return.
	for _, pattern := range b.endWithRules {
		if strings.HasSuffix(uaNormalized, pattern) {
			if b.cache != nil {
				b.cache.Add(ua, true)
			}
			return true
		}
	}

	// O1+O2: substring match — sequential slice scan with early return.
	for _, pattern := range b.containsRules {
		if strings.Contains(uaNormalized, pattern) {
			if b.cache != nil {
				b.cache.Add(ua, true)
			}
			return true
		}
	}

	if b.cache != nil {
		b.cache.Add(ua, false)
	}
	return false
}

// normalize converts a user-agent string to lowercase and strips substrings that
// would otherwise cause false-positive or false-negative bot detections.
// For example, "cubot" is a brand name that collides with common bot patterns, so
// it is removed before matching.
func normalize(userAgent string) string {
	uaNormalized := strings.ToLower(userAgent)

	if strings.HasPrefix(uaNormalized, "lynx/") {
		uaNormalized = strings.ReplaceAll(uaNormalized, "libwww-fm", "")
		return uaNormalized
	}

	toReplace := []string{
		"cubot",
		"; m bot",
		"; crono",
		"; b bot",
		"; idbot",
		"; id bot",
		"; power bot",
		"yandexsearch/",
		"amigavoyager",
	}
	for _, rep := range toReplace {
		uaNormalized = strings.ReplaceAll(uaNormalized, rep, "")
	}

	return uaNormalized
}
