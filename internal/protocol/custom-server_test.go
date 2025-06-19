package protocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/cqfn/refrax/internal/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testAgentCard = AgentCard{
	Name:        "TestAgent",
	Description: "This is a test agent",
	URL:         "http://testagent.example.com",
	Version:     "1.0.0",
}

func TestCustomServer_AgentCard(t *testing.T) {
	card := testAgentCard
	port, err := freePort()
	require.NoError(t, err, "Failed to get a free port")
	server, err := NewCustomServer(card, jokeHandler, port)
	require.NoError(t, err, "Failed to create custom server")
	ready := make(chan struct{})
	go startServer(server, ready, &err)
	defer closeResource(server, &err)
	<-ready

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/.well-known/agent-card.json", port))

	require.NoError(t, err)
	defer closeResource(resp.Body, &err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var result AgentCard
	err = json.NewDecoder(resp.Body).Decode(&result)
	log.Debug("Agent card response: %v", result)
	require.NoError(t, err)
	require.Equal(t, "TestAgent", result.Name, "Agent name does not match")
}

func TestCustomServer_SendsMessage(t *testing.T) {
	port, err := freePort()
	require.NoError(t, err, "Failed to get a free port")
	server, err := NewCustomServer(testAgentCard, jokeHandler, port)
	require.NoError(t, err, "Failed to create custom server")
	ready := make(chan struct{})
	go startServer(server, ready, &err)
	defer closeResource(server, &err)
	<-ready

	payload := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "message/send",
		"params": map[string]any{
			"message": map[string]any{
				"role": "user",
				"parts": []map[string]any{
					{
						"kind": "text",
						"text": "tell me a joke",
					},
				},
				"messageId": "9229e770-767c-417b-a0b0-f0741243c589",
			},
			"metadata": map[string]any{},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(fmt.Sprintf("http://localhost:%d/", port), "application/json", bytes.NewBuffer(body))

	require.NoError(t, err)
	defer closeResource(resp.Body, &err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type header should be application/json")
	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")
	jsonresp := string(bodyBytes)
	require.NoError(t, err, "Failed to decode JSON-RPC response")
	log.Debug("JSON-RPC response: %v", jsonresp)
	require.NoError(t, err, "Failed to decode JSON-RPC response")
	assert.Contains(t, jsonresp, `"jsonrpc":"2.0"`, "Response should contain JSON-RPC version")
	assert.Contains(t, jsonresp, `"text":"Why did the chicken cross the road? To get to the other side!"`, "Response should contain joke text")
	assert.Contains(t, jsonresp, `"messageId":"363422be-b0f9-4692-a24d-278670e7c7f1"`, "Response should contain message ID")
}

func startServer(server Server, ready chan struct{}, err *error) {
	if cerr := server.Start(ready); cerr != nil && *err == nil {
		*err = cerr
	}
}

func closeResource(resource io.Closer, err *error) {
	if cerr := resource.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}

func jokeHandler(msg *Message) (*Message, error) {
	log.Debug("Received message: %s", msg.MessageID)
	if len(msg.Parts) == 0 || msg.Parts[0]["text"] != "tell me a joke" {
		return nil, fmt.Errorf("unexpected message content")
	}
	response := &Message{
		Role: "agent",
		Parts: []map[string]any{
			{
				"kind": "text",
				"text": "Why did the chicken cross the road? To get to the other side!",
			},
		},
		MessageID: "363422be-b0f9-4692-a24d-278670e7c7f1",
	}
	return response, nil
}

func freePort() (int, error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	defer closeResource(l, &err)
	port := l.Addr().(*net.TCPAddr).Port
	return port, nil
}
