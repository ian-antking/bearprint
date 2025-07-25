package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ian-antking/bearprint/shared/printer"
	"github.com/stretchr/testify/assert"
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
	t.Run("GET /health returns 200 OK with body 'ok'", func(t *testing.T) {
		app := NewApp(nil)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		app.healthHandler(w, req)

		resp := w.Result()
		body := w.Body.String()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "ok\n", body)
	})
}

func TestPrintHandler(t *testing.T) {
	t.Run("POST /api/v1/print prints successfully and closes printer", func(t *testing.T) {
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
		assert.Contains(t, mockPrinter.Written.String(), "Hello world")
	})

	t.Run("GET /api/v1/print returns 405 Method Not Allowed", func(t *testing.T) {
		app := NewApp(nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/print", nil)
		w := httptest.NewRecorder()

		app.printHandler(w, req)

		resp := w.Result()
		body := w.Body.String()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
		assert.Contains(t, body, "method not allowed")
	})

	t.Run("POST /api/v1/print with invalid JSON returns 400 Bad Request", func(t *testing.T) {
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
	})

	t.Run("POST /api/v1/print fails to open printer and returns 500 Internal Server Error", func(t *testing.T) {
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
	})

	t.Run("POST /api/v1/print with empty items fails validation", func(t *testing.T) {
	mockPrinter := &mockWriteCloser{}

	app := NewApp(func() (io.WriteCloser, error) {
		return mockPrinter, nil
	})

	printReq := printer.PrintRequest{
		Items: []printer.PrintItem{},
	}

	reqBody, err := json.Marshal(printReq)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/print", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	app.printHandler(w, req)

	resp := w.Result()
	body := w.Body.String()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, body, "validation error")
})
}
