package btcturk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

type OpenOrderModel struct {
	ID                   int64           `json:"id"`
	Price                string          `json:"price"`
	Amount               string          `json:"amount"`
	Quantity             string          `json:"quantity"`
	StopPrice            string          `json:"stopPrice"`
	PairSymbol           string          `json:"pairSymbol"`
	PairSymbolNormalized string          `json:"pairSymbolNormalized"`
	Type                 SideType        `json:"type"`
	Method               OrderType       `json:"method"`
	OrderClientID        string          `json:"orderClientId"`
	Time                 int64           `json:"time"`
	UpdateTime           int64           `json:"updateTime"`
	Status               OrderStatusType `json:"status"`
	LeftAmount           string          `json:"leftAmount"`
}

type OpenOrderResult struct {
	Asks []OpenOrderModel `json:"asks"`
	Bids []OpenOrderModel `json:"bids"`
}

type OrderResult struct {
	ID                   int64           `json:"id"`
	Price                string          `json:"price"`
	Amount               string          `json:"amount"`
	Quantity             string          `json:"quantity"`
	PairSymbol           string          `json:"pairsymbol"`
	PairSymbolNormalized string          `json:"pairsymbolnormalized"`
	Type                 SideType        `json:"type"`
	Method               OrderType       `json:"method"`
	OrderClientID        string          `json:"orderClientId"`
	Time                 int64           `json:"time"`
	UpdateTime           int64           `json:"updateTime"`
	Status               OrderStatusType `json:"status"`
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

func (c *Client) NewOrder(o *OrderInput) (OrderResult, error) {
	jsonString, err := json.Marshal(o)
	if err != nil {
		return OrderResult{}, err
	}

	req, err := c.newRequest("POST", "/api/v1/order", bytes.NewBuffer(jsonString))
	if err != nil {
		return OrderResult{}, err
	}
	if err := c.auth(req); err != nil {
		return OrderResult{}, err
	}

	var response OrderResult
	if _, err = c.do(req, &response); err != nil {
		return OrderResult{}, err
	}

	return response, nil
}

func (c *Client) OpenOrders() (OpenOrderResult, error) {
	jsonString, err := json.Marshal(c.params)
	if err != nil {
		return OpenOrderResult{}, err
	}
	req, err := c.newRequest("GET", "/api/v1/openOrders", bytes.NewBuffer(jsonString))
	if err != nil {
		return OpenOrderResult{}, err
	}
	if err := c.auth(req); err != nil {
		return OpenOrderResult{}, err
	}

	var response OpenOrderResult
	if _, err = c.do(req, &response); err != nil {
		return OpenOrderResult{}, err
	}

	return response, nil
}

func (c *Client) AllOrders() ([]OrderResult, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/api/v1/allOrders?%s", c.params.Encode()), c.body)
	if err != nil {
		return make([]OrderResult, 0), err
	}
	if err := c.auth(req); err != nil {
		return make([]OrderResult, 0), err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return make([]OrderResult, 0), err
	}

	defer func() {
		_, err := io.Copy(ioutil.Discard, resp.Body)
		if err != nil {
			return
		}
		err = resp.Body.Close()
		if err != nil {
			return
		}
		c.clearRequest()
	}()

	var response []OrderResult
	if _, err = c.do(req, &response); err != nil {
		return make([]OrderResult, 0), err
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

	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}

	defer func() {
		_, err := io.Copy(ioutil.Discard, resp.Body)
		if err != nil {
			return
		}
		err = resp.Body.Close()
		if err != nil {
			return
		}
		c.clearRequest()
	}()

	var response = &GeneralResponse{}

	if json.NewDecoder(resp.Body).Decode(response) != nil {
		return false, err
	}

	if response.Success == true {
		return true, nil
	} else {
		return false, errors.New(*response.Message)
	}
}
