package leaf

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/niklasfasching/go-org/org"
)

// OutputFormat defines output type produces during org file parsing.
type OutputFormat int

const (
	// OutputFormatOrg defines pretty printed org output.
	OutputFormatOrg OutputFormat = iota
	// OutputFormatHTML defines html output.
	OutputFormatHTML
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

	format   OutputFormat
	modtime  time.Time
	filename string
}

// OpenDeck loads deck from an org file. File format is:
// * Deck Name
// ** Question
// side 1
// side 2
func OpenDeck(filename string, format OutputFormat) (*Deck, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("file: %s", err)
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("file: %s", err)
	}

	deck := &Deck{modtime: stat.ModTime(), filename: filename, format: format}
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
	doc := org.New().Parse(f, "./")
	if len(doc.Nodes) == 0 {
		return fmt.Errorf("empty or invalid org-file")
	}

	root, ok := doc.Nodes[0].(org.Headline)
	if !ok {
		return fmt.Errorf("org-file doesn't start with a headline")
	}
	deck.Name = org.String(root.Title)
	deck.Cards = make([]Card, 0)

	for _, node := range root.Children {
		headline, ok := node.(org.Headline)
		if !ok || len(headline.Children) == 0 {
			continue
		}

		var w org.Writer
		if deck.format == OutputFormatHTML {
			w = org.NewHTMLWriter()
		} else {
			w = org.NewOrgWriter()
		}

		org.WriteNodes(w, headline.Title...)

		var answer string
		if block, ok := headline.Children[0].(org.Block); ok && block.Name == "SRC" {
			org.WriteNodes(w, block)
			answer = strings.TrimSpace(org.String(headline.Children[1:]))
		} else {
			answer = strings.TrimSpace(org.String(headline.Children))
		}

		deck.Cards = append(deck.Cards, Card{w.String(), []string{answer}})
	}

	return nil
}
