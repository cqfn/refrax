package protocol

type TaskStatus struct {
	State     TaskState `json:"state"`               // Required
	Message   *Message  `json:"message,omitempty"`   // Optional: additional status message
	Timestamp *string   `json:"timestamp,omitempty"` // Optional: ISO 8601 datetime string
}
