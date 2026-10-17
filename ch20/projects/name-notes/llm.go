package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Client talks to a model server running on your own machine that speaks the OpenAI
// chat API, which Ollama, llama.cpp's server and LM Studio all do. Nothing leaves the
// machine.
type Client struct {
	BaseURL string // e.g. http://localhost:11434/v1
	Model   string // e.g. qwen2.5:1.5b
	HTTP    *http.Client
}

// FromEnv reads LLM_URL and LLM_MODEL, with Ollama's defaults.
func FromEnv() Client {
	c := Client{BaseURL: "http://localhost:11434/v1", Model: "qwen2.5:1.5b", HTTP: &http.Client{Timeout: 2 * time.Minute}}
	if v := os.Getenv("LLM_URL"); v != "" {
		c.BaseURL = v
	}
	if v := os.Getenv("LLM_MODEL"); v != "" {
		c.Model = v
	}
	return c
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type chatRequest struct {
	Model          string          `json:"model"`
	Temperature    float64         `json:"temperature"`
	Messages       []message       `json:"messages"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Chat sends one system prompt and one user prompt and returns the model's text. With
// jsonMode the server is asked for a JSON object; the caller still validates it.
func (c Client) Chat(ctx context.Context, system, user string, jsonMode bool) (string, error) {
	req := chatRequest{Model: c.Model, Temperature: 0, Messages: []message{{"system", system}, {"user", user}}}
	if jsonMode {
		req.ResponseFormat = &responseFormat{Type: "json_object"}
	}
	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimSuffix(c.BaseURL, "/")+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("model server at %s: %w (is it running? try: ollama serve, then ollama pull %s)", c.BaseURL, err, c.Model)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}
	var out chatResponse
	if err := json.Unmarshal(raw, &out); err != nil || resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("model server: %s: %s", resp.Status, strings.TrimSpace(string(raw)))
	}
	if out.Error != nil {
		return "", errors.New("model server: " + out.Error.Message)
	}
	if len(out.Choices) == 0 {
		return "", errors.New("model server returned no choices")
	}
	return out.Choices[0].Message.Content, nil
}

// ExtractJSON returns the first {...} object in s, for models that wrap their JSON in
// prose or code fences despite being asked not to.
func ExtractJSON(s string) (string, bool) {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start < 0 || end <= start {
		return "", false
	}
	return s[start : end+1], true
}
