package brain

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type qa struct {
	Question string `yaml:"question"`
	Answer   string `yaml:"answer"`
}

type yamlPlaybook struct {
	Name string `yaml:"name"`
	QA   []qa   `yaml:"qa"`
}

// YAMLPlaybook represents a playbook loaded from a YAML file.
type YAMLPlaybook struct {
	data map[string]string
}

// NewYAMLPlaybook loads a YAML playbook from the specified file path and returns a YAMLPlaybook instance.
func NewYAMLPlaybook(filePath string) (*YAMLPlaybook, error) {
	content, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}
	var playbook yamlPlaybook
	err = yaml.Unmarshal(content, &playbook)
	if err != nil {
		return nil, err
	}
	data := make(map[string]string)
	for _, qa := range playbook.QA {
		data[normalise(qa.Question)] = strings.TrimSpace(qa.Answer)
	}
	return &YAMLPlaybook{data: data}, nil
}

// Ask retrieves the answer to a given question from the playbook.
// If the question is not found, it returns a default "not found" message.
func (p *YAMLPlaybook) Ask(question string) string {
	nq := normalise(question)
	for k, v := range p.data {
		if strings.Contains(nq, k) || strings.Contains(k, nq) {
			return v
		}
	}
	return "Question not found in the playbook"
}

func normalise(raw string) string {
	return strings.Join(strings.Fields(raw), "")
}
