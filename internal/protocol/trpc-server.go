package protocol

import (
	"context"
	"fmt"

	"trpc.group/trpc-go/trpc-a2a-go/protocol"
	"trpc.group/trpc-go/trpc-a2a-go/server"
	"trpc.group/trpc-go/trpc-a2a-go/taskmanager"
)

type TrpcServer struct {
	server server.A2AServer
	port   int
	card   AgentCard
}

func NewTrpcServer(card AgentCard, port int) (Server, error) {
	processor := &myTaskProcessor{}
	taskManager, err := taskmanager.NewMemoryTaskManager(processor)
	if err != nil {
		return nil, fmt.Errorf("failed to create task manager: %w", err)
	}
	serv, err := server.NewA2AServer(agentCard(card), taskManager)
	if err != nil {
		return nil, fmt.Errorf("failed to create A2A server: %w", err)
	}
	return &TrpcServer{card: card, server: *serv, port: port}, nil
}

func (s *TrpcServer) Start() error {
	port := fmt.Sprintf(":%d\n", s.port)
	if err := s.server.Start(port); err != nil {
		return fmt.Errorf("failed to start A2A server: %w", err)
	}
	return nil
}

func agentCard(card AgentCard) server.AgentCard {
	return server.AgentCard{
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
			for _, skill := range card.Skills {
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

type myTaskProcessor struct {
}

func (p *myTaskProcessor) Process(
	ctx context.Context,
	taskID string,
	message protocol.Message,
	handle taskmanager.TaskHandle,
) error {
	return nil
}
