package facilitator

import (
	"fmt"

	"github.com/cqfn/refrax/internal/protocol"
)

type RefraxServer struct {
	port int
}

func FacilitatorServer(port int) *RefraxServer {
	return &RefraxServer{
		port: port,
	}
}

func StartServer(port int) error {
	server, error := protocol.NewCustomServer(agentCard(), port)
	if error != nil {
		return fmt.Errorf("failed to create A2A server: %w", error)
	}
	return server.Start()
}

func agentCard() protocol.AgentCard {
	return protocol.Card().
		Name("Refrax Agent").
		Description("A test agent for unit tests").
		URL("http://localhost:8080").
		Version("0.0.1").
		Skill("refactor-java", "Refactor Java Projects", "Refrax can refactor java projects").
		Build()
}
