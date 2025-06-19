package protocol

type Message struct {
	Role             string           `json:"role"`                       // "user" or "agent"
	Parts            []map[string]any `json:"parts"`                      // Required message content
	Metadata         map[string]any   `json:"metadata,omitempty"`         // Optional extension metadata
	Extensions       []string         `json:"extensions,omitempty"`       // Optional list of extension URIs
	ReferenceTaskIDs []string         `json:"referenceTaskIds,omitempty"` // Optional task references
	MessageID        string           `json:"messageId"`                  // Required message ID
	TaskID           *string          `json:"taskId,omitempty"`           // Optional task ID
	ContextID        *string          `json:"contextId,omitempty"`        // Optional context ID
	Kind             string           `json:"kind"`                       // Must be "message"
}

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

type PushNotificationConfig struct {
	ID             *string                             `json:"id,omitempty"`             // Optional
	URL            string                              `json:"url"`                      // Required
	Token          *string                             `json:"token,omitempty"`          // Optional
	Authentication *PushNotificationAuthenticationInfo `json:"authentication,omitempty"` // Optional
}

type PushNotificationAuthenticationInfo struct {
	Schemes     []string `json:"schemes"`               // Required
	Credentials *string  `json:"credentials,omitempty"` // Optional
}
