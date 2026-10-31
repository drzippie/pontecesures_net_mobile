package wp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/drzippie/pontecesures_net_mobile/internal/cache"
)

// DefaultBaseURL is the production WordPress site.
const DefaultBaseURL = "https://pontecesures.net"

// ErrNotFound is returned when a post does not exist.
var ErrNotFound = errors.New("not found")

// thumbWidth is the preferred srcset width for list thumbnails.
const thumbWidth = 400

// Client is a read-only WordPress REST client with an offline fallback cache.
type Client struct {
	base      string
	http      *http.Client
	cache     *cache.Store
	userAgent string
	pcn       pcnState
}

// Options configures a Client.
type Options struct {
	BaseURL   string
	CacheDir  string // Empty disables the offline cache.
	UserAgent string
	Timeout   time.Duration
}

// NewClient builds a Client.
func NewClient(o Options) *Client {
	if o.BaseURL == "" {
		o.BaseURL = DefaultBaseURL
	}
	if o.Timeout == 0 {
		o.Timeout = 15 * time.Second
	}
	if o.UserAgent == "" {
		o.UserAgent = "PontecesuresNetMobile/0.1"
	}
	c := &Client{
		base:      strings.TrimRight(o.BaseURL, "/"),
		http:      &http.Client{Timeout: o.Timeout},
		cache:     cache.New(o.CacheDir),
		userAgent: o.UserAgent,
	}
	go c.cache.Prune(30 * 24 * time.Hour)
	return c
}

// fetch performs a GET against /wp-json and returns the body plus pagination
// headers. On network/server failure it falls back to the cached copy and
// reports stale=true.
func (c *Client) fetch(ctx context.Context, path string, q url.Values) (body []byte, total, pages int, stale bool, err error) {
	u := c.base + "/wp-json" + path
	if len(q) > 0 {
		u += "?" + q.Encode()
	}

	body, total, pages, err = c.do(ctx, u)
	if err == nil {
		c.cache.Put(u, cache.Entry{StoredAt: time.Now(), Body: body, Total: total, TotalPages: pages})
		return body, total, pages, false, nil
	}
	if errors.Is(err, ErrNotFound) {
		return nil, 0, 0, false, err
	}
	if e, ok := c.cache.Get(u); ok {
		return e.Body, e.Total, e.TotalPages, true, nil
	}
	return nil, 0, 0, false, err
}

func (c *Client) do(ctx context.Context, u string) ([]byte, int, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, 0, 0, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, 0, 0, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(io.LimitReader(res.Body, 16<<20))
	if err != nil {
		return nil, 0, 0, err
	}
	switch {
	case res.StatusCode == http.StatusNotFound:
		return nil, 0, 0, ErrNotFound
	case res.StatusCode == http.StatusBadRequest && strings.Contains(string(b), "rest_post_invalid_page_number"):
		// Asking past the last page: treat as an empty page.
		return []byte("[]"), 0, 0, nil
	case res.StatusCode >= 400:
		return nil, 0, 0, fmt.Errorf("wordpress: HTTP %d", res.StatusCode)
	}
	total, _ := strconv.Atoi(res.Header.Get("X-WP-Total"))
	pages, _ := strconv.Atoi(res.Header.Get("X-WP-TotalPages"))
	return b, total, pages, nil
}

// Categories returns the categories that have at least minCount posts,
// ordered by number of posts (most used first).
func (c *Client) Categories(ctx context.Context, minCount int) ([]Category, error) {
	q := url.Values{
		"per_page":   {"100"},
		"orderby":    {"count"},
		"order":      {"desc"},
		"hide_empty": {"true"},
		"_fields":    {"id,name,slug,count,parent"},
	}
	body, _, _, _, err := c.fetch(ctx, "/wp/v2/categories", q)
	if err != nil {
		return nil, err
	}
	var raw []Category
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	out := make([]Category, 0, len(raw))
	for _, cat := range raw {
		if cat.Count < minCount || cat.Slug == "sincategoria" {
			continue
		}
		cat.Name = html.UnescapeString(cat.Name)
		cat.Illustration = IllustrationURL(IllustrationSlug([]Category{cat}))
		out = append(out, cat)
	}
	return out, nil
}

