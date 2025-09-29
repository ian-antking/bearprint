package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/ian-antking/bearprint/bearprint-api/localprinter"
	"github.com/ian-antking/bearprint/shared/printer"
	"gopkg.in/ini.v1"
)

var version = "dev"

//go:embed README.md
var readmeContent []byte

var configPath = "/etc/bearprint/config.ini"

func getPrinterDevice() string {
    device := "/dev/usb/lp0"

    cfg, err := ini.Load(configPath)
    if err != nil {
        fmt.Println("⚠️  Warning: could not load /etc/bearprint/config.ini, using default", device)
        return device
    }

    device = cfg.Section("printer").Key("device").MustString(device)
    fmt.Println("✅ Using printer device:", device)
    return device
}

type App struct {
	printerWriterFactory func() (io.WriteCloser, error)
}

func NewApp(printerWriterFactory func() (io.WriteCloser, error)) *App {
	return &App{printerWriterFactory: printerWriterFactory}
}

func (a *App) printStartupInfo() {
	f, err := a.printerWriterFactory()
	if err != nil {
		return
	}
	defer f.Close()

	p := localprinter.NewPrinter(f)

	ip, err := getLocalIP()
	if err != nil {
		ip = "unknown"
	}

	items := []printer.PrintItem{
		{Type: "text", Content: "BearPrint", Align: "center"},
		{Type: "text", Content: fmt.Sprintf("Version: %s", version), Align: "center"},
		{Type: "text", Content: fmt.Sprintf("Server address: %s:8080", ip), Align: "center"},
		{Type: "blank", Count: 1},
		{Type: "line"},
		{Type: "text", Content: "API Endpoints:", Align: "left"},
		{Type: "text", Content: "GET  /api/v1/health  - health check", Align: "left"},
		{Type: "text", Content: "POST /api/v1/print   - print jobs", Align: "left"},
		{Type: "line"},
		{Type: "blank", Count: 3},
		{Type: "cut"},
	}

	if err := p.PrintJob(items); err != nil {
		fmt.Printf("print error: %v", err)
		return
	}
}

func (a *App) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	rendered, err := renderMarkdown(readmeContent)
	if err != nil {
		http.Error(w, "Failed to render markdown", http.StatusInternalServerError)
		return
	}

	if err := writeHTML(w, rendered); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
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

	if err := printer.ValidatePrintRequest(req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
				http.Error(w, "validation error: "+errs.Error(), http.StatusBadRequest)
		} else {
				http.Error(w, "validation error: "+err.Error(), http.StatusBadRequest)
		}
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
		return os.OpenFile(getPrinterDevice(), os.O_WRONLY, 0)
	})

	app.printStartupInfo()

	http.HandleFunc("/", app.rootHandler)
	http.HandleFunc("/api/v1/health", app.healthHandler)
	http.HandleFunc("/api/v1/print", app.printHandler)

	fmt.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
