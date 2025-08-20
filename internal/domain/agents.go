package domain

// Facilitator represents an interface to refactor tasks into multiple classes.
type Facilitator interface {
	Refactor(task Task) ([]Class, error)
}

// Critic represents an interface to review a class and provide suggestions.
type Critic interface {
	Review(job *Job) ([]Suggestion, error)
}

// Fixer represents an interface to fix a class based on suggestions and an example.
type Fixer interface {
	Fix(job *Job) (Class, error)
}

// Reviewer represents an interface for a reviewer that can review changes made.
type Reviewer interface {
	Review() ([]Suggestion, error)
}

type Job struct {
	Descr       *Description
	Classes     []Class
	Suggestions []Suggestion
	Examples    []Class
}

type Description struct {
	Text string
	meta map[string]any
}

func (j *Job) FirstClass() Class {
	if len(j.Classes) == 0 {
		return nil
	}
	return j.Classes[0]
}
