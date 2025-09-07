// internal/handler/write.go
package handler

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"lilith/internal/domain/dto"
	"lilith/internal/infrastructure/session"
)

// WriteLatestFrom writes the latest assistant response from a given session mode
// into a target markdown file.
func WriteLatestFrom(mode, target string) error {
	var history []dto.Message
	if err := session.Load(mode, &history); err != nil {
		return err
	}
	if len(history) == 0 {
		return errors.New("no history found for mode " + mode)
	}

	// find last assistant message
	var latest string
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role == "assistant" {
			latest = history[i].Content
			break
		}
	}
	if latest == "" {
		return errors.New("no assistant response found in history for " + mode)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	path := filepath.Join(home, "Documents", mode, target)
	fmt.Println(path)
	// write to file
	if err := os.WriteFile(path, []byte(latest), 0o644); err != nil {
		return err
	}
	fmt.Printf("Wrote %s from %s session\n", path, mode)
	return nil
}
