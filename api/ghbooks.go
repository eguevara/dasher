package api

import (
	"log"
	"net/http"
	"time"

	"github.com/eguevara/dasher/common"
	"github.com/eguevara/dasher/config"
	"github.com/eguevara/go-books"
	"github.com/gorilla/mux"
)

const (
	serviceBooks = "books"
)

// BooksResponse stores the Books handler response.
type BooksResponse struct {
	Items []books.Volume `json:"items"`
}

// NotesResponse stores the Annotation handler response.
type NotesResponse struct {
	Items []books.Annotation `json:"items"`
}

// BooksHandler stores data for bookHandler
type BooksHandler struct {
	Client *books.Client
	Config *config.AppConfig
}

// NewBooksHandler returns an instance of BookHandler.
func NewBooksHandler(cfg *config.AppConfig) *BooksHandler {
	client := common.GetOAuthClient(cfg.BooksOAuth)
	booksClient, err := books.New(client)
	if err != nil {
		log.Fatal("error on creating Books client")
	}
	return &BooksHandler{
		Config: cfg,
		Client: booksClient,
	}
}

// ListNotesByVolumeID will list all of the Annotations for a volumeId.
func (b *BooksHandler) ListNotesByVolumeID(w http.ResponseWriter, r *http.Request) {
	defer func(begin time.Time) {
		requestCount(serviceBooks)
		requestLatency(serviceBooks, begin)
	}(time.Now())

	vars := mux.Vars(r)
	volumeID := vars["id"]

	response, err := b.getNotesByVolumeID(volumeID)
	common.Respond(w, response, err)
}

// List will call books.mylibrary.bookshelves.volumes.list API.
func (b *BooksHandler) List(w http.ResponseWriter, r *http.Request) {
	defer func(begin time.Time) {
		requestCount(serviceBooks)
		requestLatency(serviceBooks, begin)
	}(time.Now())

	response, err := b.getBooks()
	if err != nil {
		reportRequestFailed(*r, err)
	}

	common.Respond(w, response, err)
}

// getNotesByVolumeID will return a list of Notes for a volumeId
func (b *BooksHandler) getNotesByVolumeID(volumeID string) (*NotesResponse, error) {
	opts := &books.AnnotationsListOptions{
		VolumeID:       volumeID,
		ContentVersion: "full-1.0.0",
		MaxResults:     40,
		Fields:         "items(selectedText,volumeId),totalItems",
	}

	notes, _, err := b.Client.Annotations.List(opts)
	if err != nil {
		return nil, err
	}

	response := &NotesResponse{
		Items: notes,
	}

	return response, nil
}

// getBooks returns a list of Google Books you've read.
func (b *BooksHandler) getBooks() (*BooksResponse, error) {
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
