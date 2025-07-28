package protocol

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

// a2aClient represents a client for interacting with a custom API.
type a2aClient struct {
	url    string
	client *http.Client
}

// NewClient creates a new instance of CustomClient with a specified URL.
func NewClient(url string) Client {
	return &a2aClient{
		url: url,
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

// SendMessage sends a message using the custom API and returns the JSON-RPC response.
func (c *a2aClient) SendMessage(params MessageSendParams) (*JSONRPCResponse, error) {
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "1", // Static ID for simplicity, can be changed to a unique ID generator
		Method:  "message/send",
		Params:  params,
	}
	var resp JSONRPCResponse
	if err := c.doRequest(req, &resp); err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("message/send error: '%s' (code: %d)", resp.Error.Message, resp.Error.Code)
	}

	return &resp, nil
}

// CancelTask is a placeholder for canceling a task.
func (c *a2aClient) CancelTask() {
	panic("unimplemented")
}

// GetTask is a placeholder for retrieving a task.
func (c *a2aClient) GetTask() {
	panic("unimplemented")
}

// StreamMessage is a placeholder for streaming a message.
func (c *a2aClient) StreamMessage() {
	panic("unimplemented")
}

// doRequest sends a JSON-RPC request to the server and decodes the response.
func (c *a2aClient) doRequest(req any, resp *JSONRPCResponse) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request %v: %w", req, err)
	}
	httpReq, err := http.NewRequest("POST", c.url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create POST request for %s: %w", c.url, err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return fmt.Errorf("request to %s timed out after %s", c.url, c.client.Timeout)
		}
		return fmt.Errorf("failed to send request '%v': %w", req, err)
	}
	defer func() {
		if err := httpResp.Body.Close(); err != nil {
			panic(fmt.Errorf("failed to close response body: %w", err))
		}
	}()
	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}
	var rawResp JSONRPCResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&rawResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	copyResp(&rawResp, resp)
	return nil
}

// copyResp copies the content of one JSONRPCResponse to another.
func copyResp(from, to *JSONRPCResponse) {
	to.JSONRPC = from.JSONRPC
	to.ID = from.ID
	to.Result = from.Result
	if from.Error != nil {
		to.Error = &JSONRPCError{
			Code:    from.Error.Code,
			Message: from.Error.Message,
			Data:    from.Error.Data,
		}
	} else {
		to.Error = nil
	}
}
