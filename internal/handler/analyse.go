package handler

import (
	"context"
	"strings"

	"lilith/internal/domain/dto"
	"lilith/internal/domain/prompts"
	"lilith/internal/infrastructure/adapters"
	"lilith/internal/infrastructure/session"
)

const analyseMode = "analyse"

// AnalyseChat: multi-turn analysis; persists history in ~/.lilith/analyse.json.
// - diagram: allow Mermaid in prompt
// - fresh:   clear session before running (equivalent to --cleanup)
func AnalyseChat(ctx context.Context, completionAdapter adapters.ICompletionAdapter, user string, diagram, fresh bool) (string, error) {
	if fresh {
		_ = session.ResetMode(analyseMode)
	}

	var msgs []dto.Message
	_ = session.Load(analyseMode, &msgs)

	if len(msgs) == 0 || msgs[0].Role != "system" {
		sys := prompts.GlobalRules() + "\n\n" + prompts.AnalysePrompt(diagram)
		msgs = append([]dto.Message{{Role: "system", Content: sys}}, msgs...)
	}

	msgs = append(msgs, dto.Message{Role: "user", Content: user})

	out, err := completionAdapter.CompleteMessages(ctx, msgs)
	if err != nil {
		return "", err
	}

	reply := strings.TrimSpace(out)
	msgs = append(msgs, dto.Message{Role: "assistant", Content: reply})
	_ = session.Save(analyseMode, msgs)
	return reply, nil
}
