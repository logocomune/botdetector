package botdetector

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"strings"
)

const (
	strict = iota
	startWith
	endWith
	contains
)

type expressionInfo struct {
	expressionType int
	source         string
	detector       string
}

type BotDetector struct {
	expressions map[string]expressionInfo
	cache       *lru.Cache[string, bool]
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

func (b *BotDetector) importRules(r []string) {
	if r == nil || len(r) == 0 {
		b.expressions = make(map[string]expressionInfo)
		return
	}
	b.expressions = make(map[string]expressionInfo, len(r))
	for _, s := range r {
		b.addExpression(s)
	}
}

// NewWithRules initializes a new instance of BotDetector with provided rules.
func NewWithRules(rules []string) *BotDetector {
	uBot := BotDetector{expressions: make(map[string]expressionInfo)}

	for _, s := range rules {
		uBot.addExpression(s)
	}

	return &uBot
}

func (b *BotDetector) addExpression(original string) {
	lowered := strings.ToLower(original)
	isStrict := strings.HasPrefix(lowered, "^") && strings.HasSuffix(lowered, "$")
	isStartWith := strings.HasPrefix(lowered, "^")
	isEndWith := strings.HasSuffix(lowered, "$")

	switch {
	case isStrict:
		b.addExpressionInfo(original, strict, lowered[1:len(lowered)-1])
	case isStartWith:
		b.addExpressionInfo(original, startWith, lowered[1:])
	case isEndWith:
		b.addExpressionInfo(original, endWith, lowered[:len(lowered)-1])
	default:
		b.addExpressionInfo(original, contains, lowered)
	}
}

func (b *BotDetector) addExpressionInfo(source string, exprType int, detector string) {
	b.expressions[source] = expressionInfo{
		source:         source,
		expressionType: exprType,
		detector:       detector,
	}
}

// IsBot tests whether the useragent is a bot, crawler or a spider.
func (b *BotDetector) IsBot(ua string) bool {
	uaNormalized := normalize(ua)
	if b.cache != nil {
		if ret, ok := b.cache.Get(uaNormalized); ok {
			return ret
		}
	}
	ret := false
	for _, exp := range b.expressions {
		switch exp.expressionType {
		case strict:
			if uaNormalized == exp.detector {
				ret = true
			}
		case startWith:
			if strings.HasPrefix(uaNormalized, exp.detector) {
				ret = true
			}
		case endWith:
			if strings.HasSuffix(uaNormalized, exp.detector) {

				ret = true
			}
		case contains:
			if strings.Contains(uaNormalized, exp.detector) {

				ret = true
			}
		}
	}
	if b.cache != nil {
		b.cache.Add(uaNormalized, ret)
	}

	return ret
}

func normalize(userAgent string) string {
	uaNormalized := strings.ToLower(userAgent)

	if strings.HasPrefix(uaNormalized, "lynx/") {
		uaNormalized = strings.Replace(uaNormalized, "libwww-fm", "", 1)
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
		uaNormalized = strings.Replace(uaNormalized, rep, "", -1)
	}

	return uaNormalized
}
