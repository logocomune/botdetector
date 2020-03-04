package botdetector

import (
	"strings"
)

const (
	strict = iota
	startWidth
	contains
	)

type expressionInfo struct {
	expressionType int
	source         string
	detector       interface{}
}

type BotDetector struct {
	expression map[string]expressionInfo
}

func New() *BotDetector {
	return newDetector(rules)
}

func newDetector(rules []string) *BotDetector {
	uBot := BotDetector{expression: make(map[string]expressionInfo)}

	for _, s := range rules {
		uBot.addExpression(s)
	}

	return &uBot
}

func (b *BotDetector) addExpression(original string) {
	e := expressionInfo{
		source: original,
	}

	s := strings.ToLower(original)
	if strings.HasPrefix(s, "^") && strings.HasSuffix(s, "$") {
		b.expression[original] = expressionInfo{
			source:         original,
			expressionType: strict,
			detector:       s[1 : len(s)-1],
		}

		return
	}

	if strings.HasPrefix(s, "^") {
		b.expression[original] = expressionInfo{
			source:         original,
			expressionType: startWidth,
			detector:       s[1:],
		}

		return
	}

	e.expressionType = contains
	b.expression[original] = expressionInfo{
		source:         original,
		expressionType: contains,
		detector:       s,
	}
}

// IsBot tests whether the useragent is a bot, crawler or a spider.
func (b *BotDetector) IsBot(ua string) bool {
	uaNormalized := normalize(ua)

	for _, exp := range b.expression {
		switch exp.expressionType {
		case strict:
			if uaNormalized == exp.detector {
				//fmt.Printf("%+v\n",exp)
				return true
			}
		case startWidth:
			if strings.HasPrefix(uaNormalized, exp.detector.(string)) {
				//fmt.Printf("%+v\n",exp)
				return true
			}
		case contains:
			if strings.Contains(uaNormalized, exp.detector.(string)) {
				//fmt.Printf("%+v\n",exp)
				return true
			}
		}
	}

	return false
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
