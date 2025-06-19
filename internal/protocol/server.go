package protocol

type Server interface {
	Start(ready chan<- struct{}) error
	Close() error
}

type MessageHandler func(message *Message) (*Message, error)
