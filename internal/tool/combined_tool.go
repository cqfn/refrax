package tool

import "strings"

// CombinedTool represents a tool that combines multiple tools into one.
type CombinedTool struct {
	tools []Tool
}

// NewCombined creates a new CombinedTool instance with the provided tools.
func NewCombined(tls ...Tool) Tool {
	return &CombinedTool{tls}
}

// Imperfections gathers and returns the imperfections from all combined tools.
func (c *CombinedTool) Imperfections() string {
	var result strings.Builder
	for pos, t := range c.tools {
		i := strings.TrimSpace(t.Imperfections())
		if i == "" || i == "\n" {
			continue
		}
		result.WriteString(t.Imperfections())
		if pos < len(c.tools)-1 {
			result.WriteString("\n")
		}
	}
	return result.String()
}
