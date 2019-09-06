package leaf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReviewSession(t *testing.T) {
	cards := []*CardWithStats{
		{Card{"foo", []string{"bar"}}, NewStats(SRSSupermemo2PlusCustom)},
		{Card{"bar", []string{"baz"}}, NewStats(SRSSupermemo2PlusCustom)},
	}

	stats := make(map[string]*Stats)
	s := NewReviewSession(cards, HarshRater{}, func(card *CardWithStats) error {
		stats[card.Question] = card.Stats
		return nil
	})

	t.Run("StartedAt", func(t *testing.T) {
		assert.NotNil(t, s.StartedAt())
	})

	t.Run("Total", func(t *testing.T) {
		assert.Equal(t, 2, s.Total())
	})

	t.Run("Left", func(t *testing.T) {
		assert.Equal(t, 2, s.Left())
	})

	t.Run("Next", func(t *testing.T) {
		assert.Equal(t, "foo", s.Next())
	})

	t.Run("CorrectAnswer", func(t *testing.T) {
		assert.Equal(t, "bar", s.CorrectAnswer())
	})

	t.Run("Answer - incorrect", func(t *testing.T) {
		correct, err := s.Answer("123")
		require.NoError(t, err)
		assert.False(t, correct)
		assert.Equal(t, 2, s.Left())
		assert.Equal(t, "bar", s.Next())
	})

	t.Run("Answer - correct", func(t *testing.T) {
		correct, err := s.Answer("baz")
		require.NoError(t, err)
		assert.True(t, correct)
		assert.Equal(t, 1, s.Left())
		assert.Equal(t, "foo", s.Next())
	})

	t.Run("Answer - multiple incorrect", func(t *testing.T) {
		for i := 0; i < 4; i++ {
			s.Answer("123")
		}
		assert.Equal(t, 1, s.Left())
	})

	t.Run("Answer - finish session", func(t *testing.T) {
		correct, err := s.Answer("bar")
		require.NoError(t, err)
		assert.True(t, correct)
		assert.Equal(t, 0, s.Left())
	})

	fooStats := stats["foo"].SRSAlgorithm.(*Supermemo2PlusCustom)
	assert.InDelta(t, 0.52, fooStats.Difficulty, 0.01)
	assert.InDelta(t, 0.2, fooStats.Interval, 0.01)

	barStats := stats["bar"].SRSAlgorithm.(*Supermemo2PlusCustom)
	assert.InDelta(t, 0.28, barStats.Difficulty, 0.01)
	assert.InDelta(t, 0.37, barStats.Interval, 0.01)
}
