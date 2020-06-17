package leaf

import (
	"errors"
	"time"
)

// StatsSaveFunc persists stats updates.
type StatsSaveFunc func(card *CardWithStats) error

// ReviewSession contains parameters for a Deck review sessions.
type ReviewSession struct {
	statsSaver StatsSaveFunc
	cards      []CardWithStats
	queue      []string
	sides      []string
	startedAt  time.Time
	ratingType RatingType
}

// NewReviewSession constructs a new ReviewSession for a given set of cards.
// Rating calculation will be performed using provided rater.
// Provided StatsSaveFunc will be used for stats updates post review.
func NewReviewSession(cards []CardWithStats, sides []string, rt RatingType, statsSaver StatsSaveFunc) *ReviewSession {
	queue := make([]string, len(cards))
	for idx, card := range cards {
		queue[idx] = card.Question
	}

	return &ReviewSession{statsSaver, cards, queue, sides, time.Now(), rt}
}

// Sides returns side names from the reviewed deck.
func (s *ReviewSession) Sides() []string {
	return s.sides
}

// StartedAt returns start time of the review session.
func (s *ReviewSession) StartedAt() time.Time {
	return s.startedAt
}

// RatingType returns type of rating to be used for the review session.
func (s *ReviewSession) RatingType() RatingType {
	return s.ratingType
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

// Again re-queues current card back for review.
func (s *ReviewSession) Again() error {
	card := s.currentCard()
	if card == nil {
		return errors.New("no cards in queue")
	}

	s.queue = s.queue[1:]
	s.queue = append(s.queue, card.Question)
	return nil
}

// Rate assign rating to a current card and removes it from the queue if rating > 0.
func (s *ReviewSession) Rate(rating float64) error {
	card := s.currentCard()
	if card == nil {
		return errors.New("no cards in queue")
	}

	s.queue = s.queue[1:]
	card.Advance(rating)
	return s.statsSaver(card)
}

func (s *ReviewSession) currentCard() *CardWithStats {
	question := s.Next()
	for _, c := range s.cards {
		if c.Question == question {
			return &c
		}
	}

	return nil
}
