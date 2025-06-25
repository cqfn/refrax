package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToken_LoadsTokenFromEnvFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, ".env")
	content := "TOKEN=test-token"
	err := os.WriteFile(file, []byte(content), 0644)
	require.NoError(t, err)

	token := Token(file)

	assert.Equal(t, "test-token", token)
}

func TestToken_EnvFileNotFoundUsesDefaultEnv(t *testing.T) {
	token := Token("nonexistent-file.env")

	assert.Equal(t, "", token)
}

func TestToken_EnvVariableNotSetReturnsEmptyString(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, ".env")
	content := "OTHER_VAR=other-value"
	err := os.WriteFile(file, []byte(content), 0644)
	require.NoError(t, err)

	token := Token(file)

	assert.Equal(t, "", token)
}

func TestToken_HandlesInvalidEnvFileGracefully(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, ".env")
	content := "\x00TOKEN=test-token"
	err := os.WriteFile(file, []byte(content), 0644)
	require.NoError(t, err)

	token := Token(file)

	assert.Equal(t, "", token)
}

func TestToken_EmptyEnvFileUsesDefaultEnv(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, ".env")
	err := os.WriteFile(file, []byte(""), 0644)
	require.NoError(t, err)

	token := Token(file)

	assert.Equal(t, "", token)
}
