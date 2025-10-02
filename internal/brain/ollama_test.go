package brain

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOllama(t *testing.T) {
	u := "http://example.com"
	tk := "dummy-token"
	m := "gemma3"
	sys := "system-message"

	b := NewOllama(u, m, tk, sys)
	ob, ok := b.(*ollamaBrain)
	require.True(t, ok, "expected ollamaBrain type")
	assert.Equal(t, u, ob.url)
	assert.NotNil(t, ob.httpCient)
	assert.Equal(t, "gemma3", ob.model)
	assert.Equal(t, sys, ob.system)
}

func TestAsk_InvalidURL(t *testing.T) {
	b := &ollamaBrain{
		url:       "://invalid-url",
		httpCient: nil,
		model:     "gemma3",
		system:    "system-message",
	}
	ans, err := b.Ask("test-question")
	assert.Empty(t, ans)
	assert.Error(t, err)
}

type MockTransport struct {
	RoundTripFunc func(req *http.Request) *http.Response
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req), nil
}

func NeoMockClient(resp string, status int) *http.Client {
	return &http.Client{
		Transport: &MockTransport{
			RoundTripFunc: func(req *http.Request) *http.Response {
				body := io.NopCloser(strings.NewReader(resp))
				return &http.Response{
					StatusCode: status,
					Body:       body,
					Header:     make(http.Header),
				}
			},
		},
	}
}

func TestAsk_ClientError(t *testing.T) {
	client := NeoMockClient(`{"error":"something went wrong"}`, 500)
	ollamaBrain := &ollamaBrain{
		url:       "http://example.com",
		httpCient: client,
		model:     "gemma3",
		system:    "system-message",
	}
	_, err := ollamaBrain.Ask("test-question")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error from Ollama API")
}

func TestAsk_Success(t *testing.T) {
	client := NeoMockClient(`{"model":"gemma3","created_at":"2025-10-02T12:00:00Z","message":{"role":"assistant","content":"Hello"},"done":true}`, 200)

	ollamaBrain := &ollamaBrain{
		url:       "http://example.com",
		httpCient: client,
		model:     "llama3.1",
		system:    "system-message",
	}
	answ, err := ollamaBrain.Ask("test-question")
	require.NoError(t, err)
	assert.Equal(t, "Hello", answ)
}
