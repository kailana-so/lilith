package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lilith/internal/domain/dto"
	"net/http"
	"os"
	"strings"
)

type anthropicClient struct {
	apiKey string
	model  string
}

func NewAnthropicClient() *anthropicClient {
	return &anthropicClient{
		apiKey: os.Getenv("ANTHROPIC_API_KEY"),
		model:  os.Getenv("ANTHROPIC_MODEL"),
	}
}

func (c *anthropicClient) Complete(ctx context.Context, system, user string) (string, error) {
	return c.CompleteMessages(ctx, []dto.Message{
		{Role: "system", Content: system},
		{Role: "user", Content: user},
	})
}

func (c *anthropicClient) CompleteMessages(ctx context.Context, msgs []dto.Message) (string, error) {
	if c.apiKey == "" {
		return "", errors.New("ANTHROPIC_API_KEY not set")
	}

	// Split out system prompt(s)
	var systemPrompt string
	var amsgs []map[string]any
	for _, m := range msgs {
		if m.Role == "system" {
			systemPrompt += m.Content + "\n"
			continue
		}
		amsgs = append(amsgs, map[string]any{
			"role":    m.Role,
			"content": m.Content,
		})
	}

	body := map[string]any{
		"model":      c.model,
		"system":     systemPrompt,
		"messages":   amsgs,
		"max_tokens": 1000,
	}
	bs, _ := json.Marshal(body)

	req, _ := http.NewRequestWithContext(ctx, "POST",
		"https://api.anthropic.com/v1/messages", bytes.NewReader(bs))
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		raw, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("anthropic: %s\n%s", res.Status, raw)
	}

	var jr struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.NewDecoder(res.Body).Decode(&jr); err != nil {
		return "", err
	}

	var b strings.Builder
	for _, c := range jr.Content {
		if c.Type == "text" {
			b.WriteString(c.Text)
		}
	}
	return b.String(), nil
}
