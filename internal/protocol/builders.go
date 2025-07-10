package protocol

// CardBuilder is a builder for constructing AgentCard objects.
type CardBuilder struct {
	card *AgentCard
}

// Card creates and returns a new cardBuilder with default AgentCard values.
func Card() *CardBuilder {
	return &CardBuilder{card: &AgentCard{
		Capabilities: AgentCapabilities{},
		Provider: &AgentProvider{
			Organization: "Default Organization",
			URL:          "",
		},
	}}
}

// Name sets the name of the agent card.
func (b *CardBuilder) Name(name string) *CardBuilder {
	b.card.Name = name
	return b
}

// Description sets the description of the agent card.
func (b *CardBuilder) Description(desc string) *CardBuilder {
	b.card.Description = desc
	return b
}

// URL sets the URL of the agent.
func (b *CardBuilder) URL(url string) *CardBuilder {
	b.card.URL = url
	return b
}

// Provider sets the provider information of the agent.
func (b *CardBuilder) Provider(provider AgentProvider) *CardBuilder {
	b.card.Provider = &provider
	return b
}

// Version sets the version string of the agent.
func (b *CardBuilder) Version(version string) *CardBuilder {
	b.card.Version = version
	return b
}

// DocumentationURL sets the optional documentation URL for the agent.
func (b *CardBuilder) DocumentationURL(url string) *CardBuilder {
	b.card.DocumentationURL = &url
	return b
}

// Capabilities sets the agent's declared capabilities.
func (b *CardBuilder) Capabilities(capabilities AgentCapabilities) *CardBuilder {
	b.card.Capabilities = capabilities
	return b
}

// DefaultInputModes sets the default input modes supported by the agent.
func (b *CardBuilder) DefaultInputModes(modes []string) *CardBuilder {
	b.card.DefaultInputModes = modes
	return b
}

// DefaultOutputModes sets the default output modes supported by the agent.
func (b *CardBuilder) DefaultOutputModes(modes []string) *CardBuilder {
	b.card.DefaultOutputModes = modes
	return b
}

// Skill appends a single skill to the agent's skill list.
func (b *CardBuilder) Skill(_, name, description string) *CardBuilder {
	skill := AgentSkill{
		Name:        name,
		Description: description,
	}
	b.card.Skills = append(b.card.Skills, skill)
	return b
}

// Skills sets the list of skills for the agent.
func (b *CardBuilder) Skills(skills []AgentSkill) *CardBuilder {
	b.card.Skills = skills
	return b
}

// Build returns the constructed AgentCard.
func (b *CardBuilder) Build() *AgentCard {
	return b.card
}
