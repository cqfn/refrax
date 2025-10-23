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
	MaxSize     int
	Log         io.Writer
	Checks      []string
	Colorless   bool
	Model       string
	Attempts    int
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
		MaxSize:     200,
		Log:         io.Discard,
		Checks:      []string{"mvn clean test"},
		Colorless:   false,
		Model:       "gpt-3.5-turbo",
		Attempts:    3,
	}
}
