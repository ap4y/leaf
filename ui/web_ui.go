//go:generate esc -o ui/static.go -prefix ui/static -pkg ui ui/static
package ui

import (
	"encoding/json"
	"log"
	"net/http"
)

// DevMode makes WebUI to use local static files.
var DevMode = false

// WebUI implements web UI.
type WebUI struct {
	addr string
}

// NewWebUI construct a new WebUI instance.
func NewWebUI(addr string) *WebUI {
	return &WebUI{addr: addr}
}

// Handler returns net.Handler for provided session state.
func (ui *WebUI) Handler(s *SessionState) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(FS(DevMode)))
	mux.HandleFunc("/next", nextHandler(s))
	mux.HandleFunc("/resolve", resolveHandler(s))

	return mux
}

// Render renders current ui state using termbox.
func (ui *WebUI) Render(s *SessionState) error {
	mux := ui.Handler(s)
	log.Println("Serving HTTP on", ui.addr)
	return http.ListenAndServe(ui.addr, mux)
}

func nextHandler(s *SessionState) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, "invalid method", http.StatusBadRequest)
			return
		}

		s.Advance()
		if err := json.NewEncoder(w).Encode(s); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
	}
}

func resolveHandler(s *SessionState) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "invalid method", http.StatusBadRequest)
			return
		}

		data := map[string]string{}
		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		result, correct := s.ResolveAnswer(data["answer"])
		res := map[string]interface{}{"is_correct": result, "correct": correct}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
	}
}
