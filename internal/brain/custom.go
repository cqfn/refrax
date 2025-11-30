package brain

func NewCustom(token, url, model, system string) (Brain, error) {
	return NewOpenAI(token, url, model, system)
}
