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

func TestRelatedFallsBackWithoutPlugin(t *testing.T) {
	var pluginHits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		switch {
		case strings.HasPrefix(r.URL.Path, "/wp-json/pcn/v1/"):
			pluginHits++
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"code":"rest_no_route"}`))
		case r.URL.Path == "/wp-json/wp/v2/posts/10":
			_, _ = w.Write([]byte(`{"tags":[1,2,3],"categories":[24]}`))
		case r.URL.Path == "/wp-json/wp/v2/tags":
			// Tag 1 is generic and must be dropped; 3 is rarer than 2.
			_, _ = w.Write([]byte(`[{"id":1,"count":3906},{"id":2,"count":42},{"id":3,"count":5}]`))
		case r.URL.Path == "/wp-json/wp/v2/posts" && q.Get("tags") == "3":
			// Rarest tag first.
			if q.Get("exclude") != "10" || q.Get("per_page") != "4" {
				t.Errorf("rare tag query = %v", q)
			}
			_, _ = w.Write([]byte(`[{"id":11,"title":{"rendered":"Etiqueta rara"}}]`))
		case r.URL.Path == "/wp-json/wp/v2/posts" && q.Get("tags") == "2":
			if q.Get("exclude") != "10,11" || q.Get("per_page") != "3" {
				t.Errorf("second tag query = %v", q)
			}
			_, _ = w.Write([]byte(`[{"id":13,"title":{"rendered":"Outra etiqueta"}}]`))
		case r.URL.Path == "/wp-json/wp/v2/posts" && q.Get("categories") != "":
			if q.Get("exclude") != "10,11,13" || q.Get("per_page") != "2" {
				t.Errorf("category query = %v", q)
			}
			_, _ = w.Write([]byte(`[{"id":12,"title":{"rendered":"Por sección"}}]`))
		default:
			t.Errorf("unexpected request %s", r.URL)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	c := NewClient(Options{BaseURL: srv.URL})
	got, err := c.Related(context.Background(), 10, 4)
	if err != nil || len(got) != 3 || got[0].ID != 11 || got[1].ID != 13 || got[2].ID != 12 {
		t.Fatalf("Related = %+v, %v", got, err)
	}
	// The missing plugin route is remembered.
	_, _ = c.Related(context.Background(), 10, 4)
	if pluginHits != 1 {
		t.Fatalf("plugin endpoint hit %d times, want 1", pluginHits)
	}
}

func TestRelatedUsesPluginOrder(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/wp-json/pcn/v1/posts/10/related":
			_, _ = w.Write([]byte(`{"ids":[30,20]}`))
		case "/wp-json/wp/v2/posts":
			q := r.URL.Query()
			if q.Get("include") != "30,20" || q.Get("orderby") != "include" {
				t.Errorf("include query = %v", q)
			}
			_, _ = w.Write([]byte(`[{"id":30},{"id":20}]`))
		default:
			t.Errorf("unexpected request %s", r.URL)
		}
	}))
	defer srv.Close()

	got, err := NewClient(Options{BaseURL: srv.URL}).Related(context.Background(), 10, 4)
	if err != nil || len(got) != 2 || got[0].ID != 30 {
		t.Fatalf("Related = %+v, %v", got, err)
	}
}
