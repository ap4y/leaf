package ui

import (
	"github.com/ap4y/leaf"
)

// SessionState state holds public state of the ReviewSession.
type SessionState struct {
	Total      int             `json:"total"`
	Left       int             `json:"left"`
	Question   string          `json:"question"`
	AnswerLen  int             `json:"answer_length"`
	RatingType leaf.RatingType `json:"rating_type"`

	session *leaf.ReviewSession
	rater   leaf.Rater
}

// NewSessionState constructs a new SessionState.
func NewSessionState(session *leaf.ReviewSession) *SessionState {
	var rater leaf.Rater
	if session.RatingType() == leaf.RatingTypeSelf {
		rater = leaf.TableRater()
	} else {
		rater = leaf.HarshRater()
	}

	s := &SessionState{
		Total:      session.Total(),
		Left:       session.Left(),
		Question:   session.Next(),
		AnswerLen:  len(session.CorrectAnswer()),
		RatingType: session.RatingType(),
		session:    session,
		rater:      rater,
	}

	return s
}

// ResolveAnswer submits answer to a session.
func (s *SessionState) ResolveAnswer() (correctAnswer string) {
	return s.session.CorrectAnswer()
}

// Advance fetches next question if available or sets session to finished otherwise.
func (s *SessionState) Advance(score leaf.ReviewScore) {
	rating := s.rater.Rate(s.Question, score)

	if score == leaf.ReviewScoreAgain {
		s.session.Again() // nolint: errcheck
	} else {
		s.session.Rate(rating) // nolint: errcheck
		s.Left = s.session.Left()
	}

	s.Question = s.session.Next()
	s.AnswerLen = len(s.session.CorrectAnswer())
}
