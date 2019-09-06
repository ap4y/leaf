package leaf

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHarshRater(t *testing.T) {
	rater := HarshRater()

	assert.InDelta(t, 1.0, rater.Rate("foo", ReviewScoreEasy), 0.01)

	tcs := []float64{0.59, 0.39, 0.19, 0, 0}
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("%f", tc), func(t *testing.T) {
			rater.Rate("foo", ReviewScoreAgain)
			assert.InDelta(t, tc, rater.Rate("foo", ReviewScoreEasy), 0.01)
		})
	}
}

func TestTableRater(t *testing.T) {
	rater := TableRater()

	assert.InDelta(t, 0, rater.Rate("foo", ReviewScoreAgain), 0.01)
	assert.InDelta(t, 0.2, rater.Rate("foo", ReviewScoreHard), 0.01)
	assert.InDelta(t, 0.6, rater.Rate("foo", ReviewScoreGood), 0.01)
	assert.InDelta(t, 1.0, rater.Rate("foo", ReviewScoreEasy), 0.01)
}
