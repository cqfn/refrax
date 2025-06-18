package protocol

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cqfn/refrax/internal/log"
)

type CustomServer struct {
	mux    *http.ServeMux
	card   AgentCard
	port   int
	server *http.Server
}

func NewCustomServer(card AgentCard, port int) (Server, error) {
	mux := http.NewServeMux()
	server := &CustomServer{
		mux:  mux,
		card: card,
		port: port,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
	mux.HandleFunc("/.well-known/agent-card.json", server.handleAgentCard)
	return server, nil
}

func (c *CustomServer) Start() error {
	log.Info("starting custom a2a server on port %d...", c.port)
	return c.server.ListenAndServe()
}

func (c *CustomServer) Stop() error {
	log.Info("stopping custom a2a server on port %d...", c.port)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.server.Shutdown(ctx)
}

func (c *CustomServer) handleAgentCard(w http.ResponseWriter, r *http.Request) {
	log.Debug("request for agent card received: %s", r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(c.card); err != nil {
		http.Error(w, "failed to encode agent card", http.StatusInternalServerError)
	}
}
