package leaf

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsReady(t *testing.T) {
	s := &Stats{&Supermemo2Plus{LastReviewedAt: time.Now(), Interval: 1}, nil}
	assert.False(t, s.IsReady())

	s = &Stats{&Supermemo2Plus{LastReviewedAt: time.Now().Add(-24 * time.Hour), Interval: 1}, nil}
	assert.True(t, s.IsReady())

	s = DefaultStats()
	assert.True(t, s.IsReady())
	s.Advance(5)
	assert.False(t, s.IsReady())
}

func TestRating(t *testing.T) {
	tcs := []struct {
		mistakes int
		rating   float64
	}{{0, 1}, {1, 0.59}, {2, 0.39}, {3, 0.19}, {4, 0}, {5, 0}, {6, 0}}
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("%d", tc.mistakes), func(t *testing.T) {
			s := DefaultStats()
			assert.InDelta(t, tc.rating, s.Rate(tc.mistakes), 0.01)
		})
	}
}
