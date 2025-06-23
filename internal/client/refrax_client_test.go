package client

import (
	"testing"

	"github.com/cqfn/refrax/internal/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefraxClient_Creates_Successfully(t *testing.T) {
	client := NewRefraxClient("none", "none")
	assert.NotNil(t, client, "Refrax client should not be nil")
}

func TestRefraxClient_Refactors_Successfully(t *testing.T) {
	client := NewRefraxClient("none", "none")

	result, err := client.Refactor(SingleClassProject("Main.java", before))

	log.Debug("Refactored project: %s", result)
	require.NoError(t, err, "Expected no error during refactoring")
	classes, err := result.Classes()
	require.NoError(t, err, "Expected no error retrieving classes from refactored project")
	class := classes[0]
	assert.Equal(t, "Main.java", class.Name(), "Class name should match")
	assert.Equal(t, after, class.Content(), "Class content should match expected refactored content")
}

const (
	before = `public class Main {
	public static void main(String[] args) {
		String m = "Hello, World";
		System.out.println(m);
	}
}
`

	after = `public class Main {
	public static void main(String[] args) {
		System.out.println("Hello, World");
	}
`
)
