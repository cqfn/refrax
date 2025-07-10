package protocol

import (
	"context"
	"fmt"

	"trpc.group/trpc-go/trpc-a2a-go/protocol"
	"trpc.group/trpc-go/trpc-a2a-go/server"
	"trpc.group/trpc-go/trpc-a2a-go/taskmanager"
)

type trpcServer struct {
	server server.A2AServer
	port   int
	card   *AgentCard
}

// NewTrpcServer is an interface for a TRPC server that handles A2A communication.
func NewTrpcServer(card *AgentCard, port int) (Server, error) {
	processor := &myTaskProcessor{}
	taskManager, err := taskmanager.NewMemoryTaskManager(processor)
	if err != nil {
		return nil, fmt.Errorf("failed to create task manager: %w", err)
	}
	serv, err := server.NewA2AServer(*agentCard(card), taskManager)
	if err != nil {
		return nil, fmt.Errorf("failed to create A2A server: %w", err)
	}
	return &trpcServer{card: card, server: *serv, port: port}, nil
}

// SetHandler sets the message handler for the TRPC server.
func (s *trpcServer) SetHandler(_ MessageHandler) {
	panic("unimplemented")
}

// Start starts the TRPC server and listens on the specified port, while signaling readiness.
func (s *trpcServer) Start(ready chan<- struct{}) error {
	port := fmt.Sprintf(":%d\n", s.port)
	if ready != nil {
		close(ready)
	}
	if err := s.server.Start(port); err != nil {
		return fmt.Errorf("failed to start A2A server: %w", err)
	}
	return nil
}

// Close stops the TRPC server gracefully.
func (s *trpcServer) Close() error {
	if err := s.server.Stop(context.Background()); err != nil {
		return fmt.Errorf("failed to stop A2A server: %w", err)
	}
	return nil
}

func agentCard(card *AgentCard) *server.AgentCard {
	return &server.AgentCard{
		Name:             card.Name,
		Description:      &card.Description,
		URL:              card.URL,
		Version:          card.Version,
		DocumentationURL: card.DocumentationURL,
		Provider: &server.AgentProvider{
			Organization: card.Provider.Organization,
			URL:          &card.Provider.URL,
		},
		Capabilities: server.AgentCapabilities{
			Streaming:              boolean(card.Capabilities.Streaming),
			PushNotifications:      boolean(card.Capabilities.PushNotifications),
			StateTransitionHistory: boolean(card.Capabilities.StateTransitionHistory),
		},
		DefaultInputModes:  card.DefaultInputModes,
		DefaultOutputModes: card.DefaultOutputModes,
		Skills: func() []server.AgentSkill {
			var skills []server.AgentSkill
			for i := range card.Skills {
				skill := &card.Skills[i]
				skills = append(skills, server.AgentSkill{
					Name:        skill.Name,
					Description: &skill.Description,
					Tags:        skill.Tags,
					Examples:    skill.Examples,
					InputModes:  skill.InputModes,
					OutputModes: skill.OutputModes,
				})
			}
			return skills
		}(),
	}
}

func boolean(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

type myTaskProcessor struct{}

func (p *myTaskProcessor) Process(
	_ context.Context,
	_ string,
	_ protocol.Message,
	_ taskmanager.TaskHandle,
) error {
	return nil
}
