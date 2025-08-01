package protocol

import (
	"encoding/json"
	"fmt"

	"github.com/cqfn/refrax/internal/log"
)

// JSON-RPC 2.0 standard error codes (per A2A specification)
const (
	// -32700: Invalid JSON payload
	ErrCodeParseError = -32700

	// -32600: Invalid JSON-RPC Request
	ErrCodeInvalidRequest = -32600

	// -32601: Method not found
	ErrCodeMethodNotFound = -32601

	// -32602: Invalid method parameters
	ErrCodeInvalidParams = -32602

	// -32603: Internal server error
	ErrCodeInternalError = -32603

	// -32000 to -32099: Reserved for server-defined errors (A2A-specific)
	ErrCodeServerErrorStart = -32099
	ErrCodeServerErrorEnd   = -32000
)

// JSONRPCRequest represents a request in the JSON-RPC protocol.
type JSONRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`          // MUST be "2.0"
	Method  string `json:"method"`           // e.g., "message/send"
	Params  any    `json:"params,omitempty"` // Can be any structured value (typically an object)
	ID      any    `json:"id,omitempty"`     // string, number (int), or nil
}

// JSONRPCResponse represents a response to a JSON-RPC request.
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`          // MUST be "2.0"
	ID      any           `json:"id"`               // Same type as in request (string, int, or null)
	Result  any           `json:"result,omitempty"` // Present only on success
	Error   *JSONRPCError `json:"error,omitempty"`  // Present only on failure
}

// JSONRPCError represents an error in a JSON-RPC response.
type JSONRPCError struct {
	Code    int    `json:"code"`           // Error code indicating the type of error
	Message string `json:"message"`        // Short description of the error
	Data    any    `json:"data,omitempty"` // Optional additional information (any type)
}

func (r *JSONRPCResponse) String() string {
	return fmt.Sprintf("JSONRPCResponse{jsonrpc: %s, id: %v, error: %v}", r.JSONRPC, r.ID, r.Error)
}

// UnmarshalJSON implements custom unmarshalling for JSON-RPC responses.
func (r *JSONRPCResponse) UnmarshalJSON(data []byte) error {
	type alias JSONRPCResponse
	type temp struct {
		Result json.RawMessage `json:"result,omitempty"`
		*alias
	}
	aux := temp{
		alias: (*alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if len(aux.Result) == 0 || string(aux.Result) == "null" {
		r.Result = nil
		return nil
	}
	var kind struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(aux.Result, &kind); err != nil {
		return fmt.Errorf("failed to detect kind: %w", err)
	}
	log.Debug("Detected kind: %s", kind.Kind)
	switch kind.Kind {
	case "message":
		var msg Message
		if err := json.Unmarshal(aux.Result, &msg); err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}
		r.Result = &msg
	case "task":
		var task Task
		if err := json.Unmarshal(aux.Result, &task); err != nil {
			return fmt.Errorf("failed to unmarshal task: %w", err)
		}
		r.Result = &task
	default:
		var generic map[string]any
		if err := json.Unmarshal(aux.Result, &generic); err != nil {
			return fmt.Errorf("failed to unmarshal unknown result type: %w", err)
		}
		r.Result = generic
	}
	return nil
}
