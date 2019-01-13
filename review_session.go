package leaf

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"time"
)

type ReviewSession struct {
	deck      *Deck
	stats     map[string]*Stats
	queue     []string
	mistakes  map[string]int
	db        *StatsDB
	total     int
	startedAt time.Time
}

type reviewedCard struct {
	card           string
	percentOverdue float64
}

func NewReviewSession(deck *Deck, db *StatsDB, total int) (*ReviewSession, error) {
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

func (s *ReviewSession) DeckName() string {
	return s.deck.Name
}

func (s *ReviewSession) StartedAt() time.Time {
	return s.startedAt
}

func (s *ReviewSession) Total() int {
	return s.total
}

func (s *ReviewSession) Left() int {
	return len(s.queue)
}

func (s *ReviewSession) Next() string {
	if len(s.queue) == 0 {
		return ""
	}

	return s.queue[0]
}

func (s *ReviewSession) CorrectAnswer() string {
	question := s.Next()
	card := s.deck.Cards[question]
	if card == nil {
		return ""
	}

	return card.Answer()
}

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
	} else {
		s.mistakes[question]++
		s.queue = append(s.queue, question)
		return false, nil
	}
}

func (s *ReviewSession) rating(question string) float64 {
	miss := float64(s.mistakes[question])
	return math.Max(0, 5*(1-miss/5))
}
