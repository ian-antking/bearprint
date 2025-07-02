package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"bufio"
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

	url := "http://192.168.0.191:8080/api/v1/print"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
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
