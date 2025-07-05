package critic

type MockTool struct {
	data string
}

func NewMockTool(data string) Tool {
	return &MockTool{data}
}

func NewMockToolEmpty() Tool {
	return &MockTool{string([]byte{})}
}

func (a *MockTool) Imperfections() string {
	return a.data
}
