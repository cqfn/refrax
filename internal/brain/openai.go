package brain

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// OpenAI represents a client for interacting with the OpenAI API
type openAI struct {
	token string
	url   string
	model string
}

type openaiReq struct {
	Model    string      `json:"model"`
	Messages []openaiMsg `json:"messages"`
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

// NewOpenAI creates a new OpenAI instance
func NewOpenAI(apiKey string) Brain {
	return &openAI{
		token: apiKey,
		url:   "https://api.openai.com/v1/chat/completions",
		model: "gpt-3.5-turbo", // Default model
	}
}

// Ask sends a question to the OpenAI API
func (o *openAI) Ask(question string) (string, error) {
	return o.send("You are a helpful assistant.", question)
}

func (o *openAI) send(system, user string) (answer string, err error) {
	body := openaiReq{
		Model: o.model,
		Messages: []openaiMsg{
			{Role: "system", Content: system},
			{Role: "user", Content: strings.TrimSpace(user)},
		},
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
	req.Header.Set("Authorization", "Bearer "+o.token)
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
		return "", fmt.Errorf("API error: %s", body)
	}
	var response openaiResp
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}
	if len(response.Choices) == 0 {
		return "", errors.New("no choices in response")
	}
	return strings.TrimSpace(response.Choices[0].Message.Content), nil
}
