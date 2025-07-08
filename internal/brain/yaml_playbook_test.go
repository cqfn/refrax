package brain

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const content = `
name: qa-basic-test
qa:
  - question: |
      What should I rename the variable ` + "`tmp`" + ` to?
    answer: |
      You can rename ` + "`tmp`" + ` to ` + "`userID`" + ` for better clarity.
  - question: |
      Should I add a null check for ` + "`user.getName()`" + `?
    answer: |
      Yes, you should check if ` + "`user`" + ` is null before calling ` + "`getName()`" + ` to avoid a NullPointerException.
  - question: |
      Is there a simpler way to write this loop?
    answer: |
      Yes, consider using a ` + "`forEach`" + ` method or a stream to simplify the loop.
`

func TestNewYAMLPlaybook_Success(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "playbook_example.yml")
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	pb, err := NewYAMLPlaybook(filePath)
	require.NoError(t, err)
	require.NotNil(t, pb)

	fmt.Println("Playbook Name:", pb)
	assert.Equal(t, "You can rename `tmp` to `userID` for better clarity.", pb.Ask("What should I rename the variable `tmp` to?"))
	assert.Equal(t, "Yes, you should check if `user` is null before calling `getName()` to avoid a NullPointerException.", pb.Ask("Should I add a null check for `user.getName()`?"))
	assert.Equal(t, "Yes, consider using a `forEach` method or a stream to simplify the loop.", pb.Ask("Is there a simpler way to write this loop?"))
}

func TestNewYAMLPlaybook_FileNotFound(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "nonexistent.yml")

	pb, err := NewYAMLPlaybook(filePath)
	require.Error(t, err)
	assert.Nil(t, pb)
}

func TestNewYAMLPlaybook_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "invalid.yml")
	content := strings.Join([]string{
		"name: qa-basic-test",
		"qa:",
		"  - question: \"What should I rename the variable 'tmp' to?\"",
		"    answer: [\"This should fail because it's not a string\"]",
	}, "\n")
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	pb, err := NewYAMLPlaybook(filePath)
	require.Error(t, err)
	assert.Nil(t, pb)
}

func TestYAMLPlaybook_Ask_QuestionNotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "playbook_example.yml")
	content := "name: qa-basic-test\n" +
		"qa:\n" +
		"  - question: |\n" +
		"      What should I rename the variable `tmp` to?\n" +
		"    answer: |\n" +
		"      You can rename `tmp` to `userID` for better clarity.\n"
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)

	pb, err := NewYAMLPlaybook(path)
	require.NoError(t, err)
	require.NotNil(t, pb)

	assert.Equal(t, "Question not found in the playbook", pb.Ask("This question does not exist"))
}

func TestYAMLPlaybook_ReadsExistingFile(t *testing.T) {
	pb, err := NewYAMLPlaybook("playbook_example.yml")
	require.NoError(t, err, "Failed to read the playbook file")
	require.NotNil(t, pb, "Playbook should not be nil")

	answer := pb.Ask("Should I add a null check for `user.getName()`?")

	assert.Equal(t, "Yes, you should check if `user` is null before calling `getName()` to avoid a NullPointerException.", answer, "The answer should match the expected response")
}
