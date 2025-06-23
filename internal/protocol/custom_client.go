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

type CustomClient struct {
	url    string
	client *http.Client
}

func NewCustomClient(url string) Client {
	return &CustomClient{
		url: url,
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (c *CustomClient) SendMessage(params MessageSendParams) (*JSONRPCResponse, error) {
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

func (c *CustomClient) CancelTask() {
	panic("unimplemented")
}

func (c *CustomClient) GetTask() {
	panic("unimplemented")
}

func (c *CustomClient) StreamMessage() {
	panic("unimplemented")
}

func (client *CustomClient) doRequest(req any, resp *JSONRPCResponse) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request %v: %w", req, err)
	}
	httpReq, err := http.NewRequest("POST", client.url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create POST request for %s: %w", client.url, err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := client.client.Do(httpReq)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return fmt.Errorf("request to %s timed out after %s", client.url, client.client.Timeout)
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
	copy(&rawResp, resp)
	return nil
}

func copy(from *JSONRPCResponse, to *JSONRPCResponse) {
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
