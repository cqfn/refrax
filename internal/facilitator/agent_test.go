package facilitator

import (
	"fmt"
	"sync"
	"testing"

	"github.com/cqfn/refrax/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockCritic is a mock implementation of the domain.Critic interface
type mockCritic struct {
	reviewFunc func(job *domain.Job) (*domain.Artifacts, error)
}

func (m *mockCritic) Review(job *domain.Job) (*domain.Artifacts, error) {
	if m.reviewFunc != nil {
		return m.reviewFunc(job)
	}
	return &domain.Artifacts{
		Descr:       &domain.Description{Text: "mock critique"},
		Suggestions: []domain.Suggestion{},
	}, nil
}

// mockLogger is a mock implementation of the log.Logger interface
type mockLogger struct {
	mu   sync.Mutex
	logs []string
}

func (m *mockLogger) Info(msg string, args ...any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logs = append(m.logs, fmt.Sprintf("INFO: "+msg, args...))
}

func (m *mockLogger) Debug(msg string, args ...any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logs = append(m.logs, fmt.Sprintf("DEBUG: "+msg, args...))
}

func (m *mockLogger) Warn(msg string, args ...any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logs = append(m.logs, fmt.Sprintf("WARN: "+msg, args...))
}

func (m *mockLogger) Error(msg string, args ...any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logs = append(m.logs, fmt.Sprintf("ERROR: "+msg, args...))
}

func (m *mockLogger) GetLogs() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	logs := make([]string, len(m.logs))
	copy(logs, m.logs)
	return logs
}

func TestCriticizeAll_WithValidClasses(t *testing.T) {
	// Arrange
	logger := &mockLogger{}
	critic := &mockCritic{
		reviewFunc: func(job *domain.Job) (*domain.Artifacts, error) {
			return &domain.Artifacts{
				Descr: &domain.Description{Text: "mock critique"},
				Suggestions: []domain.Suggestion{
					*domain.NewSuggestion("This is a suggestion", job.Classes[0].Path()),
				},
			}, nil
		},
	}

	agent := &agent{
		log:    logger,
		critic: critic,
	}

	classes := []domain.Class{
		domain.NewInMemoryClass("Class1.java", "/path/Class1.java", "public class Class1 {}"),
		domain.NewInMemoryClass("Class2.java", "/path/Class2.java", "public class Class2 {}"),
	}

	// Act
	result, err := agent.criticizeAll(classes, 1000)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Len(t, result[0].suggestions, 1)
	assert.Len(t, result[1].suggestions, 1)
	assert.Equal(t, "Class1.java", result[0].class.Name())
	assert.Equal(t, "Class2.java", result[1].class.Name())
}

func TestCriticizeAll_WithNoSuggestions(t *testing.T) {
	// Arrange
	logger := &mockLogger{}
	critic := &mockCritic{
		reviewFunc: func(job *domain.Job) (*domain.Artifacts, error) {
			return &domain.Artifacts{
				Descr:       &domain.Description{Text: "mock critique"},
				Suggestions: []domain.Suggestion{},
			}, nil
		},
	}

	agent := &agent{
		log:    logger,
		critic: critic,
	}

	classes := []domain.Class{
		domain.NewInMemoryClass("Class1.java", "/path/Class1.java", "public class Class1 {}"),
	}

	// Act
	result, err := agent.criticizeAll(classes, 1000)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Len(t, result[0].suggestions, 0)
	assert.Equal(t, "Class1.java", result[0].class.Name())
}

func TestCriticizeAll_WithErrorInCritic(t *testing.T) {
	// Arrange
	logger := &mockLogger{}
	critic := &mockCritic{
		reviewFunc: func(job *domain.Job) (*domain.Artifacts, error) {
			return nil, fmt.Errorf("critic error")
		},
	}

	agent := &agent{
		log:    logger,
		critic: critic,
	}

	classes := []domain.Class{
		domain.NewInMemoryClass("Class1.java", "/path/Class1.java", "public class Class1 {}"),
	}

	// Act
	result, err := agent.criticizeAll(classes, 1000)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to review class")
}

func TestCriticizeAll_WithLargeClass_SkipsReview(t *testing.T) {
	// Arrange
	logger := &mockLogger{}
	critic := &mockCritic{
		reviewFunc: func(job *domain.Job) (*domain.Artifacts, error) {
			return &domain.Artifacts{
				Descr: &domain.Description{Text: "mock critique"},
				Suggestions: []domain.Suggestion{
					*domain.NewSuggestion("This is a suggestion", job.Classes[0].Path()),
				},
			}, nil
		},
	}

	agent := &agent{
		log:    logger,
		critic: critic,
	}

	// Create a large class content that exceeds MAX_TOKENS
	largeContent := ""
	for i := 0; i < 10000; i++ {
		largeContent += fmt.Sprintf("line %d: some content\n", i)
	}

	classes := []domain.Class{
		domain.NewInMemoryClass("LargeClass.java", "/path/LargeClass.java", largeContent),
		domain.NewInMemoryClass("SmallClass.java", "/path/SmallClass.java", "public class SmallClass {}"),
	}

	// Act
	result, err := agent.criticizeAll(classes, 1000)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "SmallClass.java", result[0].class.Name())

	// Check that the large class was skipped
	logs := logger.GetLogs()
	found := false
	for _, log := range logs {
		if contains(log, "has too many tokens") {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected log message about too many tokens")
}

func TestCriticizeAll_EmptyClasses(t *testing.T) {
	// Arrange
	logger := &mockLogger{}
	critic := &mockCritic{}

	agent := &agent{
		log:    logger,
		critic: critic,
	}

	classes := []domain.Class{}

	// Act
	result, err := agent.criticizeAll(classes, 1000)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 0)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) != -1
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
