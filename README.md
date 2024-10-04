# BotDetector

[![Build Status](https://app.travis-ci.com/logocomune/botdetector.svg?branch=master)](https://app.travis-ci.com/logocomune/botdetector)
[![Go Report Card](https://goreportcard.com/badge/github.com/logocomune/botdetector)](https://goreportcard.com/report/github.com/logocomune/botdetector)
[![codecov](https://codecov.io/gh/logocomune/botdetector/branch/master/graph/badge.svg)](https://codecov.io/gh/logocomune/botdetector)

BotDetector is a Go library that detects bots, spiders, and crawlers from user agents.

## Installation

`go get -u github.com/logocomune/botdetector/v2`

## Usage

### Simple usage

```go
   userAgent := req.Header.Get("User-Agent")

detector, _ := botdetector.New()
isBot := detector.IsBot(userAgnet)

if isBot {
log.Println("Bot, Spider or Crawler detected")
}

```

### Adding Custom Rules

You can add custom detection rules with the `WithRules` method. For example:

```go
userAgent := req.Header.Get("User-Agent")

detector, _ := botdetector.New(WithRules([]string{"my rule", "^test"}))
isBot := detector.IsBot(userAgent)

if isBot {
log.Println("Bot, Spider or Crawler detected")
}

```

Custom Rule Patterns:

| pattern | description                                               |
|---------|-----------------------------------------------------------|
| "..."   | Checks if the string contains the specified pattern.      |
| "^..."  | Checks if the string starts with the specified pattern.   |
| "...$"  | Checks if the string ends with the specified pattern.     |
| "^...$" | Checks if the string strictly matches the entire pattern. |

In this example, the custom rules "my rule" and "^test" are added to the existing detection rules.

### Adding Cache
You can add a lru cache rules with the `WithCache` method. For example:

```go
userAgent := req.Header.Get("User-Agent")

detector, _ := botdetector.New(WithCache(1000))
isBot := detector.IsBot(userAgent)

if isBot {
log.Println("Bot, Spider or Crawler detected")
}

```



### Example

[Simple example](_example/main.go)

## Inspiration

BotSeeker is inspired by [CrawlerDetect](https://github.com/JayBizzle/Crawler-Detect), an awesome PHP project.
 
