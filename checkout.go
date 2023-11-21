package melhorenvio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CheckoutRequest struct {
	Orders []string `json:"orders"`
}

type CheckoutResponsePurchase struct {
	Id       string  `json:"id"`
	Protocol string  `json:"protocol"`
	Total    float64 `json:"total"`
	Discount float64 `json:"discount"`
	Status   string  `json:"status"`
	Orders   []struct {
		Id string `json:"id"`
		// TODO
	} `json:"orders"`
	// TODO
}

type CheckoutResponse struct {
	Purchase CheckoutResponsePurchase `json:"purchase"`
}

type CheckoutError struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

func (ce *CheckoutError) Error() string {
	return "melhor envio: checkout: " + ce.Message + ": " + fmt.Sprintf("%v", ce.Errors)
}

func (c *Client) Checkout(req *CheckoutRequest) (*CheckoutResponse, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(c.context, "POST", c.config.ApiUrl+"/api/v2/me/shipment/checkout", buf)
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
		var resp *CheckoutResponse
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return nil, fmt.Errorf("melhor envio: checkout: unrecognized response: %v %v", httpResp.StatusCode, string(body))
		}

		return resp, nil
	case http.StatusUnprocessableEntity:
		ret := &CheckoutError{}
		err = json.Unmarshal(body, ret)
		if err != nil {
			return nil, fmt.Errorf("melhor envio: checkout: unrecognized response: %v %v", httpResp.StatusCode, string(body))
		}
		return nil, ret

	case http.StatusUnauthorized:
		return nil, ErrInvalidToken
	default:
		return nil, fmt.Errorf("melhor envio: checkout: unrecognized response: %v %v", httpResp.StatusCode, string(body))
	}
}
