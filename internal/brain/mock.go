package brain

import (
	"fmt"
)

// mockBrain represents a mock implementation of the Brain interface,
// used for testing purposes with predefined playbooks.
type mockBrain struct {
	playbooks []string
}

// NewMock creates a new instance of MockBrain with the provided playbooks.
// If no playbooks are provided, a default echo playbook will be used.
func NewMock(playbooks ...string) Brain {
	return &mockBrain{playbooks: playbooks}
}

// Ask processes a given question and provides a response based on the active playbook.
// Returns an error if the question is empty or if the playbook fails to respond.
func (b *mockBrain) Ask(question string) (string, error) {
	if question == "" {
		return "", fmt.Errorf("question cannot be empty")
	}
	book, err := b.Playbook()
	if err != nil {
		return "", fmt.Errorf("failed to get playbook: %w", err)
	}
	return book.Ask(question), nil
}

// Playbook retrieves the active playbook for the MockBrain.
// Supports either an echo playbook or a single YAML-based playbook.
// Returns an error if multiple playbooks are provided.
func (b *mockBrain) Playbook() (Playbook, error) {
	switch len(b.playbooks) {
	case 0:
		return &EchoPlaybook{}, nil
	case 1:
		if b.playbooks[0] != "" {
			return NewYAMLPlaybook(b.playbooks[0])
		}
		return &EchoPlaybook{}, nil
	default:
		return nil, fmt.Errorf("mock brain supports only one playbook at a time")
	}
}
