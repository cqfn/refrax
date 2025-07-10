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

// MessageBuilder provides a fluent API for constructing Message objects.
type MessageBuilder struct {
	msg *Message
}

// NewMessageBuilder creates and returns a new MessageBuilder instance.
func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{
		msg: &Message{
			Kind: KindMessage,
		},
	}
}

// Role sets the role of the message sender (e.g., "user" or "agent").
func (b *MessageBuilder) Role(role string) *MessageBuilder {
	b.msg.Role = role
	return b
}

// MessageID sets the unique identifier for the message.
func (b *MessageBuilder) MessageID(id string) *MessageBuilder {
	b.msg.MessageID = id
	return b
}

// Part appends a message part to the message content.
func (b *MessageBuilder) Part(part Part) *MessageBuilder {
	b.msg.Parts = append(b.msg.Parts, part)
	return b
}

// MetadataField sets a key-value pair in the message metadata.
func (b *MessageBuilder) MetadataField(key string, value any) *MessageBuilder {
	if b.msg.Metadata == nil {
		b.msg.Metadata = make(map[string]any)
	}
	b.msg.Metadata[key] = value
	return b
}

// Extensions appends one or more extension URIs to the message.
func (b *MessageBuilder) Extensions(ext ...string) *MessageBuilder {
	b.msg.Extensions = append(b.msg.Extensions, ext...)
	return b
}

// ReferenceTaskIDs appends one or more task IDs to the reference list.
func (b *MessageBuilder) ReferenceTaskIDs(ids ...string) *MessageBuilder {
	b.msg.ReferenceTaskIDs = append(b.msg.ReferenceTaskIDs, ids...)
	return b
}

// TaskID sets the optional task ID associated with the message.
func (b *MessageBuilder) TaskID(taskID string) *MessageBuilder {
	b.msg.TaskID = &taskID
	return b
}

// ContextID sets the optional context ID associated with the message.
func (b *MessageBuilder) ContextID(contextID string) *MessageBuilder {
	b.msg.ContextID = &contextID
	return b
}

// Build finalizes and returns the constructed Message.
func (b *MessageBuilder) Build() *Message {
	return b.msg
}
