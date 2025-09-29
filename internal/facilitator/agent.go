package facilitator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/domain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/prompts"
	"github.com/cqfn/refrax/internal/stats"
	"github.com/cqfn/refrax/internal/util"
)

type agent struct {
	brain    brain.Brain
	log      log.Logger
	critic   domain.Critic
	fixer    domain.Fixer
	reviewer domain.Reviewer
	frounds  int
}

func (a *agent) Refactor(job *domain.Job) (*domain.Artifacts, error) {
	size, err := maxSize(job)
	if err != nil {
		return nil, fmt.Errorf("failed to get max size limit: %w", err)
	}
	if job.Descr.Text != "refactor the project" {
		a.log.Warn("Received a message that is not related to refactoring, ignoring")
		return nil, fmt.Errorf("received a message that is not related to refactoring")
	}
	diff := 0
	result := make([]domain.Class, 0)
	for diff < size {
		improvements, err := a.criticizeAll(job.Classes, size)
		if err != nil {
			return nil, fmt.Errorf("failed to criticize classes: %w", err)
		}
		if len(improvements) == 0 {
			a.log.Warn("No improvements found, returning original classes")
			res := &domain.Artifacts{
				Descr:   &domain.Description{Text: "no improvements found"},
				Classes: job.Classes,
			}
			return res, nil
		}
		nsuggestions := 0
		for _, imp := range improvements {
			nsuggestions += len(imp.suggestions)
		}
		if nsuggestions == 0 {
			a.log.Warn("No suggestions found")
			break
		}
		mostImportant, err := a.mostFrequent(improvements)
		if err != nil {
			return nil, fmt.Errorf("failed to get most frequent suggestions: %w", err)
		}
		a.log.Info("Received %d most frequent suggestions from brain", len(mostImportant))
		refactored, changed, err := a.fixAllSuggestions(mostImportant, size)
		if err != nil {
			return nil, fmt.Errorf("failed to fix all suggestions: %w", err)
		}
		diff += changed
		err = a.fixErrors(refactored)
		result = append(result, refactored...)
		if err != nil {
			return nil, fmt.Errorf("failed to stabilize refactored classes: %w", err)
		}
	}
	res := &domain.Artifacts{
		Descr:   &domain.Description{Text: "refactored classes"},
		Classes: result,
	}
	return res, nil
}

// fixAllSuggestions processes all improvements concurrently, ensuring that the total changes do not exceed the specified size limit.
func (a *agent) fixAllSuggestions(improvements []improvement, size int) ([]domain.Class, int, error) {
	refactored := make([]domain.Class, 0)
	fixChannel := make(chan fixResult, len(improvements))
	send := make(map[string]improvement, 0)
	for _, imp := range improvements {
		send[imp.class.Path()] = imp
		go a.fixSuggestions(imp, fixChannel)
	}
	changed := 0
	for range len(send) {
		fixRes := <-fixChannel
		if fixRes.err != nil {
			panic(fmt.Sprintf("Failed to fix class: %v", fixRes.err))
		}
		path := fixRes.class.Path()
		class := send[path].class
		if changed >= size {
			a.log.Warn("Refactoring class %s would exceed max-size of %d (current %d), skipping refactoring", class.Name(), size, changed)
			continue
		}
		modified := fixRes.class
		refactored = append(refactored, modified)
		diff := util.Diff(class.Content(), modified.Content())
		a.log.Info("Fixed class %s (%s), changed content (diff %d)", modified.Name(), modified.Path(), diff)
		changed += diff
	}
	for _, c := range refactored {
		class := domain.NewFSClass(c.Name(), c.Path())
		err := class.SetContent(c.Content())
		a.log.Info("Setting content for class %s (%s)", class.Name(), class.Path())
		if err != nil {
			return nil, 0, fmt.Errorf("failed to set content for class %s: %w", class.Name(), err)
		}
	}
	return refactored, changed, nil
}

