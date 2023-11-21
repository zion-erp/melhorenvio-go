package melhorenvio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CartToFrom struct {
	Name            string `json:"name,omitempty"`
	Phone           string `json:"phone,omitempty"`
	Email           string `json:"email,omitempty"`
	Document        string `json:"document,omitempty"`
	CompanyDocument string `json:"company_document,omitempty"`
	StateRegister   string `json:"state_register,omitempty"`
	Address         string `json:"address,omitempty"`
	Complement      string `json:"complement,omitempty"`
	Number          string `json:"number,omitempty"`
	District        string `json:"district,omitempty"`
	City            string `json:"city,omitempty"`
	CountryId       string `json:"country_id,omitempty"`
	PostalCode      string `json:"postal_code,omitempty"`
	StateAbbr       string `json:"state_abbr,omitempty"`
	// Note            string `json:"note,omitempty"` // removido pois dá erro 500 se enviar
}

type CartProduct struct {
	Name         string  `json:"name"`
	Quantity     float64 `json:"quantity,omitempty"`
	UnitaryValue float64 `json:"unitary_value,omitempty"`

	Weight float64 `json:"weight,omitempty"`
}

type CartVolume struct {
	Dimensions
	Weight float64 `json:"weight"`
}

type CartOptions struct {
	Options
	Reverse       bool `json:"reverse"`
	NonCommercial bool `json:"non_commercial"`
	Invoice       struct {
		Key string `json:"key"`
	} `json:"invoice,omitempty"`
	Plataform string `json:"plataform,omitempty"`
	Tag       []struct {
		Tag string `json:"tag,omitempty"`
		Url string `json:"url,omitempty"`
	} `json:"tags,omitempty"`
}

type AddToCartRequest struct {
	Service  int32         `json:"service"`
	Agency   int32         `json:"agency,omitempty"`
	From     CartToFrom    `json:"from"`
	To       CartToFrom    `json:"to"`
	Products []CartProduct `json:"products"`
	Volumes  []CartVolume  `json:"volumes"`
	Options  CartOptions   `json:"options"`
}

type CartResponseVolume struct {
	Id        int32  `json:"id"`
	Height    string `json:"height"`
	Width     string `json:"width"`
	Length    string `json:"length"`
	Diameter  string `json:"diameter"`
	Weight    string `json:"weight"`
	Format    string `json:"format"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CartResponse struct {
	Id                 string  `json:"id"`
	Protocol           string  `json:"protocol"`
	ServiceId          int32   `json:"service_id"`
	AgencyId           int32   `json:"agency_id"`
	Contract           string  `json:"contract"`
	ServiceCode        string  `json:"service_code"`
	Quote              float64 `json:"quote"`
	Price              float64 `json:"price"`
	Coupon             string  `json:"coupon"`
	Discount           float64 `json:"discount"`
	DeliveryMin        int32   `json:"delivery_min"`
	DeliveryMax        int32   `json:"delivery_max"`
	Status             string  `json:"status"`
	Reminder           string  `json:"reminder"`
	InsuranceValue     float64 `json:"insurance_value"`
	Weight             string  `json:"weight"`
	Width              string  `json:"width"`
	Height             string  `json:"height"`
	Length             string  `json:"length"`
	Diameter           string  `json:"diameter"`
	Format             string  `json:"format"`
	BilledWeight       float64 `json:"billed_weight"`
	Receipt            bool    `json:"receipt"`
	OwnHand            bool    `json:"own_hand"`
	Collect            bool    `json:"collect"`
	CollectScheduledAt string  `json:"collect_scheduled_at"`
	// Reverse            bool                 `json:"reverse"` // removido pois vem 0 no lugar de um bool
	NonCommercial     bool                 `json:"non_commercial"`
	AuthorizationCode string               `json:"authorization_code"`
	Tracking          string               `json:"tracking"`
	SelfTracking      string               `json:"self_tracking"`
	DeliveryReceipt   string               `json:"delivery_receipt"`
	AdditionalInfo    string               `json:"additional_info"`
	CteKey            string               `json:"cte_key"`
	PaidAt            string               `json:"paid_at"`
	GeneratedAt       string               `json:"generated_at"`
	PostedAt          string               `json:"posted_at"`
	DeliveredAt       string               `json:"delivered_at"`
	CanceledAt        string               `json:"canceled_at"`
	SuspendedAt       string               `json:"suspended_at"`
	ExpiredAt         string               `json:"expired_at"`
	CreatedAt         string               `json:"created_at"`
	UpdatedAt         string               `json:"updated_at"`
	ParsePiAt         string               `json:"parse_pi_at"`
	Products          []CartProduct        `json:"products"`
	Volumes           []CartResponseVolume `json:"volumes"`
}

type CartError struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"error"`
}

func (ce *CartError) Error() string {
	return "melhor envio: cart: " + ce.Message + ": " + fmt.Sprintf("%v", ce.Errors)
}

func (c *Client) AddToCart(req *AddToCartRequest) (*CartResponse, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(c.context, "POST", c.config.ApiUrl+"/api/v2/me/cart", buf)
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
	case http.StatusCreated:
		var resp *CartResponse
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return nil, fmt.Errorf("melhor envio: cart: unrecognized response: %v %v", httpResp.StatusCode, string(body))
		}

		return resp, nil
	case http.StatusUnprocessableEntity:
		ret := &CartError{}
		err = json.Unmarshal(body, ret)
		if err != nil {
			return nil, fmt.Errorf("melhor envio: cart: unrecognized response: %v %v", httpResp.StatusCode, string(body))
		}
		return nil, ret

	case http.StatusUnauthorized:
		return nil, ErrInvalidToken
	default:
		return nil, fmt.Errorf("melhor envio: cart: unrecognized response: %v %v", httpResp.StatusCode, string(body))
	}
}

func (c *Client) RemoveFromCart(orderId string) error {
	httpReq, err := http.NewRequestWithContext(c.context, "DELETE", c.config.ApiUrl+"/api/v2/me/cart/"+orderId, nil)
	if err != nil {
		return err
	}

	httpResp, err := c.doRequest(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	body, _ := io.ReadAll(httpResp.Body)

	switch httpResp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusUnprocessableEntity, http.StatusBadRequest:
		ret := &CartError{}
		err = json.Unmarshal(body, ret)
		if err != nil {
			return fmt.Errorf("melhor envio: cart: unrecognized response: %v %v", httpResp.StatusCode, string(body))
		}
		return ret

	case http.StatusUnauthorized:
		return ErrInvalidToken
	default:
		return fmt.Errorf("melhor envio: cart: unrecognized response: %v %v", httpResp.StatusCode, string(body))
	}
}
