package tool

// Tool defines the interface for a tool that can be used to identify and report imperfections in artifacts.
type Tool interface {
	Imperfections() string
}
