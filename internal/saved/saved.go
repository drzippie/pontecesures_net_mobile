// Package saved keeps articles the user saved for later, fully on the
// device, so they can be read offline. Each article is one JSON file
// (<id>.json) holding the full post and the time it was saved.
//
// Accounts do not exist yet; when they do, this list can be synced with the
// user's profile.
package saved

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/drzippie/pontecesures_net_mobile/internal/wp"
)

// Store is a directory of saved articles.
type Store struct {
	mu  sync.Mutex
	dir string
}

type entry struct {
	SavedAt time.Time `json:"savedAt"`
	Post    wp.Post   `json:"post"`
}

// NewStore returns a store in dir (created on first write). An empty dir
// disables persistence.
func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

func (s *Store) file(id int) string {
	return filepath.Join(s.dir, strconv.Itoa(id)+".json")
}

// Put saves (or refreshes) an article.
func (s *Store) Put(p wp.Post) error {
	if s.dir == "" {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return err
	}
	e := entry{SavedAt: time.Now(), Post: p}
	// Keep the original save time when refreshing an existing entry.
	if old, ok := s.read(p.ID); ok {
		e.SavedAt = old.SavedAt
	}
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	tmp := s.file(p.ID) + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.file(p.ID))
}

// Remove deletes a saved article (no error if it was not saved).
func (s *Store) Remove(id int) error {
	if s.dir == "" {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	err := os.Remove(s.file(id))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// Get returns a saved article.
func (s *Store) Get(id int) (wp.Post, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.read(id)
	return e.Post, ok
}

func (s *Store) read(id int) (entry, bool) {
	var e entry
	if s.dir == "" {
		return e, false
	}
	b, err := os.ReadFile(s.file(id))
	if err != nil || json.Unmarshal(b, &e) != nil {
		return e, false
	}
	return e, true
}

// IDs returns the ids of all saved articles.
func (s *Store) IDs() []int {
	s.mu.Lock()
	defer s.mu.Unlock()
	var ids []int
	entries, _ := os.ReadDir(s.dir)
	for _, de := range entries {
		name := de.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}
		if id, err := strconv.Atoi(strings.TrimSuffix(name, ".json")); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

// List returns the saved articles as summaries, most recently saved first.
func (s *Store) List() []wp.PostSummary {
	type item struct {
		at time.Time
		p  wp.PostSummary
	}
	var items []item
	for _, id := range s.IDs() {
		s.mu.Lock()
		e, ok := s.read(id)
		s.mu.Unlock()
		if ok {
			items = append(items, item{e.SavedAt, e.Post.PostSummary})
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].at.After(items[j].at) })
	out := make([]wp.PostSummary, 0, len(items))
	for _, it := range items {
		out = append(out, it.p)
	}
	return out
}
