//go:generate esc -o ui/static.go -prefix ui/static -pkg ui ui/static

package ui

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/ap4y/leaf"
)

type statsResponse struct {
	Card  string      `json:"card"`
	Stats *leaf.Stats `json:"stats"`
}

// Server implements web ui for reviews.
type Server struct {
	dm    *leaf.DeckManager
	rater leaf.Rater

	cardsPerReview int
	sessionState   *SessionState
}

// NewServer construct a new Server instance.
func NewServer(dm *leaf.DeckManager, rater leaf.Rater, cardsPerReview int) *Server {
	return &Server{dm: dm, rater: rater, cardsPerReview: cardsPerReview}
}

// Serve starts listening loop on addr.
func (srv *Server) Serve(addr string, devMode bool) error {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(FS(devMode)))
	mux.HandleFunc("/decks", srv.listDecks)
	mux.HandleFunc("/start/", srv.startSession)
	mux.HandleFunc("/stats/", srv.deckStats)
	mux.HandleFunc("/next", srv.nextCard)
	mux.HandleFunc("/resolve", srv.resolveCard)

	log.Println("Serving HTTP on", addr)
	return http.ListenAndServe(addr, mux)
}

func (srv *Server) listDecks(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	decks, err := srv.dm.ReviewDecks(100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := json.NewEncoder(w).Encode(decks); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (srv *Server) startSession(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	deckName := strings.Replace(req.URL.Path, "/start/", "", -1)
	session, err := srv.dm.ReviewSession(deckName, srv.rater, srv.cardsPerReview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	srv.sessionState = NewSessionState(session)
	if err := json.NewEncoder(w).Encode(srv.sessionState); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (srv *Server) deckStats(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	deckName := strings.Replace(req.URL.Path, "/stats/", "", -1)
	stats, err := srv.dm.DeckStats(deckName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	res := make([]statsResponse, 0)
	for _, stat := range stats {
		res = append(res, statsResponse{
			stat.Question, stat.Stats,
		})
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (srv *Server) nextCard(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	if srv.sessionState == nil {
		http.Error(w, "no active sessions", http.StatusBadRequest)
		return
	}

	srv.sessionState.Advance()
	if err := json.NewEncoder(w).Encode(srv.sessionState); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (srv *Server) resolveCard(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	if srv.sessionState == nil {
		http.Error(w, "no active sessions", http.StatusBadRequest)
		return
	}

	data := map[string]string{}
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	result, correct := srv.sessionState.ResolveAnswer(data["answer"])
	res := map[string]interface{}{"is_correct": result, "correct": correct}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}
