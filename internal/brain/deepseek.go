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

// deepSeek represents a client for interacting with the deepSeek API.
type deepSeek struct {
	token string
	url   string
	model string
}

type deepseekReq struct {
	Model    string        `json:"model"`
	Messages []deepseekMsg `json:"messages"`
	Stream   bool          `json:"stream"`
}

type deepseekResp struct {
	Choices []deepseekChoice `json:"choices"`
}

type deepseekChoice struct {
	Message deepseekMsg `json:"message"`
}

type deepseekMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// NewDeepSeek creates a new deepSeek instance with the provided API key.
func NewDeepSeek(apiKey string) Brain {
	return &deepSeek{
		token: apiKey,
		url:   "https://api.deepseek.com/chat/completions",
		model: "deepseek-chat",
	}
}

// Ask sends a question to the deepSeek API and retrieves an answer.
func (d *deepSeek) Ask(question string) (string, error) {
	log.Debug("deepSeek: asking question: %s", question)
	return d.send("You are a helpful assistant.", question)
}

func (d *deepSeek) send(system, user string) (answer string, err error) {
	content := trimmed(user)
	log.Debug("deepSeek: sending request with system promt: '%s' and userPrompt: '%s'", system, content)
	body := deepseekReq{
		Model: d.model,
		Messages: []deepseekMsg{
			{Role: "system", Content: system},
			{Role: "user", Content: content},
		},
		Stream: false,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %w", err)
	}
	req, err := http.NewRequest("POST", d.url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request to deepseek api: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = fmt.Errorf("error closing response body: %w", cerr)
		}
	}()
	if resp.StatusCode != 200 {
		content, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s", content)
	}
	var parsed deepseekResp
	if err = json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}
	if len(parsed.Choices) == 0 {
		return "", errors.New("no choices in response")
	}
	answer = strings.TrimSpace(parsed.Choices[0].Message.Content)
	return answer, err
}
