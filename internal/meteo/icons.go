package meteo

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// IconPath is the URL prefix the frontend uses for sky icons. Requests are
// answered by the app's own asset server (see IconStore.Middleware), so they
// work offline once an icon has been fetched.
const IconPath = "/meteo/icons/"

const (
	iconSource  = "https://www.meteogalicia.gal/web/assets/icons/svg/"
	maxIconSize = 256 << 10
)

// IconStore keeps MeteoGalicia's official sky icons on disk. A missing icon
// is downloaded once, on first use or when Prefetch is called after a
// forecast arrives, and then served locally forever (icons are keyed by code
// and never change).
type IconStore struct {
	dirFunc   func() string
	once      sync.Once
	dir       string
	source    string
	http      *http.Client
	userAgent string
	mu        sync.Mutex // serialises downloads
}

// NewIconStore creates a store. dirFunc is resolved lazily on first use,
// because the mobile storage path is only available once the app runs.
func NewIconStore(dirFunc func() string, userAgent string) *IconStore {
	return &IconStore{
		dirFunc:   dirFunc,
		source:    iconSource,
		http:      &http.Client{Timeout: 15 * time.Second},
		userAgent: userAgent,
	}
}

func (s *IconStore) directory() string {
	s.once.Do(func() {
		if s.dirFunc == nil {
			return
		}
		if d := s.dirFunc(); d != "" && os.MkdirAll(d, 0o755) == nil {
			s.dir = d
		}
	})
	return s.dir
}

// validSkyCode accepts the codes of annex II (101-125 day, 201-225 night).
func validSkyCode(code int) bool {
	return (code >= 101 && code <= 125) || (code >= 201 && code <= 225)
}

// Middleware serves IconPath requests and passes everything else on.
func (s *IconStore) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, IconPath) {
			next.ServeHTTP(w, r)
			return
		}
		s.serve(w, r)
	})
}

func (s *IconStore) serve(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, IconPath)
	code, err := strconv.Atoi(strings.TrimSuffix(name, ".svg"))
	if err != nil || !strings.HasSuffix(name, ".svg") || !validSkyCode(code) {
		http.NotFound(w, r)
		return
	}
	body, err := s.icon(r.Context(), code)
	if err != nil {
		http.Error(w, "icon unavailable", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	_, _ = w.Write(body)
}

// icon returns the SVG for a code, downloading and storing it if needed.
func (s *IconStore) icon(ctx context.Context, code int) ([]byte, error) {
	dir := s.directory()
	file := filepath.Join(dir, strconv.Itoa(code)+".svg")
	if dir != "" {
		if b, err := os.ReadFile(file); err == nil {
			return b, nil
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if dir != "" {
		// Another request may have stored it while we waited.
		if b, err := os.ReadFile(file); err == nil {
			return b, nil
		}
	}
	b, err := s.download(ctx, code)
	if err != nil {
		return nil, err
	}
	if dir != "" {
		tmp := file + ".tmp"
		if os.WriteFile(tmp, b, 0o644) == nil {
			_ = os.Rename(tmp, file)
		}
	}
	return b, nil
}

func (s *IconStore) download(ctx context.Context, code int) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.source+strconv.Itoa(code)+".svg", nil)
	if err != nil {
		return nil, err
	}
	if s.userAgent != "" {
		req.Header.Set("User-Agent", s.userAgent)
	}
	res, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meteogalicia icon %d: HTTP %d", code, res.StatusCode)
	}
	b, err := io.ReadAll(io.LimitReader(res.Body, maxIconSize+1))
	if err != nil {
		return nil, err
	}
	// Only store what really is an SVG (never an HTML error page).
	if len(b) > maxIconSize || !strings.Contains(string(b[:min(len(b), 512)]), "<svg") {
		return nil, fmt.Errorf("meteogalicia icon %d: not an SVG", code)
	}
	return b, nil
}

// Prefetch downloads the icons used by a forecast so they are available
// offline before the forecast view is opened. Errors are ignored.
func (s *IconStore) Prefetch(f Forecast) {
	seen := map[int]bool{}
	add := func(v *int) {
		if v != nil && validSkyCode(*v) {
			seen[*v] = true
		}
	}
	for _, d := range f.Days {
		add(d.Sky)
		for _, p := range d.Periods {
			add(p.Sky)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	for code := range seen {
		_, _ = s.icon(ctx, code)
	}
}
