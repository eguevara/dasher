package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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
		flagConfigPath = flag.String("config-file", defaultConfigFile, "application configuration file")
		flagBookID     = flag.String("book-id", "", "the id of the book")
		flagShowBooks  = flag.Bool("show-books", false, "show all books for a shelf")
		flagLucky      = flag.Bool("feeling-lucky", false, "feeling lucky")
	)

	flag.Parse()

	cfg := buildConfigFromFie(flagConfigPath)

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

	if *flagBookID != "" || *flagShowBooks {
		svc := api.NewBooksHandler(cfg)

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