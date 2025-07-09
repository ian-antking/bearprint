package printservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ian-antking/bearprint/shared/printer"
)

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

func (c *Client) Print(items []printer.PrintItem) error {
	reqBody := printer.PrintRequest{Items: items}

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
