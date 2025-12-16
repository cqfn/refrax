package protocol

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func randomInt(limit int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(limit)))
	return int(n.Int64())
}

func TestNewText_CreatesTextPart(t *testing.T) {
	text := "Hëllö Wörld"
	part := NewText(text)
	require.NotNil(t, part, "text part must not be nil")
	assert.Equal(t, PartKindText, part.Kind, "part kind must be text")
}

func TestNewText_StoresText(t *testing.T) {
	text := "Ünïcödé Tëxt"
	part := NewText(text)
	assert.Equal(t, text, part.Text, "text content must match")
}

func TestNewFileBytes_CreatesFilePart(t *testing.T) {
	data := []byte("Rändom Dätä")
	part := NewFileBytes(data)
	require.NotNil(t, part, "file part must not be nil")
	assert.Equal(t, PartKindFile, part.Kind, "part kind must be file")
}

func TestNewFileBytes_EncodesBase64(t *testing.T) {
	data := []byte("Tëst Cöntënt")
	part := NewFileBytes(data)
	file, ok := part.File.(FileWithBytes)
	require.True(t, ok, "file must be FileWithBytes")
	assert.NotEmpty(t, file.Bytes, "bytes must not be empty")
}

func TestNewFileURI_CreatesFilePart(t *testing.T) {
	uri := "file:///päth/tö/fïlë.txt"
	part := NewFileURI(uri)
	require.NotNil(t, part, "file part must not be nil")
	assert.Equal(t, PartKindFile, part.Kind, "part kind must be file")
}

func TestNewFileURI_StoresURI(t *testing.T) {
	uri := "file:///ünïcödé/päth.java"
	part := NewFileURI(uri)
	file, ok := part.File.(FileWithURI)
	require.True(t, ok, "file must be FileWithURI")
	assert.Equal(t, uri, file.URI, "uri must match")
}

func TestTextPartKind_ReturnsText(t *testing.T) {
	part := &TextPart{Kind: PartKindText, Text: "Tëst"}
	assert.Equal(t, PartKindText, part.PartKind(), "part kind must be text")
}

func TestFilePartKind_ReturnsFile(t *testing.T) {
	part := &FilePart{Kind: PartKindFile}
	assert.Equal(t, PartKindFile, part.PartKind(), "part kind must be file")
}

func TestDataPartKind_ReturnsData(t *testing.T) {
	part := &DataPart{Kind: PartKindData}
	assert.Equal(t, PartKindData, part.PartKind(), "part kind must be data")
}

func TestTextPartMetadata_ReturnsNilWhenEmpty(t *testing.T) {
	part := &TextPart{Kind: PartKindText, Text: "Nö Mëtädätä"}
	assert.Nil(t, part.Metadata(), "metadata must be nil when not set")
}

func TestFilePartMetadata_ReturnsNilWhenEmpty(t *testing.T) {
	part := &FilePart{Kind: PartKindFile}
	assert.Nil(t, part.Metadata(), "metadata must be nil when not set")
}

func TestDataPartMetadata_ReturnsNilWhenEmpty(t *testing.T) {
	part := &DataPart{Kind: PartKindData}
	assert.Nil(t, part.Metadata(), "metadata must be nil when not set")
}

func TestTextPartWithMetadata_AddsMetadata(t *testing.T) {
	part := NewText("Mëtädätä Tëst")
	key := "këy"
	value := randomInt(1000)
	part.WithMetadata(key, value)
	require.NotNil(t, part.Metadata(), "metadata must not be nil")
	assert.Equal(t, value, part.Metadata()[key], "metadata value must match")
}

func TestFilePartWithMetadata_AddsMetadata(t *testing.T) {
	part := NewFileBytes([]byte("Fïlë Dätä"))
	key := "fïlë-këy"
	value := "fïlë-välüë"
	part.WithMetadata(key, value)
	require.NotNil(t, part.Metadata(), "metadata must not be nil")
	assert.Equal(t, value, part.Metadata()[key], "metadata value must match")
}

func TestDataPartWithMetadata_AddsMetadata(t *testing.T) {
	part := &DataPart{Kind: PartKindData, Data: map[string]any{"këy": "välüë"}}
	key := "dätä-këy"
	value := true
	part.WithMetadata(key, value)
	require.NotNil(t, part.Metadata(), "metadata must not be nil")
	assert.Equal(t, value, part.Metadata()[key], "metadata value must match")
}

