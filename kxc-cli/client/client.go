package client

import (
	"net/http"
	"time"

	"github.com/didil/kubexcloud/kxc-cli/config"
)

type Client struct {
	apiURL     string
	authToken  string
	httpClient *http.Client
}

func NewClient() *Client {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	cl := &Client{
		httpClient: httpClient,
	}

	if apiURL := config.GetApiUrl(); apiURL != "" {
		cl.apiURL = apiURL
	}

	if authToken := config.GetAuthToken(); authToken != "" {
		cl.authToken = authToken
	}

	return cl
}
