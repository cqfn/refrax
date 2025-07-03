package aibolit

import (
	"regexp"
	"strings"
)

type SanitizedAibolit struct {
	origin Aibolit
}

func NewSanitizedAibolit(aibolit Aibolit) Aibolit {
	return &SanitizedAibolit{aibolit}
}

func (r *SanitizedAibolit) Imperfections() string {
	complaint := regexp.MustCompile(`^[^:]+\.java\[\d+\]: .+`)
	var complaints []string
	for line := range strings.SplitSeq(r.origin.Imperfections(), "\n") {
		if complaint.MatchString(line) {
			complaints = append(complaints, line)
		}
	}
	return strings.Join(complaints, "\n")
}
