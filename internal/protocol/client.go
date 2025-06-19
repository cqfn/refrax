package protocol

type Client interface {
	SendMessage(question MessageSendParams) (*JSONRPCResponse, error)
	StreamMessage()
	GetTask()
	CancelTask()
}
