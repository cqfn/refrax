package client

import (
	"bytes"
	"io"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefraxClient_Creates_Successfully(t *testing.T) {
	client := NewRefraxClient(NewMockParams())
	assert.NotNil(t, client, "Refrax client should not be nil")
}

func TestRefraxClient_Refactors_EmptyProject(t *testing.T) {
	client := NewRefraxClient(NewMockParams())
	origin := NewInMemoryProject(map[string]string{})

	proj, err := client.Refactor(origin)

	assert.Equal(t, origin, proj, "Refactoring an empty project should return the same project")
	assert.Error(t, err, "Expected an error when refactoring an empty project")
	assert.Equal(t, "no java classes found in the project [empty project], add java files to the appropriate directory", err.Error(), "Error message should indicate no classes found")
}

func TestRefraxClient_PrintsStatsIfEnabled(t *testing.T) {
	params := NewMockParams()
	params.Stats = true
	out := bytes.Buffer{}
	params.Log = NewSyncWriter(io.Writer(&out))
	client := NewRefraxClient(params)
	_, err := client.Refactor(SingleClassProject("Foo.java", "abstract class Foo {}"))
	assert.NoError(t, err)
	assert.Contains(t, out.String(), "Total messages asked", "Expected total messages asked to be logged")
	assert.Contains(t, out.String(), "Brain finished asking", "Expected interaction stats to be logged")
}

type SyncWriter struct {
	mu sync.Mutex
	w  io.Writer
}

func NewSyncWriter(w io.Writer) *SyncWriter {
	return &SyncWriter{
		mu: sync.Mutex{},
		w:  w,
	}
}

func (s *SyncWriter) Write(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.w.Write(p)
}
