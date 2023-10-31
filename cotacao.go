package melhorenvio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ToFrom struct {
	PostalCode string `json:"postal_code"`
}

type Dimensions struct {
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
	Length float64 `json:"length"`
}

type Product struct {
	ID string `json:"id"`
	Dimensions
	Weight         float64 `json:"weight"`
	InsuranceValue float64 `json:"insurance_value"`
	Quantity       int32   `json:"quantity"`
}

type Volume struct {
	Dimensions
	Weight         float64 `json:"weight"`
	InsuranceValue float64 `json:"insurance_value"`
}

type Options struct {
	Receipt bool `json:"receipt"`
	OwnHand bool `json:"own_hand"`
}

type CotacaoRequest struct {
	From     ToFrom    `json:"from"`
	To       ToFrom    `json:"to"`
	Products []Product `json:"products,omitempty"`
	Volumes  []Volume  `json:"volumes,omitempty"`
	Options  Options   `json:"options,omitempty"`
}

type DeliveryRange struct {
	Min int32 `json:"min"`
	Max int32 `json:"max"`
}

type Package struct {
	Price          string     `json:"price"`
	Discount       string     `json:"discount"`
	Format         string     `json:"format"`
	Dimensions     Dimensions `json:"dimensions"`
	Weight         string     `json:"weight"`
	InsuranceValue string     `json:"insurance_value"`
	Products       []Product  `json:"products"`
}

type AdditionalService struct {
	Receipt bool `json:"receipt"`
	OwnHand bool `json:"own_hand"`
	Collect bool `json:"collect"`
}

type Company struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type CotacaoResponse struct {
	ID                  int32             `json:"id"`
	Name                string            `json:"name"`
	Price               string            `json:"price"`
	CustomPrice         string            `json:"custom_price"`
	Discount            string            `json:"discount"`
	Currency            string            `json:"currency"`
	DeliveryTime        int32             `json:"delivery_time"`
	DeliveryRange       DeliveryRange     `json:"delivery_range"`
	CustomDeliveryTime  int32             `json:"custom_delivery_time"`
	CustomDeliveryRange DeliveryRange     `json:"custom_delivery_range"`
	Packages            []Package         `json:"packages"`
	AdditionalServices  AdditionalService `json:"additional_services"`
	Company             Company           `json:"company"`
}

type CotacaoError struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

func (ce *CotacaoError) Error() string {
	return "melhor envio: cotacao: " + ce.Message + ": " + fmt.Sprintf("%v", ce.Errors)
}

func (c *Client) CotarFrete(req *CotacaoRequest) ([]*CotacaoResponse, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(c.context, "POST", c.config.ApiUrl+"/api/v2/me/shipment/calculate", buf)
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
		var resp []*CotacaoResponse
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return nil, err
		}

		return resp, nil
	case http.StatusUnprocessableEntity:
		ret := &CotacaoError{}
		err = json.Unmarshal(body, ret)
		if err != nil {
			return nil, err
		}
		return nil, ret

	case http.StatusUnauthorized:
		return nil, ErrInvalidToken
	default:
		return nil, fmt.Errorf("melhor envio: cotacao: unrecognized response: %v %v", httpResp.StatusCode, string(body))
	}
}
