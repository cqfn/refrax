// Package brain provides an interface and implementations for asking questions programmatically.
package brain

// Brain represents an interface for asking questions and receiving answers.
type Brain interface {
	Ask(question string) (string, error)
}

const deepseek = "deepseek"
const openai = "openai"

// New creates a new instance of Brain based on the provided provider and optional playbook strings.
func New(provider, token string, playbook ...string) Brain {
	switch provider {
	case deepseek:
		return NewDeepSeek(token)
	case openai:
		return NewOpenAI(token)
	default:
		if len(playbook) == 0 {
			return NewMock()
		}
		return NewMock(playbook[0]) // Outdented as per lint suggestion.
	}
}

func trimmed(prompt string) string {
	limit := 120 * 400
	runes := []rune(prompt)
	if len(runes) > limit {
		return string(runes[:limit])
	}
	return prompt
}
