package brain

type Playbook interface {
	Ask(question string) string
}
