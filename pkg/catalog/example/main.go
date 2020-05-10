package main

import (
	"crypto/md5"
	"encoding/hex"
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

	// Lets create a unique has for the note to store and lookup.
	hash := md5.Sum([]byte("testing blood did not work well at the end."))
	key := hex.EncodeToString(hash[:])

	// Append new card or skip already existing card.
	if card, found := cards.Find(key); found {
		fmt.Printf("Found card, letas update: %v - %v\n", card.Title, card.KeyHash)
	} else {
		fmt.Println("Inserting new card to catalog")
		cards = append(cards, &catalog.Card{Title: "Bad Blood", KeyHash: key})
	}

	// Writing Cards to local cache.
	err = librarian.Write(cards)
	if err != nil {
		log.Fatal(err)
	}
}
