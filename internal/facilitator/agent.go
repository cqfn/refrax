package facilitator

import (
	"fmt"
	"strconv"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/domain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/util"
)

type agent struct {
	brain  brain.Brain
	log    log.Logger
	critic domain.Critic
	fixer  domain.Fixer
}

func (a *agent) Refactor(t domain.Task) ([]domain.Class, error) {
	size, err := maxSize(t)
	if err != nil {
		return nil, fmt.Errorf("failed to get max size limit: %w", err)
	}
	if t.Description() != "refactor the project" {
		a.log.Warn("received a message that is not related to refactoring. ignoring.")
		return nil, fmt.Errorf("received a message that is not related to refactoring")
	}
	a.log.Info("received request for refactoring, number of attached files: %d, max-size: %d", len(t.Classes()), size)
	refactored := make([]domain.Class, 0, len(t.Classes()))
	var example domain.Class
	changed := 0
	for _, class := range t.Classes() {
		a.log.Info("received class for refactoring: %q", class.Name())
		if changed >= size {
			a.log.Warn("refactoring class %s would exceed max size %d, skipping refactoring", class.Name(), size)
			refactored = append(refactored, class)
			continue
		}
		suggestions, err := a.critic.Review(class)
		if err != nil {
			return nil, fmt.Errorf("failed to ask critic: %w", err)
		}
		a.log.Info("received %d suggestions from critic", len(suggestions))
		modified, err := a.fixer.Fix(class, suggestions, example)
		if err != nil {
			return nil, fmt.Errorf("failed to ask fixer: %w", err)
		}
		a.log.Info("fixed class %s, changed content", modified.Name())
		refactored = append(refactored, modified)
		diff := util.Diff(class.Content(), modified.Content())
		changed += diff
	}
	return refactored, nil
}

func maxSize(t domain.Task) (int, error) {
	size, ok := t.Param("max-size")
	if !ok {
		size = "200"
	}
	return strconv.Atoi(size)
}
