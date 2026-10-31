// Package settings stores the user's app preferences on the device.
//
// Today it only holds display options for anonymous users. When accounts are
// added, profile data will live alongside (never the auth token: that goes to
// secure storage).
package settings

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// Settings are the user's display preferences. Empty values mean "not chosen
// yet": the frontend then applies its defaults (device language, system
// theme, normal text size).
type Settings struct {
	Language string `json:"language"` // "", "es", "gl"
	Theme    string `json:"theme"`    // "", "auto", "light", "dark"
	TextSize string `json:"textSize"` // "", "normal", "large", "xlarge"
}

var allowed = map[string]map[string]bool{
	"language": {"": true, "es": true, "gl": true},
	"theme":    {"": true, "auto": true, "light": true, "dark": true},
	"textSize": {"": true, "normal": true, "large": true, "xlarge": true},
}

// ErrInvalid is returned for unknown option values.
var ErrInvalid = errors.New("settings: invalid value")

// Validate checks every field against the allowed values.
func (s Settings) Validate() error {
	if !allowed["language"][s.Language] || !allowed["theme"][s.Theme] || !allowed["textSize"][s.TextSize] {
		return ErrInvalid
	}
	return nil
}

// Store persists Settings as a JSON file.
type Store struct {
	mu   sync.Mutex
	path string
}

// NewStore returns a store writing to path. An empty path keeps settings in
// memory only.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Load returns the saved settings, or zero values if none were saved (or the
// file is unreadable/invalid).
func (s *Store) Load() Settings {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out Settings
	if s.path == "" {
		return out
	}
	b, err := os.ReadFile(s.path)
	if err != nil || json.Unmarshal(b, &out) != nil || out.Validate() != nil {
		return Settings{}
	}
	return out
}

// Save validates and writes the settings atomically.
func (s *Store) Save(v Settings) error {
	if err := v.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
