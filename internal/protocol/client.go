package protocol

// Client represents the interface for a protocol client.
type Client interface {
	SendMessage(question *MessageSendParams) (*JSONRPCResponse, error)
	StreamMessage()
	GetTask()
	CancelTask()
}
