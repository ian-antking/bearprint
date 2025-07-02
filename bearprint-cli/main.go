package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ian-antking/bear-print/bearprint-cli/config"
	"github.com/ian-antking/bear-print/bearprint-cli/printer"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	printerClient := printer.NewClient(cfg.ServerHost, cfg.ServerPort)

	var items []printer.PrintItem
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		items = append(items, printer.PrintItem{
			Type:    printer.Text,
			Content: scanner.Text(),
		})
	}
	items = append(items, printer.PrintItem{Type: "cut"})

	if err := printerClient.Print(items); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
