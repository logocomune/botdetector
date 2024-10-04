package main

import (
	"fmt"
	"net/http"

	"github.com/logocomune/botdetector/v2"
)

var detector *botdetector.BotDetector

func init() {
	var err error
	detector, err = botdetector.New()
	if err != nil {
		panic(err)
	}
}

func main() {

	http.HandleFunc("/", userAgentHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func userAgentHandler(w http.ResponseWriter, r *http.Request) {
	ua := r.Header.Get("User-Agent")
	fmt.Fprintf(w, "Your user agent --> %s\n", ua)
	fmt.Fprintf(w, "It's a bot?         %t\n", detector.IsBot(ua))
}
