package brain

type echoPlaybook struct{}

func NewEchoPlaybook() *echoPlaybook {
	return &echoPlaybook{}
}

func (e *echoPlaybook) Ask(question string) string {
	return question
}
