package common

import (
	"io/ioutil"
	"log"
	"net/http"

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
