// Package protocol A2A Agent Card.
// Current version: 0.2.2
package protocol

// AgentCard represents the A2A Agent Card as defined in the A2A specification.
// SecurityScheme isn't implemented yet, but should be defined as per the A2A spec.
// Check docs for it:
// https://google-a2a.github.io/A2A/latest/specification/#553-securityscheme-object
type AgentCard struct {
	Name                              string                `json:"name"`
	Description                       string                `json:"description"`
	URL                               string                `json:"url"`
	IconURL                           *string               `json:"iconUrl,omitempty"`
	Provider                          *AgentProvider        `json:"provider,omitempty"`
	Version                           string                `json:"version"`
	DocumentationURL                  *string               `json:"documentationUrl,omitempty"`
	Capabilities                      AgentCapabilities     `json:"capabilities"`
	SecuritySchemes                   map[string]string     `json:"securitySchemes,omitempty"`
	Security                          []map[string][]string `json:"security,omitempty"`
	DefaultInputModes                 []string              `json:"defaultInputModes"`
	DefaultOutputModes                []string              `json:"defaultOutputModes"`
	Skills                            []AgentSkill          `json:"skills"`
	SupportsAuthenticatedExtendedCard *bool                 `json:"supportsAuthenticatedExtendedCard,omitempty"`
}

// AgentProvider represents the organization and URL of the agent provider.
type AgentProvider struct {
	Organization string `json:"organization"`
	URL          string `json:"url"`
}

// AgentCapabilities defines the capabilities of the agent as per the A2A specification.
type AgentCapabilities struct {
	Streaming              *bool            `json:"streaming,omitempty"`
	PushNotifications      *bool            `json:"pushNotifications,omitempty"`
	StateTransitionHistory *bool            `json:"stateTransitionHistory,omitempty"`
	Extensions             []AgentExtension `json:"extensions,omitempty"`
}

// AgentExtension represents an extension that can be used by the agent.
type AgentExtension struct {
	URI         string         `json:"uri"`
	Description *string        `json:"description,omitempty"`
	Required    *bool          `json:"required,omitempty"`
	Params      map[string]any `json:"params,omitempty"`
}

// AgentSkill represents a skill that the agent can perform.
type AgentSkill struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Examples    []string `json:"examples,omitempty"`
	InputModes  []string `json:"inputModes,omitempty"`
	OutputModes []string `json:"outputModes,omitempty"`
}
