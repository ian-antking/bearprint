package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ian-antking/bear-print/shared/printer"
)

type mockWriteCloser struct {
	Written bytes.Buffer
	Closed  bool
}

func (m *mockWriteCloser) Write(p []byte) (int, error) {
	return m.Written.Write(p)
}

func (m *mockWriteCloser) Close() error {
	m.Closed = true
	return nil
}

func TestHealthHandler(t *testing.T) {
	app := NewApp(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	app.healthHandler(w, req)

	resp := w.Result()
	body := w.Body.String()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "ok\n", body)
}

func TestPrintHandler_Success(t *testing.T) {
	mockPrinter := &mockWriteCloser{}

	app := NewApp(func() (io.WriteCloser, error) {
		return mockPrinter, nil
	})

	printReq := printer.PrintRequest{
		Items: []printer.PrintItem{
			{Type: "text", Content: "Hello world", Align: "left"},
		},
	}

	reqBody, err := json.Marshal(printReq)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/print", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	app.printHandler(w, req)

	resp := w.Result()
	body := w.Body.String()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, body, "printed")
	assert.True(t, mockPrinter.Closed, "printer writer should be closed")

	assert.True(t, strings.Contains(mockPrinter.Written.String(), "Hello world"), "printed output should contain 'Hello world'")
}

func TestPrintHandler_MethodNotAllowed(t *testing.T) {
	app := NewApp(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/print", nil)
	w := httptest.NewRecorder()

	app.printHandler(w, req)

	resp := w.Result()
	body := w.Body.String()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	assert.Contains(t, body, "method not allowed")
}

func TestPrintHandler_BadRequest(t *testing.T) {
	mockPrinter := &mockWriteCloser{}

	app := NewApp(func() (io.WriteCloser, error) {
		return mockPrinter, nil
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/print", strings.NewReader("invalid-json"))
	w := httptest.NewRecorder()

	app.printHandler(w, req)

	resp := w.Result()
	body := w.Body.String()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, body, "invalid request body")
}

func TestPrintHandler_PrinterOpenError(t *testing.T) {
	app := NewApp(func() (io.WriteCloser, error) {
		return nil, io.ErrClosedPipe
	})

	printReq := printer.PrintRequest{
		Items: []printer.PrintItem{
			{Type: "text", Content: "Hello", Align: "left"},
		},
	}

	reqBody, err := json.Marshal(printReq)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/print", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	app.printHandler(w, req)

	resp := w.Result()
	body := w.Body.String()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Contains(t, body, "failed to open printer device")
}

