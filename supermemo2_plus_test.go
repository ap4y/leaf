package leaf

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSM2PlusNextReviewAt(t *testing.T) {
	sm := &Supermemo2Plus{LastReviewedAt: time.Unix(100, 0), Interval: 1}
	assert.Equal(t, int64(86500), sm.NextReviewAt().Unix())

	sm = &Supermemo2Plus{LastReviewedAt: time.Unix(100, 0).Add(-24 * time.Hour), Interval: 1}
	assert.Equal(t, int64(100), sm.NextReviewAt().Unix())

	sm = NewSupermemo2Plus()
	assert.InDelta(t, time.Since(sm.NextReviewAt()), time.Hour, float64(time.Minute))
	interval := sm.Advance(1)
	assert.InDelta(t, time.Duration(24*interval)*time.Hour, time.Until(sm.NextReviewAt()), float64(time.Minute))
}

func TestSM2PlusPercentOverdue(t *testing.T) {
	sm := Supermemo2Plus{LastReviewedAt: time.Now().Add(-time.Hour), Interval: 1}
	assert.InDelta(t, 0.04, sm.PercentOverdue(), 0.01)

	sm = Supermemo2Plus{LastReviewedAt: time.Now().Add(-48 * time.Hour), Interval: 1}
	assert.InDelta(t, 2.0, sm.PercentOverdue(), 0.01)
}

func TestSM2PlusRecord(t *testing.T) {
	results := [][]float64{
		{0.04, 0.01, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
		{0.45, 0.88, 1.51, 2.17, 2.82, 3.67, 4.77, 6.21, 8.07},
		{0.53, 1.43, 4.0, 11.59, 34.71, 104.13, 312.41, 937.23, 2811.7},
	}
	for idx, rating := range []float64{0.5, 0.6, 1.0} {
		t.Run(fmt.Sprintf("%f", rating), func(t *testing.T) {
			sm := NewSupermemo2Plus()
			intervals := []float64{}
			for i := 0; i < 9; i++ {
				interval := sm.Advance(rating)
				intervals = append(intervals, interval)

				curInterval := sm.Interval * 24 * float64(time.Hour)
				sm.LastReviewedAt = time.Now().Add(-time.Duration(curInterval))
			}

			assert.InDeltaSlice(t, results[idx], intervals, 0.01)
		})
	}

	t.Run("sequence", func(t *testing.T) {
		sm := NewSupermemo2Plus()
		intervals := []float64{}
		for _, rating := range []float64{1, 1, 1, 1, 0.5, 1} {
			interval := sm.Advance(rating)
			intervals = append(intervals, interval)

			curInterval := sm.Interval * 24 * float64(time.Hour)
			sm.LastReviewedAt = time.Now().Add(-time.Duration(curInterval))
		}

		assert.InDeltaSlice(t, []float64{0.53, 1.43, 4.0, 11.59, 1.79, 4.73}, intervals, 0.01)

		historical := []float64{}
		for _, snap := range sm.Historical {
			assert.NotNil(t, snap.Timestamp)
			historical = append(historical, snap.Interval)
		}
		assert.InDeltaSlice(t, []float64{0.2, 0.53, 1.43, 4.0, 11.59, 1.79}, historical, 0.01)
	})
}

func TestSM2PlusJsonMarshalling(t *testing.T) {
	sm := &Supermemo2Plus{LastReviewedAt: time.Unix(100, 0), Interval: 1, Difficulty: 0.2}
	res, err := json.Marshal(sm)
	require.NoError(t, err)

	newSM := new(Supermemo2Plus)
	require.NoError(t, json.Unmarshal(res, newSM))
	assert.Equal(t, sm, newSM)
}
