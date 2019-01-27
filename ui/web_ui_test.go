package ui

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ap4y/leaf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebUIMux(t *testing.T) {
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
	ui := NewWebUI(":8080")
	mux := ui.Handler(state)

	req := httptest.NewRequest("GET", "http://example.com/next", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	resState := &SessionState{}
	require.NoError(t, json.NewDecoder(res.Body).Decode(resState))
	assert.Equal(t, "test", state.DeckName)
	assert.Equal(t, 2, state.Total)
	assert.Equal(t, 2, state.Left)
	assert.Len(t, state.Question, 3)
	assert.Equal(t, 3, state.AnswerLen)

	body := new(bytes.Buffer)
	require.NoError(t, json.NewEncoder(body).Encode(map[string]string{"answer": "123"}))
	req = httptest.NewRequest("POST", "http://example.com/resolve", body)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	res = w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	result := make(map[string]interface{})
	require.NoError(t, json.NewDecoder(res.Body).Decode(&result))
	assert.False(t, result["is_correct"].(bool))
	assert.Len(t, result["correct"], 3)
}
