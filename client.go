package melhorenvio

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"
)

type Environment string

const (
	SandboxApiUrl    = "https://sandbox.melhorenvio.com.br"
	ProductionApiUrl = "https://melhorenvio.com.br"
)

type CredentialsChangedCallback = func(credentials Credentials) error

type Credentials struct {
	ClientId     int32
	ClientSecret string

	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time

	Code string
}

type Config struct {
	Credentials Credentials

	ApiUrl      string
	RedirectUri string

	ApplicationName string
	Email           string

	CredentialsChangedCallback CredentialsChangedCallback
}

type Client struct {
	context context.Context

	config Config

	httpClient  *http.Client
	initialized bool

	mutex sync.Mutex
}

func NewClient(ctx context.Context, config Config) *Client {
	c := &Client{}
	c.context = ctx
	c.config = config
	if c.config.ApiUrl == "" {
		c.config.ApiUrl = SandboxApiUrl
	}
	c.httpClient = http.DefaultClient
	c.initialized = true

	return c
}

func (c *Client) injectDefaultHeaders(req *http.Request) {
	if req == nil {
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.config.ApplicationName+" ("+c.config.Email+")")
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	// faz a requisição, já injetando a autenticação e gerenciando o processo de refresh de token
	// ao dar retry por conta da autenticação, pode dar erro se o body do request não for um dos tipos
	// que é possível fazer o retry (ex: bytes.Buffer)
	c.injectDefaultHeaders(req)

	if c.config.Credentials.ExpiresAt.Before(time.Now()) {
		err := c.RefreshToken()
		if err != nil {
			return nil, err
		}
	}

	req.Header.Set("Authorization", "Bearer "+c.config.Credentials.AccessToken)

	response, err := c.httpClient.Do(req)
	if err != nil {
		return response, err
	}

	switch response.StatusCode {
	case http.StatusUnauthorized:
		io.Copy(io.Discard, response.Body)
		response.Body.Close()

		err = c.RefreshToken()
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+c.config.Credentials.AccessToken)

		response, err = c.httpClient.Do(req)
		if err != nil {
			return response, err
		}
		// don't close body

		// probably not needed, but I'll check anyway
		switch response.StatusCode {
		case http.StatusUnauthorized:
			return nil, ErrInvalidToken
		}
	}
	return response, err
}
