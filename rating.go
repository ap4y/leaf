package leaf

import "math"

const ratingSuccess = 0.6

// RatingType defines types of review rating options.
type RatingType string

const (
	// RatingTypeAuto defines auto rated review option.
	RatingTypeAuto RatingType = "auto"
	// RatingTypeSelf defines self rated review option.
	RatingTypeSelf RatingType = "self"
)

// ReviewScore defines grade for review attempts. Rater uses scores to
// calculate rating in range from [0, 1].
type ReviewScore int

const (
	// ReviewScoreAgain defines "again" score.
	ReviewScoreAgain ReviewScore = iota
	// ReviewScoreHard defines "hard" score.
	ReviewScoreHard
	// ReviewScoreGood defines "good" score.
	ReviewScoreGood
	// ReviewScoreEasy defines "easy" score.
	ReviewScoreEasy
)

// Rater rates review attempt based on amount of mistakes. Rating
// should be within [0, 1] range.
type Rater interface {
	Rate(question string, score ReviewScore) float64
}

type harshRater struct {
	mistakes map[string]int
}

// HarshRater returns miss count based Rater. Miss counter will
// increase for each "again" score. Rating declines really fast and
// even 1 mistake results in 0.59 rating.
func HarshRater() Rater {
	return &harshRater{make(map[string]int)}
}

func (rater harshRater) Rate(question string, score ReviewScore) float64 {
	if score == ReviewScoreAgain {
		rater.mistakes[question]++
		return 0
	}

	mistakes := rater.mistakes[question]
	if mistakes == 0 {
		return 1
	}

	return math.Max(0, 0.79-float64(mistakes)/5)
}

type tableRater struct {
}

// TableRater returns Rater implementation with following the conversion table:
// again => 0
// hard => 0.2
// good => 0.6
// easy => 1.0
func TableRater() Rater {
	return &tableRater{}
}

func (rater tableRater) Rate(question string, score ReviewScore) float64 {
	switch score {
	case ReviewScoreHard:
		return 0.2
	case ReviewScoreGood:
		return 0.6
	case ReviewScoreEasy:
		return 1.0
	default:
		return 0
	}
}
