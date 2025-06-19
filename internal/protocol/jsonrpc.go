package protocol

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

type JSONRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`          // MUST be "2.0"
	Method  string `json:"method"`           // e.g., "message/send"
	Params  any    `json:"params,omitempty"` // Can be any structured value (typically an object)
	ID      any    `json:"id,omitempty"`     // string, number (int), or nil
}

type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`          // MUST be "2.0"
	ID      any           `json:"id"`               // Same type as in request (string, int, or null)
	Result  any           `json:"result,omitempty"` // Present only on success
	Error   *JSONRPCError `json:"error,omitempty"`  // Present only on failure
}

type JSONRPCError struct {
	Code    int    `json:"code"`           // Error code indicating the type of error
	Message string `json:"message"`        // Short description of the error
	Data    any    `json:"data,omitempty"` // Optional additional information (any type)
}
