package leaf

import (
	"errors"
	"path/filepath"
	"sort"
	"time"
)

// ErrNotFound represents error returned for requests for non-existing deck.
var ErrNotFound = errors.New("deck not found")

// DeckStats stores overview stats for a Deck.
type DeckStats struct {
	Name         string    `json:"name"`
	CardsReady   int       `json:"cards_ready"`
	NextReviewAt time.Time `json:"next_review_at"`
}

// DeckManager manages set of decks.
type DeckManager struct {
	db    StatsStore
	decks []*Deck
}

// NewDeckManager constructs a new DeckManager by reading all decks
// from a given folder using provided store.
func NewDeckManager(path string, db StatsStore, outFormat OutputFormat) (*DeckManager, error) {
	files, err := filepath.Glob(path + "/*.org")
	if err != nil {
		return nil, err
	}

	decks := make([]*Deck, 0, len(files))
	for _, file := range files {
		deck, err := OpenDeck(file, outFormat)
		if err != nil {
			return nil, err
		}
		decks = append(decks, deck)
	}

	return &DeckManager{db, decks}, nil
}

// ReviewDecks returns stats for available decks.
func (dm DeckManager) ReviewDecks() ([]DeckStats, error) {
	result := make([]DeckStats, 0, len(dm.decks))
	for _, deck := range dm.decks {
		nextReviewAt, reviewDeck, err := dm.reviewDeck(deck, -1)
		if err != nil {
			return nil, err
		}

		result = append(result, DeckStats{deck.Name, len(reviewDeck), nextReviewAt})
	}

	return result, nil
}

// ReviewSession initiates a new ReviewSession for a given deck name.
func (dm DeckManager) ReviewSession(deckName string) (*ReviewSession, error) {
	var deck *Deck
	for _, d := range dm.decks {
		if d.Name == deckName {
			deck = d
			break
		}
	}

	if deck == nil {
		return nil, ErrNotFound
	}

	_, cards, err := dm.reviewDeck(deck, deck.PerReview)
	if err != nil {
		return nil, err
	}

	return NewReviewSession(cards, deck.RatingType, func(card *CardWithStats) error {
		return dm.db.SaveStats(deckName, card.Question, card.Stats)
	}), nil
}

// DeckStats returns card stats for a given deck name.
func (dm DeckManager) DeckStats(deckName string) ([]CardWithStats, error) {
	var deck *Deck
	for _, d := range dm.decks {
		if d.Name == deckName {
			deck = d
			break
		}
	}

	if deck == nil {
		return nil, ErrNotFound
	}

	return dm.deckStats(deck)
}

func (dm DeckManager) deckStats(deck *Deck) ([]CardWithStats, error) {
	stats := make(map[string]*Stats)
	err := dm.db.RangeStats(deck.Name, deck.Algorithm, func(card string, s *Stats) bool {
		stats[card] = s
		return true
	})
	if err != nil {
		return nil, err
	}

	result := make([]CardWithStats, 0, len(deck.Cards))
	for _, card := range deck.Cards {
		if stats[card.Question] != nil {
			result = append(result, CardWithStats{card, stats[card.Question]})
		} else {
			result = append(result, CardWithStats{card, NewStats(deck.Algorithm)})
		}
	}

	return result, nil
}

func (dm DeckManager) reviewDeck(deck *Deck, total int) (nextReviewAt time.Time, cards []CardWithStats, err error) {
	if fErr := deck.Reload(); fErr != nil {
		err = fErr
		return
	}

	stats, sErr := dm.deckStats(deck)
	if sErr != nil {
		err = sErr
		return
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[j].SRSAlgorithm.Less(stats[i].Stats.SRSAlgorithm)
	})

	if len(stats) > 0 {
		nextReviewAt = stats[0].NextReviewAt()
	}

	cards = make([]CardWithStats, 0, len(stats))
	for _, s := range stats {
		if total > 0 && len(cards) == total {
			break
		}

		if !s.IsReady() {
			continue
		}

		cards = append(cards, s)
	}

	return
}
