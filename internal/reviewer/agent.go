// Package reviewer is for the reviewer.
package reviewer

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/domain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/prompts"
)

type agent struct {
	logger log.Logger
	cmds   []string
	ai     brain.Brain
}

// promptData holds all inputs for the template in one place.
type promptData struct {
	Command string
	WorkDir string
	Error   string
	Stderr  string
	Stdout  string
}

func (a *agent) Review() (*domain.Artifacts, error) {
	var res []domain.Suggestion
	a.logger.Info("Starting review using %d commands, %s", len(a.cmds), strings.Join(a.cmds, ", "))
	for _, cmd := range a.cmds {
		suggestions, err := a.runCmd(cmd)
		if err != nil {
			return nil, fmt.Errorf("Failed to run command %s: %w", cmd, err)
		}
		res = append(res, suggestions...)
	}
	artifacts := &domain.Artifacts{
		Descr:       &domain.Description{Text: "suggestions based on command outputs"},
		Suggestions: res,
	}
	return artifacts, nil
}

func (a *agent) runCmd(cmd string) ([]domain.Suggestion, error) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	root, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Failed to get current working directory: %w", err)
	}
	parts := strings.Split(cmd, " ")
	command := exec.Command(parts[0], parts[1:]...) // #nosec G204
	command.Stdout = &out
	command.Stderr = &errOut
	command.Dir = root
	a.logger.Info("Running review command: %s in %s", cmd, root)
	err = command.Run()
	if err == nil {
		a.logger.Info("Review command completed successfully: %s", cmd)
		return make([]domain.Suggestion, 0), nil
	}
	a.logger.Info("Failed to run review command: %s, error: %v", cmd, err)
	a.logger.Info("Asking AI to form suggestions based on the error output")
	outb := out.Bytes()
	errb := errOut.Bytes()
	data := promptData{
		Command: cmd,
		WorkDir: root,
		Error:   err.Error(),
		Stderr:  string(errb),
		Stdout:  string(outb),
	}
	prompt := prompts.User{
		Data: data,
		Name: "reviewer/review.md.tmpl",
	}
	raw, err := a.ai.Ask(prompt.String())
	if err != nil {
		return nil, fmt.Errorf("Failed to ask AI for suggestions: %w", err)
	}
	parsed := a.parseSuggestions(raw)
	return parsed, nil
}

func (a *agent) parseSuggestions(output string) []domain.Suggestion {
	lines := strings.Split(output, "\n")
	res := make([]domain.Suggestion, 0, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			a.logger.Warn("Skipping malformed suggestion line: %q", line)
			continue
		}
		path := strings.TrimSpace(parts[0])
		text := strings.TrimSpace(strings.Join(parts[1:], ":"))
		res = append(res, *domain.NewSuggestion(text, path))
	}
	return res
}
