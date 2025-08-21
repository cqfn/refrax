// Package reviewer provides an implementation of a code review agent that uses AI to suggest fixes based on command output.
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
)

type agent struct {
	logger log.Logger
	cmds   []string
	ai     brain.Brain
}

func (a *agent) Review() (*domain.Artifacts, error) {
	var res []domain.Suggestion
	for _, cmd := range a.cmds {
		suggestions, err := a.runCmd(cmd)
		if err != nil {
			return nil, fmt.Errorf("failed to run command %s: %w", cmd, err)
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
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	parts := strings.Split(cmd, " ")
	command := exec.Command(parts[0], parts[1:]...) // #nosec G204
	command.Stdout = &out
	command.Stderr = &errOut
	command.Dir = root
	a.logger.Info("running review command: %s in %s", cmd, root)
	err = command.Run()
	if err == nil {
		a.logger.Info("review command completed successfully: %s", cmd)
		return make([]domain.Suggestion, 0), nil
	}
	a.logger.Info("failed to run review command: %s, error: %v", cmd, err)
	a.logger.Info("asking AI to form suggestions based on the error output")
	outb := out.Bytes()
	errb := errOut.Bytes()
	raw, err := a.ai.Ask(prompt(cmd, root, err, string(outb), string(errb)))
	if err != nil {
		return nil, fmt.Errorf("failed to ask AI for suggestions: %w", err)
	}
	parsed := parseSuggestions(raw)
	return parsed, nil
}

func prompt(cmd, cwd string, runErr error, stdout, stderr string) string {
	return fmt.Sprintf(
		`You are an experienced software engineer. The following command failed during compilation or build:

Command: %q
Working directory: %s
Error: %s

--- STDERR ---
%s
______________

--- STDOUT ---
%s
______________

Please suggest specific, actionable steps to fix the problem.
- Focus only on the issues visible in the provided output.
- Keep the suggestions short and practical.
- Number the suggestions starting from 1.
- If there are multiple possible causes, list them in priority order.
- Do not suggest general fixes or unrelated changes.
- Do not suggest new libraries or tools.

Answer in the following format:
		<java class name>: <suggestion 1>
		<java class name>: <suggestion 2>
		<java class name>: <suggestion 3>

Use single-line suggestions, each starting with a class name followed by a colon and the suggestion text.
`,
		cmd, cwd, runErr.Error(), stderr, stdout,
	)
}

func parseSuggestions(output string) []domain.Suggestion {
	lines := strings.Split(output, "\n")
	res := make([]domain.Suggestion, 0, len(lines))
	for _, line := range lines {
		res = append(res, domain.NewSuggestion(line))
	}
	return res
}
