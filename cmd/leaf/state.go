package main

import "github.com/ap4y/leaf"

// UI renders session state.
type UI interface {
	Render(state *SessionState)
}

// Step defines step of the ReviewSession.
type Step int

const (
	// StepAnswering defines step at which Session is awaiting for answer.
	StepAnswering Step = iota
	// StepScore defines step at which Session is showing score for an answer.
	StepScore
	// StepFinished defines step for a finished Session.
	StepFinished
)

// SessionState state holds public state of the ReviewSession.
type SessionState struct {
	Step     Step
	DeckName string
	Total    int
	Left     int
	Question string
	Answer   string
	Result   bool

	session *leaf.ReviewSession
}

// NewSessionState constructs a new SessionState.
func NewSessionState(session *leaf.ReviewSession) *SessionState {
	return &SessionState{
		DeckName: session.DeckName(),
		Total:    session.Total(),
		Left:     session.Left(),
		Question: session.Next(),
		Answer:   session.CorrectAnswer(),
		session:  session,
	}
}

func (s *SessionState) resolveAnswer(userInput string) {
	s.Result, _ = s.session.Answer(userInput)
	s.Step = StepScore
	s.Left = s.session.Left()
}

func (s *SessionState) advance() {
	if s.session.Left() == 0 {
		s.Step = StepFinished
		return
	}

	s.Step = StepAnswering
	s.Question = s.session.Next()
	s.Answer = s.session.CorrectAnswer()
	return
}
