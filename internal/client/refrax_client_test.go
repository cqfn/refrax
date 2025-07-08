package client

import (
	"strings"
	"testing"

	"github.com/cqfn/refrax/internal/log"
	"github.com/stretchr/testify/assert"
)

func TestRefraxClient_Creates_Successfully(t *testing.T) {
	client := NewRefraxClient("none", "none", "")
	assert.NotNil(t, client, "Refrax client should not be nil")
}

func TestRefraxClient_Refactors_EmptyProject(t *testing.T) {
	client := NewRefraxClient("none", "none", "")
	origin := NewInMemoryProject(map[string]string{})

	proj, err := client.Refactor(origin, false, log.NewMock())

	assert.Equal(t, origin, proj, "Refactoring an empty project should return the same project")
	assert.Error(t, err, "Expected an error when refactoring an empty project")
	assert.Equal(t, "no java classes found in the project [empty project], add java files to the appropriate directory", err.Error(), "Error message should indicate no classes found")
}

func TestRefraxClient_PrintsStatsIfEnabled(t *testing.T) {
	client := NewRefraxClient("mock", "ABC", "")
	logger := log.NewMock()
	_, err := client.Refactor(SingleClassProject("Foo.java", "abstract class Foo {}"), true, logger)
	assert.NoError(t, err)
	assert.True(t, logMessageFoundWithText(logger, "Total messages asked"), "Expected total messages asked to be logged")
	assert.True(t, logMessageFoundWithText(logger, "Brain finished asking"), "Expected interaction stats to be logged")
}

func logMessageFoundWithText(logger *log.Mock, contains string) bool {
	found := false
	for _, msg := range logger.Messages {
		if strings.Contains(msg, contains) {
			found = true
			break
		}
	}
	return found
}
