package brain

type Brain interface {
	Ask(question string) (string, error)
}

func New(provider, token string) Brain {
	switch provider {
	case "deepseek":
		return NewDeepSeek(token)
	default:
		return NewMock()
	}
}