// ListQuery selects which posts to list.
type ListQuery struct {
	Page       int
	PerPage    int
	CategoryID int    // 0 = all
	Search     string // empty = no search
	// Order: "relevance" (only with Search) or "date" (default).
	Order string
}

const postFields = "id,date,link,title,excerpt,content,_links,_embedded"

// Posts returns a page of post summaries.
func (c *Client) Posts(ctx context.Context, lq ListQuery) (PostPage, error) {
	if lq.Page < 1 {
		lq.Page = 1
	}
	if lq.PerPage <= 0 || lq.PerPage > 50 {
		lq.PerPage = 12
	}
	q := url.Values{
		"page":     {strconv.Itoa(lq.Page)},
		"per_page": {strconv.Itoa(lq.PerPage)},
		"_embed":   {"author,wp:term"},
		"_fields":  {postFields},
	}
	if lq.CategoryID > 0 {
		q.Set("categories", strconv.Itoa(lq.CategoryID))
	}
	if s := strings.TrimSpace(lq.Search); s != "" {
		q.Set("search", s)
		if lq.Order == "relevance" {
			// Title matches first, then by date (WordPress core ranking).
			q.Set("orderby", "relevance")
		}
	}

	body, total, pages, stale, err := c.fetch(ctx, "/wp/v2/posts", q)
	if err != nil {
		return PostPage{}, err
	}
	var raw []rawPost
	if err := json.Unmarshal(body, &raw); err != nil {
		return PostPage{}, err
	}
	items := make([]PostSummary, 0, len(raw))
	for _, r := range raw {
		items = append(items, r.summary())
	}
	return PostPage{Items: items, Page: lq.Page, TotalPages: pages, Total: total, Stale: stale}, nil
}

// Post returns a full post by ID.
func (c *Client) Post(ctx context.Context, id int) (Post, error) {
	q := url.Values{"_embed": {"author,wp:term"}, "_fields": {postFields}}
	body, _, _, _, err := c.fetch(ctx, "/wp/v2/posts/"+strconv.Itoa(id), q)
	if err != nil {
		return Post{}, err
	}
	var r rawPost
	if err := json.Unmarshal(body, &r); err != nil {
		return Post{}, err
	}
	return r.full(), nil
}

// PostBySlug returns a full post by its slug.
func (c *Client) PostBySlug(ctx context.Context, slug string) (Post, error) {
	q := url.Values{"slug": {slug}, "_embed": {"author,wp:term"}, "_fields": {postFields}}
	body, _, _, _, err := c.fetch(ctx, "/wp/v2/posts", q)
	if err != nil {
		return Post{}, err
	}
	var raw []rawPost
	if err := json.Unmarshal(body, &raw); err != nil {
		return Post{}, err
	}
	if len(raw) == 0 {
		return Post{}, ErrNotFound
	}
	return raw[0].full(), nil
}

func (r rawPost) categories() []Category {
	var out []Category
	for _, group := range r.Embedded.Terms {
		for _, t := range group {
			if t.Taxonomy == "category" {
				out = append(out, Category{ID: t.ID, Name: html.UnescapeString(t.Name), Slug: t.Slug})
			}
		}
	}
	return out
}

func (r rawPost) summary() PostSummary {
	cats := r.categories()
	s := PostSummary{
		ID:         r.ID,
		Date:       r.Date,
		Title:      PlainText(r.Title.Rendered),
		Excerpt:    Excerpt(r.Excerpt.Rendered, 220),
		Link:       r.Link,
		Categories: cats,
		Image:      PostImage(r.Content.Rendered, cats, thumbWidth),
	}
	if len(r.Embedded.Author) > 0 {
		s.Author = html.UnescapeString(r.Embedded.Author[0].Name)
	}
	return s
}

func (r rawPost) full() Post {
	s := r.summary()
	// Larger hero image for the article view.
	s.Image = PostImage(r.Content.Rendered, s.Categories, 800)
	return Post{PostSummary: s, Content: SanitizeHTML(r.Content.Rendered)}
}
