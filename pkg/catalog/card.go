package catalog

// Card holds information stored in a catalog.
type Card struct {
	// Title is the title of the book.
	Title string

	// KeyHash is the unique hash of the annotation.
	KeyHash string
}

// Cards are a list of Card.
type Cards []*Card

// Find will look through a list of cards.
func (c Cards) Find(key string) (*Card, bool) {
	for _, card := range c {
		if card.KeyHash == key {
			return card, true
		}
	}

	return nil, false
}
