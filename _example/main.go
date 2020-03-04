package main

import (
	"fmt"
	"net/http"

	"github.com/logocomune/botdetector"
)

var detector = botdetector.New()

func main() {

	http.HandleFunc("/", userAgentHandler)

	http.ListenAndServe(":8080", nil)
}

func userAgentHandler(w http.ResponseWriter, r *http.Request) {
	ua := r.Header.Get("User-Agent")
	fmt.Fprintf(w, "Your user agent --> %s\n", ua)
	fmt.Fprintf(w, "It's a bot?         %t\n", detector.IsBot(ua))
}
