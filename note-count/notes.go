package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		for _, challenge := range os.Args[1:] {
			showChallenge(challenge)
		}
	} else {
		challenges := getChallenges()
		if len(challenges) == 4 {
			showChallenge(challenges[0])
			showChallenge(challenges[2])
		} else {
			for _, challenge := range challenges {
				showChallenge(challenge)
			}
		}
	}
}

func showChallenge(challenge string) {
	resp, err := http.Get("http://rocksmithchallenge.com/challenges/" + challenge)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: http error: %v", challenge, err)
		return
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		fmt.Fprintf(os.Stderr, "%s: http read error: %v", challenge, err)
		return
	}
	body := buf.String()
	artist, title := ScrapeSong(body)
	scores := ScrapeScores(body)
	ComputeScores(scores)
	fmt.Printf("%s: %s - %s\n", challenge, artist, title)
	PrintScores(os.Stdout, scores)
}

func getChallenges() []string {
	resp, err := http.Get("http://rocksmithchallenge.com/")
	if err != nil {
		fmt.Fprintf(os.Stderr, "http error: %v", err)
		return nil
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		fmt.Fprintf(os.Stderr, "http read error: %v", err)
		return nil
	}
	return ScrapeChallenges(buf.String())
}
