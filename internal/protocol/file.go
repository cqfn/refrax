package protocol

import (
	"encoding/json"
	"fmt"
)

type FileWithBytes struct {
	Bytes string `json:"bytes"` // Required: base64-encoded content
}

type FileWithURI struct {
	URI string `json:"uri"` // Required: URI to the file
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
