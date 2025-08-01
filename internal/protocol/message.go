package protocol

// Kind represents the type of a protocol entity, such as a task or a message.
type Kind string

const (
	// KindTask represents a task entity.
	KindTask Kind = "task"

	// KindMessage represents a message entity.
	KindMessage Kind = "message"
)

// Message represents a communication unit between a user and an agent.
type Message struct {
	Role             string         `json:"role"`
	Parts            Parts          `json:"parts"`
	Metadata         map[string]any `json:"metadata,omitempty"`
	Extensions       []string       `json:"extensions,omitempty"`
	ReferenceTaskIDs []string       `json:"referenceTaskIds,omitempty"`
	MessageID        string         `json:"messageId"`
	TaskID           *string        `json:"taskId,omitempty"`
	ContextID        *string        `json:"contextId,omitempty"`
	Kind             Kind           `json:"kind"`
}

// NewMessage creates and returns a new Message instance with Kind set to KindMessage.
func NewMessage() *Message {
	return &Message{
		Kind: KindMessage,
	}
}

// WithRole sets the Role field of the Message and returns the updated Message instance.
func (m *Message) WithRole(role string) *Message {
	m.Role = role
	return m
}

// WithMessageID sets the MessageID field of the Message and returns the updated Message instance.
func (m *Message) WithMessageID(id string) *Message {
	m.MessageID = id
	return m
}

// AddPart appends a Part to the Parts slice in the Message and returns the updated Message instance.
func (m *Message) AddPart(part Part) *Message {
	m.Parts = append(m.Parts, part)
	return m
}

// AddMetadata adds a key-value pair to the Metadata map of the Message and returns the updated Message instance.
func (m *Message) AddMetadata(key string, value any) *Message {
	if m.Metadata == nil {
		m.Metadata = make(map[string]any)
	}
	m.Metadata[key] = value
	return m
}

// AddExtensions appends one or more extensions to the Extensions slice of the Message and returns the updated Message instance.
func (m *Message) AddExtensions(ext ...string) *Message {
	m.Extensions = append(m.Extensions, ext...)
	return m
}

// AddReferenceTaskIDs appends one or more task IDs to the ReferenceTaskIDs slice of the Message and returns the updated Message instance.
func (m *Message) AddReferenceTaskIDs(ids ...string) *Message {
	m.ReferenceTaskIDs = append(m.ReferenceTaskIDs, ids...)
	return m
}

// WithTaskID sets the TaskID field of the Message and returns the updated Message instance.
func (m *Message) WithTaskID(taskID string) *Message {
	m.TaskID = &taskID
	return m
}

// WithContextID sets the ContextID field of the Message and returns the updated Message instance.
func (m *Message) WithContextID(contextID string) *Message {
	m.ContextID = &contextID
	return m
}
