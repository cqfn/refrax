package brain

import (
	"fmt"
	"regexp"
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
		return "mock response to: " + question, nil
	} else if blocks[0] == known {
		return refactored, nil
	} else {
		return blocks[0], nil
	}
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