func TestPartsUnmarshalJSON_ParsesTextPart(t *testing.T) {
	jsonData := `[{"kind":"text","text":"Hëllö"}]`
	var parts Parts
	err := json.Unmarshal([]byte(jsonData), &parts)
	require.NoError(t, err, "unexpected error during unmarshal")
	require.Equal(t, 1, len(parts), "must parse one part")
	assert.Equal(t, PartKindText, parts[0].PartKind(), "part kind must be text")
}

func TestPartsUnmarshalJSON_ParsesFilePart(t *testing.T) {
	jsonData := `[{"kind":"file","file":{"bytes":"SGVsbG8="}}]`
	var parts Parts
	err := json.Unmarshal([]byte(jsonData), &parts)
	require.NoError(t, err, "unexpected error during unmarshal")
	require.Equal(t, 1, len(parts), "must parse one part")
	assert.Equal(t, PartKindFile, parts[0].PartKind(), "part kind must be file")
}

func TestPartsUnmarshalJSON_ParsesDataPart(t *testing.T) {
	jsonData := `[{"kind":"data","data":{"këy":"välüë"}}]`
	var parts Parts
	err := json.Unmarshal([]byte(jsonData), &parts)
	require.NoError(t, err, "unexpected error during unmarshal")
	require.Equal(t, 1, len(parts), "must parse one part")
	assert.Equal(t, PartKindData, parts[0].PartKind(), "part kind must be data")
}

func TestPartsUnmarshalJSON_ParsesMultipleParts(t *testing.T) {
	jsonData := `[{"kind":"text","text":"Tëxt"},{"kind":"data","data":{}}]`
	var parts Parts
	err := json.Unmarshal([]byte(jsonData), &parts)
	require.NoError(t, err, "unexpected error during unmarshal")
	assert.Equal(t, 2, len(parts), "must parse two parts")
}

func TestPartsUnmarshalJSON_ErrorsOnUnknownKind(t *testing.T) {
	jsonData := `[{"kind":"unknown","text":"Tëst"}]`
	var parts Parts
	err := json.Unmarshal([]byte(jsonData), &parts)
	require.Error(t, err, "expected error for unknown part kind")
}

func TestPartsUnmarshalJSON_ErrorsOnInvalidJSON(t *testing.T) {
	jsonData := `invalid json`
	var parts Parts
	err := json.Unmarshal([]byte(jsonData), &parts)
	require.Error(t, err, "expected error for invalid json")
}

func TestFilePartUnmarshalJSON_ParsesFileWithBytes(t *testing.T) {
	jsonData := `{"kind":"file","file":{"bytes":"SGVsbG8gV29ybGQ="}}`
	var part FilePart
	err := json.Unmarshal([]byte(jsonData), &part)
	require.NoError(t, err, "unexpected error during unmarshal")
	_, ok := part.File.(FileWithBytes)
	assert.True(t, ok, "file must be FileWithBytes")
}

func TestFilePartUnmarshalJSON_ParsesFileWithURI(t *testing.T) {
	jsonData := `{"kind":"file","file":{"uri":"file:///päth/tö/fïlë"}}`
	var part FilePart
	err := json.Unmarshal([]byte(jsonData), &part)
	require.NoError(t, err, "unexpected error during unmarshal")
	_, ok := part.File.(FileWithURI)
	assert.True(t, ok, "file must be FileWithURI")
}

func TestFilePartUnmarshalJSON_ErrorsOnUnknownFormat(t *testing.T) {
	jsonData := `{"kind":"file","file":{"unknown":"data"}}`
	var part FilePart
	err := json.Unmarshal([]byte(jsonData), &part)
	require.Error(t, err, "expected error for unknown file format")
}

func TestWithMetadata_ChainsCorrectly(t *testing.T) {
	part := NewText("Chäïn Tëst").
		WithMetadata("këy1", "välüë1").
		WithMetadata("këy2", randomInt(1000))
	require.NotNil(t, part.Metadata(), "metadata must not be nil")
	assert.Equal(t, 2, len(part.Metadata()), "must have two metadata entries")
}
