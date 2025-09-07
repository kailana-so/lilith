package handler

import (
	"context"
	"strings"

	"lilith/internal/domain/dto"
	"lilith/internal/domain/prompts"
	"lilith/internal/infrastructure/adapters"
	"lilith/internal/infrastructure/session"
)

const chatMode = "chat"

func Chat(ctx context.Context, completionAdapter adapters.ICompletionAdapter, user string) (string, error) {
	var msgs []dto.Message
	_ = session.Load(chatMode, &msgs)

	if len(msgs) == 0 || msgs[0].Role != "system" {
		sys := prompts.GlobalRules() + "\n\n" + prompts.ChatPrompt()
		msgs = append([]dto.Message{{Role: "system", Content: sys}}, msgs...)
	}

	msgs = append(msgs, dto.Message{Role: "user", Content: user})

	out, err := completionAdapter.CompleteMessages(ctx, msgs)
	if err != nil {
		return "", err
	}

	reply := strings.TrimSpace(out)
	msgs = append(msgs, dto.Message{Role: "assistant", Content: reply})
	_ = session.Save(chatMode, msgs)
	return reply, nil
}
