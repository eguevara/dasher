package catalog

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

// Options holds the options for a Catalog.
type Options struct {
	Name string
}

// Librarian is a general interface used to manage catalog operations.
type Librarian interface {
	// Add will add a new annotation to the catalog.
	Write(Cards) error

	// ReadCards will return all of the annotations available in the catalog.
	ReadCards() (Cards, error)
}

type librarian struct {
	catalog Options
}

var _ Librarian = librarian{}

// New returns a new catalog.
func New(catalog Options) Librarian {
	return librarian{
		catalog: catalog,
	}
}

// Write will take a list of annotations and write to a file.
func (l librarian) Write(cards Cards) error {
	file, err := os.OpenFile(l.catalog.Name, os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	// Write the cards to the file in json.
	if err := json.NewEncoder(file).Encode(cards); err != nil {
		return err
	}

	return nil
}

// Get will return all cards in a catalog.
func (l librarian) ReadCards() (cards Cards, err error) {
	file, err := os.OpenFile(l.catalog.Name, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	if err := json.NewDecoder(file).Decode(&cards); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		return nil, err
	}

	return cards, nil
}
