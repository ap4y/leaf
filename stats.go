package leaf

import (
	"encoding/json"
	"math"
	"time"
)

const ratingSuccess = 0.6

// Supermemo calculates review intervals
type Supermemo interface {
	json.Marshaler
	json.Unmarshaler

	// Advance advances supermemo state for a card.
	Advance(rating float64) (interval float64)
	// NextReviewAt returns next review timestamp for a card.
	NextReviewAt() time.Time
	// SortParam returns values that should used as a review order for cards
	SortParam() float64
}

// Rater rates review attempt based on amount of mistakes. Rating
// should be within [0, 1] range.
type Rater interface {
	Rate(mistakes int) float64
}

// HarshRater implements offset based rating. Rating declines really
// fast and even 1 mistakes results in 0.59 rating.
type HarshRater struct{}

// Rate implements Rater for a HarshRater.
func (hr HarshRater) Rate(mistakes int) float64 {
	if mistakes == 0 {
		return 1
	}

	return math.Max(0, 0.79-float64(mistakes)/5)
}

// Stats store SM2+ parameters for a Card.
type Stats struct {
	Supermemo
	Rater
}

// CardWithStats joins Stats to a Card
type CardWithStats struct {
	Card
	*Stats
}

// DefaultStats returns a new Stats initialized with default values.
func DefaultStats() *Stats {
	return &Stats{NewSupermemo2PlusCustom(), &HarshRater{}}
}

// IsReady signals whether card is read for review.
func (s *Stats) IsReady() bool {
	return s.NextReviewAt().Before(time.Now())
}
