package ui

import (
	"github.com/ap4y/leaf"
)

// UI renders session state.
type UI interface {
	Render(state *SessionState) error
}

// SessionState state holds public state of the ReviewSession.
type SessionState struct {
	DeckName  string `json:"deck"`
	Total     int    `json:"total"`
	Left      int    `json:"left"`
	Question  string `json:"question"`
	AnswerLen int    `json:"answerLen"`

	session *leaf.ReviewSession
}

// NewSessionState constructs a new SessionState.
func NewSessionState(session *leaf.ReviewSession) *SessionState {
	s := &SessionState{
		DeckName:  session.DeckName(),
		Total:     session.Total(),
		Left:      session.Left(),
		Question:  session.Next(),
		AnswerLen: len(session.CorrectAnswer()),
		session:   session,
	}

	return s
}

// ResolveAnswer submits answer to a session.
func (s *SessionState) ResolveAnswer(userInput string) (isCorrect bool, correctAnswer string) {
	correctAnswer = s.session.CorrectAnswer()
	isCorrect, _ = s.session.Answer(userInput)
	s.Left = s.session.Left()
	return isCorrect, correctAnswer
}

// Advance fetches next question if available or sets session to finished otherwise.
func (s *SessionState) Advance() {
	s.Question = s.session.Next()
	s.AnswerLen = len(s.session.CorrectAnswer())
}
