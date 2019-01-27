package leaf

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
		{0.2, 0.43, 0.92, 1.94, 4.05, 8.32, 16.87, 33.74, 66.55},
		{0.2, 0.44, 1, 2.3, 5.35, 12.63, 30.25, 73.45, 180.83},
	}
	for idx, rating := range []float64{0.6, 0.8, 1.0} {
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
}
