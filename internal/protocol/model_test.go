package protocol

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func secureRandomInt(limit int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(limit)))
	return int(n.Int64())
}

func TestNewAgentCard_CreatesValidInstance(t *testing.T) {
	card := NewAgentCard()
	require.NotNil(t, card, "agent card must not be nil")
}

func TestNewAgentCard_HasDefaultProvider(t *testing.T) {
	card := NewAgentCard()
	require.NotNil(t, card.Provider, "provider must not be nil")
	assert.Equal(t, "Default Organization", card.Provider.Organization, "default organization must be set")
}

func TestWithName_SetsName(t *testing.T) {
	name := "Tëst Ägënt"
	card := NewAgentCard().WithName(name)
	assert.Equal(t, name, card.Name, "name must match")
}

func TestWithDescription_SetsDescription(t *testing.T) {
	desc := "Ünïcödé dëscrïptïön"
	card := NewAgentCard().WithDescription(desc)
	assert.Equal(t, desc, card.Description, "description must match")
}

func TestWithURL_SetsURL(t *testing.T) {
	url := "http://ëxämplë.cöm:8080"
	card := NewAgentCard().WithURL(url)
	assert.Equal(t, url, card.URL, "url must match")
}

func TestWithProvider_SetsProvider(t *testing.T) {
	provider := AgentProvider{Organization: "Tëst Örg", URL: "http://örg.cöm"}
	card := NewAgentCard().WithProvider(provider)
	require.NotNil(t, card.Provider, "provider must not be nil")
	assert.Equal(t, "Tëst Örg", card.Provider.Organization, "organization must match")
}

func TestWithVersion_SetsVersion(t *testing.T) {
	version := "1.2.3-bëtä"
	card := NewAgentCard().WithVersion(version)
	assert.Equal(t, version, card.Version, "version must match")
}

func TestWithDocumentationURL_SetsDocURL(t *testing.T) {
	docURL := "http://döcs.ëxämplë.cöm"
	card := NewAgentCard().WithDocumentationURL(docURL)
	require.NotNil(t, card.DocumentationURL, "documentation url must not be nil")
	assert.Equal(t, docURL, *card.DocumentationURL, "documentation url must match")
}

func TestWithCapabilities_SetsCapabilities(t *testing.T) {
	streaming := true
	caps := AgentCapabilities{Streaming: &streaming}
	card := NewAgentCard().WithCapabilities(caps)
	require.NotNil(t, card.Capabilities.Streaming, "streaming capability must be set")
	assert.True(t, *card.Capabilities.Streaming, "streaming must be true")
}

func TestWithDefaultInputModes_SetsModes(t *testing.T) {
	modes := []string{"tëxt", "fïlë"}
	card := NewAgentCard().WithDefaultInputModes(modes)
	assert.Equal(t, modes, card.DefaultInputModes, "input modes must match")
}

func TestWithDefaultOutputModes_SetsModes(t *testing.T) {
	modes := []string{"jSön", "xml"}
	card := NewAgentCard().WithDefaultOutputModes(modes)
	assert.Equal(t, modes, card.DefaultOutputModes, "output modes must match")
}

func TestAddSkill_AppendsSkill(t *testing.T) {
	card := NewAgentCard().AddSkill("ïd", "Skïll Nämë", "Skïll Dëscrïptïön")
	require.Equal(t, 1, len(card.Skills), "must have one skill")
	assert.Equal(t, "Skïll Nämë", card.Skills[0].Name, "skill name must match")
}

func TestAddSkill_AppendsMultipleSkills(t *testing.T) {
	card := NewAgentCard().
		AddSkill("ïd1", "Skïll Önë", "Dësc Önë").
		AddSkill("ïd2", "Skïll Twö", "Dësc Twö")
	assert.Equal(t, 2, len(card.Skills), "must have two skills")
}

func TestWithSkills_SetsSkillsList(t *testing.T) {
	skills := []AgentSkill{
		{ID: "ïd", Name: "Nämë", Description: "Dësc", Tags: []string{"täg1"}},
	}
	card := NewAgentCard().WithSkills(skills)
	assert.Equal(t, skills, card.Skills, "skills must match")
}

func TestChainedMethods_BuildCompleteCard(t *testing.T) {
	port := secureRandomInt(10000) + 20000
	card := NewAgentCard().
		WithName("Cömplëtë Ägënt").
		WithDescription("Ä cömplëtë ägënt cärd").
		WithURL("http://löcälhöst:" + string(rune(port))).
		WithVersion("0.0.1")
	require.NotNil(t, card, "card must not be nil")
	assert.NotEmpty(t, card.Name, "name must be set")
}

func TestAgentSkill_HasRequiredFields(t *testing.T) {
	skill := AgentSkill{
		ID:          "tëst-ïd",
		Name:        "Tëst Skïll",
		Description: "Ä tëst skïll",
		Tags:        []string{"täg1", "täg2"},
	}
	assert.Equal(t, "tëst-ïd", skill.ID, "id must match")
	assert.Equal(t, "Tëst Skïll", skill.Name, "name must match")
}

func TestAgentExtension_HasRequiredFields(t *testing.T) {
	desc := "Ëxtënsïön dësc"
	required := true
	ext := AgentExtension{
		URI:         "urn:ëxämplë:ëxt",
		Description: &desc,
		Required:    &required,
		Params:      map[string]any{"pärämëtër": "välüë"},
	}
	assert.Equal(t, "urn:ëxämplë:ëxt", ext.URI, "uri must match")
	assert.True(t, *ext.Required, "required must be true")
}
