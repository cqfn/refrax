package protocol

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCustomClient_SendsMessage(t *testing.T) {
	var err error
	serv, port, ready := testServer(t)
	<-ready
	client := NewCustomClient(fmt.Sprintf("http://localhost:%d", port))
	message := MessageSendParams{
		Message: askJoke(),
	}

	resp, err := client.SendMessage(message)

	require.NoError(t, err, "Failed to send message")
	err = serv.Close()
	require.NoError(t, err, "Failed to close server")
	expected := &JSONRPCResponse{
		ID:      "1",
		JSONRPC: "2.0",
		Result:  *tellJoke(),
	}
	require.NoError(t, err, "Failed to send message")
	require.NotNil(t, resp, "Response should not be nil")
	require.Equal(t, expected, resp, "Response text should match")
}
