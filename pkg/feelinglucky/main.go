package feelinglucky

import (
	"fmt"
	"log"

	"github.com/eguevara/dasher/common"
	"github.com/eguevara/dasher/config"
	books "github.com/eguevara/go-books"
	"k8s.io/apimachinery/pkg/util/rand"
)

// LuckyService defines the behavior required by types that implement a new
// feeling lucky type.
type LuckyService interface {
	Annotation() (*LuckyAnnotation, error)
}

// Request stores options needed to make feelinglucky requests.
type Request struct {
	// Config holds the oauth credential used by the api.
	Config *config.AppConfig

	// Shelf is the bookshelf to look for volumes/annotations.
	Shelf string
}

type lucky struct {
	// Options holds all stateful fields added to the request.
	Options *Request

	// Client holds the book service client to communicate to the books api.
	Client *books.Client

	Title string
}

// LuckyAnnotation holds annotation information to send as a response.
type LuckyAnnotation struct {
	// A wrapper field to hold Annotation fields.
	*books.Annotation

	// The title of the book.
	Title string
}

// NewService returns a new instance of the lucky type.
func NewService(opts *Request) LuckyService {
	client := common.GetOAuthClient(opts.Config.BooksOAuth)
	booksClient, err := books.New(client)
	if err != nil {
		log.Fatal("error on creating Books client")
	}
	return &lucky{
		Options: opts,
		Client:  booksClient,
	}
}

// Validates that we implement the LuckyService interface.
var _ LuckyService = &lucky{}

// Annotation returns a random luckyannotation of any of the volumes available.
func (l *lucky) Annotation() (*LuckyAnnotation, error) {
	book, err := l.book()
	if err != nil {
		log.Fatalf("could not find volume: %v", err)
	}

	opts := &books.AnnotationsListOptions{
		VolumeID:       *book.ID,
		ContentVersion: *book.Info.ContentVersion,
		MaxResults:     40,
		Fields:         "items(layerId,selectedText,volumeId),totalItems",
	}

	list, _, err := l.Client.Annotations.List(opts)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, fmt.Errorf("no annotations found: %v", l.Title)
	}

	randomIndex := rand.Intn(len(list))

	annotation := &LuckyAnnotation{
		Title:      l.Title,
		Annotation: &list[randomIndex],
	}
	return annotation, nil
}

// Book returns a random volume from the bookshelves avaialble.
func (l *lucky) book() (*books.Volume, error) {
	opts := &books.VolumesListOptions{
		Fields:     "items(id,volumeInfo(contentVersion,title)),totalItems",
		MaxResults: 40,
	}

	volumes, _, err := l.Client.Volumes.List(l.Options.Shelf, opts)
	if err != nil {
		return nil, err
	}

	randomIndex := rand.Intn(len(volumes))

	l.Title = *volumes[randomIndex].Info.Title

	return &volumes[randomIndex], nil
}
