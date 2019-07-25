package leaf

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEbisu(t *testing.T) {
	cases := []struct {
		model [3]float64
		op    [2]float64
		post  [3]float64
	}{
		{[3]float64{3.3, 4.4, 1.0}, [2]float64{0.0, 0.1}, [3]float64{3.0, 5.49, 1.0}},
		{[3]float64{3.3, 4.4, 1.0}, [2]float64{1.0, 0.1}, [3]float64{3.4, 4.4, 1.0}},
		{[3]float64{3.3, 4.4, 1.0}, [2]float64{0.0, 1.0}, [3]float64{3.3, 5.4, 1.0}},
		{[3]float64{3.3, 4.4, 1.0}, [2]float64{1.0, 1.0}, [3]float64{4.3, 4.4, 1.0}},
		{[3]float64{34.4, 3.4, 1.0}, [2]float64{0.0, 1.0}, [3]float64{3.1, 4.64, 8.33}},
		{[3]float64{34.4, 3.4, 1.0}, [2]float64{1.0, 1.0}, [3]float64{3.47, 3.5, 8.33}},
		{[3]float64{34.4, 3.4, 1.0}, [2]float64{0.0, 5.5}, [3]float64{3.29, 4.59, 8.33}},
		{[3]float64{34.4, 3.4, 1.0}, [2]float64{1.0, 5.5}, [3]float64{3.98, 3.48, 8.33}},
		{[3]float64{34.4, 3.4, 1.0}, [2]float64{0.0, 50.0}, [3]float64{3.55, 3.91, 8.33}},
		{[3]float64{34.4, 3.4, 1.0}, [2]float64{1.0, 50.0}, [3]float64{2.89, 3.53, 22.64}},
	}

	for idx, tc := range cases {
		t.Run(fmt.Sprintf("Advance %d", idx), func(t *testing.T) {
			eb := &Ebisu{
				time.Now().Add(toHourDuration(-1 * tc.op[1])),
				tc.model[0],
				tc.model[1],
				tc.model[2],
				make([]IntervalSnapshot, 0),
			}
			eb.Advance(tc.op[0])
			assert.InDelta(t, tc.post[0], eb.Alpha, 0.01)
			assert.InDelta(t, tc.post[1], eb.Beta, 0.01)
			assert.InDelta(t, tc.post[2], eb.Interval, 0.01)
		})
	}
}

func TestEbisuRecord(t *testing.T) {
	results := [][]float64{
		{1.0, 1.0, 0.41, 0.41, 0.41, 0.41, 0.41, 0.41, 0.41},
		{1.0, 1.0, 1.0, 3.06, 3.06, 3.06, 3.06, 3.06, 9.39},
		{1.0, 1.0, 1.0, 3.06, 3.06, 3.06, 3.06, 3.06, 9.39},
	}
	for idx, rating := range []float64{0.5, 0.6, 1.0} {
		t.Run(fmt.Sprintf("%f", rating), func(t *testing.T) {
			srs := NewEbisu()
			intervals := []float64{}
			for i := 0; i < 9; i++ {
				interval := srs.Advance(rating)
				intervals = append(intervals, interval/24)

				curInterval := srs.Interval * float64(time.Hour)
				srs.LastReviewedAt = time.Now().Add(-time.Duration(curInterval))
			}

			assert.InDeltaSlice(t, results[idx], intervals, 0.01)
		})
	}

	t.Run("sequence", func(t *testing.T) {
		srs := NewEbisu()
		intervals := []float64{}
		for _, rating := range []float64{1, 1, 1, 1, 0.5, 1, 1, 1} {
			interval := srs.Advance(rating)
			intervals = append(intervals, interval/24)

			curInterval := srs.Interval * float64(time.Hour)
			srs.LastReviewedAt = time.Now().Add(-time.Duration(curInterval))
		}

		assert.InDeltaSlice(t, []float64{1, 1, 1, 3.06, 1.27, 1.27, 1.27, 3.89}, intervals, 0.01)
	})
}

func TestEbisuPredictRecall(t *testing.T) {
	eb := &Ebisu{LastReviewedAt: time.Now().Add(-1 * time.Hour), Alpha: 4, Beta: 4, Interval: 24}
	assert.InDelta(t, 0.96, eb.predictRecall(), 0.01)

	eb = &Ebisu{LastReviewedAt: time.Now().Add(-1 * time.Hour), Alpha: 2, Beta: 4, Interval: 24}
	assert.InDelta(t, 0.94, eb.predictRecall(), 0.01)
}

func TestEbisuLess(t *testing.T) {
	eb1 := &Ebisu{LastReviewedAt: time.Now().Add(-1 * time.Hour), Alpha: 4, Beta: 4, Interval: 24}
	eb2 := &Ebisu{LastReviewedAt: time.Now().Add(-1 * time.Hour), Alpha: 2, Beta: 4, Interval: 24}

	slice := []Supermemo{eb1, eb2}
	sort.Slice(slice, func(i, j int) bool { return slice[j].Less(slice[i]) })
	assert.Equal(t, []Supermemo{eb2, eb1}, slice)
}

func TestBetaln(t *testing.T) {
	assert.InDelta(t, -70.97, betaln(99, 30.25), 0.01)
	assert.InDelta(t, -28.23, betaln(13, 35), 0.01)
	assert.InDelta(t, -84.94, betaln(47.25, 80.5), 0.01)
	assert.InDelta(t, -59.9, betaln(79.75, 26.25), 0.01)
}

func toHourDuration(hours float64) time.Duration {
	return time.Duration(hours * float64(time.Hour))
}
