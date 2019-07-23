package leaf

import (
	"io/ioutil"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeck(t *testing.T) {
	t.Run("OpenDeck", func(t *testing.T) {
		deck, err := OpenDeck("./fixtures/hiragana.org")
		require.NoError(t, err)
		assert.Equal(t, "Hiragana", deck.Name)
		require.Len(t, deck.Cards, 46)

		cards := deck.Cards
		sort.Slice(cards, func(i, j int) bool {
			return cards[i].Question > cards[j].Question
		})
		assert.Equal(t, "ん", cards[0].Question)
		assert.Equal(t, "n", cards[0].Answer())
	})

	t.Run("Reload", func(t *testing.T) {
		deckfile, err := ioutil.TempFile("", "deck.org")
		require.NoError(t, err)
		defer os.Remove(deckfile.Name())

		_, err = deckfile.Write([]byte("* Test\n** foo\nbar\n"))
		require.NoError(t, err)
		require.NoError(t, deckfile.Sync())

		deck, err := OpenDeck(deckfile.Name())
		require.NoError(t, err)
		require.Len(t, deck.Cards, 1)

		require.NoError(t, deck.Reload())
		require.Len(t, deck.Cards, 1)

		time.Sleep(100 * time.Millisecond)
		_, err = deckfile.Write([]byte("** bar\nbaz\n"))
		require.NoError(t, err)
		require.NoError(t, deckfile.Sync())

		require.NoError(t, deck.Reload())
		require.Len(t, deck.Cards, 2)
	})
}
