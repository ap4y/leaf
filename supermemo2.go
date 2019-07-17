package leaf

import (
	"encoding/json"
	"math"
	"time"
)

// Supermemo2 calculates review intervals using SM2 algorithm
type Supermemo2 struct {
	LastReviewedAt time.Time
	Interval       float64
	Easiness       float64
	Correct        int
	Total          int
	Historical     []IntervalSnapshot
}

// NewSupermemo2 returns a new Supermemo2 instance
func NewSupermemo2() *Supermemo2 {
	return &Supermemo2{
		LastReviewedAt: time.Now(),
		Interval:       0,
		Easiness:       2.5,
		Correct:        0,
		Total:          0,
	}
}

// NextReviewAt returns next review timestamp for a card.
func (sm *Supermemo2) NextReviewAt() time.Time {
	return sm.LastReviewedAt.Add(time.Duration(24*sm.Interval) * time.Hour)
}

// SortParam returns values that should used as a review order for cards
func (sm *Supermemo2) SortParam() float64 {
	return sm.Interval
}

// Advance advances supermemo state for a card.
func (sm *Supermemo2) Advance(rating float64) float64 {
	sm.Total++
	sm.LastReviewedAt = time.Now()

	sm.Easiness += 0.1 - (1-rating)*(0.4+(1-rating)*0.5)
	sm.Easiness = math.Max(sm.Easiness, 1.3)

	interval := 1.0
	if rating >= ratingSuccess {
		if sm.Total == 2 {
			interval = 6
		} else if sm.Total > 2 {
			interval = math.Round(sm.Interval * sm.Easiness)
		}
		sm.Correct++
	} else {
		sm.Correct = 0
	}

	if sm.Historical == nil {
		sm.Historical = make([]IntervalSnapshot, 0)
	}
	sm.Historical = append(
		sm.Historical,
		IntervalSnapshot{time.Now().Unix(), sm.Interval, sm.Easiness},
	)

	sm.Interval = interval
	return interval
}

// MarshalJSON implements json.Marshaller for Supermemo2
func (sm *Supermemo2) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		LastReviewedAt time.Time
		Interval       float64
		Easiness       float64
		Correct        int
		Total          int
		Historical     []IntervalSnapshot
	}{sm.LastReviewedAt, sm.Interval, sm.Easiness, sm.Correct, sm.Total, sm.Historical})
}

// UnmarshalJSON implements json.Unmarshaller for Supermemo2
func (sm *Supermemo2) UnmarshalJSON(b []byte) error {
	payload := &struct {
		LastReviewedAt time.Time
		Interval       float64
		Easiness       float64
		Correct        int
		Total          int
		Historical     []IntervalSnapshot
	}{}

	if err := json.Unmarshal(b, payload); err != nil {
		return err
	}

	sm.LastReviewedAt = payload.LastReviewedAt
	sm.Easiness = payload.Easiness
	sm.Interval = payload.Interval
	sm.Correct = payload.Correct
	sm.Total = payload.Total
	sm.Historical = payload.Historical
	return nil
}
