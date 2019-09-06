package leaf

import (
	"encoding/json"
	"time"
)

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
