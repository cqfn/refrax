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

type customServer struct {
	mux     *http.ServeMux
	card    AgentCard
	handler MessageHandler
	port    int
	server  *http.Server
}

// NewCustomServer creates a new instance of a custom server that handles A2A requests
func NewCustomServer(card *AgentCard, port int) Server {
	mux := http.NewServeMux()
	server := &customServer{
		mux:     mux,
		card:    *card,
		port:    port,
		handler: logRequest,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
	mux.HandleFunc("/.well-known/agent-card.json", server.handleAgentCard)
	mux.HandleFunc("/", server.handleJSONRPC)
	return server
}

// SetHandler sets the message handler for the custom server.
func (serv *customServer) SetHandler(handler MessageHandler) {
	serv.handler = handler
}

// Start starts the custom server and listens on the specified port, while signaling readiness.
func (serv *customServer) Start(ready chan<- struct{}) error {
	log.Debug("starting custom a2a server on port %d...", serv.port)
	address := fmt.Sprintf(":%d", serv.port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", serv.port, err)
	}
	close(ready)
	if err := http.Serve(l, serv.mux); err != nil {
		return fmt.Errorf("failed to start server on port %d: %w", serv.port, err)
	}
	return nil
}

// Close stops the custom server gracefully, allowing for a timeout.
func (serv *customServer) Close() error {
	log.Debug("stopping custom a2a server on port %d...", serv.port)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return serv.server.Shutdown(ctx)
}

func (serv *customServer) handleAgentCard(w http.ResponseWriter, r *http.Request) {
	log.Debug("request for agent card received: %s", r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(serv.card); err != nil {
		http.Error(w, "failed to encode agent card", http.StatusInternalServerError)
	}
}

func (serv *customServer) handleJSONRPC(w http.ResponseWriter, r *http.Request) {
	log.Debug("JSON-RPC request received: %s", r.URL.Path)
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "", ErrCodeInvalidRequest, "invalid JSON payload")
		return
	}
	pbytes, err := json.Marshal(req.Params)
	var params MessageSendParams
	if err != nil {
		msg := fmt.Sprintf("failed to marshal params '%v': %v", req.Params, err)
		sendError(w, toString(req.ID), ErrCodeInvalidRequest, msg)
		return
	}
	if err := json.Unmarshal(pbytes, &params); err != nil {
		msg := fmt.Sprintf("failed to unmarshal params '%v': : %v", req.Params, err)
		sendError(w, toString(req.ID), ErrCodeInvalidRequest, msg)
		return
	}
	log.Debug("handling JSON-RPC request: %s, params: %v", req.Method, params)
	switch req.Method {
	case "message/send":
		if err := serv.handleMessageSend(w, &params, toString(req.ID)); err != nil {
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

func (serv *customServer) handleMessageSend(w http.ResponseWriter, params *MessageSendParams, id string) error {
	msg := params.Message
	udpadted, err := serv.handler(msg)
	if err != nil {
		return err
	}
	log.Debug("message handler returned: %s", udpadted.MessageID)
	return serv.sendResponse(w, id, udpadted)
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

func (serv *customServer) sendResponse(w http.ResponseWriter, id string, result any) error {
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

func logRequest(message *Message) (*Message, error) {
	log.Debug("server received the following message: %s", message.MessageID)
	return message, nil
}
