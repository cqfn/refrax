package protocol

type AgentCardBuilder struct {
	card *AgentCard
}

func Card() *AgentCardBuilder {
	return &AgentCardBuilder{card: &AgentCard{
		Capabilities: AgentCapabilities{},
		Provider: &AgentProvider{
			Organization: "Default Organization",
			URL:          "",
		},
	}}
}

func (b *AgentCardBuilder) Name(name string) *AgentCardBuilder {
	b.card.Name = name
	return b
}

func (b *AgentCardBuilder) Description(desc string) *AgentCardBuilder {
	b.card.Description = desc
	return b
}

func (b *AgentCardBuilder) URL(url string) *AgentCardBuilder {
	b.card.URL = url
	return b
}

func (b *AgentCardBuilder) Provider(provider AgentProvider) *AgentCardBuilder {
	b.card.Provider = &provider
	return b
}

func (b *AgentCardBuilder) Version(version string) *AgentCardBuilder {
	b.card.Version = version
	return b
}

func (b *AgentCardBuilder) DocumentationURL(url string) *AgentCardBuilder {
	b.card.DocumentationURL = &url
	return b
}

func (b *AgentCardBuilder) Capabilities(capabilities AgentCapabilities) *AgentCardBuilder {
	b.card.Capabilities = capabilities
	return b
}

func (b *AgentCardBuilder) DefaultInputModes(modes []string) *AgentCardBuilder {
	b.card.DefaultInputModes = modes
	return b
}

func (b *AgentCardBuilder) DefaultOutputModes(modes []string) *AgentCardBuilder {
	b.card.DefaultOutputModes = modes
	return b
}

func (b *AgentCardBuilder) Skill(id, name, description string) *AgentCardBuilder {
	var skill = AgentSkill{
		Name:        name,
		Description: description,
	}
	b.card.Skills = append(b.card.Skills, skill)
	return b
}

func (b *AgentCardBuilder) Skills(skills []AgentSkill) *AgentCardBuilder {
	b.card.Skills = skills
	return b
}

func (b *AgentCardBuilder) Build() AgentCard {
	return *b.card
}
