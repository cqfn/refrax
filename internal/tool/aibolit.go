// Package tool provides tools for identifying refactoring opportunities in code.
package tool

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/cqfn/refrax/internal/log"
)

// Aibolit is a tool for identifying refactoring opportunities in Java code.
// @todo #2:45min Implement support for Aibolit with multiple classes.
// As for now, we check only the first class, and return imperfections result. Instead, we need to support
// multiple files instead. Let's implement such Aibolit struct, that will be able to manage whole project, instead
// of single Java file. Also see this related issue: https://github.com/cqfn/refrax/issues/28.
type Aibolit struct {
	filename string
	executor runner
}

// NewAibolit creates a new instance of Aibolit for analyzing a single Java file.
func NewAibolit(filename string) *Aibolit {
	return &Aibolit{filename, &exexRunner{}}
}

// Imperfections analyzes the Java code and identifies refactoring opportunities.
func (a *Aibolit) Imperfections() string {
	opportunities, _ := a.executor.Run("aibolit", "check", "--filenames", "Foo.java")
	log.Debug("Identified refactoring opportunities with aibolit: \n%s", opportunities)
	return sanitized(string(opportunities))
}

func sanitized(raw string) string {
	complaint := regexp.MustCompile(`^[^:]+\.java\[\d+\]: .+`)
	var complaints []string
	for _, line := range strings.Split(raw, "\n") {
		if complaint.MatchString(line) {
			complaints = append(complaints, line)
		}
	}
	return strings.Join(complaints, "\n")
}

type runner interface {
	Run(name string, args ...string) ([]byte, error)
}

type exexRunner struct{}

func (e *exexRunner) Run(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}
