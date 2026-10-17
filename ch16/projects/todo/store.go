package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Item is one line of the list.
type Item struct {
	Text    string     `json:"text"`
	Created time.Time  `json:"created"`
	Done    *time.Time `json:"done,omitempty"` // nil until it is done
}

// Store is the list plus where it lives.
type Store struct {
	Path  string
	Items []Item
}

// Load reads the file, or starts an empty list when there is none yet.
func Load(path string) (*Store, error) {
	s := &Store{Path: path}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &s.Items); err != nil {
		return nil, fmt.Errorf("%s is not a todo file: %w", path, err)
	}
	return s, nil
}

// Save writes the list back, creating the folder if needed. It writes to a temporary
// file and renames, so a crash mid-write cannot leave you with half a list.
func (s *Store) Save() error {
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.Items, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.Path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.Path)
}

// Add appends an item and returns its 1-based number.
func (s *Store) Add(text string, now time.Time) int {
	s.Items = append(s.Items, Item{Text: text, Created: now})
	return len(s.Items)
}

// MarkDone sets the done time on item n (1-based).
func (s *Store) MarkDone(n int, now time.Time) error {
	if n < 1 || n > len(s.Items) {
		return fmt.Errorf("no item %d", n)
	}
	if s.Items[n-1].Done != nil {
		return fmt.Errorf("item %d is already done", n)
	}
	s.Items[n-1].Done = &now
	return nil
}

// Remove deletes item n (1-based).
func (s *Store) Remove(n int) error {
	if n < 1 || n > len(s.Items) {
		return fmt.Errorf("no item %d", n)
	}
	s.Items = append(s.Items[:n-1], s.Items[n:]...)
	return nil
}

// DefaultPath is $TODO_FILE, or todo/todo.json under the user's config folder.
func DefaultPath() (string, error) {
	if p := os.Getenv("TODO_FILE"); p != "" {
		return p, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "todo", "todo.json"), nil
}
