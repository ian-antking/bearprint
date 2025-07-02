package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"bufio"
	"net/url"

	"github.com/ian-antking/bear-print/bearprint-cli/config"
)
type PrintItem struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	Align   string `json:"align,omitempty"`
	Count   int    `json:"count,omitempty"`
}

type PrintRequest struct {
	Items []PrintItem `json:"items"`
}

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
	}

	url := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort),
		Path:   "/api/v1/print",
  }
	
	scanner := bufio.NewScanner(os.Stdin)

	printReq := PrintRequest{
	Items: []PrintItem{},
}

for scanner.Scan() {
	line := scanner.Text()

	printReq.Items = append(printReq.Items, PrintItem{
		Type:    "text",
		Content: line,
	})
}

	printReq.Items = append(printReq.Items, PrintItem{Type: "cut"})

	data, err := json.Marshal(printReq)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode request: %v\n", err)
		os.Exit(1)
	}

	resp, err := http.Post(url.String(), "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to send request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "printer returned error: %s\n", resp.Status)
		os.Exit(1)
	}

	fmt.Println("✅ Print sent successfully!")
}
