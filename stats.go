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
	// Less defines card order for the review.
	Less(other Supermemo) bool
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
}

// CardWithStats joins Stats to a Card
type CardWithStats struct {
	Card
	*Stats
}

// SupermemoAlgorithm defines supported supermemo algorithms.
type SupermemoAlgorithm string

const (
	// SM2 represents Supermemo2 algorithm
	SM2 SupermemoAlgorithm = "sm2"
	// SM2Plus represents Supermemo2Plus algorithm
	SM2Plus = "sm2+"
	// SM2PlusCustom represents Supermemo2PlusCustom algorithm
	SM2PlusCustom = "sm2+c"
)

// NewStats returns a new Stats initialized with provided algorithm
// with default values. Supported values: sm2, sm2+, sm2+c. If smAlgo
// is missing or unknown will default to Supermemo2PlusCustom.
func NewStats(smAlgo SupermemoAlgorithm) *Stats {
	var sm Supermemo
	switch smAlgo {
	case SM2:
		sm = NewSupermemo2()
	case SM2Plus:
		sm = NewSupermemo2Plus()
	default:
		sm = NewSupermemo2PlusCustom()
	}
	return &Stats{sm}
}

// IsReady signals whether card is read for review.
func (s *Stats) IsReady() bool {
	return s.NextReviewAt().Before(time.Now())
}
