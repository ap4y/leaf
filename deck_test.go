package leaf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenDeck(t *testing.T) {
	deck, err := OpenDeck("hiragana.org")
	require.NoError(t, err)
	assert.Equal(t, "Hiragana", deck.Name)
	require.Len(t, deck.Cards, 46)
	assert.Equal(t, Card{"a"}, deck.Cards["あ"])
	assert.Equal(t, Card{"n"}, deck.Cards["ん"])
}
