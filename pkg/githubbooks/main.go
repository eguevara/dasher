package githubbooks

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"

	"github.com/eguevara/dasher/common"
	"github.com/eguevara/dasher/config"
)

// BookService is an interface to represents a book model.
type BookService interface {
	// AddBook creates a new book.
	AddBook(request *github.IssueRequest) (*github.Issue, error)

	// BookExists returns the book found if it exists.
	BookExists(string) (*github.Issue, bool, error)

	// List will list out all books found.
	List() (*BooksResponse, error)

	// AddNote adds a book annotation.
	AddNote(int, string) (*github.IssueComment, error)
}

// BooksResponse represents what the book service returns.
type BooksResponse struct {
	Issues []*github.Issue
}

// Request stores options needed to make books requests.
type Request struct {
	// Config holds the oauth credential used by the api.
	Config *config.AppConfig
}

type githubbook struct {
	// Options holds all stateful fields added to the request.
	Options *Request

	// Client holds the book service client to communicate to the github api.
	Client *github.Client

	Title string
}

// NewService returns a new instance of the lucky type.
func NewService(opts *Request) BookService {
	client := common.GetGithubClient(opts.Config.GitHub)
	return &githubbook{
		Options: opts,
		Client:  client,
	}
}

// Validates that we implement the LuckyService interface.
var _ BookService = &githubbook{}

// AddBook will add a new issue to an existing GitHub repo to mask as a Book.
func (b *githubbook) AddNote(issue int, body string) (*github.IssueComment, error) {
	issueComment := &github.IssueComment{Body: &body}
	comment, _, err := b.Client.Issues.CreateComment(
		context.Background(),
		b.Options.Config.GitHub.Owner,
		b.Options.Config.GitHub.Repo,
		issue,
		issueComment)

	if err != nil {
		return nil, err
	}
	return comment, nil
}

// AddBook will add a new issue to an existing GitHub repo to mask as a Book.
func (b *githubbook) AddBook(opt *github.IssueRequest) (*github.Issue, error) {
	issue, _, err := b.Client.Issues.Create(
		context.Background(),
		b.Options.Config.GitHub.Owner,
		b.Options.Config.GitHub.Repo, opt)

	if err != nil {
		return nil, err
	}
	return issue, nil
}

// List returns a list of books from github masked as issues.
func (b *githubbook) List() (*BooksResponse, error) {
	ghIssues, _, err := b.Client.Issues.ListByRepo(
		context.Background(),
		b.Options.Config.GitHub.Owner,
		b.Options.Config.GitHub.Repo,
		nil)
	if err != nil {
		return nil, err
	}

	issues := &BooksResponse{
		Issues: ghIssues,
	}

	return issues, nil
}

// bookExist checks whether the book already exists in the repo as an issue.
// Uses the search api to do a search on the :title and is:issue qualifiers.
//
// Returns The github.Issue and boolean.
func (b *githubbook) BookExists(bookTitle string) (*github.Issue, bool, error) {

	query := fmt.Sprintf("repo:%v/%v is:issue \"%v\" in:title", b.Options.Config.GitHub.Owner, b.Options.Config.GitHub.Repo, bookTitle)
	result, _, err := b.Client.Search.Issues(context.Background(), query, nil)

	if err != nil {
		return nil, true, err
	}

	if len(result.Issues) > 0 {
		return result.Issues[0], true, nil
	}

	return nil, false, nil
}
