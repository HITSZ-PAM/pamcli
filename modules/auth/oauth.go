package auth

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
	cc "golang.org/x/oauth2/clientcredentials"
)

type Auth interface {
	// Token is used to obtain auth token
	Token() (string, error)
}

// auth is the object used to get token
// See Token method bellow
type auth struct {
	// for authentication
	Client *http.Client // Used to interact with the center, implicit token authentication and update logic
}

// Config is the configuration struct of auth object
type Config struct {
	// ClientID is the application's ID.
	ClientID string

	// ClientSecret is the application's secret.
	ClientSecret string

	// TokenURL is the resource server's token endpoint
	// URL. This is a constant specific to each server.
	TokenURL string

	JwkURL string
}

// Token get current token. if token is expired, we will refresh it and return new one.
func (a *auth) Token() (string, error) {
	trans, ok := a.Client.Transport.(*oauth2.Transport)
	if !ok {
		return "", errors.New("wrong type of client transport of auth")
	}

	t, err := trans.Source.Token()
	if err != nil {
		return "", err
	}

	return t.AccessToken, nil
}

// NewAuth generate a auth object
func NewAuth(ctx context.Context, c Config) (Auth, error) {
	//The new client is used to ignore illegal server certificates
	realClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	ctxWithHTTPClient := context.WithValue(ctx, oauth2.HTTPClient, realClient)

	config := &cc.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		TokenURL:     c.TokenURL,
	}

	client := config.Client(ctxWithHTTPClient)

	return &auth{
		Client: client,
	}, nil
}
