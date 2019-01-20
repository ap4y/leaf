//go:generate esc -o static.go -prefix static static
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// WebUI implements web UI.
type WebUI struct {
	addr  string
	state *SessionState
}

// NewWebUI construct a new WebUI instance.
func NewWebUI(addr string) *WebUI {
	return &WebUI{addr: addr}
}

// Render renders current ui state using termbox.
func (ui *WebUI) Render(s *SessionState) error {
	ui.state = s

	http.Handle("/", http.FileServer(FS(false)))
	http.HandleFunc("/next", ui.nextHandler)
	http.HandleFunc("/resolve", ui.resolveHandler)

	log.Println("Serving HTTP on", ui.addr)
	return http.ListenAndServe(ui.addr, nil)
}

func (ui *WebUI) nextHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "invalid method", http.StatusBadRequest)
		return
	}

	ui.state.Advance()
	if err := json.NewEncoder(w).Encode(ui.state); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (ui *WebUI) resolveHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusBadRequest)
		return
	}

	data := map[string]string{}
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	result, correct := ui.state.ResolveAnswer(data["answer"])
	res := map[string]interface{}{"is_correct": result, "correct": correct}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}
