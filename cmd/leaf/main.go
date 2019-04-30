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
	db    = flag.String("db", "leaf.db", "database location")
	count = flag.Int("count", 20, "cards to review")
)

func main() {
	flag.Parse()

	deckName := flag.Arg(1)
	if deckName == "" {
		log.Fatal("Missing deck name")
	}

	db, err := leaf.OpenBoltStore(*db)
	if err != nil {
		log.Fatal("Failed to open stats DB: ", err)
	}

	defer db.Close()

	dm, err := leaf.NewDeckManager("./", db)
	if err != nil {
		log.Fatal("Failed to initialise deck manager: ", err)
	}

	switch flag.Arg(0) {
	case "stats":
		stats, err := dm.DeckStats(deckName)
		if err != nil {
			log.Fatal("Failed to get card stats: ", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 5, 5, 5, ' ', 0)
		fmt.Fprintln(w, "Card\tDifficulty\tInterval")
		for _, s := range stats {
			fmt.Fprintf(w, "%s\t%.1f\t%.1f\n", s.Question, s.Difficulty, s.Interval)
		}
		w.Flush()
	case "review":
		session, err := dm.ReviewSession(deckName, *count)
		if err != nil {
			log.Fatal("Failed to create review session: ", err)
		}

		if err := termbox.Init(); err != nil {
			log.Fatal("Failed to initialise tui: ", err)
		}
		defer termbox.Close()

		u := ui.NewTUI(deckName)

		if err := u.Render(ui.NewSessionState(session)); err != nil {
			log.Fatal("Failed to render: ", err)
		}
	default:
		log.Fatal("unknown command")
	}
}
