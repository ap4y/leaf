package leaf

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsReady(t *testing.T) {
	s := &Stats{&Supermemo2Plus{LastReviewedAt: time.Now(), Interval: 1}}
	assert.False(t, s.IsReady())

	s = &Stats{&Supermemo2Plus{LastReviewedAt: time.Now().Add(-24 * time.Hour), Interval: 1}}
	assert.True(t, s.IsReady())

	s = DefaultStats()
	assert.True(t, s.IsReady())
	s.Advance(5)
	assert.False(t, s.IsReady())
}
