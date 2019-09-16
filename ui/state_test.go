package ui

import (
	"testing"

	"github.com/ap4y/leaf"
	"github.com/stretchr/testify/assert"
)

func TestSessionState(t *testing.T) {
	cards := []leaf.CardWithStats{
		{leaf.Card{"foo", []string{"bar"}}, leaf.NewStats(leaf.SRSSupermemo2Plus)},
		{leaf.Card{"bar", []string{"baz"}}, leaf.NewStats(leaf.SRSSupermemo2Plus)},
	}

	stats := make(map[string]*leaf.Stats)
	s := leaf.NewReviewSession(cards, leaf.RatingTypeAuto, func(card *leaf.CardWithStats) error {
		stats[card.Question] = card.Stats
		return nil
	})

	state := NewSessionState(s)
	t.Run("state", func(t *testing.T) {
		assert.Equal(t, 2, state.Total)
		assert.Equal(t, 2, state.Left)
		assert.Equal(t, "foo", state.Question)
		assert.Equal(t, 3, state.AnswerLen)
	})

	t.Run("ResolveAnswer", func(t *testing.T) {
		answer := state.ResolveAnswer()
		assert.Equal(t, "bar", answer)
		assert.Equal(t, 2, state.Left)
	})

	t.Run("Advance - incorrect", func(t *testing.T) {
		state.Advance(leaf.ReviewScoreAgain)
		assert.Equal(t, "bar", state.Question)
		assert.Equal(t, 3, state.AnswerLen)
		assert.Equal(t, 2, state.Left)
	})

	t.Run("Advance - correct", func(t *testing.T) {
		state.Advance(leaf.ReviewScoreGood)
		assert.Equal(t, 1, state.Left)
	})
}
