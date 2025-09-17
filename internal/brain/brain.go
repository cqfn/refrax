package brain

import "fmt"

// Brain represents an interface for asking questions and receiving answers.
type Brain interface {
	Ask(question string) (string, error)
}

const deepseek = "deepseek"

const openai = "openai"

const mock = "mock"

// New creates a new instance of Brain based on the provided provider and optional playbook strings.
func New(provider, token, system string, playbook ...string) (Brain, error) {
	switch provider {
	case deepseek:
		return NewDeepSeek(token, system), nil
	case openai:
		return NewOpenAI(token, system), nil
	case mock:
		if len(playbook) == 0 {
			return NewMock(), nil
		}
		return NewMock(playbook[0]), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
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
