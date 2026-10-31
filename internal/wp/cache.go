package wp

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// fileCache is a tiny persistent key/value cache used for offline reading.
// Each entry is a JSON file named after the SHA-1 of its key.
type fileCache struct {
	dir string
}

type cacheEntry struct {
	StoredAt   time.Time       `json:"storedAt"`
	Body       json.RawMessage `json:"body"`
	Total      int             `json:"total"`
	TotalPages int             `json:"totalPages"`
}

func newFileCache(dir string) *fileCache {
	if dir == "" {
		return &fileCache{}
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return &fileCache{}
	}
	return &fileCache{dir: dir}
}

func (c *fileCache) path(key string) string {
	sum := sha1.Sum([]byte(key))
	return filepath.Join(c.dir, hex.EncodeToString(sum[:])+".json")
}

func (c *fileCache) get(key string) (*cacheEntry, bool) {
	if c.dir == "" {
		return nil, false
	}
	b, err := os.ReadFile(c.path(key))
	if err != nil {
		return nil, false
	}
	var e cacheEntry
	if json.Unmarshal(b, &e) != nil {
		return nil, false
	}
	return &e, true
}

func (c *fileCache) put(key string, e cacheEntry) {
	if c.dir == "" {
		return
	}
	b, err := json.Marshal(e)
	if err != nil {
		return
	}
	tmp := c.path(key) + ".tmp"
	if os.WriteFile(tmp, b, 0o644) == nil {
		_ = os.Rename(tmp, c.path(key))
	}
}

// prune removes entries older than maxAge so the cache does not grow forever.
func (c *fileCache) prune(maxAge time.Duration) {
	if c.dir == "" {
		return
	}
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return
	}
	cutoff := time.Now().Add(-maxAge)
	for _, e := range entries {
		if info, err := e.Info(); err == nil && info.ModTime().Before(cutoff) {
			_ = os.Remove(filepath.Join(c.dir, e.Name()))
		}
	}
}
