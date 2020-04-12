package api

import (
	"github.com/eguevara/dasher/common"

	"github.com/eguevara/dasher/config"
	"github.com/google/go-github/github"
)

const (
	serviceGitHubBooks = "books"
)

// GitHubBooksHandler stores data for bookHandler
type GitHubBooksHandler struct {
	Client *github.Client
	Config *config.AppConfig
}

// NewGitHubBooksHandler returns an instance of BookHandler.
func NewGitHubBooksHandler(cfg *config.AppConfig) *GitHubBooksHandler {
	client := common.GetGithubClient(cfg.GitHub)

	return &GitHubBooksHandler{
		Config: cfg,
		Client: client,
	}
}
