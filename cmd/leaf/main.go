package main

import (
	"fmt"
	"log"

	"github.com/ap4y/leaf"
)

func main() {
	deck, err := leaf.OpenDeck("hiragana.org")
	if err != nil {
		log.Fatal("Failed to open deck: ", err)
	}

	db, err := leaf.OpenStatsDB("leaf.db")
	if err != nil {
		log.Fatal("Failed to open stats DB: ", err)
	}

	defer db.Close()

	session, err := leaf.NewReviewSession(deck, db, 20)
	if err != nil {
		log.Fatal("Failed to create review session: ", err)
	}

	for {
		question := session.Next()
		if question == "" {
			break
		}

		fmt.Println(question)
		var answer string
		if _, err := fmt.Scanln("%s", &answer); err != nil {
			log.Fatal("Failed to read answer: ", err)
		}

		result, err := session.Answer(answer)
		if err != nil {
			log.Fatal("Failed to save answer: ", err)
		}

		if result {
			log.Println("correct!")
		} else {
			log.Println(deck.Cards[question].Answer())
		}
	}
}
