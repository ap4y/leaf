package leaf

import (
	"encoding/json"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSM2NextReviewAt(t *testing.T) {
	sm := &Supermemo2{LastReviewedAt: time.Unix(100, 0), Interval: 1}
	assert.Equal(t, int64(86500), sm.NextReviewAt().Unix())

	sm = &Supermemo2{LastReviewedAt: time.Unix(100, 0).Add(-24 * time.Hour), Interval: 1}
	assert.Equal(t, int64(100), sm.NextReviewAt().Unix())

	sm = NewSupermemo2()
	assert.InDelta(t, time.Since(sm.NextReviewAt()), time.Duration(0), float64(time.Minute))
	interval := sm.Advance(1)
	assert.InDelta(t, time.Duration(24*interval)*time.Hour, time.Until(sm.NextReviewAt()), float64(time.Minute))
}

func TestSM2Less(t *testing.T) {
	sm1 := &Supermemo2{Interval: 1}
	sm2 := &Supermemo2{Interval: 0.2}

	slice := []SRSAlgorithm{sm1, sm2}
	sort.Slice(slice, func(i, j int) bool { return slice[j].Less(slice[i]) })
	assert.Equal(t, []SRSAlgorithm{sm2, sm1}, slice)
}

func TestSM2Record(t *testing.T) {
	results := [][]float64{
		{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0},
		{1.0, 6.0, 12.0, 23.0, 41.0, 68.0, 103.0, 142.0, 185.0},
		{1.0, 6.0, 17.0, 49.0, 147.0, 456.0, 1459.0, 4815.0, 16371.0},
	}
	for idx, rating := range []float64{0.5, 0.6, 1.0} {
		t.Run(fmt.Sprintf("%f", rating), func(t *testing.T) {
			sm := NewSupermemo2()
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
		sm := NewSupermemo2()
		intervals := []float64{}
		for _, rating := range []float64{1, 1, 1, 1, 0.5, 1} {
			interval := sm.Advance(rating)
			intervals = append(intervals, interval)

			curInterval := sm.Interval * 24 * float64(time.Hour)
			sm.LastReviewedAt = time.Now().Add(-time.Duration(curInterval))
		}

		assert.InDeltaSlice(t, []float64{1.0, 6.0, 17.0, 49.0, 1.0, 3.0}, intervals, 0.01)

		historical := []float64{}
		for _, snap := range sm.Historical {
			assert.NotNil(t, snap.Timestamp)
			historical = append(historical, snap.Interval)
		}

		assert.InDeltaSlice(t, []float64{0.0, 1.0, 6.0, 17.0, 49.0, 1.0}, historical, 0.01)
	})
}

func TestSM2JsonMarshalling(t *testing.T) {
	sm := &Supermemo2{LastReviewedAt: time.Unix(100, 0), Interval: 1, Easiness: 2.5}
	res, err := json.Marshal(sm)
	require.NoError(t, err)

	newSM := new(Supermemo2)
	require.NoError(t, json.Unmarshal(res, newSM))
	assert.Equal(t, int64(100), newSM.LastReviewedAt.Unix())
	assert.Equal(t, 1.0, newSM.Interval)
	assert.Equal(t, 2.5, newSM.Easiness)
}
