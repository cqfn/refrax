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

// NewMessageSendParams creates a new instance of MessageSendParams for building.
func NewMessageSendParams() *MessageSendParams {
	return &MessageSendParams{}
}

// WithMessage sets the message to be sent in the MessageSendParams.
func (p *MessageSendParams) WithMessage(msg *Message) *MessageSendParams {
	p.Message = msg
	return p
}

// WithConfiguration sets the configuration for sending the message in the MessageSendParams.
func (p *MessageSendParams) WithConfiguration(cfg *MessageSendConfiguration) *MessageSendParams {
	p.Configuration = cfg
	return p
}

// WithMetadata sets the metadata for the message in the MessageSendParams.
func (p *MessageSendParams) WithMetadata(meta map[string]any) *MessageSendParams {
	p.Metadata = meta
	return p
}
