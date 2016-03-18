package main

import (
	"strconv"
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

func scrapeLine(page string, i int) (int, string) {
	if i >= len(page) {
		return i, ""
	}
	for j := i; j < len(page); j++ {
		if page[j] == '\n' {
			return j + 1, page[i:j]
		}
	}
	return len(page), page[i:]
}

func ScrapeChallenge(page string) Challenge {
	var challenge Challenge
	scoreIndex := -1
	for i := 0; i < len(page); i++ {
		if page[i] != '<' {
			continue
		}
		if ok, index, name := scrapeValue(page, i, "<td class=\"user\">", '<'); ok {
			i = index
			scoreIndex++
			challenge.Scores = append(challenge.Scores, Score{})
			challenge.Scores[scoreIndex].Name = name
			continue
		}
		if scoreIndex < 0 {
			if ok, index, value := scrapeValue(page, i, "<span class=\"title\">", '<'); ok {
				i = index
				challenge.Title = value
			} else if ok, index, value := scrapeValue(page, i, "<span class=\"artist\">", '<'); ok {
				i = index
				challenge.Artist = value
			} else if ok, index, value := scrapeValue(page, i, "<h1 class=\"challenge-title ", '"'); ok {
				i = index
				challenge.Arrangement = value
			}
		} else if ok, index, difficulty := scrapeValue(page, i, "<td class=\"difficulty\" data-difficulty=\"", '"'); ok {
			i = index
			challenge.Scores[scoreIndex].Difficulty = difficulty
		} else if ok, index, score := scrapeValue(page, i, "<td class=\"score\" data-column=\"Score\">", '<'); ok {
			i = index
			challenge.Scores[scoreIndex].Score = score
		} else if ok, index, accuracy := scrapeValue(page, i, "<td class=\"accuracy\" data-column=\"Accuracy\">", '<'); ok {
			i = index
			challenge.Scores[scoreIndex].Accuracy = accuracy
		} else if ok, index, noteStreak := scrapeValue(page, i, "<td class=\"note-streak\" data-column=\"NoteStreak\">", '<'); ok {
			i = index
			challenge.Scores[scoreIndex].NoteStreak = noteStreak
		}
	}
	return challenge
}

func ScrapeChallenges(page string) []Challenge {
	var result []Challenge
	challengeIndex := -1
	inResults := false
	division := ""
	for i := 0; i < len(page); i++ {
		if page[i] != '<' {
			continue
		}
		if ok, index, value := scrapeValue(page, i, "<a href=\"/challenges/", '"'); ok {
			i = index
			if strings.HasPrefix(page[i:], ">") {
				challengeIndex++
				result = append(result, Challenge{})
				result[challengeIndex].ChallengeId = value
			}
		}
		if challengeIndex < 0 {
			continue
		}
		if ok, index, value := scrapeValue(page, i, "<img class=\"instrument\" src=\"/bundles/rsui/img/", '.'); ok {
			i = index
			result[challengeIndex].Arrangement = value
		} else if ok, index, value := scrapeValue(page, i, "<span class=\"title\"", '>'); ok {
			i = index
			i, value = scrapeLine(page, i)
			i, value = scrapeLine(page, i)
			value = strings.TrimSpace(value)
			if strings.HasPrefix(value, "<") {
				i, value = scrapeLine(page, i)
				value = strings.TrimSpace(value)
			}
			result[challengeIndex].Title = value
		} else if ok, index, value := scrapeValue(page, i, "<span class=\"artist\"", '>'); ok {
			i = index
			i, value = scrapeLine(page, i)
			i, value = scrapeLine(page, i)
			value = strings.TrimSpace(value)
			if strings.HasPrefix(value, "<") {
				i, value = scrapeLine(page, i)
				value = strings.TrimSpace(value)
			}
			result[challengeIndex].Artist = value
		}
		if !inResults {
			if ok, index, _ := scrapeValue(page, i, "<div class=\"result-content\"", '>'); ok {
				i = index
				inResults = true
			}
		} else {
			if ok, index, _ := scrapeValue(page, i, "<div class=\"result-link\"", '>'); ok {
				i = index
				inResults = false
			} else if ok, index, value := scrapeValue(page, i, "<h4>Division ", '<'); ok {
				i = index
				division = value
			} else if ok, index, value := scrapeValue(page, i, "<td>", '<'); ok {
				i = index
				result[challengeIndex].Scores = append(result[challengeIndex].Scores, Score{})
				result[challengeIndex].Scores[len(result[challengeIndex].Scores)-1].Name = value
				result[challengeIndex].Scores[len(result[challengeIndex].Scores)-1].Difficulty = division
			} else if ok, index, value := scrapeValue(page, i, "<td class=\"score\">", '<'); ok {
				i = index
				result[challengeIndex].Scores[len(result[challengeIndex].Scores)-1].Score = value
			} else if ok, index, value := scrapeValue(page, i, "<td class=\"accuracy\">", '<'); ok {
				i = index
				result[challengeIndex].Scores[len(result[challengeIndex].Scores)-1].Accuracy = value
			} else if ok, index, value := scrapeValue(page, i, "<td class=\"note-streak\">", '<'); ok {
				i = index
				result[challengeIndex].Scores[len(result[challengeIndex].Scores)-1].NoteStreak = value
			}
		}
	}

	for i := range result {
		lastDifficulty := ""
		lastPct100 := 10000
		difficulty := 1
		for j := range result[i].Scores {
			if lastDifficulty != result[i].Scores[j].Difficulty {
				lastPct100 = 10000
				lastDifficulty = result[i].Scores[j].Difficulty
			}
			pct100 := ComputePct100(result[i].Scores[j].Accuracy)
			if pct100 > lastPct100 {
				difficulty++
			}
			lastPct100 = pct100
			result[i].Scores[j].Difficulty += "-" + strconv.Itoa(difficulty)
		}
		ComputeScores(result[i].Scores)
	}
	return result
}
