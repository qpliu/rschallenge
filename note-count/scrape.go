package main

import (
	"strings"
)

func scrapeValue(page string, i int, element string, end byte) (bool, int, string) {
	if strings.HasPrefix(page[i:], element) {
		start := i + len(element)
		for j := start; j < len(page); j++ {
			if page[j] == end {
				return true, j + 1, page[start:j]
			}
		}
	}
	return false, i, ""
}

func ScrapeSong(page string) (string, string) {
	artist := ""
	title := ""
	for i := 0; i < len(page); i++ {
		if page[i] != '<' {
			continue
		}
		if ok, index, value := scrapeValue(page, i, "<span class=\"title\">", '<'); ok {
			i = index
			title = value
		} else if ok, index, value := scrapeValue(page, i, "<span class=\"artist\">", '<'); ok {
			i = index
			artist = value
			break
		}
	}
	return artist, title
}

func ScrapeScores(page string) []Score {
	var scores []Score
	scoreIndex := -1
	for i := 0; i < len(page); i++ {
		if page[i] != '<' {
			continue
		}
		if ok, index, name := scrapeValue(page, i, "<td class=\"user\">", '<'); ok {
			i = index
			scoreIndex++
			scores = append(scores, Score{})
			scores[scoreIndex].Name = name
			continue
		}
		if scoreIndex < 0 {
			continue
		}
		if ok, index, difficulty := scrapeValue(page, i, "<td class=\"difficulty\" data-difficulty=\"", '"'); ok {
			i = index
			scores[scoreIndex].Difficulty = difficulty
		}
		if ok, index, score := scrapeValue(page, i, "<td class=\"score\" data-column=\"Score\">", '<'); ok {
			i = index
			scores[scoreIndex].Score = score
			continue
		}
		if ok, index, accuracy := scrapeValue(page, i, "<td class=\"accuracy\" data-column=\"Accuracy\">", '<'); ok {
			i = index
			scores[scoreIndex].Accuracy = accuracy
			continue
		}
		if ok, index, noteStreak := scrapeValue(page, i, "<td class=\"note-streak\" data-column=\"NoteStreak\">", '<'); ok {
			i = index
			scores[scoreIndex].NoteStreak = noteStreak
			continue
		}
	}
	return scores
}
