package api

import (
	"log"
	"net/http"
	"time"

	"github.com/eguevara/dasher/common"
	"github.com/eguevara/dasher/config"
	"github.com/eguevara/go-books"
)

const (
	serviceBooks = "books"
)

// BooksResponse stores the Books handler response.
type BooksResponse struct {
	Items []books.Volume `json:"items"`
}

type booksHandler struct {
	Client *books.Client
	Config *config.AppConfig
}

// BooksHandler handles http requests for the go-books api.
func BooksHandler(cfg *config.AppConfig) http.Handler {
	client := common.GetOAuthClient(cfg.BooksOAuth)
	booksClient, err := books.New(client)
	if err != nil {
		log.Fatal("error on creating Books client")
	}
	return &booksHandler{
		Config: cfg,
		Client: booksClient,
	}
}

func (b *booksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Call Prometheus Collector to instrument requests.
	reportRequestReceived(serviceBooks)
	response, err := b.GetBooks()
	if err != nil {
		reportRequestFailed(*r, err)
	}
	common.Respond(w, response, err)

	// Call Prometheus Collector to instrument service duration.
	reportServiceCompleted(serviceBooks, startTime)
}

// GetBooks returns a list of Google Books you've read.
func (b *booksHandler) GetBooks() (*BooksResponse, error) {
	opts := &books.VolumesListOptions{
		Fields:     *b.Config.BooksVolumesFields,
		MaxResults: *b.Config.BooksVolumesMax,
	}

	volumes, _, err := b.Client.Volumes.List(*b.Config.BooksShelf, opts)
	if err != nil {
		return nil, err
	}

	response := &BooksResponse{
		Items: volumes,
	}

	return response, nil
}
