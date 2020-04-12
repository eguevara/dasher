package github_books

import (
	"github.com/eguevara/dasher/common"
	"github.com/eguevara/dasher/config"
	"log"
	"github.com/eguevara/go-books"
)

type GitHubBookService interface {
	AddBook()(error)
}

// Request stores options needed to make feelinglucky requests.
type Request struct {
	// Config holds the oauth credential used by the api.
	Config *config.AppConfig

	// Shelf is the bookshelf to look for volumes/annotations.
	Shelf string
}

type githubbook struct {
	// Options holds all stateful fields added to the request.
	Options *Request

	// Client holds the book service client to communicate to the books api.
	Client *books.Client

	Title string
}

// NewService returns a new instance of the lucky type.
func NewService(opts *Request) GitHubBookService {
	client := common.GetOAuthClient(opts.Config.BooksOAuth)
	booksClient, err := books.New(client)
	if err != nil {
		log.Fatal("error on creating Books client")
	}
	return &githubbook{
		Options: opts,
		Client:  booksClient,
	}
}

// Validates that we implement the LuckyService interface.
var _ GitHubBookService = &githubbook{}


func (b *githubbook) AddBook()(error) {
	log.Fatalf("cound not add book")
}

