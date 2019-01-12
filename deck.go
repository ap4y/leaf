package leaf

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Card []string

func (c Card) Answer() string {
	return strings.Join(c, " ")
}

type Stack map[string]Card

type Deck struct {
	Name  string
	Cards Stack
}

func OpenDeck(filename string) (*Deck, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("file: %s", err)
	}

	deck := &Deck{Cards: make(Stack)}
	scanner := bufio.NewScanner(f)
	var question string
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "* ") {
			deck.Name = strings.Replace(text, "* ", "", -1)
		} else if strings.HasPrefix(text, "** ") {
			question = strings.Replace(text, "** ", "", -1)
			deck.Cards[question] = make([]string, 0)
		} else {
			deck.Cards[question] = append(deck.Cards[question], text)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner: %s", err)
	}

	return deck, nil
}
