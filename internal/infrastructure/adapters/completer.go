package adapters

import (
	"context"
	"lilith/internal/domain/dto"
)

// ILLMAdapter defines a common adapter interface for any LLM provider.
// that can generate completions from user/system messages.
type ICompletionAdapter interface {
	Complete(ctx context.Context, system, user string) (string, error)
	CompleteMessages(ctx context.Context, msgs []dto.Message) (string, error)
}
