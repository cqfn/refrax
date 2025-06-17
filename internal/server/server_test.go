package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cqfn/refrax/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_AgentCard(t *testing.T) {
	server := NewA2AServer(8080)
	req, err := http.NewRequest(http.MethodGet, "/.well-known/agent-card", nil)
	require.NoError(t, err, "could not create request")
	rec := httptest.NewRecorder()

	server.mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Expected status OK")
	var agentCard models.AgentCard
	err = json.NewDecoder(rec.Body).Decode(&agentCard)
	require.NoError(t, err, "expected no error decoding response")
	assert.Equal(t, server.agentCard, agentCard, "Expected agent card to match")
}

func TestServer_AgentCard_MethodNotAllowed(t *testing.T) {
	server := NewA2AServer(8080)
	req, err := http.NewRequest(http.MethodPost, "/.well-known/agent-card", nil)
	require.NoError(t, err, "expected no error creating request")
	rec := httptest.NewRecorder()

	server.mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code, "Expected status Method Not Allowed")
}
