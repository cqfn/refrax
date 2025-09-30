package domain

import (
	"fmt"
	"strconv"
)

// Facilitator represents an interface to refactor tasks into multiple classes.
type Facilitator interface {
	Refactor(job *Job) (*Artifacts, error)
}

// Critic represents an interface to review a class and provide suggestions.
type Critic interface {
	Review(job *Job) (*Artifacts, error)
}

// Fixer represents an interface to fix a class based on suggestions and an example.
type Fixer interface {
	Fix(job *Job) (*Artifacts, error)
}

// Reviewer represents an interface for a reviewer that can review changes made.
type Reviewer interface {
	Review() (*Artifacts, error)
}

type Job struct {
	Descr       *Description
	Classes     []Class
	Suggestions []Suggestion
	Examples    []Class
}

func (j *Job) Param(key string) (any, bool) {
	if j.Descr == nil || j.Descr.Meta == nil {
		return nil, false
	}
	res, ok := j.Descr.Meta[key]
	return res, ok
}

func (j *Job) MaxSize() (int, error) {
	size, ok := j.Param("max-size")
	ssize := fmt.Sprintf("%v", size)
	if !ok {
		ssize = "200"
	}
	return strconv.Atoi(ssize)
}

type Artifacts struct {
	Descr       *Description
	Classes     []Class
	Suggestions []Suggestion
}

type Description struct {
	Text string
	Meta map[string]any
}

type Suggestion struct {
	ClassPath string
	Text      string
}

func (s *Suggestion) String() string {
	return fmt.Sprintf("suggestion for %s: %s", s.ClassPath, s.Text)
}
