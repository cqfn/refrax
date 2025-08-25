package domain

// OldSuggestion represents a proposed improvement or fix for a class.
type OldSuggestion interface {
	Text() string
}

type oldSuggestion struct {
	text string
}

// NewOldSuggestion creates a new Suggestion instance with the given text.
func NewOldSuggestion(text string) OldSuggestion {
	return &oldSuggestion{
		text: text,
	}
}

func NewSuggestion(text, path string) *Suggestion {
	return &Suggestion{
		Text:      text,
		ClassPath: path,
	}
}

// Text implements Suggestion.
func (a *oldSuggestion) Text() string {
	return a.text
}
