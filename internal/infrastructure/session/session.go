package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Path returns where session files live (~/.lilith/<mode>.json)
func Path(mode string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".lilith")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, fmt.Sprintf("%s.json", mode)), nil
}

// Save persists the session messages
func Save(mode string, msgs any) error {
	path, err := Path(mode)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(msgs)
}

// Load reads messages (or returns nil if no session yet)
func Load(mode string, out any) error {
	path, err := Path(mode)
	if err != nil {
		return err
	}
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil // no history yet
	}
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(out)
}

// ResetMode clears the stored session
func ResetMode(mode string) error {
	path, err := Path(mode)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
