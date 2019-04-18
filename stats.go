package leaf

import (
	"math"
	"time"
)

const ratingSuccess = 0.6

// Stats store SM2+ parameters for a Card.
type Stats struct {
	LastReviewedAt time.Time
	Difficulty     float64
	Interval       float64

	initial bool
}

// CardWithStats joins Stats to a Card
type CardWithStats struct {
	Card
	*Stats
}

// DefaultStats returns a new Stats initialized with default values.
func DefaultStats() *Stats {
	return &Stats{time.Now(), 0.3, 0.2, true}
}

// NextReviewAt returns next review timestamp for a card.
func (s *Stats) NextReviewAt() time.Time {
	if s.initial {
		return time.Now()
	}

	return s.LastReviewedAt.Add(time.Duration(24*s.Interval) * time.Hour)
}

// IsReady signals whether card is read for review.
func (s *Stats) IsReady() bool {
	if s.initial {
		return true
	}

	return s.NextReviewAt().Before(time.Now())
}

// PercentOverdue returns corresponding SM2+ value for a Card.
func (s *Stats) PercentOverdue() float64 {
	percentOverdue := time.Now().Sub(s.LastReviewedAt).Hours() / float64(24*s.Interval)
	return math.Min(2, percentOverdue)
}

// Record advances SM2+ state for a card.
func (s *Stats) Record(rating float64) float64 {
	s.initial = false
	success := rating >= ratingSuccess
	percentOverdue := float64(1)
	if success {
		percentOverdue = s.PercentOverdue()
	}

	s.Difficulty += percentOverdue / 50 * (8 - 9*rating)
	s.Difficulty = math.Max(0, math.Min(1, s.Difficulty))
	difficultyWeight := 3.5 - 1.7*s.Difficulty

	minInterval := math.Min(1.0, s.Interval)
	factor := minInterval / math.Pow(difficultyWeight, 2)
	if success {
		minInterval = 0.2
		factor = minInterval + (difficultyWeight-1)*percentOverdue
	}

	s.LastReviewedAt = time.Now()
	s.Interval = math.Max(minInterval, math.Min(s.Interval*factor, 300))
	return s.Interval
}
