package ui

import (
	"github.com/ap4y/leaf"
)

// RatingType defines types of review rating options.
type RatingType string

const (
	// RatingTypeAuto defines auto rated review option.
	RatingTypeAuto RatingType = "auto"
	// RatingTypeSelf defines self rated review option.
	RatingTypeSelf RatingType = "self"
)

// SessionState state holds public state of the ReviewSession.
type SessionState struct {
	Total      int        `json:"total"`
	Left       int        `json:"left"`
	Question   string     `json:"question"`
	AnswerLen  int        `json:"answer_length"`
	RatingType RatingType `json:"rating_type"`

	session *leaf.ReviewSession
	rater   leaf.Rater
}

// NewSessionState constructs a new SessionState.
func NewSessionState(session *leaf.ReviewSession, rt RatingType) *SessionState {
	var rater leaf.Rater
	if rt == RatingTypeSelf {
		rater = leaf.TableRater()
	} else {
		rater = leaf.HarshRater()
	}

	s := &SessionState{
		Total:      session.Total(),
		Left:       session.Left(),
		Question:   session.Next(),
		AnswerLen:  len(session.CorrectAnswer()),
		RatingType: rt,
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
		s.session.Again()
	} else {
		s.session.Rate(rating)
		s.Left = s.session.Left()
	}

	s.Question = s.session.Next()
	s.AnswerLen = len(s.session.CorrectAnswer())
}
