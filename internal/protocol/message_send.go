package protocol

type MessageSendParams struct {
	Message       Message                   `json:"message"`                 // Required
	Configuration *MessageSendConfiguration `json:"configuration,omitempty"` // Optional
	Metadata      map[string]any            `json:"metadata,omitempty"`      // Optional key-value extension metadata
}

type MessageSendConfiguration struct {
	AcceptedOutputModes    []string                `json:"acceptedOutputModes"`              // Required
	HistoryLength          *int                    `json:"historyLength,omitempty"`          // Optional
	PushNotificationConfig *PushNotificationConfig `json:"pushNotificationConfig,omitempty"` // Optional
	Blocking               *bool                   `json:"blocking,omitempty"`               // Optional
}
