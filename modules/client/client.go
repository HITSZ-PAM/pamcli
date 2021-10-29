package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/HITSZ-PAM/pamcli/models"
	"github.com/HITSZ-PAM/pamcli/modules/auth"
	"github.com/go-resty/resty/v2"
)

type Client interface {
	Resolve(accountID string) (string, string, error)
}

type client struct {
	client      *resty.Client
	config      Config
	oauthClient auth.Auth
}
type Config struct {
	ServerAddr   string
	ClientID     string
	ClientSecret string
}

// NewClient create a new client
// require context
func NewClient(ctx context.Context, c *Config) (Client, error) {

	// Start OAuth client
	authCfg := auth.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		TokenURL:     c.ServerAddr + "/PAM-OAuth/oauth2/token",
	}
	oauthClient, err := auth.NewAuth(ctx, &authCfg)

	// Deal with error
	if err != nil {
		return nil, err
	}
	apiClient := resty.New()
	apiClient.SetRetryCount(3)
	apiClient.SetRetryMaxWaitTime(20 * time.Second)
	apiClient.SetHostURL(c.ServerAddr)
	apiClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	// return if success
	return &client{
		client:      apiClient,
		config:      *c,
		oauthClient: oauthClient,
	}, nil
}

func (c *client) Resolve(accountID string) (string, string, error) {
	endpoint := "/PAM-Privilege/v1/accounts/check-out/" + accountID
	token, err := c.oauthClient.Token()
	if err != nil {
		return "", "", err
	}
	resp, err := c.client.R().
		SetResult(&models.AccountCheckoutResp{}).
		SetAuthToken(token).
		Put(endpoint)
	if err != nil {
		return "", "", err
	}
	if resp.StatusCode() != http.StatusOK {
		return "", "", fmt.Errorf("recieved http error code: %v", resp.StatusCode())
	}

	return resp.Result().(*models.AccountCheckoutResp).Data.Username,
		resp.Result().(*models.AccountCheckoutResp).Data.Password,
		nil
}
