package main

import (
	"bytes"
	"testing"
)

const (
	EXPECTED_PRINT = `The Black Keys - Little Black Submarines - guitar
                User Difficulty       Score     %   NS Notes
        elemenohpenc       hard     221,835 91.38   51 743/813 (-70)
     RollingStone222       hard     376,390 93.84  105 763/813 (-50)
             Ezzy911       hard     261,882 91.75   62 746/813 (-67)
               qpliu       hard     181,361 91.51   76 744/813 (-69)
          the_xandos       hard     202,436 87.57  132 712/813 (-101)
               jaq-b     medium     603,897 99.47  150 382/384 (-2)
         Sl1mehunter     medium     547,410 99.21  206 381/384 (-3)
          Z1ronJones     medium     380,836 97.91  153 376/384 (-8)
             Vecco34     medium     259,707 94.53   94 363/384 (-21)
          advalencia       easy      82,580 94.92   97 131/138 (-7)
`
)

func TestScore(t *testing.T) {
	check := func(score Score, expectedName, expectedDifficulty string, expectedPct, expectedStreak, expectedNotes, expectedHits int) {
		if score.Name != expectedName {
			t.Errorf("Expected name %s, got %s", expectedName, score.Name)
		}
		if score.Difficulty != expectedDifficulty {
			t.Errorf("Expected difficulty %s, got %s", expectedDifficulty, score.Difficulty)
		}
		if score.Pct100 != expectedPct {
			t.Errorf("Expected percentage %d, got %d", expectedPct, score.Pct100)
		}
		if score.Streak != expectedStreak {
			t.Errorf("Expected streak %d, got %d", expectedStreak, score.Streak)
		}
		if score.Notes != expectedNotes {
			t.Errorf("Expected notes %d, got %d", expectedNotes, score.Notes)
		}
		if score.Hits != expectedHits {
			t.Errorf("Expected hits %d, got %d", expectedHits, score.Hits)
		}
	}
	challenge := ScrapeChallenge(CHALLENGE_PAGE)
	if len(challenge.Scores) != 10 {
		t.Errorf("Expected 10 scores, got %d", len(challenge.Scores))
	}
	ComputeScores(challenge.Scores)
	check(challenge.Scores[0], "elemenohpenc", "hard", 9138, 51, 813, 743)
	check(challenge.Scores[1], "RollingStone222", "hard", 9384, 105, 813, 763)
	check(challenge.Scores[2], "Ezzy911", "hard", 9175, 62, 813, 746)
	check(challenge.Scores[3], "qpliu", "hard", 9151, 76, 813, 744)
	check(challenge.Scores[4], "the_xandos", "hard", 8757, 132, 813, 712)
	check(challenge.Scores[5], "jaq-b", "medium", 9947, 150, 384, 382)
	check(challenge.Scores[6], "Sl1mehunter", "medium", 9921, 206, 384, 381)
	check(challenge.Scores[7], "Z1ronJones", "medium", 9791, 153, 384, 376)
	check(challenge.Scores[8], "Vecco34", "medium", 9453, 94, 384, 363)
	check(challenge.Scores[9], "advalencia", "easy", 9492, 97, 138, 131)
	var buf bytes.Buffer
	PrintChallenge(&buf, challenge)
	if buf.String() != EXPECTED_PRINT {
		t.Errorf("Expected:\n%sGot:\n%s", EXPECTED_PRINT, buf.String())
	}
}
