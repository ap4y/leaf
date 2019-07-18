package ui

import (
	"testing"

	"github.com/ap4y/leaf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionState(t *testing.T) {
	cards := []*leaf.CardWithStats{
		{leaf.Card{"foo", []string{"bar"}}, leaf.NewStats(leaf.SM2Plus)},
		{leaf.Card{"bar", []string{"baz"}}, leaf.NewStats(leaf.SM2Plus)},
	}

	stats := make(map[string]*leaf.Stats)
	s := leaf.NewReviewSession(cards, &leaf.HarshRater{}, func(question string, s *leaf.Stats) error {
		stats[question] = s
		return nil
	})

	state := NewSessionState(s)
	t.Run("state", func(t *testing.T) {
		assert.Equal(t, 2, state.Total)
		assert.Equal(t, 2, state.Left)
		assert.Equal(t, "foo", state.Question)
		assert.Equal(t, 3, state.AnswerLen)
	})

	t.Run("ResolveAnswer - incorrect", func(t *testing.T) {
		res, answer := state.ResolveAnswer("123")
		require.False(t, res)
		assert.Equal(t, "bar", answer)
		assert.Equal(t, 2, state.Left)
	})

	t.Run("ResolveAnswer - incorrect", func(t *testing.T) {
		state.Advance()
		assert.Equal(t, "bar", state.Question)
		assert.Equal(t, 3, state.AnswerLen)
	})

	t.Run("ResolveAnswer - correct", func(t *testing.T) {
		res, _ := state.ResolveAnswer("baz")
		require.True(t, res)
		assert.Equal(t, 1, state.Left)
	})
}
