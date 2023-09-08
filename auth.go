package melhorenvio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// preciso de uma função pra obter o primeiro access token
// de uma função pra fazer o refresh do token
// de uma função pra injetar os headers em um certo request

type authRequest struct {
	GrantType    string `json:"grant_type"`
	ClientId     int32  `json:"client_id"`
	ClientSecret string `json:"client_secret"`

	RedirectUri string `json:"redirect_uri,omitempty"`
	Code        string `json:"code,omitempty"`

	RefreshToken string `json:"refresh_token,omitempty"`
}

type authResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int32  `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (c *Client) AutenticateByCode() error {
	if !c.initialized {
		return ErrClientNotInitialized
	}

	if c.config.Credentials.Code == "" {
		return ErrInvalidToken
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	buf := &bytes.Buffer{}
	aReq := &authRequest{
		GrantType:    "authorization_code",
		ClientId:     c.config.Credentials.ClientId,
		ClientSecret: c.config.Credentials.ClientSecret,
		RedirectUri:  c.config.RedirectUri,
		Code:         c.config.Credentials.Code,
	}
	err := json.NewEncoder(buf).Encode(aReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(c.context, "POST", c.config.ApiUrl+"/oauth/token", buf)
	if err != nil {
		return err
	}

	c.injectDefaultHeaders(req)

	response, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return c.parseAuthResponse(response)
}

func (c *Client) RefreshToken() error {
	if !c.initialized {
		return ErrClientNotInitialized
	}

	if c.config.Credentials.RefreshToken == "" {
		return ErrInvalidToken
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	buf := &bytes.Buffer{}
	aReq := &authRequest{
		GrantType:    "refresh_token",
		ClientId:     c.config.Credentials.ClientId,
		ClientSecret: c.config.Credentials.ClientSecret,
		RefreshToken: c.config.Credentials.RefreshToken,
	}
	err := json.NewEncoder(buf).Encode(aReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(c.context, "POST", c.config.ApiUrl+"/oauth/token", buf)
	if err != nil {
		return err
	}

	c.injectDefaultHeaders(req)

	response, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return c.parseAuthResponse(response)
}

func (c *Client) parseAuthResponse(response *http.Response) error {
	body, _ := io.ReadAll(response.Body)

	switch response.StatusCode {
	case http.StatusOK:
		aResp := &authResponse{}
		err := json.Unmarshal(body, aResp)
		if err != nil {
			return err
		}

		c.config.Credentials.AccessToken = aResp.AccessToken
		c.config.Credentials.RefreshToken = aResp.RefreshToken
		c.config.Credentials.ExpiresAt = time.Now().Add(time.Duration(aResp.ExpiresIn) * time.Second)
		c.config.Credentials.Code = ""

		if c.config.CredentialsChangedCallback != nil {
			err = c.config.CredentialsChangedCallback(c.config.Credentials)
			if err != nil {
				return err
			}
		}

		return nil
	case http.StatusUnauthorized:
		return ErrInvalidToken
	default:
		return fmt.Errorf("melhor envio: auth: unrecognized response: %v %v", response.StatusCode, string(body))
	}
}
