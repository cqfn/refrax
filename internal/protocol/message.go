package protocol

type Kind string

const (
	KindTask    Kind = "task"
	KindMessage Kind = "message"
)

type Message struct {
	Role             string         `json:"role"`                       // "user" or "agent"
	Parts            Parts          `json:"parts"`                      // Required message content
	Metadata         map[string]any `json:"metadata,omitempty"`         // Optional extension metadata
	Extensions       []string       `json:"extensions,omitempty"`       // Optional list of extension URIs
	ReferenceTaskIDs []string       `json:"referenceTaskIds,omitempty"` // Optional task references
	MessageID        string         `json:"messageId"`                  // Required message ID
	TaskID           *string        `json:"taskId,omitempty"`           // Optional task ID
	ContextID        *string        `json:"contextId,omitempty"`        // Optional context ID
	Kind             Kind           `json:"kind"`                       // Must be "message"
}

type MessageBuilder struct {
	msg Message
}

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{
		msg: Message{
			Kind: KindMessage,
		},
	}
}

func (b *MessageBuilder) Role(role string) *MessageBuilder {
	b.msg.Role = role
	return b
}

func (b *MessageBuilder) MessageID(id string) *MessageBuilder {
	b.msg.MessageID = id
	return b
}

func (b *MessageBuilder) Part(part Part) *MessageBuilder {
	b.msg.Parts = append(b.msg.Parts, part)
	return b
}

func (b *MessageBuilder) MetadataField(key string, value any) *MessageBuilder {
	if b.msg.Metadata == nil {
		b.msg.Metadata = make(map[string]any)
	}
	b.msg.Metadata[key] = value
	return b
}

func (b *MessageBuilder) Extensions(ext ...string) *MessageBuilder {
	b.msg.Extensions = append(b.msg.Extensions, ext...)
	return b
}

func (b *MessageBuilder) ReferenceTaskIDs(ids ...string) *MessageBuilder {
	b.msg.ReferenceTaskIDs = append(b.msg.ReferenceTaskIDs, ids...)
	return b
}

func (b *MessageBuilder) TaskID(taskID string) *MessageBuilder {
	b.msg.TaskID = &taskID
	return b
}

func (b *MessageBuilder) ContextID(contextID string) *MessageBuilder {
	b.msg.ContextID = &contextID
	return b
}

func (b *MessageBuilder) Build() Message {
	return b.msg
}
