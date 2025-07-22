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
	mux        *http.ServeMux
	card       AgentCard
	msgHandler MsgHandler
	port       int
	server     *http.Server
	handler    Handler
}

type (
	// Handler handles all incoming requests on JSONRPC level
	Handler func(next Handler, r *JSONRPCRequest) (*JSONRPCResponse, error)
	// MsgHandler handles messages received from the A2A server on Message level
	MsgHandler func(message *Message) (*Message, error)
)

// NewCustomServer creates a new instance of a custom server that handles A2A requests
func NewCustomServer(card *AgentCard, port int) Server {
	mux := http.NewServeMux()
	server := &customServer{
		mux:        mux,
		card:       *card,
		port:       port,
		msgHandler: record,
		server: &http.Server{
			Addr:              fmt.Sprintf(":%d", port),
			Handler:           mux,
			ReadHeaderTimeout: 20 * time.Second,
		},
	}
	mux.HandleFunc("/.well-known/agent-card.json", server.handleAgentCard)
	mux.HandleFunc("/", server.handleRequest)
	return server
}

// MsgHandler sets the message handler for the custom server.
func (serv *customServer) MsgHandler(handler MsgHandler) {
	serv.msgHandler = handler
}

// Handler sets the handler function for processing requests.
func (serv *customServer) Handler(handler Handler) {
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
	if err = serv.server.Serve(l); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server on port %d: %w", serv.port, err)
	}
	return err
}

// Close stops the custom server gracefully, allowing for a timeout.
func (serv *customServer) Close() error {
	log.Debug("stopping custom a2a server on port %d...", serv.port)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
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

func (serv *customServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	err := serv.handleJSONRPC(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to handle request: %v", err), http.StatusInternalServerError)
		return
	}
}

func (serv *customServer) handleJSONRPC(w http.ResponseWriter, r *http.Request) error {
	log.Debug("JSON-RPC request received: %s", r.URL.Path)
	if r.Method != http.MethodPost {
		return fmt.Errorf("method not allowed: %s", r.Method)
	}
	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := failure("", ErrCodeInvalidRequest, "invalid JSON payload")
		return send(w, &resp)
	}
	var resp *JSONRPCResponse
	var err error
	if serv.handler != nil {
		start := serv.handler
		resp, err = start(basic(serv.msgHandler), &req)
	} else {
		start := basic(serv.msgHandler)
		resp, err = start(nil, &req)
	}
	if err != nil {
		resp := failure(str(req.ID), ErrCodeInternalError, fmt.Sprintf("failed to handle request: %v", err))
		return send(w, &resp)
	}
	return send(w, resp)
}

func basic(mh MsgHandler) Handler {
	return func(_ Handler, r *JSONRPCRequest) (*JSONRPCResponse, error) {
		id := str(r.ID)
		switch r.Method {
		case "message/send":
			pbytes, err := json.Marshal(r.Params)
			var params MessageSendParams
			if err != nil {
				msg := fmt.Sprintf("failed to marshal params '%v': %v", r.Params, err)
				resp := failure(id, ErrCodeInvalidRequest, msg)
				return &resp, nil
			}
			if err = json.Unmarshal(pbytes, &params); err != nil {
				msg := fmt.Sprintf("failed to unmarshal params '%v': : %v", r.Params, err)
				resp := failure(id, ErrCodeInvalidRequest, msg)
				return &resp, nil
			}
			log.Debug("handling JSON-RPC request: %s, params: %v", r.Method, params)
			msg := params.Message
			msg, err = mh(msg)
			if err != nil {
				resp := failure(id, ErrCodeInternalError, fmt.Sprintf("failed to handle message send: %v", err))
				return &resp, nil
			}
			resp := success(id, msg)
			return &resp, nil
		case "message/stream":
			panic("message/stream is not implemented yet")
		case "tasks/get":
			panic("tasks/get is not implemented yet")
		case "tasks/cancel":
			panic("tasks/get is not implemented yet")
		default:
			resp := failure(id, ErrCodeMethodNotFound, "method not found")
			return &resp, nil
		}
	}
}

// str converts various types to a string representation.
func str(v any) string {
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

func success(id string, result any) JSONRPCResponse {
	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

func failure(id string, code int, message string) JSONRPCResponse {
	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
		},
	}
}

func send(w http.ResponseWriter, r *JSONRPCResponse) error {
	w.Header().Set("Content-Type", "application/json")
	log.Debug("sending response: %v", r)
	return json.NewEncoder(w).Encode(r)
}

func record(message *Message) (*Message, error) {
	log.Debug("server received the following message: %s", message.MessageID)
	return message, nil
}
