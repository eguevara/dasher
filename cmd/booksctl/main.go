package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"

	"github.com/eguevara/dasher/pkg/githubbooks"

	"github.com/eguevara/dasher/api"
	"github.com/eguevara/dasher/config"
	"github.com/eguevara/dasher/pkg/feelinglucky"
	books "github.com/eguevara/go-books"
)

const (
	defaultDotFolder  = ".booksctl"
	defaultConfigFile = "config.json"
	defaultPEMFile    = "booksctl.pem"
)

var defaultHomeDirectory = os.Getenv("HOME")

func main() {

	var (
		flagConfigPath      = flag.String("config-file", defaultConfigFile, "application configuration file")
		flagBookID          = flag.String("book-id", "", "the id of the book")
		flagShowBooks       = flag.Bool("show-books", false, "show all books for a shelf")
		flagGitHubShowBooks = flag.Bool("show-ghbooks", false, "show all github books")
		flagSyncGitHubBook  = flag.Bool("sync-books", false, "Syncs google books to github issues")
		flagSyncGitHubNotes = flag.String("sync-notes", "", "Syncs google book notes to github comments")
		flagShowBookShelf   = flag.Bool("show-bookshelf", false, "show all bookshelfs")
		flagLucky           = flag.Bool("feeling-lucky", false, "feeling lucky")
	)

	flag.Parse()

	cfg := buildConfigFromFie(flagConfigPath)

	if *flagGitHubShowBooks {
		opts := &githubbooks.Request{
			Config: cfg,
		}

		githubService := githubbooks.NewService(opts)

		resp, err := githubService.List()
		if err != nil {
			log.Fatalf("error in getting github books: %v", err)
		}

		for _, b := range resp.Issues {
			fmt.Printf("Adding book: %v\n", *b.Title)
		}

	}

	if *flagSyncGitHubBook {
		opts := &githubbooks.Request{
			Config: cfg,
		}
		githubService := githubbooks.NewService(opts)
		googleBookService := api.NewBooksHandler(cfg)

		volumeOpts := &books.VolumesListOptions{
			MaxResults: 100,
		}
		books, _, err := googleBookService.Client.Volumes.List("1", volumeOpts)
		if err != nil {
			fmt.Printf("could not find any books in this shelf: %v", err)
		}

		for _, book := range books {
			_, found, err := githubService.BookExists(*book.Info.Title)
			if err != nil || found {
				log.Printf("Skipping book: %v", err)
				continue
			}

			bookOpt := &github.IssueRequest{
				Title: book.Info.Title,

				//TODO: make this configurable.
				Labels: &[]string{"google-books"},
			}

			ghIssue, err := githubService.AddBook(bookOpt)

			if err != nil {
				log.Fatalf("error in List(): %v", err)
			}
			fmt.Printf("Adding book: %v, %v - %v\n", *book.Info.Title, bookOpt, *ghIssue.Title)

		}
		return
	}
	if *flagSyncGitHubNotes != "" {
		opts := &githubbooks.Request{
			Config: cfg,
		}
		githubService := githubbooks.NewService(opts)
		googleBookService := api.NewBooksHandler(cfg)

		volumeOpts := &books.VolumesListOptions{
			MaxResults: 100,
		}
		resp, _, err := googleBookService.Client.Volumes.List("1", volumeOpts)
		if err != nil {
			fmt.Printf("could not find any books in this shelf: %v", err)
		}

		for _, book := range resp {
			if *flagSyncGitHubNotes != *book.ID {
				continue
			}
			githubBook, notFound, err := githubService.BookExists(*book.Info.Title)
			if err != nil || !notFound {
				log.Printf("Skipping: %v", err)
				continue
			}

			opts := &books.AnnotationsListOptions{
				VolumeID:       *book.ID,
				ContentVersion: "full-1.0.0",
				MaxResults:     40,
				Fields:         "items(layerId,selectedText,volumeId,pageIds),totalItems,nextPageToken",
			}

			totalNotes := 0
			for {
				list, resp, err := googleBookService.Client.Annotations.List(opts)
				if err != nil {
					log.Fatalf("error in list(): %v ", err)
				}

				for _, note := range list {
					if *note.SelectedText == "" {
						totalNotes--
						continue
					}

					ghComment := buildComment(note)
					_, err := githubService.AddNote(*githubBook.Number, ghComment)
					if err != nil {
						log.Fatalf("error: %v", err)
					}

				}

				totalNotes += len(list)

				// If there is no more nextPage tokens, then you are done!
				if resp.NextPageToken == "" {
					break
				}

				// Update the token to the next page to fetch the rest.
				opts.PageToken = resp.NextPageToken
			}

			//fmt.Printf(resp.s)
			fmt.Printf("Count: %v\n\n", totalNotes)
		}
		return
	}

	if *flagLucky == true {
		opts := &feelinglucky.Request{
			Config: cfg,
			Shelf:  "1",
		}

		svc := feelinglucky.NewService(opts)

		quote, err := svc.Annotation()
		if err != nil {
			log.Fatalf("error in getting annotation: %v", err)
		}
		fmt.Printf("\"%v\n\n\n%v\n", *quote.SelectedText, quote.Title)
	}

	if *flagBookID != "" || *flagShowBooks || *flagShowBookShelf {
		svc := api.NewBooksHandler(cfg)

		if *flagShowBookShelf == true {
			opts := &books.ShelvesListOptions{}

			shelves, _, err := svc.Client.Shelves.List(opts)
			if err != nil {
				log.Fatalf("error in List(): %v", err)
			}

			for _, v := range shelves {
				fmt.Printf("Id: %d, Title: %v \n", *v.ID, *v.Title)
			}
			return
		}

		if *flagShowBooks {
			opts := &books.VolumesListOptions{
				MaxResults: 100,
			}

			books, _, err := svc.Client.Volumes.List("1", opts)
			if err != nil {
				fmt.Printf("could not find any books in this shelf: %v", err)
			}

			for _, book := range books {
				fmt.Printf("Book: %v, ID: %v\n\n", *book.Info.Title, *book.ID)
			}

			return
		}
		opts := &books.AnnotationsListOptions{
			VolumeID:       *flagBookID,
			ContentVersion: "full-1.0.0",
			MaxResults:     40,
			Fields:         "items(layerId,selectedText,volumeId),totalItems",
		}

		list, _, err := svc.Client.Annotations.List(opts)
		if err != nil {
			log.Fatalf("error in list(): %v ", err)
		}

		fmt.Printf("Count: %v\n\n", len(list))

		for _, quote := range list {
			fmt.Printf("...%v...\n\n", *quote.SelectedText)
		}

	}

}

// buildComment takes a note and builds a comment used as an annotation.
// Returns a String
func buildComment(note books.Annotation) string {
	pages := strings.Join(note.PageIds, ", ")
	pages = strings.Replace(pages, "PA", "page: ", -1)
	pages = strings.Replace(pages, "PT", "page: ", -1)
	pages = strings.Replace(pages, "PR", "page: ", -1)

	comment := fmt.Sprintf("%v\n%v", pages, *note.SelectedText)

	return comment
}

func buildConfigFromFie(file *string) *config.AppConfig {
	configFile := filepath.Join(defaultHomeDirectory, defaultDotFolder, *file)
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("DASHER error reading config.json: %v", err)
	}

	config := &config.AppConfig{}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("DASHER setting up app configuration: %v", err)
	}

	if config.BooksOAuth.PrivateFilePath == "" {
		config.BooksOAuth.PrivateFilePath = filepath.Join(defaultHomeDirectory, defaultDotFolder, defaultPEMFile)
	}

	return config
}
