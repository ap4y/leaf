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

func TestSM2PlusNextReviewAt(t *testing.T) {
	sm := &Supermemo2Plus{LastReviewedAt: time.Unix(100, 0), Interval: 1}
	assert.Equal(t, int64(86500), sm.NextReviewAt().Unix())

	sm = &Supermemo2Plus{LastReviewedAt: time.Unix(100, 0).Add(-24 * time.Hour), Interval: 1}
	assert.Equal(t, int64(100), sm.NextReviewAt().Unix())

	sm = NewSupermemo2Plus()
	assert.InDelta(t, time.Since(sm.NextReviewAt()), time.Duration(0), float64(time.Minute))
	interval := sm.Advance(1)
	assert.InDelta(t, time.Duration(24*interval)*time.Hour, time.Until(sm.NextReviewAt()), float64(time.Minute))
}

func TestSM2PlusPercentOverdue(t *testing.T) {
	sm := Supermemo2Plus{LastReviewedAt: time.Now().Add(-time.Hour), Interval: 1}
	assert.InDelta(t, 0.04, sm.PercentOverdue(), 0.01)

	sm = Supermemo2Plus{LastReviewedAt: time.Now().Add(-48 * time.Hour), Interval: 1}
	assert.InDelta(t, 2.0, sm.PercentOverdue(), 0.01)
}

func TestSM2PlusLess(t *testing.T) {
	sm1 := &Supermemo2Plus{LastReviewedAt: time.Now().Add(-time.Hour), Interval: 1}
	sm2 := &Supermemo2Plus{LastReviewedAt: time.Now().Add(-48 * time.Hour), Interval: 1}

	slice := []SRSAlgorithm{sm1, sm2}
	sort.Slice(slice, func(i, j int) bool { return slice[j].Less(slice[i]) })
	assert.Equal(t, []SRSAlgorithm{sm2, sm1}, slice)
}

func TestSM2PlusRecord(t *testing.T) {
	results := [][]float64{
		{0.04, 0.01, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
		{0.41, 0.82, 1.45, 2.17, 2.82, 3.67, 4.77, 6.21, 8.07},
		{0.46, 1.23, 3.42, 9.84, 29.27, 87.83, 263.49, 790.49, 2371.48},
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

		assert.InDeltaSlice(t, []float64{0.46, 1.23, 3.42, 9.84, 1.54, 4.05}, intervals, 0.01)

		historical := []float64{}
		for _, snap := range sm.Historical {
			assert.NotNil(t, snap.Timestamp)
			historical = append(historical, snap.Interval)
		}
		assert.InDeltaSlice(t, []float64{0.2, 0.46, 1.23, 3.42, 9.84, 1.54}, historical, 0.01)
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
