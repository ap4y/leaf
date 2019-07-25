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

func TestNextReviewAt(t *testing.T) {
	sm := &Supermemo2PlusCustom{Supermemo2Plus{LastReviewedAt: time.Unix(100, 0), Interval: 1}}
	assert.Equal(t, int64(86500), sm.NextReviewAt().Unix())

	sm = &Supermemo2PlusCustom{Supermemo2Plus{LastReviewedAt: time.Unix(100, 0).Add(-24 * time.Hour), Interval: 1}}
	assert.Equal(t, int64(100), sm.NextReviewAt().Unix())

	sm = NewSupermemo2PlusCustom()
	assert.InDelta(t, time.Since(sm.NextReviewAt()), time.Duration(0), float64(time.Minute))
	interval := sm.Advance(1)
	assert.InDelta(t, time.Duration(24*interval)*time.Hour, time.Until(sm.NextReviewAt()), float64(time.Minute))
}

func TestPercentOverdue(t *testing.T) {
	sm := &Supermemo2PlusCustom{Supermemo2Plus{LastReviewedAt: time.Now().Add(-time.Hour), Interval: 1}}
	assert.InDelta(t, 0.04, sm.PercentOverdue(), 0.01)

	sm = &Supermemo2PlusCustom{Supermemo2Plus{LastReviewedAt: time.Now().Add(-48 * time.Hour), Interval: 1}}
	assert.InDelta(t, 2.0, sm.PercentOverdue(), 0.01)
}

func TestLess(t *testing.T) {
	sm1 := &Supermemo2PlusCustom{Supermemo2Plus{LastReviewedAt: time.Now().Add(-time.Hour), Interval: 1}}
	sm2 := &Supermemo2PlusCustom{Supermemo2Plus{LastReviewedAt: time.Now().Add(-48 * time.Hour), Interval: 1}}

	slice := []Supermemo{sm1, sm2}
	sort.Slice(slice, func(i, j int) bool { return slice[j].Less(slice[i]) })
	assert.Equal(t, []Supermemo{sm2, sm1}, slice)
}

func TestRecord(t *testing.T) {
	results := [][]float64{
		{0.2, 0.2, 0.2, 0.2, 0.2, 0.2, 0.2, 0.2, 0.2},
		{0.35, 0.69, 1.27, 2.16, 3.42, 4.97, 6.60, 7.92, 8.52},
		{0.37, 0.86, 2.0, 4.76, 11.56, 28.60, 72.14, 185.44, 300},
	}
	for idx, rating := range []float64{0.5, 0.6, 1.0} {
		t.Run(fmt.Sprintf("%f", rating), func(t *testing.T) {
			sm := NewSupermemo2PlusCustom()
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
		sm := NewSupermemo2PlusCustom()
		intervals := []float64{}
		for _, rating := range []float64{1, 1, 1, 1, 0.5, 1} {
			interval := sm.Advance(rating)
			intervals = append(intervals, interval)

			curInterval := sm.Interval * 24 * float64(time.Hour)
			sm.LastReviewedAt = time.Now().Add(-time.Duration(curInterval))
		}

		assert.InDeltaSlice(t, []float64{0.37, 0.86, 2.00, 4.76, 1, 2.25}, intervals, 0.01)

		historical := []float64{}
		for _, snap := range sm.Historical {
			assert.NotNil(t, snap.Timestamp)
			historical = append(historical, snap.Interval)
		}
		assert.InDeltaSlice(t, []float64{0.2, 0.37, 0.86, 2.00, 4.76, 1}, historical, 0.01)
	})
}

func TestJsonMarshalling(t *testing.T) {
	sm := &Supermemo2PlusCustom{Supermemo2Plus{LastReviewedAt: time.Unix(100, 0), Interval: 1, Difficulty: 0.2}}
	res, err := json.Marshal(sm)
	require.NoError(t, err)

	newSM := new(Supermemo2PlusCustom)
	require.NoError(t, json.Unmarshal(res, newSM))
	assert.Equal(t, sm, newSM)
}
