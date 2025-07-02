package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ian-antking/bear-print/bearprint-cli/config"
	"github.com/ian-antking/bear-print/bearprint-cli/printer"
)

func main() {
	host := flag.String("host", "", "Server host")
	port := flag.String("port", "", "Server port")
	q_flag := flag.Bool("q", false, "Treat input as a single QR code")
	qrcode_flag := flag.Bool("qrcode", false, "Treat input as a single QR code")
	flag.Parse()
	
	cfg, err := config.NewConfig(*host, *port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	printerClient := printer.NewClient(cfg.ServerHost, cfg.ServerPort)

	var items []printer.PrintItem

	if *q_flag || *qrcode_flag {
		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}

		content := strings.TrimSpace(string(stdin))
		if len(content) > 0 {
			items = append(items, printer.PrintItem{
				Type:    printer.QRCode,
				Content: content,
				Align: printer.AlignCenter,
			})
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			items = append(items, printer.PrintItem{
				Type:    printer.Text,
				Content: scanner.Text(),
			})
		}
	}

	if len(items) > 0 {
		items = append(items, printer.PrintItem{Type: printer.Cut})
	}

	if err := printerClient.Print(items); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
