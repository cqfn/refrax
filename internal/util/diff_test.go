package util

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiff_InsertionsInMiddle(t *testing.T) {
	before := "ab\nb\nd\ne"
	after := "ab\nb\nc\nd\ne"
	result := Diff(before, after)
	assert.Equal(t, 1, result, "One line inserted in the middle")
}

func TestDiff_DeletionsInMiddle(t *testing.T) {
	before := "ad\nb\nc\nd\ne"
	after := "ad\nb\nd\ne"
	result := Diff(before, after)
	assert.Equal(t, 1, result, "One line deleted in the middle")
}

func TestDiff_CommonPrefixOnly(t *testing.T) {
	before := "ah\nb\nx\ny"
	after := "ah\nb\nc\nd"
	result := Diff(before, after)
	assert.Equal(t, 4, result, "Two lines were removed and two lines were added")
}

func TestDiff_CommonSuffixOnly(t *testing.T) {
	before := "x\ny\na\nb"
	after := "z\nw\na\nb"
	result := Diff(before, after)
	assert.Equal(t, 4, result, "Two lines were removed and two lines were added")
}

func TestDiff_InterleavedDifferences(t *testing.T) {
	before := "ai\nx\nb\ny\nc"
	after := "ai\nb\nc"
	result := Diff(before, after)
	assert.Equal(t, 2, result, "Interleaved differences should be counted correctly")
}

func TestDiff_MovedLines(t *testing.T) {
	before := "a\nb\nc\nd"
	after := "c\na\nd\nb"
	result := Diff(before, after)
	assert.Equal(t, 4, result, "LCS doesn't match moved lines unless in same order")
}

func TestDiff_OneLineChangedOutOfMany(t *testing.T) {
	before := "a\nb\nc\nd\ne"
	after := "a\nb\nX\nd\ne"
	result := Diff(before, after)
	assert.Equal(t, 2, result, "One line was added and one was removed")
}

func TestDiff_AllDifferentLines(t *testing.T) {
	before := "a\nb\nc"
	after := "x\ny\nz"
	result := Diff(before, after)
	assert.Equal(t, 6, result, "All lines changed: 3 removed, 3 added")
}

func TestDiff_OneLineAdded(t *testing.T) {
	before := "line1\nline2\nline3"
	after := "line1\nline2\nnew line\nline3"
	result := Diff(before, after)
	assert.Equal(t, 1, result, "Only one line was added")
}

func TestDiff_AddOneLineToEmptyBefore(t *testing.T) {
	before := ""
	after := "hello"
	result := Diff(before, after)
	assert.Equal(t, 2, result, "Adding one line to an empty 'before' should count as exactly one change")
}

func TestDiff_RemoveOneLineToEmptyAfter(t *testing.T) {
	before := "bye"
	after := ""
	result := Diff(before, after)
	assert.Equal(t, 2, result, "Removing the only line should count as exactly one change")
}

func TestDiff_RealisticCodeChange(t *testing.T) {
	before := `func add(a, b int) int {
    return a + b
}`

	after := `func add(a, b int) int {
    result := a + b
    return result
}`

	result := Diff(before, after)
	assert.Equal(t, 3, result, "Two lines were added inside function body, and one line was changed")
}

func TestDiff_FromTestDataFiles(t *testing.T) {
	root := "test_data"
	beforePath := filepath.Clean(filepath.Join(root, "Before.java"))
	afterPath := filepath.Clean(filepath.Join(root, "After.java"))
	beforeBytes, err := os.ReadFile(beforePath)
	if err != nil {
		t.Fatalf("failed to read %q: %v", beforePath, err)
	}
	afterBytes, err := os.ReadFile(afterPath)
	if err != nil {
		t.Fatalf("failed to read %q: %v", afterPath, err)
	}
	before := string(beforeBytes)
	after := string(afterBytes)
	result := Diff(before, after)
	assert.Equal(t, 7, result, "Expected exactly 6 changed lines between %q and %q", beforePath, afterPath)
}
