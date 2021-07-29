package btcturk

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type OrderResponse struct {
	ID                   int64           `json:"id"`
	Price                string          `json:"price"`
	Amount               string          `json:"amount"`
	Quantity             string          `json:"quantity"`
	StopPrice            string          `json:"stopPrice"`
	PairSymbol           string          `json:"pairsymbol"`
	PairSymbolNormalized string          `json:"pairsymbolnormalized"`
	Type                 SideType        `json:"type"`
	Method               OrderType       `json:"method"`
	OrderClientID        string          `json:"orderClientId"`
	Time                 int64           `json:"time"`
	UpdateTime           int64           `json:"updateTime"`
	Status               OrderStatusType `json:"status"`
	LeftAmount           string          `json:"leftAmount"`
}

type NewOrderResponse struct {
	ID                   int64     `json:"id"`
	Quantity             string    `json:"quantity"`
	Price                string    `json:"price"`
	StopPrice            string    `json:"stopPrice"`
	NewOrderClientID     string    `json:"newOrderClientId"`
	Type                 SideType  `json:"type"`
	Method               OrderType `json:"method"`
	PairSymbol           string    `json:"pairSymbol"`
	PairSymbolNormalized string    `json:"pairSymbolNormalized"`
	Datetime             int64     `json:"datetime"`
}

type OpenOrderResult struct {
	Asks []OrderResponse `json:"asks"`
	Bids []OrderResponse `json:"bids"`
}

type OrderInput struct {
	Quantity         float64   `json:"quantity"`
	Price            float64   `json:"price"`
	StopPrice        float64   `json:"stopPrice"`
	NewOrderClientId string    `json:"newOrderClientId"`
	OrderMethod      OrderType `json:"orderMethod"`
	OrderType        SideType  `json:"orderType"`
	PairSymbol       string    `json:"pairSymbol"`
}

func (c *Client) NewOrder(o *OrderInput) (*NewOrderResponse, error) {
	jsonString, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest("POST", "/api/v1/order", bytes.NewBuffer(jsonString))
	if err != nil {
		return nil, err
	}
	if err := c.auth(req); err != nil {
		return nil, err
	}

	var response NewOrderResponse
	if _, err = c.do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) OpenOrders() (*OpenOrderResult, error) {
	jsonString, err := json.Marshal(c.params)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest("GET", "/api/v1/openOrders", bytes.NewBuffer(jsonString))
	if err != nil {
		return nil, err
	}

	if err := c.auth(req); err != nil {
		return nil, err
	}

	var response OpenOrderResult
	if _, err = c.do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) AllOrders() ([]OrderResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/api/v1/allOrders?%s", c.params.Encode()), c.body)
	if err != nil {
		return nil, err
	}
	if err := c.auth(req); err != nil {
		return nil, err
	}

	var response []OrderResponse
	if _, err = c.do(req, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) CancelOrder() (bool, error) {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/api/v1/order?%s", c.params.Encode()), c.body)
	if err != nil {
		return false, err
	}
	if err := c.auth(req); err != nil {
		return false, err
	}

	var response = &GeneralResponse{}

	if _, err = c.do(req, &response); err != nil {
		return false, err
	} else {
		return true, nil
	}
}
