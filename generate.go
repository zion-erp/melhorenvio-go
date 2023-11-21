package melhorenvio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GenerateRequest struct {
	Orders []string `json:"orders"`
}

type GenerateResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type GenerateError struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

func (ge *GenerateError) Error() string {
	return "melhor envio: generate: " + ge.Message + ": " + fmt.Sprintf("%v", ge.Errors)
}

func (c *Client) Generate(req *GenerateRequest) (map[string]*GenerateResponse, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(c.context, "POST", c.config.ApiUrl+"/api/v2/me/shipment/generate", buf)
	if err != nil {
		return nil, err
	}

	httpResp, err := c.doRequest(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	body, _ := io.ReadAll(httpResp.Body)

	switch httpResp.StatusCode {
	case http.StatusOK:
		var resp map[string]*GenerateResponse
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return nil, err
		}

		return resp, nil
	case http.StatusUnprocessableEntity:
		ret := &GenerateError{}
		err = json.Unmarshal(body, ret)
		if err != nil {
			return nil, err
		}
		return nil, ret

	case http.StatusUnauthorized:
		return nil, ErrInvalidToken
	default:
		return nil, fmt.Errorf("melhor envio: generate: unrecognized response: %v %v", httpResp.StatusCode, string(body))
	}
}
