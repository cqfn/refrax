package brain

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cqfn/refrax/internal/log"
)

// OpenAI represents a client for interacting with the OpenAI API
type openAI struct {
	token  string
	url    string
	model  string
	system string
}

type openaiReq struct {
	Model       string      `json:"model"`
	Messages    []openaiMsg `json:"messages"`
	Stream      bool        `json:"stream"`
	Temperature *float64    `json:"temperature,omitempty"`
}

type openaiResp struct {
	Choices []openaiChoice `json:"choices"`
}

type openaiChoice struct {
	Message openaiMsg `json:"message"`
}

type openaiMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// NewOpenAIDefault creates a new OpenAI instance with default settings
func NewOpenAIDefault(token, system string) (Brain, error) {
	return NewOpenAI(token, "https://api.openai.com/v1/chat/completions", "gpt-3.5-turbo", system)
}

// NewOpenAI creates a new OpenAI instance with the provided settings
func NewOpenAI(token, url, model, system string) (Brain, error) {
	err := verifyUrl(url)
	if err != nil {
		return nil, err
	}
	return &openAI{
		token:  token,
		url:    url,
		model:  model,
		system: system,
	}, nil
}

// Ask sends a question to the OpenAI API
func (o *openAI) Ask(question string) (string, error) {
	return o.send(o.system, question)
}

func (o *openAI) send(system, user string) (answer string, err error) {
	log.Debug("sending request to '%s', model '%s', and prompt: '%s'", o.url, o.model, user)
	temp := float64(0.0)
	body := openaiReq{
		Model: o.model,
		Messages: []openaiMsg{
			{Role: "system", Content: system},
			{Role: "user", Content: strings.TrimSpace(user)},
		},
		Stream:      false,
		Temperature: &temp,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}
	req, err := http.NewRequest("POST", o.url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = fmt.Errorf("error closing response body: %w", cerr)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (code: %d): %s", resp.StatusCode, body)
	}
	var response openaiResp
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}
	if len(response.Choices) == 0 {
		return "", errors.New("no choices in response")
	}
	return response.Choices[0].Message.Content, nil
}

func verifyUrl(url string) error {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("invalid URL: must start with http:// or https://")
	}
	if !strings.HasSuffix(url, "/v1/chat/completions") {
		return fmt.Errorf("invalid URL: must end with /v1/chat/completions")
	}
	return nil
}
