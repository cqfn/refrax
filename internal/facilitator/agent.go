package facilitator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/domain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/stats"
	"github.com/cqfn/refrax/internal/util"
)

type agent struct {
	brain    brain.Brain
	log      log.Logger
	critic   domain.Critic
	fixer    domain.Fixer
	reviewer domain.Reviewer
}

func (a *agent) Refactor(job *domain.Job) (*domain.Artifacts, error) {
	size, err := maxSize(job)
	if err != nil {
		return nil, fmt.Errorf("failed to get max size limit: %w", err)
	}
	if job.Descr.Text != "refactor the project" {
		a.log.Warn("received a message that is not related to refactoring. ignoring.")
		return nil, fmt.Errorf("received a message that is not related to refactoring")
	}
	classes := job.Classes
	nclasses := len(classes)
	a.log.Info("received request for refactoring, number of attached files: %d, max-size: %d", nclasses, size)
	var example domain.Class
	improvements := make([]improvement, 0, nclasses)
	ch := make(chan improvementResult, nclasses)
	nreviewed := 0
	untouched := make(map[string]domain.Class, 0)
	for _, class := range classes {
		untouched[class.Path()] = class
		tokens, _ := stats.Tokens(class.Content())
		a.log.Info("class %s has %d tokens", class.Path(), tokens)
		if tokens < 2_000 {
			nreviewed++
			go a.review(class, ch)
		} else {
			a.log.Warn("class %s (%s) has too many tokens (%d), skipping review", class.Name(), class.Path(), tokens)
		}
	}
	a.log.Info("number of classes to review: %d, untouched: %d", nreviewed, len(untouched))
	for range nreviewed {
		impr := <-ch
		if impr.err != nil {
			return nil, fmt.Errorf("failed to review class: %w", impr.err)
		}
		improvements = append(improvements, impr.important)
	}
	if len(improvements) == 0 {
		a.log.Warn("no improvements found, returning original classes")
		res := &domain.Artifacts{
			Descr:   &domain.Description{Text: "no improvements found"},
			Classes: classes,
		}
		return res, nil
	}
	mostImportant, err := a.mostFrequent(improvements)
	if err != nil {
		return nil, fmt.Errorf("failed to get most frequent suggestions: %w", err)
	}
	a.log.Info("received %d most frequent suggestions from brain", len(mostImportant))
	refactored := make([]domain.Class, 0)
	fixChannel := make(chan fixResult, len(mostImportant))
	send := make(map[string]improvement, 0)
	for _, imp := range mostImportant {
		send[imp.class.Path()] = imp
		delete(untouched, imp.class.Path())
		go a.fix(imp, example, fixChannel)
	}
	changed := 0
	for range len(send) {
		fixRes := <-fixChannel
		if fixRes.err != nil {
			panic(fmt.Sprintf("failed to fix class: %v", fixRes.err))
		}
		path := fixRes.class.Path()
		class := send[path].class
		if changed >= size {
			a.log.Warn("refactoring class %s would exceed max-size is %d (current %d), skipping refactoring", class.Name(), size, changed)
			continue
		}
		modified := fixRes.class
		refactored = append(refactored, modified)
		diff := util.Diff(class.Content(), modified.Content())
		a.log.Info("fixed class %s (%s), changed content (diff %d)", modified.Name(), modified.Path(), diff)
		changed += diff
	}
	for _, class := range untouched {
		refactored = append(refactored, class)
	}
	for _, c := range refactored {
		class := domain.NewFSClass(c.Name(), c.Path())
		err = class.SetContent(c.Content())
		a.log.Info("setting content for class %s (%s)", class.Name(), class.Path())
		if err != nil {
			return nil, fmt.Errorf("failed to set content for class %s: %w", class.Name(), err)
		}
	}
	err = a.stabilize(refactored)
	if err != nil {
		return nil, fmt.Errorf("failed to stabilize refactored classes: %w", err)
	}
	res := &domain.Artifacts{
		Descr:   &domain.Description{Text: "refactored classes"},
		Classes: refactored,
	}
	return res, nil
}

func (a *agent) stabilize(refactored []domain.Class) error {
	a.log.Info("stabilizing refactored classes, number of classes: %d", len(refactored))
	artifacts, err := a.reviewer.Review()
	improvements := artifacts.Suggestions
	a.log.Info("received %d suggestions from reviewer", len(improvements))
	for _, improvement := range improvements {
		a.log.Info("received suggestion: %s", improvement)
	}
	counter := 3
	for len(improvements) > 0 && counter > 0 {
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
			a.log.Info("updating class %s (%s) with new content", class.Name(), class.Path())
			uerr = class.SetContent(updated.Content())
			if uerr != nil {
				return fmt.Errorf("failed to set content for class %s: %w", class.Name(), uerr)
			}
		}
		counter--
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
				a.log.Info("associating suggestion %q with class %s", s.Text, c.Path())
				res[c] = append(res[c], s)
				break
			}
		}
	}
	return res
}

