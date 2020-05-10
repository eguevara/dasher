package main

import (
	"fmt"
	"log"

	"github.com/eguevara/dasher/pkg/catalog"
)

func main() {

	opts := catalog.Options{
		Name: ".catalog",
	}

	librarian := catalog.New(opts)

	//notes := catalog.Cards{
	//	{
	//		Title:        "Dark Matter",
	//		KeyHash:      "1234",
	//	},
	//	{
	//		Title:        "Origin",
	//		KeyHash:      "7890",
	//	},
	//}

	// Reading local cards into Cards.
	cards, err := librarian.ReadCards()
	if err != nil {
		log.Fatal(err)
	}

	// Append new card or skip already existing card.
	if card, found := cards.Find("Bad Blood"); found {
		fmt.Printf("Found card, letas update: %v\n", card.Title)
	} else {
		fmt.Println("Inserting new card to catalog")
		cards = append(cards, &catalog.Card{Title: "Bad Blood"})
	}

	// Writing Cards to local cache.
	err = librarian.Write(cards)
	if err != nil {
		log.Fatal(err)
	}
}
