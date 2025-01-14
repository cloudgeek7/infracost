package extclient

import (
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// AuthedAPIClient represents an API client for authorized requests.
type AuthedAPIClient struct {
	host   string
	token  string
	client *retryablehttp.Client
}

// NewAuthedAPIClient returns a new API client.
func NewAuthedAPIClient(host, token string) *AuthedAPIClient {
	client := retryablehttp.NewClient()
	client.HTTPClient.Timeout = time.Second * 30
	return &AuthedAPIClient{
		host:   host,
		token:  token,
		client: client,
	}
}

// SetHost sets the host for base host for the authed API client.
func (a *AuthedAPIClient) SetHost(host string) {
	a.host = host
}

// Get performs a GET request to provided endpoint.
func (a *AuthedAPIClient) Get(path string) ([]byte, error) {
	url := fmt.Sprintf("https://%s%s", a.host, path)
	log.Debugf("Calling Terraform Cloud API: %s", url)
	req, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.token))

	resp, err := a.client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return []byte{}, errors.New("the provided token is invalid")
	} else if resp.StatusCode != 200 {
		return []byte{}, errors.Errorf("invalid response: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
