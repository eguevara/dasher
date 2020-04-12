package common

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-github/github"

	"github.com/eguevara/dasher/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

// GetOAuthClient creates the oauth client to be used for googl analytics access.
func GetOAuthClient(cfg *config.OAuthConfig) *http.Client {

	data, err := ioutil.ReadFile(cfg.PrivateFilePath)
	if err != nil {
		log.Fatalf("could not read private key file (PEM): %v", err)
	}

	conf := &jwt.Config{
		Email:      cfg.ServiceEmail,
		PrivateKey: data,
		Scopes:     cfg.Scopes,
		TokenURL:   google.JWTTokenURL,
	}

	if cfg.ImpersonateEmail != nil {
		conf.Subject = *cfg.ImpersonateEmail
	}

	// Return an OAuth http client based on private key from service account.
	return conf.Client(oauth2.NoContext)

}

// GetGithubClient returns a github client using access token.
func GetGithubClient(cfg *config.GitHubConfig) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token})
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