func (a *agent) fix(imp improvement, example domain.Class, ch chan<- fixResult) {
	class := imp.class
	suggestions := imp.suggestions
	job := domain.Job{
		Descr: &domain.Description{
			Text: "fix the class",
		},
		Classes:     []domain.Class{class},
		Suggestions: suggestions,
		Examples:    []domain.Class{example},
	}
	modified, err := a.fixer.Fix(&job)
	if err != nil {
		ch <- fixResult{fmt.Errorf("failed to ask fixer: %w", err), nil}
		return
	}
	ch <- fixResult{nil, modified.Classes[0]}
}

func (a *agent) review(class domain.Class, ch chan<- improvementResult) {
	a.log.Info("received class for refactoring: %q", class.Path())
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
	a.log.Info("received %d suggestions from critic", len(suggestions.Suggestions))
	if len(suggestions.Suggestions) == 0 {
		a.log.Info("no suggestions found for class %s", class.Path())
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

func (a *agent) mostFrequent(improvements []improvement) ([]improvement, error) {
	a.log.Info("finding most imprtant suggestions from all the possible improvements...")
	all := make([]domain.Suggestion, 0, len(improvements))
	for _, imp := range improvements {
		all = append(all, imp.suggestions...)
	}
	prompt := "Here is a list of code improvement suggestions: \n" +
		"```\n%s\n```\n" +
		"Your task:" +
		"1. Group similar suggestions by their topic — for example: comments, naming, error handling, formatting, etc." +
		"2. If a group contains multiple similar suggestions, return all suggestions from that group." +
		"3. If no group has more than one suggestion, return just one representative suggestion from the list." +
		"4. Do not change the text of any suggestion." +
		"5. Do not explain or comment. Output only the selected suggestion(s), each on its own line." +
		"6. Do not remove or modify class names in any suggestion." +
		"Important! Return suggestions as they are. Literally!" +
		" Answer in the following format: " +
		"<java class path>: <suggestion 1> " +
		"<java class path>: <suggestion 2> " +
		"<java class path>: <suggestion 3> " +
		"Example:  " +
		"		src/test/java/com/example/service/Example.java: Fix the typo in the class comment"

	var summary string
	for _, s := range all {
		summary += fmt.Sprintf("%s: %s\n", s.ClassPath, s.Text)
	}
	prompt = fmt.Sprintf(prompt, summary)
	important, err := a.brain.Ask(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to ask brain for most frequent suggestion: %w", err)
	}
	followup := "Given the grouped suggestions below, identify the largest group (the one with the most suggestions).\n" +
		"Return only the suggestions from that group.\n" +
		"Do not include group names, headers, or any explanations. Only output the suggestions.\n" +
		"Do not modify any suggestion. Keep class names intact.\n\n" +
		"Important! Return suggestions as they are. Literally!" +
		"```\n%s\n```" +
		" Answer in the following format: " +
		"<java class path>: <suggestion 1> " +
		"<java class path>: <suggestion 2> " +
		"<java class path>: <suggestion 3> " +
		"Example:  " +
		"		src/test/java/com/example/service/Example.java: Fix the typo in the class comment"
	prompt = fmt.Sprintf(followup, important)
	important, err = a.brain.Ask(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to ask brain for most frequent suggestion: %w", err)
	}
	classSuggestions := make(map[string][]string, 0)
	for s := range strings.SplitSeq(strings.ReplaceAll(important, "\r\n", "\n"), "\n") {
		a.log.Info("suggestion to consider: %s", s)
		if !strings.Contains(s, ":") {
			a.log.Warn("can't find a delimeter (:)")
			continue
		}
		split := strings.Split(s, ":")
		if len(split) < 2 {
			a.log.Warn("skipping suggestion without class name: %s", s)
			continue
		}
		className := strings.TrimSpace(split[0])
		classSuggestion := strings.TrimSpace(split[1])
		classSuggestions[className] = append(classSuggestions[className], classSuggestion)
	}
	a.log.Info("received %d suggestions from brain", len(classSuggestions))
	ires := make([]improvement, 0)
	for k, v := range classSuggestions {
		class, err := findClass(improvements, k)
		if err != nil {
			return nil, fmt.Errorf("failed to find class %s in improvements: %w", k, err)
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
	return nil, fmt.Errorf("class %s not found in improvements", path)
}

func maxSize(t *domain.Job) (int, error) {
	size, ok := t.Param("max-size")
	ssize := fmt.Sprintf("%v", size)
	if !ok {
		ssize = "200"
	}
	return strconv.Atoi(ssize)
}
