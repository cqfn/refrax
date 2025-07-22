package client

import "io"

// Params holds the configuration parameters for Refrax commands.
type Params struct {
	Provider    string
	Token       string
	Playbook    string
	MockProject bool
	Debug       bool
	Stats       bool
	Format      string
	Soutput     string
	Input       string
	Output      string
	Log         io.Writer
}

// NewMockParams creates a new Params object with mock settings.
func NewMockParams() *Params {
	return &Params{
		Provider:    "mock",
		Token:       "ABC",
		Playbook:    "",
		MockProject: true,
		Debug:       false,
		Stats:       false,
		Format:      "std",
		Soutput:     "stats",
		Input:       "",
		Output:      "",
		Log:         io.Discard, // Default to discard if no logging is needed
	}
}
