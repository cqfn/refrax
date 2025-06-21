package protocol

type Task struct {
	ID        string         `json:"id"`                  // Required: task ID
	ContextID string         `json:"contextId"`           // Required: contextual alignment
	Status    TaskStatus     `json:"status"`              // Required: current status
	History   []Message      `json:"history,omitempty"`   // Optional: message history
	Artifacts []Artifact     `json:"artifacts,omitempty"` // Optional: artifacts created
	Metadata  map[string]any `json:"metadata,omitempty"`  // Optional: extension metadata
	Kind      Kind           `json:"kind"`                // Required: must be "task"
}
