package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ian-antking/bear-print/bearprint-api/localprinter"
	"github.com/ian-antking/bear-print/shared/printer"
)

//go:embed README.md
var readmeContent []byte

type App struct {
	printerWriterFactory func() (io.WriteCloser, error)
}

func NewApp(printerWriterFactory func() (io.WriteCloser, error)) *App {
	return &App{printerWriterFactory: printerWriterFactory}
}

func (a *App) rootHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    if err := writeHTMLHeader(w); err != nil {
        http.Error(w, "Failed to write response header", http.StatusInternalServerError)
        return
    }

    rendered, err := renderMarkdown(readmeContent)
    if err != nil {
        http.Error(w, "Failed to render markdown", http.StatusInternalServerError)
        return
    }

    if _, err := w.Write(rendered); err != nil {
        http.Error(w, "Failed to write rendered markdown", http.StatusInternalServerError)
        return
    }

    if err := writeHTMLFooter(w); err != nil {
        http.Error(w, "Failed to write response footer", http.StatusInternalServerError)
        return
    }
}

func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

func (a *App) printHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	f, err := a.printerWriterFactory()
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
}

func main() {
	app := NewApp(func() (io.WriteCloser, error) {
		return os.OpenFile("/dev/usb/lp0", os.O_WRONLY, 0)
	})

	http.HandleFunc("/", app.rootHandler)
	http.HandleFunc("/api/v1/health", app.healthHandler)
	http.HandleFunc("/api/v1/print", app.printHandler)

	fmt.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