func (a *agent) criticizeAll(classes []domain.Class, size int) ([]improvement, error) {
	nclasses := len(classes)
	a.log.Info("Received request for refactoring, number of attached files: %d, max-size: %d", nclasses, size)
	improvements := make([]improvement, 0, nclasses)
	ch := make(chan improvementResult, nclasses)
	nreviewed := 0
	for _, class := range classes {
		tokens, _ := stats.Tokens(class.Content())
		a.log.Info("Class %s has %d tokens", class.Path(), tokens)
		if tokens < 2_000 {
			nreviewed++
			go a.criticize(class, ch)
		} else {
			a.log.Warn("Class %s (%s) has too many tokens (%d), skipping review", class.Name(), class.Path(), tokens)
		}
	}
	a.log.Info("Number of classes to review: %d", nreviewed)
	for range nreviewed {
		impr := <-ch
		if impr.err != nil {
			return nil, fmt.Errorf("failed to review class: %w", impr.err)
		}
		improvements = append(improvements, impr.important)
	}
	return improvements, nil
}

// fixErrors chcks whether the refactored classes have any errors and tries to fix them if any.
func (a *agent) fixErrors(refactored []domain.Class) error {
	a.log.Info("Fixing refactored classes, number of classes: %d", len(refactored))
	artifacts, err := a.reviewer.Review()
	if err != nil {
		return fmt.Errorf("failed to review project: %w", err)
	}
	improvements := artifacts.Suggestions
	a.log.Info("Received %d suggestions from reviewer", len(improvements))
	for _, improvement := range improvements {
		a.log.Info("Received suggestion: %s", improvement)
	}
	counter := a.frounds
	for len(improvements) > 0 {
		if err != nil {
			return fmt.Errorf("failed to review project: %w", err)
		}
		perclass := a.understandClasses(refactored, improvements)
		for k, v := range perclass {
			job := domain.Job{
				Descr: &domain.Description{
					Text: "fix the class",
				},
				Classes:     []domain.Class{k},
				Suggestions: v,
			}
			fixed, uerr := a.fixer.Fix(&job)
			if uerr != nil {
				return fmt.Errorf("failed to fix project: %w", uerr)
			}
			updated := fixed.Classes[0]
			class := domain.NewFSClass(k.Name(), k.Path())
			a.log.Info("Updating class %s (%s) with new content", class.Name(), class.Path())
			uerr = class.SetContent(updated.Content())
			if uerr != nil {
				return fmt.Errorf("failed to set content for class %s: %w", class.Name(), uerr)
			}
		}
		counter--
		if counter < 0 {
			return fmt.Errorf("too many rounds of fixing errors, stopping")
		}
		artifacts, err = a.reviewer.Review()
		improvements = artifacts.Suggestions
	}
	return nil
}

func (a *agent) understandClasses(clases []domain.Class, suggestions []domain.Suggestion) map[domain.Class][]domain.Suggestion {
	res := make(map[domain.Class][]domain.Suggestion, len(clases))
	for _, s := range suggestions {
		for _, c := range clases {
			actual := s.ClassPath
			expected := c.Path()
			if actual == expected || strings.Contains(actual, expected) || strings.Contains(expected, actual) {
				a.log.Info("Associating suggestion %q with class %s", s.Text, c.Path())
				res[c] = append(res[c], s)
				break
			}
		}
	}
	return res
}

// doFixSuggestions sends a fix request to the fixer and returns the modified class or an error.
func (a *agent) fixSuggestions(imp improvement, ch chan<- fixResult) {
	class := imp.class
	suggestions := imp.suggestions
	job := domain.Job{
		Descr: &domain.Description{
			Text: "fix the class",
		},
		Classes:     []domain.Class{class},
		Suggestions: suggestions,
		Examples:    []domain.Class{nil},
	}
	modified, err := a.fixer.Fix(&job)
	if err != nil {
		ch <- fixResult{fmt.Errorf("failed to ask fixer: %w", err), nil}
		return
	}
	ch <- fixResult{nil, modified.Classes[0]}
}

