package ui

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/ap4y/leaf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebUI(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "leaf.db")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	db, err := leaf.OpenBoltStore(tmpfile.Name())
	require.NoError(t, err)

	dm, err := leaf.NewDeckManager("../fixtures", db, leaf.SRSSupermemo2PlusCustom, leaf.OutputFormatOrg)
	require.NoError(t, err)

	srv := NewServer(dm, &leaf.HarshRater{}, 20)

	t.Run("listDecks", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/decks", nil)
		w := httptest.NewRecorder()

		srv.listDecks(w, req)
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		decks := make([]*leaf.DeckStats, 0)
		require.NoError(t, json.NewDecoder(w.Body).Decode(&decks))
		assert.Len(t, decks, 2)
	})

	t.Run("deckStats", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/stats/Hiragana", nil)
		w := httptest.NewRecorder()

		srv.deckStats(w, req)
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		stats := make([]map[string]interface{}, 0)
		require.NoError(t, json.NewDecoder(w.Body).Decode(&stats))
		assert.Len(t, stats, 46)
	})

	t.Run("startReview", func(t *testing.T) {
		req := httptest.NewRequest("POST", "http://example.com/start/Hiragana", nil)
		w := httptest.NewRecorder()

		srv.startSession(w, req)
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		state := new(SessionState)
		require.NoError(t, json.NewDecoder(w.Body).Decode(state))
		assert.Equal(t, 20, state.Total)
		assert.Equal(t, 20, state.Left)
	})

	t.Run("nextCard", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/next", nil)
		w := httptest.NewRecorder()

		srv.nextCard(w, req)
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		state := new(SessionState)
		require.NoError(t, json.NewDecoder(w.Body).Decode(state))
		assert.Equal(t, 20, state.Total)
		assert.Equal(t, 20, state.Left)
	})

	t.Run("resolveCard", func(t *testing.T) {
		req := httptest.NewRequest("POST", "http://example.com/resolve", strings.NewReader("{\"answer\":\"foo\"}"))
		w := httptest.NewRecorder()

		srv.resolveCard(w, req)
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		result := make(map[string]interface{})
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Equal(t, false, result["is_correct"])
	})
}
