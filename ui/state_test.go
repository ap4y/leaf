package ui

import (
	"testing"

	"github.com/ap4y/leaf"
	"github.com/stretchr/testify/assert"
)

func TestSessionState(t *testing.T) {
	cards := []leaf.CardWithStats{
		{Card: leaf.Card{Question: "foo", Sides: []string{"bar"}}, Stats: leaf.NewStats(leaf.SRSSupermemo2Plus)},
		{Card: leaf.Card{Question: "bar", Sides: []string{"baz"}}, Stats: leaf.NewStats(leaf.SRSSupermemo2Plus)},
	}

	stats := make(map[string]*leaf.Stats)
	s := leaf.NewReviewSession(cards, []string{"answer"}, leaf.RatingTypeAuto, func(card *leaf.CardWithStats) error {
		stats[card.Question] = card.Stats
		return nil
	})

	state := NewSessionState(s)
	t.Run("state", func(t *testing.T) {
		assert.Equal(t, 2, state.Total)
		assert.Equal(t, 2, state.Left)
		assert.Equal(t, "foo", state.Question)
		assert.Equal(t, 3, state.AnswerLen)
		assert.Equal(t, []string{"answer"}, state.Sides)
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

func TestSessionStateUnicode(t *testing.T) {
	cards := []leaf.CardWithStats{
		{Card: leaf.Card{Question: "hello", Sides: []string{"おはよう"}}, Stats: leaf.NewStats(leaf.SRSSupermemo2Plus)},
	}

	stats := make(map[string]*leaf.Stats)
	s := leaf.NewReviewSession(cards, nil, leaf.RatingTypeAuto, func(card *leaf.CardWithStats) error {
		stats[card.Question] = card.Stats
		return nil
	})

	state := NewSessionState(s)
	t.Run("state", func(t *testing.T) {
		assert.Equal(t, 1, state.Total)
		assert.Equal(t, 1, state.Left)
		assert.Equal(t, "hello", state.Question)
		assert.Equal(t, 4, state.AnswerLen)
	})
}
