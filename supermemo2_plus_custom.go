package leaf

import (
	"math"
	"time"
)

// Supermemo2PlusCustom calculates review intervals using altered SM2+ algorithm
type Supermemo2PlusCustom struct {
	Supermemo2Plus
}

// NewSupermemo2PlusCustom returns a new Supermemo2PlusCustom instance
func NewSupermemo2PlusCustom() *Supermemo2PlusCustom {
	sm := NewSupermemo2Plus()
	return &Supermemo2PlusCustom{*sm}
}

// Less defines card order for the review.
func (sm *Supermemo2PlusCustom) Less(other Supermemo) bool {
	return sm.PercentOverdue() < other.(*Supermemo2PlusCustom).PercentOverdue()
}

// Advance advances supermemo state for a card.
func (sm *Supermemo2PlusCustom) Advance(rating float64) float64 {
	success := rating >= ratingSuccess
	percentOverdue := float64(1)
	if success {
		percentOverdue = sm.PercentOverdue()
	}

	sm.Difficulty += percentOverdue / 35 * (8 - 9*rating)
	sm.Difficulty = math.Max(0, math.Min(1, sm.Difficulty))
	difficultyWeight := 3.5 - 1.7*sm.Difficulty

	minInterval := math.Min(1.0, sm.Interval)
	factor := minInterval / math.Pow(difficultyWeight, 2)
	if success {
		minInterval = 0.2
		factor = minInterval + (difficultyWeight-1)*percentOverdue
	}

	sm.LastReviewedAt = time.Now()
	if sm.Historical == nil {
		sm.Historical = make([]IntervalSnapshot, 0)
	}
	sm.Historical = append(
		sm.Historical,
		IntervalSnapshot{time.Now().Unix(), sm.Interval, sm.Difficulty},
	)
	sm.Interval = math.Max(minInterval, math.Min(sm.Interval*factor, 300))
	return sm.Interval
}
