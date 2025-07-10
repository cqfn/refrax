package brain

// Playbook defines the interface for a playbook that can ask questions and return answers.
// It is used in tests mostly.
type Playbook interface {
	Ask(question string) string
}
