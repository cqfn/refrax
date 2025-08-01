package protocol

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_SendsMessage(t *testing.T) {
	var err error
	serv, port := testServer(t)
	<-serv.Ready()
	client := NewClient(fmt.Sprintf("http://localhost:%d", port)).(*a2aClient)
	client.id = func() string { return "1" }
	message := MessageSendParams{
		Message: askJoke(),
	}

	resp, err := client.SendMessage(&message)

	require.NoError(t, err, "Failed to send message")
	err = serv.Shutdown()
	require.NoError(t, err, "Failed to close server")
	expected := &JSONRPCResponse{
		ID:      "1",
		JSONRPC: "2.0",
		Result:  tellJoke(),
	}
	require.NoError(t, err, "Failed to send message")
	require.NotNil(t, resp, "Response should not be nil")
	require.Equal(t, expected, resp, "Response text should match")
}
