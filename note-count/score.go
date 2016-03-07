package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Score struct {
	Name, Difficulty, Score, Accuracy, NoteStreak string
	Pct100, Streak, Notes, Hits                   int
}

func ComputeScores(scores []Score) {
	for i := range scores {
		dot := strings.Index(scores[i].Accuracy, ".")
		if dot < 0 {
			if n, err := strconv.Atoi(scores[i].Accuracy); err == nil {
				scores[i].Pct100 = 100 * n
			}
		} else {
			if n, err := strconv.Atoi(scores[i].Accuracy[0:dot]); err == nil {
				scores[i].Pct100 = 100 * n
			}
			frac := scores[i].Accuracy[dot+1:]
			if len(frac) == 1 {
				if n, err := strconv.Atoi(frac); err == nil {
					scores[i].Pct100 += 10 * n
				}
			} else {
				if frac[0] == '0' {
					frac = frac[1:]
				}
				if n, err := strconv.Atoi(frac); err == nil {
					scores[i].Pct100 += n
				}
			}
		}
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

func PrintScores(w io.Writer, scores []Score) {
	fmt.Fprintf(w, "%20.20s %10.10s %11.11s %5.5s %4.4s Notes\n", "User", "Difficulty", "Score", "%", "NS")
	for _, score := range scores {
		fmt.Fprintf(w, "%20.20s %10.10s %11.11s %5.5s %4.4s %d/%d (%d)\n", score.Name, score.Difficulty, score.Score, score.Accuracy, score.NoteStreak, score.Hits, score.Notes, score.Hits-score.Notes)
	}
}
