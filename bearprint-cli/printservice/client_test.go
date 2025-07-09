package printservice_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ian-antking/bearprint/bearprint-cli/printservice"
	"github.com/ian-antking/bearprint/shared/printer"
	"github.com/stretchr/testify/require"
)

func TestClient_Print(t *testing.T) {
	t.Run("successful print request", func(t *testing.T) {
		var receivedBody printer.PrintRequest

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			err := json.NewDecoder(r.Body).Decode(&receivedBody)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("printed"))
		}))
		defer ts.Close()

		u := strings.TrimPrefix(ts.URL, "http://")
		hostPort := strings.Split(u, ":")
		client := printservice.NewClient(hostPort[0], hostPort[1])

		items := []printer.PrintItem{
			{Content: "hello", Type: printer.Text},
			{Content: "world", Type: printer.Text},
		}

		err := client.Print(items)
		require.NoError(t, err)
		require.Equal(t, items, receivedBody.Items)
	})

	t.Run("printer returns HTTP error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		u := strings.TrimPrefix(ts.URL, "http://")
		hostPort := strings.Split(u, ":")
		client := printservice.NewClient(hostPort[0], hostPort[1])

		err := client.Print([]printer.PrintItem{{Content: "fail", Type: printer.Text}})
		require.Error(t, err)
		require.Contains(t, err.Error(), "printer returned error")
	})
}
