package leaf

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeck(t *testing.T) {
	t.Run("OpenDeck", func(t *testing.T) {
		deck, err := OpenDeck("./fixtures/hiragana.org", OutputFormatOrg)
		require.NoError(t, err)
		assert.Equal(t, "Hiragana", deck.Name)
		assert.Equal(t, RatingTypeAuto, deck.RatingType)
		require.Len(t, deck.Cards, 46)

		cards := deck.Cards
		assert.Equal(t, "あ", cards[0].Question)
		assert.Equal(t, "a", cards[0].Answer())
	})

	t.Run("OpenRichDeck", func(t *testing.T) {
		deck, err := OpenDeck("./fixtures/org-mode.org", OutputFormatHTML)
		require.NoError(t, err)
		assert.Equal(t, "Org-mode", deck.Name)
		assert.Equal(t, RatingTypeSelf, deck.RatingType)
		require.Len(t, deck.Cards, 10)

		cards := deck.Cards
		assert.Equal(t, "<em>emphasis</em>", cards[0].Question)
		assert.Equal(t, "/emphasis/ side2", cards[0].Answer())
		assert.Equal(
			t,
			"Code sample<div class=\"src src-javascript\">\n<div class=\"highlight\">\n<pre>\nconst foo = &#34;test&#34;\n</pre>\n</div>\n</div>\n",
			cards[9].Question,
		)
		assert.Equal(t, "const foo = \"test\"", cards[9].Answer())
	})

	t.Run("Reload", func(t *testing.T) {
		deckfile, err := ioutil.TempFile("", "deck.org")
		require.NoError(t, err)
		defer os.Remove(deckfile.Name())

		_, err = deckfile.Write([]byte("* Test\n** foo\nbar\n"))
		require.NoError(t, err)
		require.NoError(t, deckfile.Sync())

		deck, err := OpenDeck(deckfile.Name(), OutputFormatOrg)
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
