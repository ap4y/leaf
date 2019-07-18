package leaf

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsReady(t *testing.T) {
	s := &Stats{&Supermemo2Plus{LastReviewedAt: time.Now(), Interval: 1}}
	assert.False(t, s.IsReady())

	s = &Stats{&Supermemo2Plus{LastReviewedAt: time.Now().Add(-24 * time.Hour), Interval: 1}}
	assert.True(t, s.IsReady())

	s = NewStats(SM2Plus)
	assert.True(t, s.IsReady())
	s.Advance(5)
	assert.False(t, s.IsReady())
}

func TestHarshRater(t *testing.T) {
	rater := &HarshRater{}

	tcs := []struct {
		mistakes int
		rating   float64
	}{{0, 1}, {1, 0.59}, {2, 0.39}, {3, 0.19}, {4, 0}, {5, 0}, {6, 0}}
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("%d", tc.mistakes), func(t *testing.T) {
			assert.InDelta(t, tc.rating, rater.Rate(tc.mistakes), 0.01)
		})
	}
}
