package protocol

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCustomClient_SendsMessage(t *testing.T) {
	var err error
	serv, port, ready := testServer(t)
	closeResource(serv, &err)
	<-ready
	client := NewCustomClient(fmt.Sprintf("http://localhost:%d", port))
	message := MessageSendParams{
		Message: tellJoke(),
	}

	resp, err := client.SendMessage(message)

	expected := &JSONRPCResponse{
		ID:      "1",
		JSONRPC: "2.0",
		Result:  joke(),
	}
	require.NoError(t, err, "Failed to send message")
	require.NotNil(t, resp, "Response should not be nil")
	require.Equal(t, expected, resp, "Response text should match")

	// response, err := client.SendMessage(message)
	// require.NoError(t, err, "Failed to send message")
	// assert.NotNil(t, response, "Response should not be nil")
	// assert.Equal(t, "Hello, world!", response.Result.Parts[0]["text"], "Response text should match")
}
