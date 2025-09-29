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
	c.log.Info("Received class %q for analysis", class.Name())
	imperfections := tool.NewCombined(c.tools...).Imperfections()
	var imp []string
	if imperfections != "" {
		imp = strings.Split(imperfections, "\n")
	}
	data := promptData{
		Code:     class.Content(),
		Defects:  imp,
		NotFound: notFound,
	}
	prompt := prompts.User{
		Data: data,
		Name: "critic/critic.md.tmpl",
	}
	p := prompt.String()
	c.log.Debug("Rendered prompt for class %s: %s", class.Name(), p)
	answer, err := c.brain.Ask(p)
	if err != nil {
		return nil, fmt.Errorf("failed to get answer from brain: %w", err)
	}
	c.log.Debug("Received answer from brain for class %s: %s", class.Name(), answer)
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

func (c *agent) associated(suggestions []string, class string) []domain.Suggestion {
	res := make([]domain.Suggestion, 0)
	for i, suggestion := range suggestions {
		if strings.EqualFold(suggestion, notFound) {
			c.log.Info("No suggestions found for the class #%d: %s", i+1, class)
		} else {
			res = append(res, *domain.NewSuggestion(suggestion, class))
		}
	}
	c.log.Info("Total suggestions associated with class %s: %d", class, len(res))
	return res
}
