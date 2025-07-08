package brain

import (
	"fmt"
)

type MockBrain struct {
	playbooks []string
}

func NewMock(playbooks ...string) Brain {
	return &MockBrain{playbooks: playbooks}
}

func (b *MockBrain) Ask(question string) (string, error) {
	if question == "" {
		return "", fmt.Errorf("question cannot be empty")
	}
	book, err := b.Playbook()
	if err != nil {
		return "", fmt.Errorf("failed to get playbook: %w", err)
	}
	return book.Ask(question), nil
}

func (b *MockBrain) Playbook() (Playbook, error) {
	if len(b.playbooks) == 0 {
		return &echoPlaybook{}, nil
	} else if len(b.playbooks) == 1 {
		if b.playbooks[0] != "" {
			return NewYAMLPlaybook(b.playbooks[0])
		} else {
			return &echoPlaybook{}, nil
		}
	} else {
		return nil, fmt.Errorf("mock brain supports only one playbook at a time")
	}
}
