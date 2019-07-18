package leaf

import (
	"errors"
	"time"
)

// StatsSaveFunc persists stats updates.
type StatsSaveFunc func(question string, stats *Stats) error

// ReviewSession contains parameters for a Deck review sessions.
type ReviewSession struct {
	rater      Rater
	statsSaver StatsSaveFunc
	cards      []*CardWithStats
	queue      []string
	mistakes   map[string]int
	startedAt  time.Time
}

// NewReviewSession constructs a new ReviewSession for a given set of cards.
// Rating calculation will be performed using provided rater.
// Provided StatsSaveFunc will be used for stats updates post review.
func NewReviewSession(cards []*CardWithStats, rater Rater, statsSaver StatsSaveFunc) *ReviewSession {
	queue := []string{}
	mistakes := make(map[string]int)
	for _, card := range cards {
		queue = append(queue, card.Question)
		mistakes[card.Question] = 0
	}

	return &ReviewSession{rater, statsSaver, cards, queue, mistakes, time.Now()}
}

// StartedAt returns start time of the review session.
func (s *ReviewSession) StartedAt() time.Time {
	return s.startedAt
}

// Total returns amount of cards in the session.
func (s *ReviewSession) Total() int {
	return len(s.cards)
}

// Left returns amount of cards left to review.
func (s *ReviewSession) Left() int {
	return len(s.queue)
}

// Next returns current card to review. Same card will be return until
// review is attempted via Answer call.
func (s *ReviewSession) Next() string {
	if len(s.queue) == 0 {
		return ""
	}

	return s.queue[0]
}

// CorrectAnswer returns correct answer for a current reviewed card.
func (s *ReviewSession) CorrectAnswer() string {
	card := s.currentCard()
	if card == nil {
		return ""
	}

	return card.Answer()
}

// Answer matches provided answer against correct and advances session.
func (s *ReviewSession) Answer(answer string) (bool, error) {
	question := s.Next()
	card := s.currentCard()
	if card == nil {
		return false, errors.New("unknown card")
	}

	s.queue = s.queue[1:]
	if answer == card.Answer() {
		rating := s.rater.Rate(s.mistakes[question])
		card.Advance(rating)
		if err := s.statsSaver(question, card.Stats); err != nil {
			return false, err
		}

		return true, nil
	}

	s.mistakes[question]++
	s.queue = append(s.queue, question)
	return false, nil
}

func (s *ReviewSession) currentCard() *CardWithStats {
	question := s.Next()
	var card *CardWithStats
	for _, c := range s.cards {
		if c.Question == question {
			card = c
			break
		}
	}

	return card
}
