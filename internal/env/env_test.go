package env

import (
	"fmt"
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
	err := os.WriteFile(file, []byte(content), 0o600)
	require.NoError(t, err)

	token := Token(file, "otherprovider")

	assert.Equal(t, "test-token", token)
}

func TestToken_EnvFileNotFoundUsesDefaultEnv(t *testing.T) {
	token := Token("nonexistent-file.env", "unknown")

	assert.Equal(t, "", token)
}

func TestToken_EnvVariableNotSetReturnsEmptyString(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, ".env")
	content := "OTHER_VAR=other-value"
	err := os.WriteFile(file, []byte(content), 0o600)
	require.NoError(t, err)

	token := Token(file, "otherprovider")

	assert.Equal(t, "", token)
}

func TestToken_HandlesInvalidEnvFileGracefully(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, ".env")
	content := "\x00TOKEN=test-token"
	err := os.WriteFile(file, []byte(content), 0o600)
	require.NoError(t, err)

	token := Token(file, "unknown")

	assert.Equal(t, "", token)
}

func TestToken_EmptyEnvFileUsesDefaultEnv(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, ".env")
	err := os.WriteFile(file, []byte(""), 0o600)
	require.NoError(t, err)

	token := Token(file, "unknown")

	assert.Equal(t, "", token)
}

func TestProviderToken_DeepseekTokenPresent(t *testing.T) {
	tmp := t.TempDir()
	deepseek := "deepseek-token-value"
	env := filepath.Join(tmp, ".env")
	err := os.WriteFile(env, fmt.Appendf(nil, "DEEPSEEK_API_KEY=%s", deepseek), 0o600)
	require.NoError(t, err)

	result := Token(env, "deepseek")

	assert.Equal(t, deepseek, result)
}

func TestProviderToken_OpenAITokenPresent(t *testing.T) {
	tmp := t.TempDir()
	openai := "openai-token-value"
	env := filepath.Join(tmp, ".env")
	err := os.WriteFile(env, fmt.Appendf(nil, "OPENAI_API_KEY=%s", openai), 0o600)
	require.NoError(t, err)

	result := Token(env, "openai")

	assert.Equal(t, openai, result)
}

func TestProviderToken_DefaultTokenPresent(t *testing.T) {
	tmp := t.TempDir()
	token := "default-token-value"
	env := filepath.Join(tmp, ".env")
	err := os.WriteFile(env, fmt.Appendf(nil, "TOKEN=%s", token), 0o600)
	require.NoError(t, err)

	result := Token(env, "otherprovider")
	assert.Equal(t, token, result)
}

func TestProviderToken_DeepseekTokenAbsent(t *testing.T) {
	t.Setenv("DEEPSEEK_API_KEY", "")

	tempDir := t.TempDir()

	result := Token(tempDir, "deepseek")
	assert.Equal(t, "", result)
}

func TestProviderToken_OpenAITokenAbsent(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")

	tempDir := t.TempDir()

	result := Token(tempDir, "openai")
	assert.Equal(t, "", result)
}

func TestProviderToken_DefaultTokenAbsent(t *testing.T) {
	tempDir := t.TempDir()

	result := Token(tempDir, "otherprovider")
	assert.Equal(t, "", result)
}

func TestProviderToken_DefaultTokenIgnoredWhenDeepseekTokenExists(t *testing.T) {
	tmp := t.TempDir()
	deepseek := "deepseek-token-val"
	token := "default-ignored-deepseek-token-val"
	env := filepath.Join(tmp, ".env")
	err := os.WriteFile(env, fmt.Appendf(nil, "DEEPSEEK_API_KEY=%s\nTOKEN=%s\n", deepseek, token), 0o600)
	require.NoError(t, err)

	result := Token(env, "deepseek")

	assert.Equal(t, deepseek, result)
}

func TestProviderToken_DefaultTokenIgnoredWhenOpenAITokenExists(t *testing.T) {
	tmp := t.TempDir()
	openai := "openai-token-val"
	token := "default-openai-token-val"
	env := filepath.Join(tmp, ".env")
	err := os.WriteFile(env, fmt.Appendf(nil, "OPENAI_API_KEY=%s\nTOKEN=%s\n", openai, token), 0o600)
	require.NoError(t, err)

	result := Token(env, "openai")

	assert.Equal(t, openai, result)
}

func TestProviderToken_DefaultTokenUsedWhenDeepseekNotRequested(t *testing.T) {
	tmp := t.TempDir()
	deepseek := "deepseek-token"
	token := "default-used-deepseek-token"
	env := filepath.Join(tmp, ".env")
	err := os.WriteFile(env, fmt.Appendf(nil, "DEEPSEEK_API_KEY=%s\nTOKEN=%s\n", deepseek, token), 0o600)
	require.NoError(t, err)

	result := Token(env, "otherprovider")

	assert.Equal(t, token, result)
}

func TestProviderToken_DefaultTokenUsedWhenOpenAINotRequested(t *testing.T) {
	tmp := t.TempDir()
	openai := "openai-token"
	token := "default-openai-token"
	env := filepath.Join(tmp, ".env")
	err := os.WriteFile(env, fmt.Appendf(nil, "OPENAI_API_KEY=%s\nTOKEN=%s\n", openai, token), 0o600)
	require.NoError(t, err)

	result := Token(env, "otherprovider")

	assert.Equal(t, token, result)
}

func TestProviderToken_ReadsMockTokenFromEnvironmentVariable(t *testing.T) {
	t.Setenv("MOCK_TOKEN", "mock-token-value")

	token := Token("", "mock")

	assert.NotEmpty(t, token, "expected mock token to be set")
	assert.Equal(t, "mock-token-value", token, "expected mock token to match environment variable")
}
