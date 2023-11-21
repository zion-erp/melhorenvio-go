package melhorenvio

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Range string
type Type string
type Status string

const (
	Range_Interstate Range = "interstate"

	Type_Normal   Type = "normal"
	Type_Express  Type = "express"
	Type_Economic Type = "economic"

	Status_Available Status = "available"
)

type MinMax struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type MinMaxMaxDec struct {
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	MaxDec float64 `json:"max_dec"`
}

type Box struct {
	Weight MinMax `json:"weight"`
	Width  MinMax `json:"width"`
	Height MinMax `json:"height"`
	Length MinMax `json:"length"`
	Sum    int32  `json:"sum"`
}
type Roll struct {
	Weight   MinMax `json:"weight"`
	Diameter MinMax `json:"diameter"`
	Length   MinMax `json:"length"`
	Sum      int32  `json:"sum"`
}
type Letter struct {
	Weight MinMax `json:"weight"`
	Width  MinMax `json:"width"`
	Length MinMax `json:"length"`
}
type Formats struct {
	Box    Box    `json:"box"`
	Roll   Roll   `json:"roll"`
	Letter Letter `json:"letter"`
}
type Restrictions struct {
	InsuranceValue MinMaxMaxDec `json:"insurance_value"`
	Formats        Formats      `json:"formats"`
}
type Company struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`

	HasGroupedVolumes int32  `json:"has_grouped_volumes"`
	Status            Status `json:"status"`
	TrackingLink      string `json:"tracking_link"`
	UseOwnContract    bool   `json:"use_own_contract"`
	BatchSize         int32  `json:"batch_size"`
}

type Service struct {
	ID           int32        `json:"id"`
	Name         string       `json:"name"`
	Status       Status       `json:"status"`
	Type         Type         `json:"type"`
	Range        Range        `json:"range"`
	Restrictions Restrictions `json:"restrictions"`
	Requirements []string     `json:"requirements"`
	Optionals    []string     `json:"optionals"`
	Company      Company      `json:"company"`
}

func (c *Client) GetServiceInfo(serviceId int32) (*Service, error) {
	httpReq, err := http.NewRequestWithContext(c.context, "GET", c.config.ApiUrl+"/api/v2/me/shipment/services/"+strconv.FormatInt(int64(serviceId), 10), nil)
	if err != nil {
		return nil, err
	}

	// TODO executar a partir de outra função, que não envie os dados de autenticação, pois esta rota é pública
	httpResp, err := c.doRequest(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	body, _ := io.ReadAll(httpResp.Body)

	switch httpResp.StatusCode {
	case http.StatusOK:
		resp := &Service{}
		err = json.Unmarshal(body, resp)
		if err != nil {
			return nil, err
		}

		return resp, nil

	// case http.StatusUnauthorized:
	// 	return nil, ErrInvalidToken
	default:
		return nil, fmt.Errorf("melhor envio: service: unrecognized response: %v %v", httpResp.StatusCode, string(body))
	}
}
