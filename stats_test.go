package leaf

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNextReviewAt(t *testing.T) {
	s := &Stats{LastReviewedAt: time.Unix(100, 0), Interval: 1}
	assert.Equal(t, int64(86500), s.NextReviewAt().Unix())

	s = &Stats{LastReviewedAt: time.Unix(100, 0).Add(-24 * time.Hour), Interval: 1}
	assert.Equal(t, int64(100), s.NextReviewAt().Unix())

	s = DefaultStats()
	assert.InDelta(t, time.Now().Unix(), s.NextReviewAt().Unix(), 100)
	s.Record(5)
	assert.InDelta(t, time.Now().Add(4*time.Hour).Unix(), s.NextReviewAt().Unix(), 100)
}

func TestIsReady(t *testing.T) {
	s := &Stats{LastReviewedAt: time.Now(), Interval: 1}
	assert.False(t, s.IsReady())

	s = &Stats{LastReviewedAt: time.Now().Add(-24 * time.Hour), Interval: 1}
	assert.True(t, s.IsReady())

	s = DefaultStats()
	assert.True(t, s.IsReady())
	s.Record(5)
	assert.False(t, s.IsReady())
}

func TestPercentOverdue(t *testing.T) {
	s := &Stats{LastReviewedAt: time.Now().Add(-time.Hour), Interval: 1}
	assert.InDelta(t, 0.04, s.PercentOverdue(), 0.01)

	s = &Stats{LastReviewedAt: time.Now().Add(-48 * time.Hour), Interval: 1}
	assert.InDelta(t, 2.0, s.PercentOverdue(), 0.01)
}

func TestRecord(t *testing.T) {
	results := [][]float64{
		{0.2, 0.2, 0.2, 0.2, 0.2, 0.2, 0.2, 0.2, 0.2},
		{0.2, 0.43, 0.84, 1.62, 2.99, 5.22, 8.67, 13.63, 20.21},
		{0.2, 0.44, 1, 2.3, 5.35, 12.63, 30.25, 73.45, 180.83},
	}
	for idx, rating := range []float64{0.5, 0.6, 1.0} {
		t.Run(fmt.Sprintf("%f", rating), func(t *testing.T) {
			s := DefaultStats()
			intervals := []float64{}
			for i := 0; i < 9; i++ {
				interval := s.Record(rating)
				intervals = append(intervals, interval)

				curInterval := s.Interval * 24 * float64(time.Hour)
				s.LastReviewedAt = time.Now().Add(-time.Duration(curInterval))
			}

			assert.InDeltaSlice(t, results[idx], intervals, 0.01)
		})
	}

	t.Run("sequence", func(t *testing.T) {
		s := DefaultStats()
		intervals := []float64{}
		for _, rating := range []float64{1, 1, 1, 1, 1, 0.5, 1} {
			interval := s.Record(rating)
			intervals = append(intervals, interval)

			curInterval := s.Interval * 24 * float64(time.Hour)
			s.LastReviewedAt = time.Now().Add(-time.Duration(curInterval))
		}

		assert.InDeltaSlice(t, []float64{0.2, 0.44, 1, 2.3, 5.35, 1, 2.24}, intervals, 0.01)

		historical := []float64{}
		for _, snap := range s.Historical {
			assert.NotNil(t, snap.Timestamp)
			historical = append(historical, snap.Interval)
		}
		assert.InDeltaSlice(t, []float64{0.2, 0.2, 0.44, 1, 2.3, 5.35, 1}, historical, 0.01)
	})
}
