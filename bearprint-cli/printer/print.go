package printer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type ItemType string

const (
	Text       ItemType = "text"
	QRCode ItemType = "qrcode"
	Blank  ItemType = "blank"
	Line   ItemType = "line"
	Cut    ItemType = "cut"
)

type Alignment string

const (
	AlignLeft   Alignment = "left"
	AlignCenter Alignment = "center"
	AlignRight  Alignment = "right"
)

type PrintItem struct {
	Type    ItemType  `json:"type"`
	Content string    `json:"content,omitempty"`
	Align   Alignment `json:"align,omitempty"`
	Count   int       `json:"count,omitempty"`
}

type PrintRequest struct {
	Items []PrintItem `json:"items"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(host, port string) *Client {
	url := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   "/api/v1/print",
	}

	return &Client{
		baseURL:    url.String(),
		httpClient: &http.Client{},
	}
}

func (c *Client) Print(items []PrintItem) error {
	reqBody := PrintRequest{Items: items}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("printer returned error: %s", resp.Status)
	}

	return nil
}
