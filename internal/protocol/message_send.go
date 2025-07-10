package protocol

// MessageSendParams defines the parameters for sending a message to an agent.
type MessageSendParams struct {
	Message       *Message                  `json:"message"`                 // Required
	Configuration *MessageSendConfiguration `json:"configuration,omitempty"` // Optional
	Metadata      map[string]any            `json:"metadata,omitempty"`      // Optional key-value extension metadata
}

// MessageSendConfiguration defines the configuration for sending a message.
type MessageSendConfiguration struct {
	AcceptedOutputModes    []string                `json:"acceptedOutputModes"`              // Required
	HistoryLength          *int                    `json:"historyLength,omitempty"`          // Optional
	PushNotificationConfig *PushNotificationConfig `json:"pushNotificationConfig,omitempty"` // Optional
	Blocking               *bool                   `json:"blocking,omitempty"`               // Optional
}

// MessageSendParamsBuilder defines a builder for constructing MessageSendParams.
type MessageSendParamsBuilder struct {
	message       *Message
	configuration *MessageSendConfiguration
	metadata      map[string]any
}

// NewMessageSendParamsBuilder creates a new instance of MessageSendParamsBuilder.
func NewMessageSendParamsBuilder() *MessageSendParamsBuilder {
	return &MessageSendParamsBuilder{}
}

// Message sets the message to be sent in the MessageSendParamsBuilder.
func (b *MessageSendParamsBuilder) Message(msg *Message) *MessageSendParamsBuilder {
	b.message = msg
	return b
}

// Configuration sets the configuration for sending the message in the MessageSendParamsBuilder.
func (b *MessageSendParamsBuilder) Configuration(cfg *MessageSendConfiguration) *MessageSendParamsBuilder {
	b.configuration = cfg
	return b
}

// Metadata sets the metadata for the message in the MessageSendParamsBuilder.
func (b *MessageSendParamsBuilder) Metadata(meta map[string]any) *MessageSendParamsBuilder {
	b.metadata = meta
	return b
}

// Build constructs the MessageSendParams from the builder.
func (b *MessageSendParamsBuilder) Build() MessageSendParams {
	return MessageSendParams{
		Message:       b.message,
		Configuration: b.configuration,
		Metadata:      b.metadata,
	}
}
