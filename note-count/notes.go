package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

func main() {
	for _, arg := range os.Args[1:] {
		resp, err := http.Get("http://rocksmithchallenge.com/challenges/" + arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: http error: %v", arg, err)
			continue
		}
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(resp.Body); err != nil {
			fmt.Fprintf(os.Stderr, "%s: http read error: %v", arg, err)
			continue
		}
		body := buf.String()
		artist, title := ScrapeSong(body)
		scores := ScrapeScores(body)
		ComputeScores(scores)
		fmt.Printf("%s: %s - %s\n", arg, artist, title)
		PrintScores(os.Stdout, scores)
	}
}
