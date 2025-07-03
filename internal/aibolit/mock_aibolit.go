package aibolit

type MockAibolit struct {
	data string
}

func NewMockAibolit(data string) Aibolit {
	return &MockAibolit{data}
}

func NewMockAibolitEmpty() Aibolit {
	return &MockAibolit{string([]byte{})}
}

func (a *MockAibolit) Imperfections() string {
	return a.data
}
