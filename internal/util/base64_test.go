package util

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeFile_ValidBase64(t *testing.T) {
	original := "Hëllö Wörld"
	encoded := base64.StdEncoding.EncodeToString([]byte(original))
	result, err := DecodeFile(encoded)
	require.NoError(t, err, "DecodeFile unexpectedly returned an error")
	assert.Equal(t, original, result, "Decoded content must match the original")
}

func TestDecodeFile_EmptyString(t *testing.T) {
	result, err := DecodeFile("")
	require.NoError(t, err, "DecodeFile unexpectedly returned an error for empty input")
	assert.Empty(t, result, "Decoded empty input must be empty")
}

func TestDecodeFile_UnicodeContent(t *testing.T) {
	original := "Привет мир 日本語 🎉"
	encoded := base64.StdEncoding.EncodeToString([]byte(original))
	result, err := DecodeFile(encoded)
	require.NoError(t, err, "DecodeFile unexpectedly returned an error")
	assert.Equal(t, original, result, "Decoded unicode content must match the original")
}

func TestDecodeFile_InvalidBase64(t *testing.T) {
	_, err := DecodeFile("nöt-välïd-bàsé64!!!")
	assert.Error(t, err, "DecodeFile must return an error for invalid base64")
}

func TestDecodeFile_MultilineContent(t *testing.T) {
	original := "Lïnë önë\nLïnë twö\nLïnë thrëë"
	encoded := base64.StdEncoding.EncodeToString([]byte(original))
	result, err := DecodeFile(encoded)
	require.NoError(t, err, "DecodeFile unexpectedly returned an error")
	assert.Equal(t, original, result, "Decoded multiline content must match the original")
}

func TestDecodeFile_BinaryContent(t *testing.T) {
	original := string([]byte{0x00, 0xFF, 0x7F, 0x80})
	encoded := base64.StdEncoding.EncodeToString([]byte(original))
	result, err := DecodeFile(encoded)
	require.NoError(t, err, "DecodeFile unexpectedly returned an error")
	assert.Equal(t, original, result, "Decoded binary content must match the original")
}
