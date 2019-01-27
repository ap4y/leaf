package ui

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/ap4y/leaf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionState(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "leaf.db")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	db, err := leaf.OpenBoltStore(tmpfile.Name())
	require.NoError(t, err)

	cards := make(leaf.Stack)
	cards["foo"] = []string{"bar"}
	cards["bar"] = []string{"baz"}
	deck := &leaf.Deck{Name: "test", Cards: cards}

	s, err := leaf.NewReviewSession(deck, db, 2)
	require.NoError(t, err)

	state := NewSessionState(s)
	assert.Equal(t, "test", state.DeckName)
	assert.Equal(t, 2, state.Total)
	assert.Equal(t, 2, state.Left)
	assert.Len(t, state.Question, 3)
	assert.Equal(t, 3, state.AnswerLen)

	correct := "bar"
	if state.Question == "bar" {
		correct = "baz"
	}

	res, answer := state.ResolveAnswer("123")
	require.False(t, res)
	assert.Equal(t, correct, answer)
	assert.Equal(t, 2, state.Left)

	state.Advance()
	correct = "bar"
	if state.Question == "bar" {
		correct = "baz"
	}

	res, _ = state.ResolveAnswer(correct)
	require.True(t, res)
	assert.Equal(t, 1, state.Left)
}
