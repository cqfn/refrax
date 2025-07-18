package brain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type aiRequest struct {
	Model    string      `json:"model"`
	Messages []aiMessage `json:"messages"`
}

type aiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// NewEchoServer creates a test server that echoes back the user message
func NewEchoServer(t *testing.T, expectedModel, expectedToken string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if expectedToken != "" {
			auth := r.Header.Get("Authorization")
			require.Equal(t, "Bearer "+expectedToken, auth, "Invalid API key")
		}
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err, "Failed to read request body")
		var request aiRequest
		err = json.Unmarshal(body, &request)
		require.NoError(t, err, "Failed to unmarshal request body")
		if expectedModel != "" {
			require.Equal(t, expectedModel, request.Model, "Unexpected model")
		}
		w.WriteHeader(http.StatusOK)
		message := strings.ReplaceAll(request.Messages[1].Content, "\n", "\\n")
		message = strings.ReplaceAll(message, "\"", "'")
		resp := fmt.Sprintf(`{
			"choices": [{
				"message": {"content": %q}
			}]
		}`, message)
		_, err = w.Write([]byte(resp))
		require.NoError(t, err, "Failed to write response")
	}))
}

// NewErrorServer creates a test server that returns an error response
func NewErrorServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error":"Internal Server Error"}`))
		require.NoError(t, err, "Failed to write error response")
	}))
}
