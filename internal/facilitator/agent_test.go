package facilitator

import (
	"testing"

	"github.com/cqfn/refrax/internal/brain"
	"github.com/cqfn/refrax/internal/domain"
	"github.com/cqfn/refrax/internal/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefactor_ReturnsOriginalClassesWhenNoImprovementsFound(t *testing.T) {
	class := domain.NewInMemoryClass("Tëst.java", "/path/Tëst.java", "class Tëst {}")
	job := &domain.Job{
		Descr:   &domain.Description{Text: "refactor the project", Meta: map[string]any{"max-size": "100"}},
		Classes: []domain.Class{class},
	}
	a := &agent{
		brain:    brain.NewMock(),
		log:      log.NewMock(),
		critic:   &mockCritic{suggestions: []domain.Suggestion{}},
		fixer:    &mockFixer{},
		reviewer: &mockReviewer{},
		attempts: 1,
		frounds:  3,
	}
	artifacts, err := a.Refactor(job)
	require.NoError(t, err, "unexpected error during refactor")
	assert.Equal(t, 1, len(artifacts.Classes), "must return original class when no improvements")
}

func TestRefactor_ReturnsErrorForNonRefactorMessage(t *testing.T) {
	class := domain.NewInMemoryClass("Ünïcödé.java", "/Ünïcödé.java", "class Ünïcödé {}")
	job := &domain.Job{
		Descr:   &domain.Description{Text: "do something else"},
		Classes: []domain.Class{class},
	}
	a := &agent{
		brain:    brain.NewMock(),
		log:      log.NewMock(),
		critic:   &mockCritic{},
		fixer:    &mockFixer{},
		reviewer: &mockReviewer{},
		attempts: 1,
		frounds:  3,
	}
	_, err := a.Refactor(job)
	require.Error(t, err, "expected error for non-refactor message")
}

func TestRefactor_SkipsRefactoringWhenAttemptsZero(t *testing.T) {
	class := domain.NewInMemoryClass("Zëro.java", "/path/Zëro.java", "class Zëro {}")
	job := &domain.Job{
		Descr:   &domain.Description{Text: "refactor the project", Meta: map[string]any{"max-size": "100"}},
		Classes: []domain.Class{class},
	}
	a := &agent{
		brain:    brain.NewMock(),
		log:      log.NewMock(),
		critic:   &mockCritic{suggestions: []domain.Suggestion{}},
		fixer:    &mockFixer{},
		reviewer: &mockReviewer{},
		attempts: 0,
		frounds:  3,
	}
	artifacts, err := a.Refactor(job)
	require.NoError(t, err, "unexpected error during refactor")
	assert.Equal(t, 0, len(artifacts.Classes), "must return empty classes when attempts is zero")
}

func TestFind_ReturnsClassWhenPathMatches(t *testing.T) {
	class := domain.NewInMemoryClass("Fïnd.java", "/path/Fïnd.java", "class Fïnd {}")
	critiques := []critique{{class: class, suggestions: []domain.Suggestion{}}}
	found, err := find(critiques, "/path/Fïnd.java")
	require.NoError(t, err, "unexpected error during find")
	assert.Equal(t, "/path/Fïnd.java", found.Path(), "found class path must match")
}

func TestFind_ReturnsErrorWhenPathNotFound(t *testing.T) {
	class := domain.NewInMemoryClass("Fïnd.java", "/path/Fïnd.java", "class Fïnd {}")
	critiques := []critique{{class: class, suggestions: []domain.Suggestion{}}}
	_, err := find(critiques, "/path/Nöt.java")
	require.Error(t, err, "expected error when path not found")
}

func TestCritiqueString_FormatsCorrectly(t *testing.T) {
	class := domain.NewInMemoryClass("Strïng.java", "/path/Strïng.java", "class Strïng {}")
	c := &critique{
		class:       class,
		suggestions: []domain.Suggestion{*domain.NewSuggestion("test", "/path/Strïng.java")},
	}
	result := c.String()
	assert.Contains(t, result, "/path/Strïng.java", "string must contain class path")
}

func TestCritiqueString_ShowsSuggestionCount(t *testing.T) {
	class := domain.NewInMemoryClass("Cöünt.java", "/path/Cöünt.java", "class Cöünt {}")
	c := &critique{
		class: class,
		suggestions: []domain.Suggestion{
			*domain.NewSuggestion("one", "/path/Cöünt.java"),
			*domain.NewSuggestion("two", "/path/Cöünt.java"),
		},
	}
	result := c.String()
	assert.Contains(t, result, "2", "string must contain suggestion count")
}

func TestUnderstandClasses_AssociatesSuggestionWithClass(t *testing.T) {
	class := domain.NewInMemoryClass("Mätch.java", "/path/Mätch.java", "class Mätch {}")
	suggestions := []domain.Suggestion{*domain.NewSuggestion("fix it", "/path/Mätch.java")}
	a := &agent{
		brain: brain.NewMock(),
		log:   log.NewMock(),
	}
	result := a.understandClasses([]domain.Class{class}, suggestions)
	require.Equal(t, 1, len(result), "must associate one class")
}

func TestUnderstandClasses_MatchesPartialPaths(t *testing.T) {
	class := domain.NewInMemoryClass("Pärtial.java", "/full/path/Pärtial.java", "class Pärtial {}")
	suggestions := []domain.Suggestion{*domain.NewSuggestion("partial match", "/path/Pärtial.java")}
	a := &agent{
		brain: brain.NewMock(),
		log:   log.NewMock(),
	}
	result := a.understandClasses([]domain.Class{class}, suggestions)
	require.Equal(t, 1, len(result), "must match partial paths")
}

func TestUnderstandClasses_ReturnsEmptyForNoMatch(t *testing.T) {
	class := domain.NewInMemoryClass("Nö.java", "/path/Nö.java", "class Nö {}")
	suggestions := []domain.Suggestion{*domain.NewSuggestion("no match", "/other/Öther.java")}
	a := &agent{
		brain: brain.NewMock(),
		log:   log.NewMock(),
	}
	result := a.understandClasses([]domain.Class{class}, suggestions)
	assert.Equal(t, 0, len(result), "must return empty for no match")
}
