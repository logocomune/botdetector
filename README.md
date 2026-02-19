# BotDetector

[![Go](https://github.com/logocomune/botdetector/actions/workflows/go.yml/badge.svg)](https://github.com/logocomune/botdetector/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/logocomune/botdetector/v2)](https://goreportcard.com/report/github.com/logocomune/botdetector/v2)
[![codecov](https://codecov.io/gh/logocomune/botdetector/branch/main/graph/badge.svg)](https://codecov.io/gh/logocomune/botdetector)

BotDetector is a Go library that detects bots, spiders, and crawlers by inspecting HTTP User-Agent strings.

It ships with **1446 built-in rules** covering all major search engine crawlers (Googlebot, Bingbot, Baiduspider, YandexBot, …), SEO tools, HTTP libraries, and other automated agents. Rules are matched case-insensitively using four strategies: exact match, prefix, suffix, and substring.

## Installation

```bash
go get -u github.com/logocomune/botdetector/v2
```

Requires **Go 1.25.6** or later.

## Quick start

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/logocomune/botdetector/v2"
)

func main() {
    detector, err := botdetector.New()
    if err != nil {
        log.Fatal(err)
    }

    ua := "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
    fmt.Println(detector.IsBot(ua)) // true

    ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
    fmt.Println(detector.IsBot(ua)) // false
}
```

## Usage in an HTTP handler

`BotDetector` is safe for concurrent use — create a single instance at startup and share it across goroutines.

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/logocomune/botdetector/v2"
)

var detector *botdetector.BotDetector

func init() {
    var err error
    detector, err = botdetector.New()
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        ua := r.Header.Get("User-Agent")
        fmt.Fprintf(w, "User-Agent: %s\n", ua)
        fmt.Fprintf(w, "Is bot:     %t\n", detector.IsBot(ua))
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Options

`New` accepts functional options. Multiple options can be combined.

### `WithCache` — LRU result cache

Enables an LRU cache of the given size. Subsequent calls with the same User-Agent string return the cached result instantly, without re-scanning the rules.

```go
// Cache up to 1 000 unique User-Agent results.
detector, err := botdetector.New(botdetector.WithCache(1000))
```

`WithCache(0)` or a negative size returns an error.

### `WithRules` — custom rule set

Replaces the built-in rules with your own list. Use this when you only want to match a specific set of agents.

> **Note:** `WithRules` *replaces* all built-in rules. If you need to keep the defaults, use `NewWithRules` only when building a fully custom detector; otherwise extend via the rule pattern format below.

```go
customRules := []string{
    "^mybot",          // prefix match
    "my-scanner$",     // suffix match
    "^internalbot$",   // exact match
    "datacollector",   // substring match
}

detector, err := botdetector.New(botdetector.WithRules(customRules))
```

### Combining options

Options are applied left to right. The most common combination is custom rules together with a cache:

```go
detector, err := botdetector.New(
    botdetector.WithRules(customRules),
    botdetector.WithCache(1000),
)
if err != nil {
    log.Fatal(err)
}
```

## Rule syntax

All patterns are matched **case-insensitively**.

| Pattern   | Match type | Example pattern | Matches                          | Does not match        |
|-----------|-----------|-----------------|----------------------------------|-----------------------|
| `"..."`   | Substring | `"googlebot"`   | `"Mozilla/.../Googlebot/2.1"`    | —                     |
| `"^..."`  | Prefix    | `"^java"`       | `"Java/11.0.2"`                  | `"not java/11"`       |
| `"...$"`  | Suffix    | `"crawler$"`    | `"my-fast-crawler"`              | `"crawler-extra"`     |
| `"^...$"` | Exact     | `"^b0t$"`       | `"b0t"`                          | `"b0t/1.0"`, `"xb0t"` |

## Standalone detector without built-in rules

`NewWithRules` creates a detector that uses *only* the provided rules (no built-in list):

```go
detector := botdetector.NewWithRules([]string{
    "^mybot$",
    "datacrawler",
})

detector.IsBot("mybot")           // true
detector.IsBot("super-crawler")   // false — "datacrawler" != "super-crawler"
detector.IsBot("datacrawler/1.0") // true
```

## Full example

See [_example/main.go](_example/main.go) for a minimal HTTP server that reports bot status for every request.

## Inspiration

BotDetector is inspired by [CrawlerDetect](https://github.com/JayBizzle/Crawler-Detect), an excellent PHP project.
