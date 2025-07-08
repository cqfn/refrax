package brain

type Brain interface {
	Ask(question string) (string, error)
}

func New(provider, token string, playbook ...string) Brain {
	switch provider {
	case "deepseek":
		return NewDeepSeek(token)
	default:
		if len(playbook) == 0 {
			return NewMock()
		} else {
			return NewMock(playbook[0])
		}
	}
}
