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
	var err error
	serv, port, ready := testServer(t)
	defer closeResource(serv, &err)
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
	var err error
	serv, port, ready := testServer(t)
	defer closeResource(serv, &err)
	<-ready
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      float64(1),
		Method:  "message/send",
		Params: map[string]any{
			"message":  tellJoke(),
			"metadata": map[string]any{},
		},
	}
	body, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal request body")

	resp, err := http.Post(fmt.Sprintf("http://localhost:%d/", port), "application/json", bytes.NewBuffer(body))

	require.NoError(t, err)
	defer closeResource(resp.Body, &err)
	var response JSONRPCResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err, "Failed to decode response body")
	expected := JSONRPCResponse{
		ID:      "1",
		JSONRPC: "2.0",
		Result:  joke(),
	}
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type header should be application/json")
	assert.Equal(t, expected, response, "Server should return the expected joke message")
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func testServer(t *testing.T) (Server, int, chan struct{}) {
	t.Helper()
	port, err := freePort()
	require.NoError(t, err, "Failed to get a free port")
	server, err := NewCustomServer(testAgentCard, jokeHandler, port)
	require.NoError(t, err, "Failed to create custom server")
	ready := make(chan struct{})
	go startServer(server, ready, &err)
	return server, port, ready
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
	// log.Debug("Received message: %s", msg.MessageID)
	// if len(msg.Parts) == 0 || msg.Parts[0]["text"] != "tell me a joke" {
	// 	return nil, fmt.Errorf("unexpected message content, we expected 'tell me a joke', got: '%v'", msg.Parts[0]["text"])
	// }
	// response := joke()
	// return &response, nil
	log.Debug("Received message: %s", msg.MessageID)
	if len(msg.Parts) == 0 || msg.Parts[0].PartKind() != PartKindText || msg.Parts[0].(*TextPart).Text != "tell me a joke" {
		return nil, fmt.Errorf("unexpected message content, we expected 'tell me a joke', got: '%v'", msg.Parts[0].(*TextPart).Text)
	}
	response := joke()
	return &response, nil
}

func tellJoke() Message {
	// return Message{
	// 	Role: "user",
	// 	Parts: []map[string]any{
	// 		{
	// 			"kind": "text",
	// 			"text": "tell me a joke",
	// 		},
	// 	},
	// 	MessageID: "9229e770-767c-417b-a0b0-f0741243c589",
	// 	Kind:  KindMessage,
	// }
	return Message{
		Role: "user",
		Parts: []Part{
			&TextPart{
				Kind: PartKindText,
				Text: "tell me a joke",
			},
		},
		MessageID: "9229e770-767c-417b-a0b0-f0741243c589",
		Kind:      KindMessage,
	}
}

func joke() Message {
	// return Message{
	// 	Role: "agent",
	// 	Parts: []map[string]any{
	// 		{
	// 			"kind": "text",
	// 			"text": "Why did the chicken cross the road? To get to the other side!",
	// 		},
	// 	},
	// 	MessageID: "363422be-b0f9-4692-a24d-278670e7c7f1",
	// 	Kind:  KindMessage,
	// }
	return Message{
		Role: "agent",
		Parts: []Part{
			&TextPart{
				Kind: PartKindText,
				Text: "Why did the chicken cross the road? To get to the other side!",
			},
		},
		MessageID: "363422be-b0f9-4692-a24d-278670e7c7f1",
		Kind:      KindMessage,
	}
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
