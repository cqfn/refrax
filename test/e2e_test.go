package test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/cqfn/refrax/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEndToEnd_Agents_FromCLI_WithoutAI_WithMockProject(t *testing.T) {
	capture := &bytes.Buffer{}
	output := io.MultiWriter(capture, os.Stdout)
	cmd := cmd.NewRootCmd(output, io.Discard)
	cmd.SetArgs([]string{"refactor", "--ai=none", "--mock"})

	err := cmd.Execute()

	require.NoError(t, err, "Expected command to execute without error")
	assert.Contains(t, capture.String(), "provider: none", "expect no AI provider to be used in output")
	assert.Contains(t, capture.String(), "System.out.println(\"Hello, World\")", "expect refactored code to contain 'Hello, World'")
}

func TestEndToEnd_JavaRefactor_InlineVariable_WithoutAI(t *testing.T) {
	const before = "public class Main {\n\tpublic static void main(String[] args) {\n\t\tString m = \"Hello, World\";\n\t\tSystem.out.println(m);\n\t}\n}\n\n"
	const expected = "public class Main {\n\tpublic static void main(String[] args) {\n\t\tSystem.out.println(\"Hello, World\");\n\t}\n"
	project := setupProject(t, before)
	capture := &bytes.Buffer{}
	output := io.MultiWriter(capture, os.Stdout)
	cmd := cmd.NewRootCmd(output, io.Discard)
	cmd.SetArgs([]string{"refactor", "--ai=none", project})

	err := cmd.Execute()

	require.NoError(t, err, "Expected command to execute without error")
	assert.Contains(t, capture.String(), "provider: none", "expect no AI provider to be used in output")
	assert.Contains(t, capture.String(), "System.out.println(\"Hello, World\")", "expect refactored code to contain inlined variable")
	assertContent(t, project, expected)
}

func setupProject(t *testing.T, code string) string {
	t.Helper()
	tmp := t.TempDir()
	java := filepath.Join(tmp, "Main.java")
	err := os.WriteFile(java, []byte(code), 0644)
	require.NoError(t, err, "Expected to write mock project file without error")
	return java
}

func assertContent(t *testing.T, path string, expected string) {
	t.Helper()
	content, err := os.ReadFile(path)
	require.NoError(t, err, "Expected to read file content without error")
	assert.Equal(t, expected, string(content), "File content does not match expected content")
}
