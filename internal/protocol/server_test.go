package protocol

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cqfn/refrax/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerMux_AgentCard(t *testing.T) {
	port, err := util.FreePort()
	require.NoError(t, err, "expected no error getting free port")
	server := NewServer(mockCard(port), port)
	server.MsgHandler(joke)
	require.NoError(t, err, "expected no error creating server")
	cserver := server.(*a2aServer)
	require.NoError(t, err, "expected no error creating server")
	req, err := http.NewRequest(http.MethodGet, "/.well-known/agent-card.json", http.NoBody)
	require.NoError(t, err, "could not create request")
	rec := httptest.NewRecorder()

	cserver.mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Expected status OK")
	var card AgentCard
	err = json.NewDecoder(rec.Body).Decode(&card)
	require.NoError(t, err, "expected no error decoding response")
	assert.Equal(t, *mockCard(port), card, "Expected agent card to match")
}

func TestServerMux_AgentCard_MethodNotAllowed(t *testing.T) {
	port, err := util.FreePort()
	require.NoError(t, err, "expected no error getting free port")
	server := NewServer(mockCard(port), port)
	server.MsgHandler(joke)
	assert.NoError(t, err, "expected no error creating server")
	cserver := server.(*a2aServer)
	req, err := http.NewRequest(http.MethodPost, "/.well-known/agent-card.json", http.NoBody)
	require.NoError(t, err, "expected no error creating request")
	rec := httptest.NewRecorder()

	cserver.mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code, "Expected status Method Not Allowed")
}

func mockCard(port int) *AgentCard {
	return NewAgentCard().
		WithName("Test Agent").
		WithDescription("A test agent for unit tests").
		WithURL(fmt.Sprintf("http://localhost:%d", port)).
		WithVersion("0.0.1").
		AddSkill("refactor-java", "Refactor Java Projects", "Refrax can refactor java projects")
}
