package wp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPlainTextAndExcerpt(t *testing.T) {
	got := PlainText("<p>Samaín &#8211; Pontecesures&nbsp;2026</p>")
	if got != "Samaín – Pontecesures 2026" {
		t.Fatalf("PlainText = %q", got)
	}
	ex := Excerpt("<p>Unha nova moi longa que segue e segue [&hellip;]</p>", 20)
	if !strings.HasSuffix(ex, "…") || len([]rune(ex)) > 21 {
		t.Fatalf("Excerpt = %q", ex)
	}
	if e := Excerpt("<p>Curta [&hellip;]</p>", 100); e != "Curta" {
		t.Fatalf("Excerpt trailer = %q", e)
	}
}

func TestFirstImagePicksSrcsetAndOrientation(t *testing.T) {
	html := `<p><img src="https://secure.gravatar.com/a.png"></p>
<img width="2048" height="1152" src="https://pontecesures.net/a.jpg" alt="Foto"
 srcset="https://pontecesures.net/a-600.jpg 600w, https://pontecesures.net/a-400.jpg 400w, https://pontecesures.net/a.jpg 2048w">`
	img := FirstImage(html, 400)
	if img == nil || img.URL != "https://pontecesures.net/a-400.jpg" || img.Width != 400 || img.Height != 225 {
		t.Fatalf("FirstImage = %+v", img)
	}
	if !img.Landscape || img.Fallback || img.Alt != "Foto" {
		t.Fatalf("flags = %+v", img)
	}
	poster := FirstImage(`<img src="http://pontecesures.net/p.jpg" width="500" height="700">`, 400)
	if poster == nil || poster.Landscape || poster.URL != "https://pontecesures.net/p.jpg" {
		t.Fatalf("poster = %+v", poster)
	}
}

func TestIllustrationFallback(t *testing.T) {
	cases := []struct {
		cats []Category
		want string
	}{
		{[]Category{{Slug: "concello"}}, "politica-v2.jpg"},
		{[]Category{{Slug: "sincategoria"}, {Slug: "cultura"}}, "cultura.jpg"},
		{[]Category{{Slug: "na-prensa"}}, "xeral.jpg"},
		{nil, "xeral.jpg"},
	}
	for _, c := range cases {
		img := PostImage("<p>sen imaxe</p>", c.cats, 400)
		if !img.Fallback || !strings.HasSuffix(img.URL, "/categorias/"+c.want) {
			t.Errorf("%v -> %s, want %s", c.cats, img.URL, c.want)
		}
	}
}

func TestSanitizeHTML(t *testing.T) {
	in := `<p onclick="x()">Ola<script>alert(1)</script></p><img class="e-vertical" src="http://pontecesures.net/i.jpg"><iframe src="https://evil.example/"></iframe><iframe src="https://www.youtube.com/embed/abc"></iframe>`
	out := SanitizeHTML(in)
	for _, bad := range []string{"onclick", "<script", "evil.example", `src="http://`} {
		if strings.Contains(out, bad) {
			t.Errorf("sanitised output contains %q: %s", bad, out)
		}
	}
	for _, good := range []string{"e-vertical", "youtube.com/embed/abc"} {
		if !strings.Contains(out, good) {
			t.Errorf("sanitised output lost %q: %s", good, out)
		}
	}
}

func TestOfflineFallback(t *testing.T) {
	up := true
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !up {
			http.Error(w, "down", http.StatusBadGateway)
			return
		}
		w.Header().Set("X-WP-Total", "1")
		w.Header().Set("X-WP-TotalPages", "1")
		_, _ = w.Write([]byte(`[{"id":7,"date":"2026-10-31T08:00:00","title":{"rendered":"Ola"},"excerpt":{"rendered":""},"content":{"rendered":""}}]`))
	}))
	defer srv.Close()

	c := NewClient(Options{BaseURL: srv.URL, CacheDir: t.TempDir()})
	ctx := context.Background()
	p, err := c.Posts(ctx, ListQuery{Page: 1})
	if err != nil || p.Stale || len(p.Items) != 1 {
		t.Fatalf("online: %+v %v", p, err)
	}
	up = false
	p, err = c.Posts(ctx, ListQuery{Page: 1})
	if err != nil || !p.Stale || p.Items[0].Title != "Ola" {
		t.Fatalf("offline: %+v %v", p, err)
	}
}
