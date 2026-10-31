package settings

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	p := filepath.Join(t.TempDir(), "sub", "settings.json")
	s := NewStore(p)
	if got := s.Load(); got != (Settings{}) {
		t.Fatalf("fresh load = %+v", got)
	}
	want := Settings{Language: "gl", Theme: "dark", TextSize: "large"}
	if err := s.Save(want); err != nil {
		t.Fatal(err)
	}
	if got := NewStore(p).Load(); got != want {
		t.Fatalf("reload = %+v, want %+v", got, want)
	}
}

func TestRejectsInvalidValues(t *testing.T) {
	p := filepath.Join(t.TempDir(), "settings.json")
	s := NewStore(p)
	for _, bad := range []Settings{{Language: "en"}, {Theme: "blue"}, {TextSize: "huge"}} {
		if err := s.Save(bad); err != ErrInvalid {
			t.Fatalf("Save(%+v) = %v", bad, err)
		}
	}
	// A tampered file falls back to defaults.
	_ = os.WriteFile(p, []byte(`{"language":"xx"}`), 0o644)
	if got := s.Load(); got != (Settings{}) {
		t.Fatalf("invalid file load = %+v", got)
	}
}
