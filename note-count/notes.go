package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		for _, challengeId := range os.Args[1:] {
			showChallenge(challengeId)
		}
	} else {
		challenges := getChallenges()
		for _, challenge := range challenges {
			PrintChallenge(os.Stdout, challenge)
		}
	}
}

func showChallenge(challengeId string) {
	resp, err := http.Get("http://rocksmithchallenge.com/challenges/" + challengeId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: http error: %v", challengeId, err)
		return
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		fmt.Fprintf(os.Stderr, "%s: http read error: %v", challengeId, err)
		return
	}
	body := buf.String()
	challenge := ScrapeChallenge(body)
	ComputeScores(challenge.Scores)
	PrintChallenge(os.Stdout, challenge)
}

func getChallenges() []Challenge {
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
