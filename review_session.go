package leaf

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"time"
)

// ReviewSession contains parameters for a Deck review sessions.
type ReviewSession struct {
	deck      *Deck
	stats     map[string]*Stats
	queue     []string
	mistakes  map[string]int
	db        StatsStore
	total     int
	startedAt time.Time
}

type reviewedCard struct {
	card           string
	percentOverdue float64
}

// NewReviewSession constructs a new ReviewSession for a given deck
// with a total amount of cards. Provided DB will be used for review stats.
func NewReviewSession(deck *Deck, db StatsStore, total int) (*ReviewSession, error) {
	rCards := []*reviewedCard{}
	s := make(map[string]*Stats)
	err := db.GetStats(deck.Name, func(card string, stats *Stats) {
		rCards = append(rCards, &reviewedCard{card, stats.PercentOverdue()})
		s[card] = stats
	})
	if err != nil {
		return nil, fmt.Errorf("db: %s", err)
	}

	for card := range deck.Cards {
		if s[card] != nil {
			continue
		}

		stats := DefaultStats()
		rCards = append(rCards, &reviewedCard{card, stats.PercentOverdue()})
		s[card] = stats
	}

	sort.Slice(rCards, func(i, j int) bool {
		return rCards[i].percentOverdue > rCards[j].percentOverdue
	})

	stack := make(Stack)
	queue := []string{}
	mistakes := make(map[string]int)
	for _, rCard := range rCards {
		if len(queue) == total {
			break
		}

		stats := s[rCard.card]
		if !stats.IsReady() {
			continue
		}

		stack[rCard.card] = deck.Cards[rCard.card]
		queue = append(queue, rCard.card)
		mistakes[rCard.card] = 0
	}

	return &ReviewSession{deck, s, queue, mistakes, db, len(queue), time.Now()}, nil
}

// DeckName returns a name of the reviewed deck.
func (s *ReviewSession) DeckName() string {
	return s.deck.Name
}

// StartedAt returns start time of the review session.
func (s *ReviewSession) StartedAt() time.Time {
	return s.startedAt
}

// Total returns amount of cards in the session.
func (s *ReviewSession) Total() int {
	return s.total
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
	question := s.Next()
	card := s.deck.Cards[question]
	if card == nil {
		return ""
	}

	return card.Answer()
}

// Answer matches provided answer against correct and advances session.
func (s *ReviewSession) Answer(answer string) (bool, error) {
	question := s.Next()
	card := s.deck.Cards[question]
	if card == nil {
		return false, errors.New("unknown card")
	}

	s.queue = s.queue[1:]
	if answer == card.Answer() {
		stats := s.stats[question]
		rating := s.rating(question)
		stats.Record(rating)
		if err := s.db.SaveStats(s.deck.Name, question, stats); err != nil {
			return false, fmt.Errorf("db: %s", err)
		}

		return true, nil
	}

	s.mistakes[question]++
	s.queue = append(s.queue, question)
	return false, nil
}

func (s *ReviewSession) rating(question string) float64 {
	miss := float64(s.mistakes[question])
	return math.Max(0, 1-miss/5)
}
