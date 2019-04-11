package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/ap4y/leaf"
	"github.com/ap4y/leaf/ui"
	termbox "github.com/nsf/termbox-go"
)

var (
	db      = flag.String("db", "leaf.db", "database location")
	count   = flag.Int("count", 20, "cards to review")
	addr    = flag.String("addr", ":8080", "addr for Web UI")
	devMode = flag.Bool("dev", false, "use local dev assets")
)

func main() {
	flag.Parse()

	ui.DevMode = *devMode

	deckFile := flag.Arg(1)
	if deckFile == "" {
		log.Fatal("Missing deck file")
	}

	deck, err := leaf.OpenDeck(deckFile)
	if err != nil {
		log.Fatal("Failed to open deck: ", err)
	}

	db, err := leaf.OpenBoltStore(*db)
	if err != nil {
		log.Fatal("Failed to open stats DB: ", err)
	}

	defer db.Close()

	switch flag.Arg(0) {
	case "stats":
		w := tabwriter.NewWriter(os.Stdout, 5, 5, 5, ' ', 0)
		fmt.Fprintln(w, "Card\tDifficulty\tInterval")
		db.GetStats(deck.Name, func(card string, stats *leaf.Stats) {
			fmt.Fprintf(w, "%s\t%.1f\t%.1f\n", card, stats.Difficulty, stats.Interval)
		})
		w.Flush()
	case "web", "review":
		session, err := leaf.NewReviewSession(deck, db, *count)
		if err != nil {
			log.Fatal("Failed to create review session: ", err)
		}

		var u ui.UI
		if flag.Arg(0) == "web" {
			u = ui.NewWebUI(*addr)
		} else {
			if err := termbox.Init(); err != nil {
				log.Fatal("Failed to initialise tui: ", err)
			}
			defer termbox.Close()

			u = ui.NewTUI()
		}

		if err := u.Render(ui.NewSessionState(session)); err != nil {
			log.Fatal("Failed to render: ", err)
		}
	default:
		log.Fatal("unknown command")
	}
}
