package meteo

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func newTestStore(t *testing.T, up *bool, hits *int) (*IconStore, string) {
	t.Helper()
	src := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*hits++
		if !*up {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		if r.URL.Path == "/999.svg" {
			_, _ = w.Write([]byte("<html>not found</html>"))
			return
		}
		_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg"></svg>`))
	}))
	t.Cleanup(src.Close)
	dir := t.TempDir()
	s := NewIconStore(func() string { return dir }, "test")
	s.source = src.URL + "/"
	return s, dir
}

func TestIconMiddlewareCachesOnDisk(t *testing.T) {
	up, hits := true, 0
	s, dir := newTestStore(t, &up, &hits)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(299) })
	h := s.Middleware(next)

	get := func(path string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		return rec
	}

	if rec := get("/index.html"); rec.Code != 299 {
		t.Fatalf("non-icon path not passed through: %d", rec.Code)
	}
	if rec := get(IconPath + "101.svg"); rec.Code != 200 || rec.Header().Get("Content-Type") != "image/svg+xml" {
		t.Fatalf("first fetch: %d %v", rec.Code, rec.Header())
	}
	if _, err := os.Stat(filepath.Join(dir, "101.svg")); err != nil {
		t.Fatalf("icon not stored: %v", err)
	}
	// Offline: served from disk.
	up = false
	if rec := get(IconPath + "101.svg"); rec.Code != 200 || hits != 1 {
		t.Fatalf("offline: code=%d hits=%d", rec.Code, hits)
	}
	// Not cached and offline.
	if rec := get(IconPath + "102.svg"); rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("missing offline: %d", rec.Code)
	}
	// Invalid codes never reach the network.
	before := hits
	for _, p := range []string{"999.svg", "abc.svg", "101.png", "../101.svg"} {
		if rec := get(IconPath + p); rec.Code != http.StatusNotFound {
			t.Fatalf("%s: %d", p, rec.Code)
		}
	}
	if hits != before {
		t.Fatalf("invalid codes hit the network")
	}
}

func TestPrefetchSkipsCachedAndRejectsNonSVG(t *testing.T) {
	up, hits := true, 0
	s, dir := newTestStore(t, &up, &hits)
	sky, night, bad := 103, 201, -1
	f := Forecast{Days: []Day{{Sky: &sky, Periods: []Period{{Sky: &sky}, {Sky: &night}, {Sky: &bad}, {Sky: nil}}}}}
	s.Prefetch(f)
	if hits != 2 {
		t.Fatalf("hits = %d, want 2", hits)
	}
	s.Prefetch(f)
	if hits != 2 {
		t.Fatalf("cached icons downloaded again: %d", hits)
	}
	for _, c := range []string{"103.svg", "201.svg"} {
		if _, err := os.Stat(filepath.Join(dir, c)); err != nil {
			t.Fatal(err)
		}
	}
	s.source = s.source[:len(s.source)-1] + "/"
	if _, err := s.download(t.Context(), 999); err == nil {
		t.Fatal("HTML accepted as SVG")
	}
}
