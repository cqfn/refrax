package brain

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/cqfn/refrax/internal/log"
	"github.com/ollama/ollama/api"
)

// ollamaBrain is the implementation of the Brain interface for interacting with the Ollama API.
type ollamaBrain struct {
	url       string       // Base URL of the Ollama API.
	httpCient *http.Client // HTTP client used for making API requests.
	token     string       // Authentication token for the Ollama API.
	model     string       // The AI model to be used in the Ollama API.
	system    string       // The system prompt or context for the AI model.
}

// NewOllama creates a new instance of the Brain implementation using the Ollama API.
//
// Parameters:
//   - address (string): The base URL of the Ollama API.
//   - model (string): The AI model to use for chat interactions.
//   - token (string): The authentication token for accessing the API.
//   - system (string): The initial system prompt or context for the AI model.
//
// Returns:
//   - Brain: An instance of the Brain interface configured to use the Ollama API.
func NewOllama(address, model, token, system string) Brain {
	return &ollamaBrain{
		url:       address,
		httpCient: &http.Client{},
		token:     token,
		model:     model,
		system:    system,
	}
}

// Ask sends a question to the configured AI model and retrieves the response.
//
// This function constructs a chat request using the configured model, system prompt,
// and the user-provided question. It sends the request to the Ollama API and waits
// for a response or an error.
//
// Parameters:
//   - question (string): The user's question to be sent to the AI model.
//
// Returns:
//   - string: The AI model's response to the question.
//   - error: An error that occurred during the API interaction, or nil if successful.
func (o *ollamaBrain) Ask(question string) (string, error) {
	address, err := url.Parse(o.url)
	if err != nil {
		return "", err
	}
	client := api.NewClient(address, o.httpCient)
	ctx := context.Background()
	req := api.ChatRequest{
		Model: o.model,
		Messages: []api.Message{
			{
				Role:    "system",
				Content: o.system,
			},
			{
				Role:    "user",
				Content: question,
			},
		},
	}
	answ := make(chan string)
	errc := make(chan error)
	res := func(r api.ChatResponse) error {
		log.Debug("Ollama response: %+v", r)
		select {
		case answ <- r.Message.Content:
			return nil
		case <-time.After(time.Minute * 5):
			return fmt.Errorf("timeout: can't write to answer channel")
		}
	}
	go func() {
		err := client.Chat(ctx, &req, res)
		if err != nil {
			errc <- err
		}
	}()
	select {
	case a := <-answ:
		return a, nil
	case e := <-errc:
		return "", fmt.Errorf("error from Ollama API: %w", e)
	}
}