// criticize sends a review request to the critic and returns the suggestions or an error.
func (a *agent) criticize(class domain.Class, ch chan<- improvementResult) {
	a.log.Info("Received class for refactoring: %q", class.Path())
	job := domain.Job{
		Descr: &domain.Description{
			Text: "refactor the class",
		},
		Classes: []domain.Class{class},
	}
	suggestions, err := a.critic.Review(&job)
	if err != nil {
		ch <- improvementResult{err: fmt.Errorf("failed to ask critic: %w", err), important: improvement{class: class}}
		return
	}
	a.log.Info("Received %d suggestions from critic", len(suggestions.Suggestions))
	if len(suggestions.Suggestions) == 0 {
		a.log.Info("No suggestions found for class %s", class.Path())
		ch <- improvementResult{err: nil, important: improvement{class: class}}
		return
	}
	ch <- improvementResult{
		err: nil,
		important: improvement{
			class:       class,
			suggestions: suggestions.Suggestions,
		},
	}
}

type fixResult struct {
	err   error
	class domain.Class
}

type improvementResult struct {
	err       error
	important improvement
}

type improvement struct {
	class       domain.Class
	suggestions []domain.Suggestion
}

type groupPromptData struct {
	Suggestions []domain.Suggestion
}

type choosePromoptData struct {
	Groupped string
}

func (a *agent) mostFrequent(improvements []improvement) ([]improvement, error) {
	a.log.Info("Grouping all suggestions...")
	all := make([]domain.Suggestion, 0, len(improvements))
	for _, imp := range improvements {
		all = append(all, imp.suggestions...)
	}
	prompt := prompts.User{
		Data: groupPromptData{
			Suggestions: all,
		},
		Name: "facilitator/group.md.tmpl",
	}
	important, err := a.brain.Ask(prompt.String())
	if err != nil {
		return nil, fmt.Errorf("failed to ask the brain to group suggestions: %w", err)
	}

	prompt = prompts.User{
		Data: choosePromoptData{
			Groupped: important,
		},
		Name: "facilitator/choose.md.tmpl",
	}
	a.log.Info("Choosing the most important suggestions...")
	important, err = a.brain.Ask(prompt.String())
	if err != nil {
		return nil, fmt.Errorf("failed to ask brain for most frequent suggestion: %w", err)
	}
	classSuggestions := make(map[string][]string, 0)
	for s := range strings.SplitSeq(strings.ReplaceAll(important, "\r\n", "\n"), "\n") {
		a.log.Info("Suggestion to consider: %s", s)
		if !strings.Contains(s, ":") {
			a.log.Warn("Can't find a delimiter ':'")
			continue
		}
		split := strings.Split(s, ":")
		if len(split) < 2 {
			a.log.Warn("Skipping suggestion without class name: %s", s)
			continue
		}
		className := strings.TrimSpace(split[0])
		classSuggestion := strings.TrimSpace(split[1])
		classSuggestions[className] = append(classSuggestions[className], classSuggestion)
	}
	a.log.Info("Received %d suggestions from brain", len(classSuggestions))
	ires := make([]improvement, 0)
	for k, v := range classSuggestions {
		class, err := findClass(improvements, k)
		if err != nil {
			a.log.Warn("Class %s not found in improvements, skipping: %v", k, err)
			continue
		}
		var suggetions []domain.Suggestion
		for _, s := range v {
			suggetions = append(suggetions, *domain.NewSuggestion(s, class.Path()))
		}
		ires = append(ires, improvement{
			class:       class,
			suggestions: suggetions,
		})
	}
	return ires, nil
}

func findClass(improvements []improvement, path string) (domain.Class, error) {
	for _, imp := range improvements {
		if imp.class.Path() == path {
			return imp.class, nil
		}
	}
	all := make([]string, 0, len(improvements))
	for _, imp := range improvements {
		all = append(all, imp.class.Path())
	}
	return nil, fmt.Errorf("class %s not found in improvements %d", path, len(all))
}

func maxSize(t *domain.Job) (int, error) {
	size, ok := t.Param("max-size")
	ssize := fmt.Sprintf("%v", size)
	if !ok {
		ssize = "200"
	}
	return strconv.Atoi(ssize)
}
