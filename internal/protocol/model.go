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

// NewAgentCard creates and returns a new AgentCard with default values.
func NewAgentCard() *AgentCard {
	return &AgentCard{
		Capabilities: AgentCapabilities{},
		Provider: &AgentProvider{
			Organization: "Default Organization",
			URL:          "",
		},
	}
}

// WithName sets the name of the agent card.
func (c *AgentCard) WithName(name string) *AgentCard {
	c.Name = name
	return c
}

// WithDescription sets the description of the agent card.
func (c *AgentCard) WithDescription(desc string) *AgentCard {
	c.Description = desc
	return c
}

// WithURL sets the URL of the agent.
func (c *AgentCard) WithURL(url string) *AgentCard {
	c.URL = url
	return c
}

// WithProvider sets the provider information of the agent.
func (c *AgentCard) WithProvider(provider AgentProvider) *AgentCard {
	c.Provider = &provider
	return c
}

// WithVersion sets the version string of the agent.
func (c *AgentCard) WithVersion(version string) *AgentCard {
	c.Version = version
	return c
}

// WithDocumentationURL sets the optional documentation URL for the agent.
func (c *AgentCard) WithDocumentationURL(url string) *AgentCard {
	c.DocumentationURL = &url
	return c
}

// WithCapabilities sets the agent's declared capabilities.
func (c *AgentCard) WithCapabilities(capabilities AgentCapabilities) *AgentCard {
	c.Capabilities = capabilities
	return c
}

// WithDefaultInputModes sets the default input modes supported by the agent.
func (c *AgentCard) WithDefaultInputModes(modes []string) *AgentCard {
	c.DefaultInputModes = modes
	return c
}

// WithDefaultOutputModes sets the default output modes supported by the agent.
func (c *AgentCard) WithDefaultOutputModes(modes []string) *AgentCard {
	c.DefaultOutputModes = modes
	return c
}

// AddSkill appends a single skill to the agent's skill list.
func (c *AgentCard) AddSkill(_, name, description string) *AgentCard {
	c.Skills = append(c.Skills, AgentSkill{
		Name:        name,
		Description: description,
	})
	return c
}

// WithSkills sets the list of skills for the agent.
func (c *AgentCard) WithSkills(skills []AgentSkill) *AgentCard {
	c.Skills = skills
	return c
}
