package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Provider interface for LLM providers
type Provider interface {
	Chat(messages []Message, opts Options) (*Response, error)
	Stream(messages []Message, opts Options) (<-chan StreamChunk, error)
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Options for LLM requests
type Options struct {
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	System      string  `json:"system,omitempty"`
}

// Response from LLM
type Response struct {
	Content    string `json:"content"`
	Model      string `json:"model"`
	StopReason string `json:"stop_reason"`
	Usage      Usage  `json:"usage"`
}

// Usage statistics
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// StreamChunk for streaming responses
type StreamChunk struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	Error   error  `json:"error,omitempty"`
	Done    bool   `json:"done"`
}

// ClaudeProvider implements Provider for Anthropic Claude
type ClaudeProvider struct {
	APIKey  string
	BaseURL string
}

// NewClaude creates a new Claude provider
func NewClaude(apiKey string) *ClaudeProvider {
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	return &ClaudeProvider{
		APIKey:  apiKey,
		BaseURL: "https://api.anthropic.com/v1",
	}
}

// ClaudeRequest is the API request format
type ClaudeRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
	System    string    `json:"system,omitempty"`
}

// ClaudeResponse is the API response format
type ClaudeResponse struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Role         string `json:"role"`
	Content      []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model        string `json:"model"`
	StopReason   string `json:"stop_reason"`
	Usage        Usage  `json:"usage"`
}

// Chat sends a chat request to Claude
func (c *ClaudeProvider) Chat(messages []Message, opts Options) (*Response, error) {
	if opts.Model == "" {
		opts.Model = "claude-sonnet-4-20250514"
	}
	if opts.MaxTokens == 0 {
		opts.MaxTokens = 4096
	}

	req := ClaudeRequest{
		Model:     opts.Model,
		MaxTokens: opts.MaxTokens,
		Messages:  messages,
		System:    opts.System,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL+"/messages", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var claudeResp ClaudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&claudeResp); err != nil {
		return nil, err
	}

	content := ""
	for _, c := range claudeResp.Content {
		if c.Type == "text" {
			content += c.Text
		}
	}

	return &Response{
		Content:    content,
		Model:      claudeResp.Model,
		StopReason: claudeResp.StopReason,
		Usage:      claudeResp.Usage,
	}, nil
}

// Stream sends a streaming chat request to Claude
func (c *ClaudeProvider) Stream(messages []Message, opts Options) (<-chan StreamChunk, error) {
	ch := make(chan StreamChunk)

	go func() {
		defer close(ch)

		// For now, fall back to non-streaming
		resp, err := c.Chat(messages, opts)
		if err != nil {
			ch <- StreamChunk{Error: err, Done: true}
			return
		}

		ch <- StreamChunk{Type: "content", Content: resp.Content}
		ch <- StreamChunk{Type: "done", Done: true}
	}()

	return ch, nil
}

// OpenRouterProvider implements Provider for OpenRouter
type OpenRouterProvider struct {
	APIKey  string
	BaseURL string
}

// NewOpenRouter creates a new OpenRouter provider
func NewOpenRouter(apiKey string) *OpenRouterProvider {
	if apiKey == "" {
		apiKey = os.Getenv("OPENROUTER_API_KEY")
	}
	return &OpenRouterProvider{
		APIKey:  apiKey,
		BaseURL: "https://openrouter.ai/api/v1",
	}
}

// Chat sends a chat request to OpenRouter
func (o *OpenRouterProvider) Chat(messages []Message, opts Options) (*Response, error) {
	if opts.Model == "" {
		opts.Model = "anthropic/claude-sonnet-4"
	}
	if opts.MaxTokens == 0 {
		opts.MaxTokens = 4096
	}

	// OpenRouter uses OpenAI-compatible format
	reqBody := map[string]interface{}{
		"model":      opts.Model,
		"max_tokens": opts.MaxTokens,
		"messages":   messages,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", o.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+o.APIKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var openAIResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Model string `json:"model"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return nil, err
	}

	content := ""
	stopReason := ""
	if len(openAIResp.Choices) > 0 {
		content = openAIResp.Choices[0].Message.Content
		stopReason = openAIResp.Choices[0].FinishReason
	}

	return &Response{
		Content:    content,
		Model:      openAIResp.Model,
		StopReason: stopReason,
		Usage: Usage{
			InputTokens:  openAIResp.Usage.PromptTokens,
			OutputTokens: openAIResp.Usage.CompletionTokens,
		},
	}, nil
}

// Stream sends a streaming chat request to OpenRouter
func (o *OpenRouterProvider) Stream(messages []Message, opts Options) (<-chan StreamChunk, error) {
	ch := make(chan StreamChunk)

	go func() {
		defer close(ch)
		resp, err := o.Chat(messages, opts)
		if err != nil {
			ch <- StreamChunk{Error: err, Done: true}
			return
		}
		ch <- StreamChunk{Type: "content", Content: resp.Content}
		ch <- StreamChunk{Type: "done", Done: true}
	}()

	return ch, nil
}
