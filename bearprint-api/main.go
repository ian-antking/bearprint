package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ian-antking/bear-print/bearprint-api/localprinter"
	"github.com/ian-antking/bear-print/shared/printer"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	http.HandleFunc("/api/v1/print", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		f, err := os.OpenFile("/dev/usb/lp0", os.O_WRONLY, 0)
		if err != nil {
			http.Error(w, "failed to open printer device", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		p := localprinter.NewPrinter(f)

		var req printer.PrintRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := p.PrintJob(req.Items); err != nil {
			http.Error(w, fmt.Sprintf("print error: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "printed")
	})

	fmt.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
