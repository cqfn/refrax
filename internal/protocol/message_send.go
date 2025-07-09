package protocol

type MessageSendParams struct {
	Message       *Message                  `json:"message"`                 // Required
	Configuration *MessageSendConfiguration `json:"configuration,omitempty"` // Optional
	Metadata      map[string]any            `json:"metadata,omitempty"`      // Optional key-value extension metadata
}

type MessageSendConfiguration struct {
	AcceptedOutputModes    []string                `json:"acceptedOutputModes"`              // Required
	HistoryLength          *int                    `json:"historyLength,omitempty"`          // Optional
	PushNotificationConfig *PushNotificationConfig `json:"pushNotificationConfig,omitempty"` // Optional
	Blocking               *bool                   `json:"blocking,omitempty"`               // Optional
}

type MessageSendParamsBuilder struct {
	message       *Message
	configuration *MessageSendConfiguration
	metadata      map[string]any
}

func NewMessageSendParamsBuilder() *MessageSendParamsBuilder {
	return &MessageSendParamsBuilder{}
}

func (b *MessageSendParamsBuilder) Message(msg *Message) *MessageSendParamsBuilder {
	b.message = msg
	return b
}

func (b *MessageSendParamsBuilder) Configuration(cfg *MessageSendConfiguration) *MessageSendParamsBuilder {
	b.configuration = cfg
	return b
}

func (b *MessageSendParamsBuilder) Metadata(meta map[string]any) *MessageSendParamsBuilder {
	b.metadata = meta
	return b
}

func (b *MessageSendParamsBuilder) Build() MessageSendParams {
	return MessageSendParams{
		Message:       b.message,
		Configuration: b.configuration,
		Metadata:      b.metadata,
	}
}
