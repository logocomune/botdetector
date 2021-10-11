# BotDetector
[![Build Status](https://app.travis-ci.com/logocomune/botdetector.svg?branch=master)](https://app.travis-ci.com/logocomune/botdetector)
[![Go Report Card](https://goreportcard.com/badge/github.com/logocomune/botdetector)](https://goreportcard.com/report/github.com/logocomune/botdetector)
[![codecov](https://codecov.io/gh/logocomune/botdetector/branch/master/graph/badge.svg)](https://codecov.io/gh/logocomune/botdetector)

BotDetector is a golang library that detects Bot/Spider/Crawler from user agent

## Installation

`go get -u github.com/logocomune/botdetector`

## Usage

```go
   userAgent :=  req.Header.Get("User-Agent")
   
   detector := botdetector.New()
   isBot := detector.IsBot(userAgnet)
   
   if isBot {
   	 log.Println("Bot, Spider or Crawler detected")
   }
    
```
 
### Example

[Simple example](_example/main.go)

## Inspiration

BotSeeker is inspired by [CrawlerDetect](https://github.com/JayBizzle/Crawler-Detect), an awesome PHP project.
 
