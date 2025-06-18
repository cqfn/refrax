package protocol

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_AgentCard(t *testing.T) {
	server, err := NewCustomServer(mockCard(), 8080)
	require.NoError(t, err, "expected no error creating server")
	cserver := server.(*CustomServer)
	require.NoError(t, err, "expected no error creating server")
	req, err := http.NewRequest(http.MethodGet, "/.well-known/agent-card.json", nil)
	require.NoError(t, err, "could not create request")
	rec := httptest.NewRecorder()

	cserver.mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Expected status OK")
	var card AgentCard
	err = json.NewDecoder(rec.Body).Decode(&card)
	require.NoError(t, err, "expected no error decoding response")
	assert.Equal(t, mockCard(), card, "Expected agent card to match")
}

func TestServer_AgentCard_MethodNotAllowed(t *testing.T) {
	server, err := NewCustomServer(mockCard(), 8080)
	assert.NoError(t, err, "expected no error creating server")
	cserver := server.(*CustomServer)
	req, err := http.NewRequest(http.MethodPost, "/.well-known/agent-card.json", nil)
	require.NoError(t, err, "expected no error creating request")
	rec := httptest.NewRecorder()

	cserver.mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code, "Expected status Method Not Allowed")
}

func mockCard() AgentCard {
	return Card().
		Name("Test Agent").
		Description("A test agent for unit tests").
		URL("http://localhost:8080").
		Version("0.0.1").
		Skill("refactor-java", "Refactor Java Projects", "Refrax can refactor java projects").
		Build()
}
