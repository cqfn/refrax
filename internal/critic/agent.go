package critic

import (
	"fmt"
	"strings"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/domain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/prompts"
	"github.com/cqfn/refrax/internal/tool"
)

type agent struct {
	brain brain.Brain
	log   log.Logger
	tools []tool.Tool
}

// notFound is the message returned when no suggestions are found.
const notFound = "No suggestions found"

// promptData holds the data to be injected into the prompt template.
type promptData struct {
	Code     string
	Defects  []string
	NotFound string
}

// Review sends the provided Java class to the Critic for analysis and returns suggested improvements.
func (c *agent) Review(job *domain.Job) (*domain.Artifacts, error) {
	class := job.Classes[0]
	c.log.Info("received class %q for analysis", class.Name())
	data := promptData{
		Code:     class.Content(),
		Defects:  []string{tool.NewCombined(c.tools...).Imperfections()},
		NotFound: notFound,
	}
	prompt := prompts.User{
		Data: data,
		Name: "critic/critic.md.tmpl",
	}
	c.log.Debug("rendered prompt for class %s: %s", class.Name(), prompt)
	answer, err := c.brain.Ask(prompt.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get answer from brain: %w", err)
	}
	suggestions := c.associated(parseAnswer(answer), class.Path())
	logSuggestions(c.log, suggestions)
	artifacts := domain.Artifacts{
		Descr: &domain.Description{
			Text: fmt.Sprintf("Critique for class %s", class.Name()),
		},
		Suggestions: suggestions,
	}
	return &artifacts, nil
}

func logSuggestions(logger log.Logger, suggestions []domain.Suggestion) {
	for i, suggestion := range suggestions {
		logger.Info("#%d: %s", i+1, suggestion)
	}
}

func parseAnswer(answer string) []string {
	lines := strings.Split(strings.TrimSpace(answer), "\n")
	var suggestions []string
	for _, line := range lines {
		suggestion := strings.TrimSpace(line)
		if suggestion != "" {
			suggestions = append(suggestions, suggestion)
		}
	}
	return suggestions
}

func (a *agent) associated(suggestions []string, class string) []domain.Suggestion {
	res := make([]domain.Suggestion, 0)
	for i, suggestion := range suggestions {
		if strings.EqualFold(suggestion, notFound) {
			a.log.Info("no suggestions found for the class #%d: %s", i+1, class)
		} else {
			res = append(res, *domain.NewSuggestion(suggestion, class))
		}
	}
	a.log.Info("total suggestions associated with class %s: %d", class, len(res))
	return res
}
