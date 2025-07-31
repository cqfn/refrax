package tool

// MockTool represents a mock implementation of the Tool interface.
type MockTool struct {
	data string
}

// NewMock creates a new instance of MockTool with the provided data.
func NewMock(data string) Tool {
	return &MockTool{data}
}

// NewEmpty creates a new instance of MockTool with empty data.
func NewEmpty() Tool {
	return &MockTool{string([]byte{})}
}

// Imperfections returns the imperfections associated with the MockTool.
func (a *MockTool) Imperfections() string {
	return a.data
}
