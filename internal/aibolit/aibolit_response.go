package aibolit

import (
	"regexp"
	"strings"
)

type AibolitResponse struct {
	data string
}

func NewAibolitResponse(data string) *AibolitResponse {
	return &AibolitResponse{data}
}

func (r *AibolitResponse) Sanitized() string {
	complaint := regexp.MustCompile(`^[^:]+\.java\[\d+\]: .+`)
	var complaints []string
	for _, line := range strings.Split(r.data, "\n") {
		if complaint.MatchString(line) {
			complaints = append(complaints, line)
		}
	}
	return strings.Join(complaints, "\n")
}
