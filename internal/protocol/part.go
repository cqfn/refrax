package protocol

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// PartKind represents the type of a part in a message.
type PartKind string

// PartBase represents the base structure for a part, containing common metadata.
type PartBase struct {
	Metadata map[string]any `json:"metadata,omitempty"`
}

const (
	// PartKindText represents a text part.
	PartKindText PartKind = "text"

	// PartKindFile represents a file part.
	PartKindFile PartKind = "file"

	// PartKindData represents a data part.
	PartKindData PartKind = "data"
)

// Part is an interface that defines the common behavior for all parts in a message.
type Part interface {
	PartKind() PartKind
	Metadata() map[string]any
}

// TextPart represents a text part in a message.
type TextPart struct {
	Kind PartKind `json:"kind"` // Must be "text"
	Text string   `json:"text"` // Required text content
	PartBase
}

// FilePart represents a file part in a message.
type FilePart struct {
	Kind PartKind `json:"kind"` // Must be "file"
	File any      `json:"file"` // Can be FileWithBytes or FileWithUri
	PartBase
}

// DataPart represents a data part in a message.
type DataPart struct {
	Kind PartKind       `json:"kind"` // Must be "data"
	Data map[string]any `json:"data"` // Required structured content
	PartBase
}

// Parts represents a collection of Part interfaces.
type Parts []Part

// PartKind returns the kind of the part.
func (p *TextPart) PartKind() PartKind {
	return p.Kind
}

// PartKind returns the kind of the part.
func (p *FilePart) PartKind() PartKind {
	return p.Kind
}

// PartKind returns the kind of the part.
func (p *DataPart) PartKind() PartKind {
	return p.Kind
}

// Metadata returns the metadata of the TextPart.
func (p *TextPart) Metadata() map[string]any {
	return p.PartBase.Metadata
}

// Metadata returns the metadata of the FilePart.
func (p *FilePart) Metadata() map[string]any {
	return p.PartBase.Metadata
}

// Metadata returns the metadata of the DataPart.
func (p *DataPart) Metadata() map[string]any {
	return p.PartBase.Metadata
}

// UnmarshalJSON implements the json.Unmarshaler interface for Parts.
func (p *Parts) UnmarshalJSON(data []byte) error {
	var parts []json.RawMessage
	if err := json.Unmarshal(data, &parts); err != nil {
		return err
	}
	for _, raw := range parts {
		var kind struct {
			Kind PartKind `json:"kind"`
		}
		if err := json.Unmarshal(raw, &kind); err != nil {
			return err
		}
		var part Part
		switch kind.Kind {
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
			return fmt.Errorf("unknown part kind: %s", kind.Kind)
		}
		*p = append(*p, part)
	}
	return nil
}

// NewFileBytes creates a new FilePart with file content as bytes.
func NewFileBytes(data []byte) *FilePart {
	return &FilePart{
		Kind: PartKindFile,
		File: FileWithBytes{
			Bytes: base64.StdEncoding.EncodeToString(data),
		},
	}
}

// NewFileURI creates a new FilePart with a file URI.
func NewFileURI(uri string) *FilePart {
	return &FilePart{
		Kind: PartKindFile,
		File: FileWithURI{
			URI: uri,
		},
	}
}

// NewText creates a new TextPart with the given text.
func NewText(text string) *TextPart {
	return &TextPart{
		Kind: PartKindText,
		Text: text,
	}
}

// WithMetadata adds a key-value pair to the metadata of the FilePart and returns the updated FilePart.
func (p *FilePart) WithMetadata(key string, val any) *FilePart {
	if p.PartBase.Metadata == nil {
		p.PartBase.Metadata = make(map[string]any)
	}
	p.PartBase.Metadata[key] = val
	return p
}

// WithMetadata adds a key-value pair to the metadata of the DataPart and returns the updated DataPart.
func (p *DataPart) WithMetadata(key string, val any) *DataPart {
	if p.PartBase.Metadata == nil {
		p.PartBase.Metadata = make(map[string]any)
	}
	p.PartBase.Metadata[key] = val
	return p
}

// WithMetadata adds a key-value pair to the metadata of the TextPart and returns the updated TextPart.
func (p *TextPart) WithMetadata(key string, val any) *TextPart {
	if p.PartBase.Metadata == nil {
		p.PartBase.Metadata = make(map[string]any)
	}
	p.PartBase.Metadata[key] = val
	return p
}
