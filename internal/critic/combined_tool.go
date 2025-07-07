package critic

import "strings"

type CombinedTool struct {
	tools []Tool
}

func NewCombinedTool(tls ...Tool) Tool {
	return &CombinedTool{tls}
}

func (c *CombinedTool) Imperfections() string {
	var result strings.Builder
	for pos, t := range c.tools {
		result.WriteString(t.Imperfections())
		if pos < len(c.tools)-1 {
			result.WriteString("\n")
		}
	}
	return result.String()
}
