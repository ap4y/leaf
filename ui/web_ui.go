//go:generate esc -o ui/static.go -prefix ui/static -pkg ui ui/static

package ui

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/ap4y/leaf"
)

// DevMode makes WebUI to use local static files.
var DevMode = false

// WebUI implements web UI.
type WebUI struct {
	addr         string
	dm           *leaf.DeckManager
	sessionState *SessionState
}

// NewWebUI construct a new WebUI instance.
func NewWebUI(addr string, dm *leaf.DeckManager) *WebUI {
	return &WebUI{addr: addr, dm: dm}
}

// Handler returns net.Handler for provided session state.
func (ui *WebUI) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(FS(DevMode)))
	mux.HandleFunc("/decks", ui.listDecks)
	mux.HandleFunc("/start/", ui.startReview)
	mux.HandleFunc("/next", ui.nextHandler)
	mux.HandleFunc("/resolve", ui.resolveHandler)

	return mux
}

// Render renders current ui state using termbox.
func (ui *WebUI) Render() error {
	mux := ui.Handler()
	log.Println("Serving HTTP on", ui.addr)
	return http.ListenAndServe(ui.addr, mux)
}

func (ui *WebUI) listDecks(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	decks, err := ui.dm.ReviewDecks(100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := json.NewEncoder(w).Encode(decks); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (ui *WebUI) startReview(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	deckName := strings.Replace(req.URL.Path, "/start/", "", -1)
	session, err := ui.dm.ReviewSession(deckName, 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	ui.sessionState = NewSessionState(session)
	if err := json.NewEncoder(w).Encode(ui.sessionState); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (ui *WebUI) nextHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	if ui.sessionState == nil {
		http.Error(w, "no active sessions", http.StatusBadRequest)
		return
	}

	ui.sessionState.Advance()
	if err := json.NewEncoder(w).Encode(ui.sessionState); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (ui *WebUI) resolveHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	if ui.sessionState == nil {
		http.Error(w, "no active sessions", http.StatusBadRequest)
		return
	}

	data := map[string]string{}
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	result, correct := ui.sessionState.ResolveAnswer(data["answer"])
	res := map[string]interface{}{"is_correct": result, "correct": correct}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}
