package protocol

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type PartKind string

const (
	PartKindText PartKind = "text"
	PartKindFile PartKind = "file"
	PartKindData PartKind = "data"
)

type Part interface {
	PartKind() PartKind
}

type TextPart struct {
	Kind PartKind `json:"kind"` // Must be "text"
	Text string   `json:"text"` // Required text content
}

type FilePart struct {
	Kind PartKind `json:"kind"` // Must be "file"
	File any      `json:"file"` // Can be FileWithBytes or FileWithUri
}

type DataPart struct {
	Kind PartKind       `json:"kind"` // Must be "data"
	Data map[string]any `json:"data"` // Required structured content
}

type Parts []Part

func (p *TextPart) PartKind() PartKind {
	return p.Kind
}

func (p *FilePart) PartKind() PartKind {
	return p.Kind
}

func (p *DataPart) PartKind() PartKind {
	return p.Kind
}

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

func NewFileBytes(data []byte) *FilePart {
	return &FilePart{
		Kind: PartKindFile,
		File: FileWithBytes{
			Bytes: base64.StdEncoding.EncodeToString(data),
		},
	}
}

func NewFileURI(uri string) *FilePart {
	return &FilePart{
		Kind: PartKindFile,
		File: FileWithURI{
			URI: uri,
		},
	}
}

func NewText(text string) *TextPart {
	return &TextPart{
		Kind: PartKindText,
		Text: text,
	}
}
