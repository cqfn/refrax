package protocol

// Artifact represents an object with various attributes and metadata.
type Artifact struct {
	ArtifactID  string           `json:"artifactId"`            // Required: unique identifier
	Name        *string          `json:"name,omitempty"`        // Optional: human-readable name
	Description *string          `json:"description,omitempty"` // Optional: human-readable description
	Parts       []map[string]any `json:"parts"`                 // Required: parts (can be refined later)
	Metadata    map[string]any   `json:"metadata,omitempty"`    // Optional: extension metadata
	Extensions  []string         `json:"extensions,omitempty"`  // Optional: contributed extension URIs
}
