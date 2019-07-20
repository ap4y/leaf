package leaf

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// Card represents a single card in a Deck. Each card may have
// multiple sides (answers).
type Card struct {
	Question string   `json:"card"`
	Sides    []string `json:"-"`
}

// Answer returns combined space separated answer for all sides of the card.
func (c *Card) Answer() string {
	return strings.Join(c.Sides, " ")
}

// Deck represents a named collection of the cards to review.
type Deck struct {
	Name  string
	Cards []Card

	modtime  time.Time
	filename string
}

// OpenDeck loads deck from an org file. File format is:
// * Deck Name
// ** Question
// side 1
// side 2
func OpenDeck(filename string) (*Deck, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("file: %s", err)
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("file: %s", err)
	}

	deck := &Deck{modtime: stat.ModTime(), filename: filename}
	if err := deck.load(f); err != nil {
		return nil, err
	}

	return deck, nil
}

// Reload compares ModTime on deck file and reloads cards if necessary.
func (deck *Deck) Reload() error {
	stat, err := os.Stat(deck.filename)
	if err != nil {
		return fmt.Errorf("file: %s", err)
	}

	if deck.modtime.UnixNano() >= stat.ModTime().UnixNano() {
		return nil
	}

	f, err := os.Open(deck.filename)
	if err != nil {
		return fmt.Errorf("file: %s", err)
	}

	if err := deck.load(f); err != nil {
		return err
	}

	deck.modtime = stat.ModTime()
	return nil
}

func (deck *Deck) load(f *os.File) error {
	cards := make(map[string][]string)
	scanner := bufio.NewScanner(f)
	var question string
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "* ") {
			deck.Name = strings.Replace(text, "* ", "", -1)
		} else if strings.HasPrefix(text, "** ") {
			question = strings.Replace(text, "** ", "", -1)
			cards[question] = make([]string, 0)
		} else {
			cards[question] = append(cards[question], text)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner: %s", err)
	}

	deck.Cards = make([]Card, 0)
	for question, sides := range cards {
		deck.Cards = append(deck.Cards, Card{question, sides})
	}

	return nil
}
