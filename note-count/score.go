package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Challenge struct {
	Artist, Title, ChallengeId, Arrangement string
	Scores                                  []Score
}

type Score struct {
	Name, Difficulty, Score, Accuracy, NoteStreak string
	Pct100, Streak, Notes, Hits                   int
}

func ComputePct100(accuracy string) int {
	pct := strings.Index(accuracy, "%")
	if pct > 0 {
		accuracy = accuracy[0:pct]
	}
	dot := strings.Index(accuracy, ".")
	if dot < 0 {
		if n, err := strconv.Atoi(accuracy); err == nil {
			return n * 100
		}
		return 0
	}
	pct100 := 0
	if n, err := strconv.Atoi(accuracy[0:dot]); err == nil {
		pct100 = n * 100
	}
	frac := accuracy[dot+1:]
	if len(frac) == 1 {
		if n, err := strconv.Atoi(frac); err == nil {
			pct100 += n * 10
		}
	} else {
		if len(frac) > 2 {
			frac = frac[0:2]
		}
		if frac[0] == '0' {
			frac = frac[1:]
		}
		if n, err := strconv.Atoi(frac); err == nil {
			pct100 += n
		}
	}
	return pct100
}

func ComputeScores(scores []Score) {
	for i := range scores {
		scores[i].Pct100 = ComputePct100(scores[i].Accuracy)
		if n, err := strconv.Atoi(scores[i].NoteStreak); err == nil {
			scores[i].Streak = n
		}
	}
	pcts := make(map[string][]int)
	notes := make(map[string]int)
	for _, score := range scores {
		diff := score.Difficulty
		if score.Streak > notes[diff] {
			notes[diff] = score.Streak
		}
		pcts[diff] = append(pcts[diff], score.Pct100)
	}
	for diff, pct := range pcts {
		if n, err := ComputeTotal(notes[diff], 5000, pct); err == nil {
			notes[diff] = n
		}
	}
	for i := range scores {
		scores[i].Notes = notes[scores[i].Difficulty]
		scores[i].Hits = (scores[i].Notes*scores[i].Pct100 + 5000) / 10000
	}
}

func PrintChallenge(w io.Writer, challenge Challenge) {
	if challenge.ChallengeId != "" {
		fmt.Fprintf(w, "%s: ", challenge.ChallengeId)
	}
	fmt.Fprintf(w, "%s - %s - %s\n", challenge.Artist, challenge.Title, challenge.Arrangement)
	fmt.Fprintf(w, "%20.20s %10.10s %11.11s %5.5s %4.4s Notes\n", "User", "Difficulty", "Score", "%", "NS")
	for _, score := range challenge.Scores {
		fmt.Fprintf(w, "%20.20s %10.10s %11.11s %5.5s %4.4s %d/%d (%d)\n", score.Name, score.Difficulty, score.Score, score.Accuracy, score.NoteStreak, score.Hits, score.Notes, score.Hits-score.Notes)
	}
}
