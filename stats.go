package leaf

import (
	"encoding/json"
	"math"
	"time"
)

const ratingSuccess = 0.6

// SRSAlgorithm calculates review intervals
type SRSAlgorithm interface {
	json.Marshaler
	json.Unmarshaler

	// Advance advances supermemo state for a card.
	Advance(rating float64) (interval float64)
	// NextReviewAt returns next review timestamp for a card.
	NextReviewAt() time.Time
	// Less defines card order for the review.
	Less(other SRSAlgorithm) bool
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
	SRSAlgorithm
}

// CardWithStats joins Stats to a Card
type CardWithStats struct {
	Card
	*Stats
}

// SRS defines supported spaced-repetiton algorithms.
type SRS string

const (
	// SRSSupermemo2 represents Supermemo2 algorithm
	SRSSupermemo2 SRS = "sm2"
	// SRSSupermemo2Plus represents Supermemo2Plus algorithm
	SRSSupermemo2Plus = "sm2+"
	// SRSSupermemo2PlusCustom represents Supermemo2PlusCustom algorithm
	SRSSupermemo2PlusCustom = "sm2+c"
	// SRSEbisu represents Ebisu algorithm
	SRSEbisu = "ebs"
)

// NewStats returns a new Stats initialized with provided algorithm
// with default values. Supported values: sm2, sm2+, sm2+c. If smAlgo
// is missing or unknown will default to Supermemo2PlusCustom.
func NewStats(srs SRS) *Stats {
	var sm SRSAlgorithm
	switch srs {
	case SRSSupermemo2:
		sm = NewSupermemo2()
	case SRSSupermemo2Plus:
		sm = NewSupermemo2Plus()
	case SRSEbisu:
		sm = NewEbisu()
	default:
		sm = NewSupermemo2PlusCustom()
	}
	return &Stats{sm}
}

// IsReady signals whether card is read for review.
func (s *Stats) IsReady() bool {
	return s.NextReviewAt().Before(time.Now())
}
