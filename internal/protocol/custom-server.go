package protocol

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/cqfn/refrax/internal/log"
)

type CustomServer struct {
	mux     *http.ServeMux
	card    AgentCard
	handler MessageHandler
	port    int
	server  *http.Server
}

func NewCustomServer(card AgentCard, handler MessageHandler, port int) (Server, error) {
	mux := http.NewServeMux()
	server := &CustomServer{
		mux:     mux,
		card:    card,
		port:    port,
		handler: handler,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
	mux.HandleFunc("/.well-known/agent-card.json", server.handleAgentCard)
	mux.HandleFunc("/", server.handleJSONRPC)
	return server, nil
}

func LogRequest(message *Message) (*Message, error) {
	log.Debug("server received the following message: %s", message.MessageID)
	return message, nil
}

func (c *CustomServer) Start(ready chan<- struct{}) error {
	log.Info("starting custom a2a server on port %d...", c.port)
	address := fmt.Sprintf(":%d", c.port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", c.port, err)
	}
	close(ready)
	if err := http.Serve(l, c.mux); err != nil {
		return fmt.Errorf("failed to start server on port %d: %w", c.port, err)
	}
	return nil
}

func (c *CustomServer) Close() error {
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

func (s *CustomServer) handleJSONRPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Debug("JSON-RPC request received: %s", r.URL.Path)
	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "", ErrCodeInvalidRequest, "invalid JSON payload")
		return
	}
	switch req.Method {
	case "message/send":
		var params MessageSendParams
		paramsBytes, err := json.Marshal(req.Params)
		if err != nil {
			msg := fmt.Sprintf("failed to marshal params: %v", err)
			sendError(w, toString(req.ID), ErrCodeInvalidRequest, msg)
			return
		}
		if err := json.Unmarshal(paramsBytes, &params); err != nil {
			msg := fmt.Sprintf("failed to unmarshal params: %v", err)
			sendError(w, toString(req.ID), ErrCodeInvalidRequest, msg)
			return
		}
		if err = s.handleMessageSend(w, &req, toString(req.ID)); err != nil {
			msg := fmt.Sprintf("failed to handle message send: %v", err)
			sendError(w, toString(req.ID), ErrCodeInternalError, msg)
			return
		}
	case "message/stream":
		panic("message/stream is not implemented in CustomServer")
	case "tasks/get":
		panic("tasks/get is not implemented in CustomServer")
	case "tasks/cancel":
		panic("tasks/get is not implemented in CustomServer")
	default:
		sendError(w, req.ID.(string), ErrCodeMethodNotFound, "method not found")
	}
}

func toString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%v", val)
	case bool:
		return fmt.Sprintf("%v", val)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", val)
	}
}

func (s *CustomServer) handleMessageSend(w http.ResponseWriter, req *JSONRPCRequest, id string) error {
	var params MessageSendParams
	paramsBytes, err := json.Marshal(req.Params)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		return err
	}
	msg := params.Message
	udpadted, err := s.handler(&msg)
	if err != nil {
		return err
	}
	log.Debug("message handler returned: %s", udpadted.MessageID)
	return s.sendResponse(w, id, udpadted)
}

func (s *CustomServer) sendResponse(w http.ResponseWriter, id string, result any) error {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}

	w.Header().Set("Content-Type", "application/json")
	log.Debug("sending response: %v", response)
	return json.NewEncoder(w).Encode(response)
}

func sendError(w http.ResponseWriter, id string, code int, message string) {
	response := errorResposne(id, code, message)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		panic(fmt.Errorf("failed to encode error response: %w", err))
	}
}

func errorResposne(id string, code int, message string) JSONRPCResponse {
	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
		},
	}
}
