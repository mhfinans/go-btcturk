package btcturk

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://api.btcturk.com"

type OrderType string
type SideType string
type OrderStatusType string

const (
	OrderTypeMarket    OrderType = "market"
	OrderTypeLimit     OrderType = "limit"
	OrderTypeStopLimit OrderType = "stoplimit"

	OrderStatusTypeUntouched       OrderStatusType = "Untouched"
	OrderStatusTypePartiallyFilled OrderStatusType = "Partial"
	OrderStatusTypeCancelled       OrderStatusType = "Canceled"
	OrderStatusTypeFilled          OrderStatusType = "Closed"

	SideTypeBuy  SideType = "buy"
	SideTypeSell SideType = "sell"
)

type Client struct {
	client  *http.Client
	baseURL *url.URL
	params  url.Values
	body    io.Reader

	publicKey  string
	privateKey string
}
type GeneralResponse struct {
	Success bool        `json:"success"`
	Message *string     `json:"message"`
	Code    int64       `json:"code"`
	Data    interface{} `json:"data"`
}

func NewBTCTurkClient() *Client {
	baseURL, _ := url.Parse(baseURL)
	return &Client{
		baseURL: baseURL,
		client:  &http.Client{},
		params:  url.Values{},
		body:    nil,
	}
}

func (c *Client) SetAuthKey(publicKey, privateKey string) {
	c.publicKey, c.privateKey = publicKey, privateKey
}

func (c *Client) SetBaseURL(basURL string) (err error) {
	c.baseURL, err = url.Parse(basURL)
	return
}

func (c *Client) newRequest(method, relURL string, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse(relURL)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, c.baseURL.ResolveReference(rel).String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(context.Background())

	return req, nil
}

func (c *Client) do(r *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	defer func() {
		cerr := resp.Body.Close()
		// Only overwrite the returned error if the original error was nil and an
		// error occurred while closing the body.
		if err == nil && cerr != nil {
			err = cerr
		}
	}()

	if resp.StatusCode >= 400 {
		return nil, errors.New(fmt.Sprintf("Request: %s, Status: %s, Response: %s", r.URL.String(), resp.Status, string(data)))
	}

	var response = &GeneralResponse{
		Data: v,
	}

	if err = json.Unmarshal(data, response); err != nil {
		return nil, err
	}

	if response.Code == 0 {
		v = response.Data
	} else {
		return nil, errors.New(*response.Message)
	}

	return resp, nil
}

func (c *Client) auth(req *http.Request) error {
	if c.privateKey == "" {
		return errors.New("private key is not set")
	}

	if c.publicKey == "" {
		return errors.New("public key is not set")
	}

	key, err := base64.StdEncoding.DecodeString(c.privateKey)
	if err != nil {
		return err
	}

	stamp := fmt.Sprint(time.Now().Unix() * 1000)
	message := c.publicKey + stamp

	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))

	req.Header.Set("X-PCK", c.publicKey)
	req.Header.Set("X-Stamp", stamp)
	req.Header.Set("X-Signature", base64.StdEncoding.EncodeToString(h.Sum(nil)))
	req.Header.Add("Content-Type", "application/json")
	return nil
}

func (c *Client) clearRequest() {
	c.params = url.Values{}
	c.body = nil
}

func (c *Client) Version() string {
	return "0.0.1"
}
