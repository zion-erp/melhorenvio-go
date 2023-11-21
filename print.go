package melhorenvio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Mode string

const (
	Mode_Private Mode = "private"
	Mode_Public  Mode = "public"
)

type PrintRequest struct {
	Mode   Mode     `json:"mode"`
	Orders []string `json:"orders"`
}

type PrintResponse struct {
	Url string `json:"url"`
}

type PrintError struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

func (pe *PrintError) Error() string {
	return "melhor envio: print: " + pe.Message
}

func (c *Client) Print(req *PrintRequest) (*PrintResponse, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(c.context, "POST", c.config.ApiUrl+"/api/v2/me/shipment/print", buf)
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
		var resp *PrintResponse
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return nil, fmt.Errorf("melhor envio: print: unrecognized response: %v %v", httpResp.StatusCode, string(body))
		}

		return resp, nil
	case http.StatusUnprocessableEntity, http.StatusBadRequest:
		ret := &PrintError{}
		err = json.Unmarshal(body, ret)
		if err != nil {
			return nil, fmt.Errorf("melhor envio: print: unrecognized response: %v %v", httpResp.StatusCode, string(body))
		}
		return nil, ret

	case http.StatusUnauthorized:
		return nil, ErrInvalidToken
	default:
		return nil, fmt.Errorf("melhor envio: print: unrecognized response: %v %v", httpResp.StatusCode, string(body))
	}
}
