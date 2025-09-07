// internal/handler/discuss_chat.go
package handler

import (
	"context"
	"strings"

	"lilith/internal/domain/dto"
	"lilith/internal/domain/prompts"
	"lilith/internal/infrastructure/adapters"
	"lilith/internal/infrastructure/session"
)

const discussMode = "discuss"

// DiscussChat: chat-like flow for discussion threads
func DiscussChat(ctx context.Context, completionAdapter adapters.ICompletionAdapter, user string, diagram, fresh bool) (string, error) {
	if fresh {
		_ = session.ResetMode(discussMode)
	}

	var msgs []dto.Message
	_ = session.Load(discussMode, &msgs)

	if len(msgs) == 0 || msgs[0].Role != "system" {
		sys := prompts.GlobalRules() + "\n\n" + prompts.DiscussPrompt(diagram)
		msgs = append([]dto.Message{{Role: "system", Content: sys}}, msgs...)
	}

	msgs = append(msgs, dto.Message{Role: "user", Content: user})

	out, err := completionAdapter.CompleteMessages(ctx, msgs)
	if err != nil {
		return "", err
	}

	reply := strings.TrimSpace(out)
	msgs = append(msgs, dto.Message{Role: "assistant", Content: reply})
	_ = session.Save(discussMode, msgs)
	return reply, nil
}
