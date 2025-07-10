package critic

// MockTool represents a mock implementation of the Tool interface.
type MockTool struct {
	data string
}

// NewMockTool creates a new instance of MockTool with the provided data.
func NewMockTool(data string) Tool {
	return &MockTool{data}
}

// NewMockToolEmpty creates a new instance of MockTool with empty data.
func NewMockToolEmpty() Tool {
	return &MockTool{string([]byte{})}
}

// Imperfections returns the imperfections associated with the MockTool.
func (a *MockTool) Imperfections() string {
	return a.data
}
