package critic

import (
	"fmt"
	"strings"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/domain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/tool"
)

type agent struct {
	brain brain.Brain
	log   log.Logger
	tools []tool.Tool
}

const notFound = "No suggestions found"

const prompt = `Analyze the following Java code:

{{code}}

Identify possible improvements or flaws such as:

* grammar and spelling mistakes in comments,
* variables that can be inlined or removed without changing functionality,
* unnecessary comments inside methods,
* redundant code.

Don't suggest any changes that would alter the functionality of the code.
Don't suggest any changes that would require moving code parts between files (like extract class or extract an interface).
Don't suggest method renaming or class renaming.

Keep in mind the following imperfections with Java code, identified by automated static analysis system:

{{imperfections}}

Respond with a few most relevant and important suggestion for improvement, formatted as a few lines of text. 
If there are no suggestions or they are insignificant, respond with "{{not-found}}".
Do not include any explanations, summaries, or extra text.
`

// Review sends the provided Java class to the Critic for analysis and returns suggested improvements.
func (c *agent) Review(job *domain.Job) (*domain.Artifacts, error) {
	class := job.Classes[0]
	c.log.Info("received class %q for analysis", class.Name())
	replacer := strings.NewReplacer(
		"{{code}}", class.Content(),
		"{{imperfections}}", tool.NewCombined(c.tools...).Imperfections(),
		"{{not-found}}", notFound,
	)
	answer, err := c.brain.Ask(replacer.Replace(prompt))
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
