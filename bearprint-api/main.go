package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ian-antking/bear-print/bearprint-api/localprinter"
	"github.com/ian-antking/bear-print/shared/printer"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
  "github.com/yuin/goldmark/renderer/html"
  "github.com/yuin/goldmark-highlighting"
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
    md := goldmark.New(
        goldmark.WithExtensions(
            extension.Table,
            highlighting.NewHighlighting(
                highlighting.WithStyle("dracula"),
                highlighting.WithFormatOptions(),
            ),
        ),
        goldmark.WithRendererOptions(
            html.WithHardWraps(),
            html.WithXHTML(),
        ),
    )
    var buf bytes.Buffer
    if err := md.Convert(readmeContent, &buf); err != nil {
        http.Error(w, "Failed to render markdown", http.StatusInternalServerError)
        return
    }

	if _, err := w.Write([]byte(`
			<!DOCTYPE html><html><head>
			<meta charset="UTF-8">
			<title>BearPrint API Docs</title>
			<style>
				body {
					font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
					margin: 4rem auto;
					max-width: 800px;
					padding: 0 1rem;
					background-color: #fff;
					color: #222;
				}
				pre {
					background-color: #2d2d2d;
					padding: 1rem;
					overflow-x: auto;
					border-radius: 5px;
				}
				table {
					border-collapse: collapse;
					width: 100%;
					margin-bottom: 1rem;
				}
				th, td {
					border: 1px solid #ddd;
					padding: 0.5rem;
					text-align: left;
				}
				th {
					background-color: #f4f4f4;
				}
			</style>
			</head><body>
	`)); err != nil {
			http.Error(w, "Failed to write response header", http.StatusInternalServerError)
			return
	}

	if _, err := w.Write(buf.Bytes()); err != nil {
			http.Error(w, "Failed to write rendered markdown", http.StatusInternalServerError)
			return
	}

	if _, err := w.Write([]byte(`</body></html>`)); err != nil {
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
