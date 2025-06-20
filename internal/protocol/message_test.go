package protocol

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilePart_UnmarshalJSON(t *testing.T) {
	before := FilePart{
		Kind: PartKindFile,
		File: FileWithBytes{
			Bytes: "c29tZSB0ZXh0IGZpbGUgY29udGVudA==", // base64 for "some text file content"
		},
	}
	var after FilePart
	data, err := json.Marshal(before)
	require.NoError(t, err, "Failed to marshal FilePart")

	err = after.UnmarshalJSON(data)

	require.NoError(t, err, "Failed to unmarshal FilePart")
	assert.Equal(t, before, after, "data structures should match")
}
