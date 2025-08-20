package domain

// Suggestion represents a proposed improvement or fix for a class.
type Suggestion interface {
	Text() string
}

type suggestion struct {
	text string
}

// NewSuggestion creates a new Suggestion instance with the given text.
func NewSuggestion(text string) Suggestion {
	return &suggestion{
		text: text,
	}
}

// Text implements Suggestion.
func (a *suggestion) Text() string {
	return a.text
}
