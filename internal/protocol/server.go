package protocol

type Server interface {
	Start(ready chan<- struct{}) error
	SetHandler(handler MessageHandler)
	Close() error
}

type MessageHandler func(message *Message) (*Message, error)
