package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cqfn/refrax/internal/log"
	"github.com/cqfn/refrax/internal/models"
)

type RefraxServer struct {
	agentCard models.AgentCard
	port      int
	mux       *http.ServeMux
}

func NewA2AServer(port int) *RefraxServer {
	mux := http.NewServeMux()
	server := &RefraxServer{
		port: port,
		mux:  mux,
		agentCard: models.AgentCard{
			Name:        "Refrax Agent",
			Description: stringPtr("A test agent for unit tests"),
			URL:         "http://localhost:8080",
			Version:     "0.0.1",
			Capabilities: models.AgentCapabilities{
				Streaming:              boolPtr(false),
				PushNotifications:      boolPtr(false),
				StateTransitionHistory: boolPtr(false),
			},
			Skills: []models.AgentSkill{
				{
					ID:          "refactor-java-skill",
					Name:        "Refactor Java Skill",
					Description: stringPtr("Capability to refactor Java code"),
				},
			},
		},
	}
	mux.HandleFunc("/.well-known/agent-card", server.handleAgentCard)
	return server
}

func StartServer(port int) error {
	log.Info("Starting HTTP server...")
	server := NewA2AServer(port)
	return server.Start()
}

func (s *RefraxServer) Start() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)
}

func (s *RefraxServer) handleAgentCard(w http.ResponseWriter, r *http.Request) {
	log.Debug("received request for agent card")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s.agentCard); err != nil {
		http.Error(w, "Failed to encode agent card", http.StatusInternalServerError)
	}
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
