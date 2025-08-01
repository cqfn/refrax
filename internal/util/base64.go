// Package util provides utility functions for various operations.
package util

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// DecodeFile decodes a base64 encoded string, trims whitespace, and returns the decoded string.
// If decoding fails, it returns an error.
func DecodeFile(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 string: %v", err)
	}
	return strings.TrimSpace(string(decoded)), err
}
