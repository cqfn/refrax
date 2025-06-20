package protocol

import (
	"encoding/json"
	"fmt"
)

type Kind string

const (
	KindTask    Kind = "task"
	KindMessage Kind = "message"
)

type PartKind string

const (
	PartKindText PartKind = "text"
	PartKindFile PartKind = "file"
	PartKindData PartKind = "data"
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

type Parts []Part

func (p *Parts) UnmarshalJSON(data []byte) error {
	var rawParts []json.RawMessage
	if err := json.Unmarshal(data, &rawParts); err != nil {
		return err
	}
	for _, raw := range rawParts {
		var kindHolder struct {
			Kind PartKind `json:"kind"`
		}
		if err := json.Unmarshal(raw, &kindHolder); err != nil {
			return err
		}

		var part Part
		switch kindHolder.Kind {
		case PartKindText:
			var tp TextPart
			if err := json.Unmarshal(raw, &tp); err != nil {
				return err
			}
			part = &tp
		case PartKindFile:
			var fp FilePart
			if err := json.Unmarshal(raw, &fp); err != nil {
				return err
			}
			part = &fp
		case PartKindData:
			var dp DataPart
			if err := json.Unmarshal(raw, &dp); err != nil {
				return err
			}
			part = &dp
		default:
			return fmt.Errorf("unknown part kind: %s", kindHolder.Kind)
		}

		*p = append(*p, part)
	}
	return nil
}

type Part interface {
	PartKind() PartKind
}

type PartBase struct {
}

type TextPart struct {
	PartBase
	Kind PartKind `json:"kind"` // Must be "text"
	Text string   `json:"text"` // Required text content
}

func (p *TextPart) PartKind() PartKind {
	return p.Kind
}

type FilePart struct {
	PartBase
	Kind PartKind `json:"kind"` // Must be "file"
	File any      `json:"file"` // Can be FileWithBytes or FileWithUri
}

func (p *FilePart) PartKind() PartKind {
	return p.Kind
}

type DataPart struct {
	PartBase
	Kind PartKind       `json:"kind"` // Must be "data"
	Data map[string]any `json:"data"` // Required structured content
}

func (p *DataPart) PartKind() PartKind {
	return p.Kind
}

type FileBase struct {
}

type FileWithBytes struct {
	FileBase
	Bytes string `json:"bytes"` // Required: base64-encoded content
}

type FileWithURI struct {
	FileBase
	URI string `json:"uri"` // Required: URI to the file
}

func (m Message) MarshalJSON() ([]byte, error) {
	type Alias Message
	m.Kind = "message"
	return json.Marshal((Alias)(m))
}

func (p *FilePart) UnmarshalJSON(data []byte) error {
	// Define an alias to avoid recursion and parse the base part
	type Alias FilePart
	aux := &struct {
		File json.RawMessage `json:"file"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Probe the file kind
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(aux.File, &probe); err != nil {
		return err
	}

	// Detect and decode the actual file content
	switch {
	case probe["bytes"] != nil:
		var fwb FileWithBytes
		if err := json.Unmarshal(aux.File, &fwb); err != nil {
			return err
		}
		p.File = fwb
	case probe["uri"] != nil:
		var fwu FileWithURI
		if err := json.Unmarshal(aux.File, &fwu); err != nil {
			return err
		}
		p.File = fwu
	default:
		return fmt.Errorf("unknown file format in FilePart")
	}

	return nil
}
