package test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/cqfn/refrax/cmd"
	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEndToEnd_Agents_FromCLI_WithoutAI_WithEmptyProject(t *testing.T) {
	command := cmd.NewRootCmd(io.Discard, io.Discard)
	command.SetArgs([]string{"refactor", "--ai=none", t.TempDir()})

	err := command.Execute()

	require.Error(t, err, "expected command to fail with an empty project")
	assert.Contains(
		t,
		err.Error(),
		"no java classes found in the project",
		"Expected the output to indicate no AI provider was used and no classes were found",
	)
}

func TestEndToEnd_Agents_FromCLI_WithoutAI_WithMockProject(t *testing.T) {
	capture := buff()
	output := io.MultiWriter(capture, os.Stdout)
	command := cmd.NewRootCmd(output, io.Discard)
	command.SetArgs([]string{"refactor", "--ai=none", "--mock", "--debug", t.TempDir()})

	err := command.Execute()

	require.NoError(t, err, "Expected command to execute without error")
	assert.Contains(t, capture.String(), "provider: none", "expect no AI provider to be used in output")
	assert.Contains(t, capture.String(), "refactor result: [Main.java]", "expect refactored code to contain list of changed files")
}

func TestEndToEnd_JavaRefactor_InlineVariable_WithoutAI(t *testing.T) {
	const before = "public class Main {\n\tpublic static void main(String[] args) {\n\t\tString m = \"Hello, World\";\n\t\tSystem.out.println(m);\n\t}\n}\n\n"
	const expected = "public class Main {\n    public static void main(String[] args) {\n        System.out.println(\"Hello, World\");\n    }\n}"
	jclass := setupJava(t, t.TempDir(), "Main.java", before)
	capture := buff()
	output := io.MultiWriter(capture, os.Stdout)
	command := cmd.NewRootCmd(output, io.Discard)
	playbook := filepath.Join("test_data", "playbooks", "plain_main.yml")
	command.SetArgs([]string{"refactor", "--ai=none", "--debug", fmt.Sprintf("--playbook=%s", playbook), jclass})

	err := command.Execute()

	require.NoError(t, err, "Expected command to execute without error")
	assert.Contains(t, capture.String(), "provider: none", "expect no AI provider to be used in output")
	assertContent(t, jclass, expected)
}

func TestEndToEnd_JavaRefactor_ManyJavaFilesProject(t *testing.T) {
	tmp := t.TempDir()
	main, err := os.ReadFile(filepath.Join("test_data", "java", "person", "src", "com", "example", "MainApp.java"))
	require.NoError(t, err, "Expected to read test file content without error")
	mainFile := setupJava(t, filepath.Join(tmp, "person", "src", "com", "example"), "MainApp.java", string(main))

	person, err := os.ReadFile(filepath.Join("test_data", "java", "person", "src", "com", "example", "model", "Person.java"))
	require.NoError(t, err, "Expected to read test file content without error")
	personFile := setupJava(t, filepath.Join(tmp, "person", "src", "com", "example", "model"), "Person.java", string(person))

	service, err := os.ReadFile(filepath.Join("test_data", "java", "person", "src", "com", "example", "service", "GreetingService.java"))
	require.NoError(t, err, "Expected to read test file content without error")
	serviseFile := setupJava(t, filepath.Join(tmp, "person", "src", "com", "example", "service"), "GreetingService.java", string(service))

	capture := buff()
	output := io.MultiWriter(capture, os.Stdout)
	command := cmd.NewRootCmd(output, io.Discard)
	playbook := filepath.Join("test_data", "playbooks", "person.yml")
	command.SetArgs([]string{"refactor", "--ai=none", "--debug", fmt.Sprintf("--playbook=%s", playbook), tmp})

	err = command.Execute()

	require.NoError(t, err, "Expected command to execute without error")
	pb, err := brain.NewYAMLPlaybook(playbook)
	require.NoError(t, err, "Expected to load playbook without error")
	assertContent(t, mainFile, clean(pb.Ask("Fix 'MainApp'")))
	assertContent(t, personFile, clean(pb.Ask("Fix 'Person'")))
	assertContent(t, serviseFile, clean(pb.Ask("Fix 'GreetingService'")))
}

func TestEndToEnd_OuputOption_CopiesProject(t *testing.T) {
	tmp := t.TempDir()
	project := filepath.Join("test_data", "java", "person")

	capture := buff()
	output := io.MultiWriter(capture, os.Stdout)
	command := cmd.NewRootCmd(output, io.Discard)
	command.SetArgs([]string{"refactor", "--ai=none", "--debug", "--output=" + tmp, project})

	err := command.Execute()

	require.NoError(t, err, "Expected command to execute without error")
	assert.FileExists(t, filepath.Join(tmp, "src", "com", "example", "MainApp.java"), "Expected MainApp.java to be copied to output directory")
	assert.FileExists(t, filepath.Join(tmp, "src", "com", "example", "model", "Person.java"), "Expected Person.java to be copied to output directory")
	assert.FileExists(t, filepath.Join(tmp, "src", "com", "example", "service", "GreetingService.java"), "Expected GreetingService.java to be copied to output directory")
}

func setupJava(t *testing.T, path, name, code string) string {
	t.Helper()
	full := filepath.Clean(path)
	err := os.MkdirAll(full, 0o700)
	log.Info("create test directories at %s", full)
	require.NoError(t, err, "Expected to create test directories correctly")
	java := filepath.Join(full, name)
	err = os.WriteFile(java, []byte(code), 0o600)
	require.NoError(t, err, "Expected to write mock project file without error")
	return java
}

func assertContent(t *testing.T, path, expected string) {
	t.Helper()
	content, err := os.ReadFile(filepath.Clean(path))
	require.NoError(t, err, "Expected to read file content without error")
	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(string(content)), "File content does not match expected content")
}

func buff() *safeBuffer {
	return &safeBuffer{buf: &bytes.Buffer{}}
}

type safeBuffer struct {
	mu  sync.Mutex
	buf *bytes.Buffer
}

func (s *safeBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *safeBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

func clean(answer string) string {
	answer = strings.ReplaceAll(answer, "```java", "")
	return strings.ReplaceAll(answer, "```", "")
}
