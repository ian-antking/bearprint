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
	var buf bytes.Buffer
	if err := goldmark.Convert(readmeContent, &buf); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error rendering README")
		fmt.Println("markdown conversion error:", err)
		return
	}

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="utf-8">
			<title>BearPrint API Docs</title>
			<style>
				body { font-family: sans-serif; max-width: 800px; margin: auto; padding: 2em; line-height: 1.6; }
				pre { background: #f4f4f4; padding: 1em; overflow-x: auto; }
				code { background: #eee; padding: 0.2em 0.4em; }
				table { border-collapse: collapse; }
				th, td { border: 1px solid #ccc; padding: 0.5em; }
			</style>
		</head>
		<body>
			%s
		</body>
		</html>`, buf.String())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(html)); err != nil {
		fmt.Println("error writing response:", err)
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
