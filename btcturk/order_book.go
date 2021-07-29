package btcturk

import "fmt"

// OrderBook https://docs.btcturk.com/#order-book
type OrderBook struct {
	TimeStamp int64      `json:"timestamp"`
	Bids      [][]string `json:"bids"`
	Asks      [][]string `json:"asks"`
}

// OrderBook GET ?pairSymbol=BTC_TRY
// Or
// GET ?pairSymbol=BTC_TRY&limit=100
func (c *Client) OrderBook() (*OrderBook, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/api/v2/orderbook?%s", c.params.Encode()), nil)

	if err != nil {
		return nil, err
	}

	var response OrderBook
	if _, err = c.do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
