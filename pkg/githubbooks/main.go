package githubbooks

import (
	"context"

	"github.com/google/go-github/github"

	"github.com/eguevara/dasher/common"
	"github.com/eguevara/dasher/config"
)

type BookService interface {
	// AddBook creates a new github issue on a repo masked as a book.
	AddBook(request *github.IssueRequest) (*github.Issue, error)

	List(request *github.IssueListByRepoOptions) (*BooksResponse, error)
}

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
func (b *githubbook) AddBook(opt *github.IssueRequest) (*github.Issue, error) {
	issue, _, err := b.Client.Issues.Create(
		context.Background(),
		b.Options.Config.GitHub.Owner,
		b.Options.Config.GitHub.Repo,
		opt,
	)

	if err != nil {
		return nil, err
	}
	return issue, nil
}

// List returns a list of books from github masked as issues.
func (b *githubbook) List(opt *github.IssueListByRepoOptions) (*BooksResponse, error) {
	ghIssues, _, err := b.Client.Issues.ListByRepo(
		context.Background(),
		b.Options.Config.GitHub.Owner,
		b.Options.Config.GitHub.Repo,
		opt,
	)

	if err != nil {
		return nil, err
	}

	issues := &BooksResponse{
		Issues: ghIssues,
	}

	return issues, nil
}
