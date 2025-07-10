package protocol

import (
	"encoding/json"
	"fmt"
)

// FileWithBytes represents a file with its content encoded in base64.
type FileWithBytes struct {
	Bytes string `json:"bytes"` // Required: base64-encoded content
}

// FileWithURI represents a file referenced by its URI.
type FileWithURI struct {
	URI string `json:"uri"` // Required: URI to the file
}

// UnmarshalJSON unmarshals a JSON representation of FilePart, detecting its file type.
func (p *FilePart) UnmarshalJSON(data []byte) error {
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
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(aux.File, &probe); err != nil {
		return err
	}
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
