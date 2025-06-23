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

func TestDeepSeek_Ask_PositiveCase(t *testing.T) {
	server := echoServer(t)
	deepseek := NewDeepSeek("test_api_key")
	deepseek.(*DeepSeek).url = server.URL

	answer, err := deepseek.Ask("This is a test question")
	require.NoError(t, err)
	require.Equal(t, "This is a test question", answer)
}

func TestDeepSeek_Ask_NegativeCase_InternalServerError(t *testing.T) {
	server := errorServer(t)
	defer server.Close()
	deepseek := NewDeepSeek("test_api_key")
	deepseek.(*DeepSeek).url = server.URL

	answer, err := deepseek.Ask("This is a test question")

	require.Error(t, err)
	require.Empty(t, answer)
}

func echoServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err, "Failed to read request body")
		var request deepseekReq
		err = json.Unmarshal(body, &request)
		require.NoError(t, err, "Failed to unmarshal request body")
		require.NoError(t, err, "Failed to read request body")
		w.WriteHeader(http.StatusOK)
		message := strings.ReplaceAll(request.Messages[1].Content, "\n", "\\n")
		message = strings.ReplaceAll(message, "\"", "'")
		resp := fmt.Sprintf(`{"choices":[{"message":{"content":"%s"}}]}`, message)
		_, err = w.Write([]byte(resp))
		require.NoError(t, err, "Failed to write response")
	}))
}

func errorServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error":"Internal Server Error"}`))
		require.NoError(t, err, "Failed to write error response")
	}))
}
