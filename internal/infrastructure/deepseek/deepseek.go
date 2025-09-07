package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"lilith/internal/domain/dto"
)

type deepSeekClient struct {
	apiKey string
	model  string
}

func NewDeepSeekClient() *deepSeekClient {
	return &deepSeekClient{
		apiKey: os.Getenv("DEEPSEEK_API_KEY"),
		model:  os.Getenv("DEEPSEEK_MODEL"),
	}
}

func (dsClient *deepSeekClient) Complete(ctx context.Context, system, user string) (string, error) {
	return dsClient.CompleteMessages(ctx, []dto.Message{
		{Role: "system", Content: system},
		{Role: "user", Content: user},
	})
}

func (dsClient *deepSeekClient) CompleteMessages(ctx context.Context, msgs []dto.Message) (string, error) {
	if dsClient.apiKey == "" {
		return "", errors.New("DEEPSEEK_API_KEY not set")
	}

	body := map[string]any{
		"model":       dsClient.model,
		"messages":    msgs,
		"temperature": 0.2,
	}
	bs, _ := json.Marshal(body)

	req, _ := http.NewRequestWithContext(ctx, "POST",
		"https://api.deepseek.com/v1/chat/completions", bytes.NewReader(bs))
	req.Header.Set("Authorization", "Bearer "+dsClient.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		raw, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("deepseek: %s\n%s", res.Status, raw)
	}

	var jr struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(res.Body).Decode(&jr); err != nil {
		return "", err
	}
	if len(jr.Choices) == 0 {
		return "", errors.New("deepseek: no choices")
	}
	return jr.Choices[0].Message.Content, nil
}
