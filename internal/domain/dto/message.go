package dto

// Message represents one chat turn.
type Message struct {
	Role    string `json:"role"`    // "system", "user", "assistant"
	Content string `json:"content"` // plain text
}
