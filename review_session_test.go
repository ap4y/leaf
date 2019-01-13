package leaf

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReviewSession(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "leaf.db")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	db, err := OpenStatsDB(tmpfile.Name())
	require.NoError(t, err)

	cards := make(Stack)
	cards["foo"] = []string{"bar"}
	cards["bar"] = []string{"baz"}
	deck := &Deck{Name: "test", Cards: cards}

	s, err := NewReviewSession(deck, db, 2)
	require.NoError(t, err)

	assert.Equal(t, "test", s.DeckName())
	assert.Equal(t, 2, s.Total())
	assert.Equal(t, 2, s.Left())

	q1 := s.Next()
	assert.Equal(t, q1, s.Next())

	if q1 == "foo" {
		assert.Equal(t, "bar", s.CorrectAnswer())
	} else {
		assert.Equal(t, "baz", s.CorrectAnswer())
	}

	correct, err := s.Answer("123")
	require.NoError(t, err)
	assert.False(t, correct)
	assert.Equal(t, 2, s.Left())

	q2 := s.Next()
	assert.NotEqual(t, q1, s.Next())
	correct, err = s.Answer(s.CorrectAnswer())
	require.NoError(t, err)
	assert.True(t, correct)
	assert.Equal(t, 1, s.Left())

	for i := 0; i < 4; i++ {
		s.Answer("123")
	}
	assert.Equal(t, 1, s.Left())

	correct, err = s.Answer(s.CorrectAnswer())
	require.NoError(t, err)
	assert.True(t, correct)
	assert.Equal(t, 0, s.Left())

	stats := make(map[string]*Stats)
	err = db.GetStats("test", func(c string, s *Stats) {
		stats[c] = s
	})
	require.NoError(t, err)

	assert.InDelta(t, 0.45, stats[q1].Difficulty, 0.01)
	assert.InDelta(t, 0.2, stats[q1].Interval, 0.01)

	assert.InDelta(t, 0.29, stats[q2].Difficulty, 0.01)
	assert.InDelta(t, 0.2, stats[q2].Interval, 0.01)

	s, err = NewReviewSession(deck, db, 2)
	require.NoError(t, err)
	assert.Equal(t, 0, s.Total())
}

func TestRating(t *testing.T) {
	tcs := []struct {
		mistakes int
		rating   float64
	}{{0, 1}, {1, 0.8}, {2, 0.6}, {3, 0.4}, {4, 0.19}, {5, 0}, {6, 0}}
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("%d", tc.mistakes), func(t *testing.T) {
			s := &ReviewSession{mistakes: map[string]int{"foo": tc.mistakes}}
			assert.InDelta(t, tc.rating, s.rating("foo"), 0.01)
		})
	}
}
