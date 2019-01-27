package leaf

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatsDB(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "leaf.db")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	db, err := OpenBoltStore(tmpfile.Name())
	require.NoError(t, err)

	s1 := Stats{Difficulty: 1}
	require.NoError(t, db.SaveStats("deck1", "foo", &s1))

	s2 := Stats{Difficulty: 2}
	require.NoError(t, db.SaveStats("deck1", "bar", &s2))

	s3 := Stats{Difficulty: 3}
	require.NoError(t, db.SaveStats("deck2", "foo", &s3))

	cards := []string{}
	stats := []Stats{}
	err = db.GetStats("deck1", func(card string, s *Stats) {
		cards = append(cards, card)
		stats = append(stats, *s)
	})
	require.NoError(t, err)

	assert.Equal(t, []string{"bar", "foo"}, cards)
	assert.Equal(t, []Stats{s2, s1}, stats)
}
