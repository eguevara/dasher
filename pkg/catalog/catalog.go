package catalog

import (
	"encoding/json"
	"os"
)

type CatalogOptions struct {
	Name string
}

// Librarian is a general interface used to manage catalog operations.
type Librarian interface {
	// Add will add a new annotation to the catalog.
	Write(Cards) error

	// GetAll will return all of the annotations available in the catalog.
	GetAll() (Cards, error)
}

type librarian struct {
	catalog CatalogOptions
}

var _ Librarian = librarian{}

func New(catalog CatalogOptions) Librarian {
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
func (l librarian) GetAll() (cards Cards, err error) {
	file, err := os.OpenFile(l.catalog.Name, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	if err := json.NewDecoder(file).Decode(&cards); err != nil {
		return nil, err
	}

	return cards, nil
}
