package brain

// EchoPlaybook represents a playbook that echoes back the question asked.
type EchoPlaybook struct{}

// NewEchoPlaybook creates a new instance of EchoPlaybook.
func NewEchoPlaybook() *EchoPlaybook {
	return &EchoPlaybook{}
}

// Ask echoes back the provided question.
func (e *EchoPlaybook) Ask(question string) string {
	return question
}
