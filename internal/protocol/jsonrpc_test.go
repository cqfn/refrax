package protocol

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONRPCResponse_UnmarshalMessage(t *testing.T) {
	before := JSONRPCResponse{
		ID: float64(1),
		Result: NewMessage().
			WithMessageID("363422be-b0f9-4692-a24d-278670e7c7f1").
			WithRole("agent").
			AddPart(NewText("Why did the chicken cross the road? To get to the other side!")),
	}
	var after JSONRPCResponse

	data, err := json.Marshal(before)

	require.NoError(t, err, "Failed to marshal JSONRPCResponse")
	err = json.Unmarshal(data, &after)
	require.NoError(t, err, "Failed to unmarshal JSONRPCResponse")
	assert.Equal(t, before, after, "data structures should match")
}

func TestJSONRPCResponse_UnmarshalTask(t *testing.T) {
	before := JSONRPCResponse{
		ID: float64(1),
		Result: &Task{
			ID:   "12345",
			Kind: KindTask,
			Status: TaskStatus{
				State: TaskStateCompleted,
			},
		},
	}
	var after JSONRPCResponse
	data, err := json.Marshal(before)
	require.NoError(t, err, "Failed to marshal JSONRPCResponse")

	err = json.Unmarshal(data, &after)

	require.NoError(t, err, "Failed to unmarshal JSONRPCResponse")
	assert.Equal(t, before, after, "data structures should match")
}
