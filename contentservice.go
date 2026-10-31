package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/drzippie/pontecesures_net_mobile/internal/saved"
	"github.com/drzippie/pontecesures_net_mobile/internal/wp"
)

// minCategoryPosts hides near-empty legacy categories from the section list.
const minCategoryPosts = 20

// ContentService exposes the public content of pontecesures.net to the
// frontend. All network access happens in Go: the WebView never talks to
// WordPress directly (no CORS issues, offline cache in one place).
type ContentService struct {
	client *wp.Client
	saved  *saved.Store
}

// ServiceStartup initialises the WordPress client once the platform storage
// path is available.
func (s *ContentService) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	s.client = wp.NewClient(wp.Options{
		CacheDir:  cacheDir("wp-cache"),
		UserAgent: "PontecesuresNetMobile/" + appVersion,
	})
	s.saved = saved.NewStore(dataPath("saved"))
	// Light status bar icons over the dark app chrome (no-op on desktop).
	application.Mobile.SetStatusBar(`{"style":"light"}`)
	return nil
}

// cacheDir returns a per-feature cache directory: the app's private storage
// on mobile, the user cache directory on desktop.
func cacheDir(name string) string {
	if dir := application.Mobile.StoragePath(); dir != "" {
		return filepath.Join(dir, name)
	}
	if dir, err := os.UserCacheDir(); err == nil {
		return filepath.Join(dir, "pontecesures_net_mobile", name)
	}
	return ""
}

// dataPath returns a path for persistent user data (not a cache): the app's
// private storage on mobile, the user config dir on desktop.
func dataPath(name string) string {
	if dir := application.Mobile.StoragePath(); dir != "" {
		return filepath.Join(dir, name)
	}
	if dir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(dir, "pontecesures_net_mobile", name)
	}
	return ""
}

func reqCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 20*time.Second)
}

// Categories returns the sections with content, most used first.
func (s *ContentService) Categories() ([]wp.Category, error) {
	ctx, cancel := reqCtx()
	defer cancel()
	return s.client.Categories(ctx, minCategoryPosts)
}

// Posts returns a page of posts, optionally filtered by section and/or a
// search term. order is "relevance" or "date" (only relevant for searches).
func (s *ContentService) Posts(page int, categoryID int, search string, order string) (wp.PostPage, error) {
	ctx, cancel := reqCtx()
	defer cancel()
	return s.client.Posts(ctx, wp.ListQuery{Page: page, CategoryID: categoryID, Search: search, Order: order})
}

// Post returns a full article. Saved articles are always readable offline,
// and their stored copy is refreshed whenever the network copy loads.
func (s *ContentService) Post(id int) (wp.Post, error) {
	ctx, cancel := reqCtx()
	defer cancel()
	p, err := s.client.Post(ctx, id)
	if err != nil {
		if sp, ok := s.saved.Get(id); ok {
			return sp, nil
		}
		return p, err
	}
	if _, ok := s.saved.Get(id); ok {
		_ = s.saved.Put(p)
	}
	return p, nil
}

// SaveArticle stores an article on the device for offline reading.
func (s *ContentService) SaveArticle(id int) error {
	ctx, cancel := reqCtx()
	defer cancel()
	p, err := s.client.Post(ctx, id)
	if err != nil {
		return err
	}
	return s.saved.Put(p)
}

// RemoveSaved deletes a saved article.
func (s *ContentService) RemoveSaved(id int) error {
	return s.saved.Remove(id)
}

// SavedIDs returns the ids of the saved articles.
func (s *ContentService) SavedIDs() []int {
	ids := s.saved.IDs()
	if ids == nil {
		ids = []int{}
	}
	return ids
}

// Saved returns the saved articles, most recently saved first.
func (s *ContentService) Saved() []wp.PostSummary {
	return s.saved.List()
}

// permalinkRe matches pontecesures.net post permalinks: /YYYY/MM/DD/slug/
var permalinkRe = regexp.MustCompile(`^/\d{4}/\d{2}/\d{2}/([^/]+)/?$`)

// ErrNotAPost is returned by PostByURL for links that are not articles.
var ErrNotAPost = errors.New("not a pontecesures.net article link")

// PostByURL resolves an internal pontecesures.net link (permalink or ?p=ID)
// to an article, so links inside posts open in the app.
func (s *ContentService) PostByURL(link string) (wp.Post, error) {
	u, err := url.Parse(link)
	if err != nil || !isSiteHost(u.Host) {
		return wp.Post{}, ErrNotAPost
	}
	ctx, cancel := reqCtx()
	defer cancel()
	if id, err := strconv.Atoi(u.Query().Get("p")); err == nil && id > 0 {
		return s.client.Post(ctx, id)
	}
	m := permalinkRe.FindStringSubmatch(u.Path)
	if m == nil {
		return wp.Post{}, ErrNotAPost
	}
	slug, _ := url.PathUnescape(m[1])
	if id, err := strconv.Atoi(slug); err == nil {
		// Some posts have a numeric slug equal to their ID.
		if p, err := s.client.Post(ctx, id); err == nil {
			return p, nil
		}
	}
	return s.client.PostBySlug(ctx, slug)
}

func isSiteHost(h string) bool {
	h = strings.ToLower(h)
	return h == "pontecesures.net" || h == "www.pontecesures.net"
}

// OpenURL opens a link in the system browser.
func (s *ContentService) OpenURL(link string) error {
	u, err := url.Parse(link)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http" && u.Scheme != "mailto" && u.Scheme != "tel") {
		return errors.New("unsupported link")
	}
	if application.System.IsMobile() {
		application.Mobile.OpenURL(u.String())
		return nil
	}
	return application.Get().Browser.OpenURL(u.String())
}

// Share opens the native share sheet on mobile. On desktop it copies the
// link to the clipboard and returns false so the UI can say so.
func (s *ContentService) Share(title, link string) bool {
	if application.System.IsMobile() {
		payload, _ := json.Marshal(map[string]string{"text": title, "url": link})
		application.Mobile.Share(string(payload))
		return true
	}
	application.Get().Clipboard.SetText(link)
	return false
}

// Insets are the system bar insets reported by the native host.
type Insets struct {
	Top    float64 `json:"top"`
	Bottom float64 `json:"bottom"`
	Left   float64 `json:"left"`
	Right  float64 `json:"right"`
}

// SafeArea returns the raw system bar insets (physical pixels on Android,
// points on iOS; zero on desktop). Android's WebView does not expose them
// through CSS env(safe-area-inset-*) when drawing edge-to-edge, so the
// frontend converts them to CSS variables.
func (s *ContentService) SafeArea() Insets {
	var in Insets
	_ = json.Unmarshal([]byte(application.Mobile.SafeAreaJSON()), &in)
	return in
}

// relatedCount is how many related posts are shown under an article (same as
// the website).
const relatedCount = 4

// Related returns the posts shown under an article ("Más sobre esto").
func (s *ContentService) Related(id int) ([]wp.PostSummary, error) {
	ctx, cancel := reqCtx()
	defer cancel()
	return s.client.Related(ctx, id, relatedCount)
}
