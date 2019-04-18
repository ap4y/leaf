package leaf

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenDeck(t *testing.T) {
	deck, err := OpenDeck("hiragana.org")
	require.NoError(t, err)
	assert.Equal(t, "Hiragana", deck.Name)
	require.Len(t, deck.Cards, 46)

	cards := deck.Cards
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Question > cards[j].Question
	})
	assert.Equal(t, "ん", cards[0].Question)
	assert.Equal(t, "n", cards[0].Answer())
}
