package catalog

type KeyHash string

type Card struct {
	// Title is the title of the book.
	Title string

	// KeyHash is the unique hash of the annotation.
	KeyHash KeyHash
}

type Cards []*Card

func (c Cards) Find(title string) (*Card, bool) {
	for _, card := range c {
		if card.Title == title {
			return card, true
		}
	}

	return nil, false
}
