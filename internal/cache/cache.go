// Package cache is a tiny persistent key/value store used for offline
// reading. Each entry is a JSON file named after the SHA-1 of its key.
package cache

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Store persists entries under a directory. The zero value (or a Store with
// an empty directory) is a valid no-op cache.
type Store struct {
	dir string
}

// Entry is a cached response body plus optional pagination metadata.
type Entry struct {
	StoredAt   time.Time       `json:"storedAt"`
	Body       json.RawMessage `json:"body"`
	Total      int             `json:"total,omitempty"`
	TotalPages int             `json:"totalPages,omitempty"`
}

// New returns a Store rooted at dir, creating it if needed. If dir is empty
// or cannot be created the store silently does nothing.
func New(dir string) *Store {
	if dir == "" || os.MkdirAll(dir, 0o755) != nil {
		return &Store{}
	}
	return &Store{dir: dir}
}

func (s *Store) path(key string) string {
	sum := sha1.Sum([]byte(key))
	return filepath.Join(s.dir, hex.EncodeToString(sum[:])+".json")
}

// Get returns the entry stored under key.
func (s *Store) Get(key string) (*Entry, bool) {
	if s == nil || s.dir == "" {
		return nil, false
	}
	b, err := os.ReadFile(s.path(key))
	if err != nil {
		return nil, false
	}
	var e Entry
	if json.Unmarshal(b, &e) != nil {
		return nil, false
	}
	return &e, true
}

// Put stores an entry atomically (write to a temp file, then rename).
func (s *Store) Put(key string, e Entry) {
	if s == nil || s.dir == "" {
		return
	}
	b, err := json.Marshal(e)
	if err != nil {
		return
	}
	tmp := s.path(key) + ".tmp"
	if os.WriteFile(tmp, b, 0o644) == nil {
		_ = os.Rename(tmp, s.path(key))
	}
}

// Prune removes entries older than maxAge so the cache does not grow forever.
func (s *Store) Prune(maxAge time.Duration) {
	if s == nil || s.dir == "" {
		return
	}
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return
	}
	cutoff := time.Now().Add(-maxAge)
	for _, e := range entries {
		if info, err := e.Info(); err == nil && info.ModTime().Before(cutoff) {
			_ = os.Remove(filepath.Join(s.dir, e.Name()))
		}
	}
}
