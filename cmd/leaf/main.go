package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ap4y/leaf"
	termbox "github.com/nsf/termbox-go"
)

var (
	db    = flag.String("db", "leaf.db", "database location")
	count = flag.Int("count", 20, "cards to review")
	stats = flag.Bool("stats", false, "show deck stats")
)

func main() {
	flag.Parse()

	deckFile := flag.Arg(0)
	if deckFile == "" {
		log.Fatal("Missing deck file")
	}

	deck, err := leaf.OpenDeck(deckFile)
	if err != nil {
		log.Fatal("Failed to open deck: ", err)
	}

	db, err := leaf.OpenStatsDB(*db)
	if err != nil {
		log.Fatal("Failed to open stats DB: ", err)
	}

	defer db.Close()

	if *stats {
		fmt.Printf("Card\tDifficulty\tInterval\n")
		db.GetStats(deck.Name, func(card string, stats *leaf.Stats) {
			fmt.Printf("%s\t%.1f\t\t%.1f\n", card, stats.Difficulty, stats.Interval)
		})
		return
	}

	session, err := leaf.NewReviewSession(deck, db, *count)
	if err != nil {
		log.Fatal("Failed to create review session: ", err)
	}

	if err := termbox.Init(); err != nil {
		log.Fatal("Failed to initialise tui: ", err)
	}
	defer termbox.Close()

	ui := &UI{
		UserInput: make([]rune, 0),
		DeckName:  session.DeckName(),
		Total:     session.Total(),
		Left:      session.Left(),
		Question:  session.Next(),
		Answer:    session.CorrectAnswer(),
	}

	if ui.Total == 0 {
		ui.State = StateFinished
	}

	ui.Render()

	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				return
			}

			if ui.State == StateFinished {
				break
			}

			if ui.State == StateResult {
				advanceState(ui, session)
				break
			}

			if ev.Key == termbox.KeyEnter {
				resolveAnswer(ui, session)
			} else if ev.Key == termbox.KeyBackspace || ev.Key == termbox.KeyBackspace2 {
				if len(ui.UserInput) > 0 {
					ui.UserInput = ui.UserInput[:len(ui.UserInput)-1]
				}
			} else {
				var ch rune
				if ev.Key == termbox.KeySpace {
					ch = ' '
				} else {
					ch = ev.Ch
				}

				ui.UserInput = append(ui.UserInput, ch)
			}
		case termbox.EventError:
			panic(ev.Err)
		}

		ui.Render()
	}
}

func advanceState(ui *UI, session *leaf.ReviewSession) {
	if session.Left() == 0 {
		ui.State = StateFinished
		return
	} else {
		ui.State = StateAnswering
		ui.Question = session.Next()
		ui.Answer = session.CorrectAnswer()
		return
	}
}

func resolveAnswer(ui *UI, session *leaf.ReviewSession) {
	ui.Result, _ = session.Answer(string(ui.UserInput))
	ui.State = StateResult
	ui.Left = session.Left()
	ui.UserInput = make([]rune, 0)
}
