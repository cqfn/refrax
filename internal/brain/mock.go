package brain

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cqfn/refrax/internal/log"
)

const known = "public class Main {\n\tpublic static void main(String[] args) {\n\t\tString m = \"Hello, World\";\n\t\tSystem.out.println(m);\n\t}\n}\n\n"
const refactored = "public class Main {\n\tpublic static void main(String[] args) {\n\t\tSystem.out.println(\"Hello, World\");\n\t}\n"

type MockBrain struct {
}

func NewMock() Brain {
	return &MockBrain{}
}

func (b *MockBrain) Ask(question string) (string, error) {
	if question == "" {
		return "", fmt.Errorf("question cannot be empty")
	}
	blocks := javaCode(question)
	if len(blocks) == 0 {
		log.Info("mock-brain: no Java code found in the question, returning mock response")
		return "mock response to: " + question, nil
	} else if strings.Contains(trimSpace(blocks[0]), trimSpace(known)) {
		log.Info("mock-brain: known Java code found, returning refactored code")
		return refactored, nil
	} else {
		log.Info("mock-brain: unknown Java code found, returning first block as response (echo)")
		return blocks[0], nil
	}
}

func trimSpace(s string) string {
	return strings.TrimSpace(s)
}

func javaCode(markdown string) []string {
	re := regexp.MustCompile("(?s)```java\\s+(.*?)```")
	matches := re.FindAllStringSubmatch(markdown, -1)
	var blocks []string
	for _, match := range matches {
		if len(match) > 1 {
			blocks = append(blocks, match[1])
		}
	}
	return blocks
}
