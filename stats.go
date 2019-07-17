package leaf

import (
	"time"
)

const ratingSuccess = 0.6

// Stats store SM2+ parameters for a Card.
type Stats struct {
	Supermemo
}

// CardWithStats joins Stats to a Card
type CardWithStats struct {
	Card
	*Stats
}

// DefaultStats returns a new Stats initialized with default values.
func DefaultStats() *Stats {
	return &Stats{NewSupermemo2Plus()}
}

// IsReady signals whether card is read for review.
func (s *Stats) IsReady() bool {
	return s.NextReviewAt().Before(time.Now())
}
