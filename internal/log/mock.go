package log

import "fmt"

type Mock struct {
	Messages []string
}

func NewMock() *Mock {
	return &Mock{Messages: []string{}}
}

func (m *Mock) Info(msg string, args ...any) {
	m.Messages = append(m.Messages, "mock info: "+fmt.Sprintf(msg, args...))
}

func (m *Mock) Debug(msg string, args ...any) {
	m.Messages = append(m.Messages, fmt.Sprintf("mock debug: %s", fmt.Sprintf(msg, args...)))
}

func (m *Mock) Warn(msg string, args ...any) {
	m.Messages = append(m.Messages, fmt.Sprintf("mock warn: %s", fmt.Sprintf(msg, args...)))
}

func (m *Mock) Error(msg string, args ...any) {
	m.Messages = append(m.Messages, fmt.Sprintf("mock error: %s", fmt.Sprintf(msg, args...)))
}
