package wp

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Related posts.
//
// The website computes "Máis sobre isto" in PHP (plugin "comments",
// includes/relacionadas.php): tags weighted by rarity plus a recency bonus,
// topped up with the latest posts of the same category. When the pcn-mobile
// plugin is installed, GET /pcn/v1/posts/{id}/related exposes exactly that
// list. Otherwise we approximate it with the core REST API.

// genericTagThreshold mirrors PONTECESURES_ETIQUETA_XENERICA: tags used by
// more posts than this say nothing about the topic.
const genericTagThreshold = 1500

// pcnRetryAfter is how long we stop asking the plugin endpoint after it was
// found missing, to avoid a wasted request per article.
const pcnRetryAfter = 6 * time.Hour

type pcnState struct {
	mu          sync.Mutex
	unavailable time.Time
}

func (s *pcnState) available() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return time.Since(s.unavailable) > pcnRetryAfter
}

func (s *pcnState) markUnavailable() {
	s.mu.Lock()
	s.unavailable = time.Now()
	s.mu.Unlock()
}

// Related returns up to n posts related to the given one.
func (c *Client) Related(ctx context.Context, id, n int) ([]PostSummary, error) {
	if n <= 0 || n > 8 {
		n = 4
	}
	if c.pcn.available() {
		ids, err := c.relatedFromPlugin(ctx, id, n)
		if err == nil {
			if len(ids) == 0 {
				return []PostSummary{}, nil
			}
			return c.postsByIDs(ctx, ids)
		}
		if errors.Is(err, ErrNotFound) {
			// Route missing: plugin not installed (or the site's related
			// posts function is not available). Stop asking for a while.
			c.pcn.markUnavailable()
		}
	}
	return c.relatedApprox(ctx, id, n)
}

func (c *Client) relatedFromPlugin(ctx context.Context, id, n int) ([]int, error) {
	q := url.Values{"per_page": {strconv.Itoa(n)}}
	body, _, _, _, err := c.fetch(ctx, "/pcn/v1/posts/"+strconv.Itoa(id)+"/related", q)
	if err != nil {
		return nil, err
	}
	var res struct {
		IDs []int `json:"ids"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res.IDs, nil
}

// relatedApprox approximates the website ranking with the core REST API:
// specific tags are queried one by one from the rarest to the most common
// (rarer tags identify the topic better), newest first within each tag, and
// the list is topped up with the latest posts of the same category.
func (c *Client) relatedApprox(ctx context.Context, id, n int) ([]PostSummary, error) {
	body, _, _, _, err := c.fetch(ctx, "/wp/v2/posts/"+strconv.Itoa(id), url.Values{"_fields": {"tags,categories"}})
	if err != nil {
		return nil, err
	}
	var meta struct {
		Tags       []int `json:"tags"`
		Categories []int `json:"categories"`
	}
	if err := json.Unmarshal(body, &meta); err != nil {
		return nil, err
	}

	out := []PostSummary{}
	exclude := []int{id}
	add := func(q url.Values) {
		if len(out) >= n {
			return
		}
		q.Set("exclude", joinInts(exclude))
		q.Set("per_page", strconv.Itoa(n-len(out)))
		items, err := c.postsWhere(ctx, q)
		if err != nil {
			return
		}
		for _, p := range items {
			out = append(out, p)
			exclude = append(exclude, p.ID)
		}
	}

	for _, tag := range c.specificTags(ctx, meta.Tags) {
		add(url.Values{"tags": {strconv.Itoa(tag)}})
	}
	if len(meta.Categories) > 0 {
		add(url.Values{"categories": {joinInts(meta.Categories)}})
	}
	return out, nil
}

// specificTags returns the post tags ordered from the rarest to the most
// common, dropping generic tags when more specific ones exist (same rule as
// the website). If every tag is generic they are all kept.
func (c *Client) specificTags(ctx context.Context, ids []int) []int {
	if len(ids) == 0 {
		return nil
	}
	body, _, _, _, err := c.fetch(ctx, "/wp/v2/tags", url.Values{
		"include":  {joinInts(ids)},
		"per_page": {"100"},
		"orderby":  {"count"},
		"order":    {"asc"},
		"_fields":  {"id,count"},
	})
	if err != nil {
		return ids
	}
	var tags []struct {
		ID    int `json:"id"`
		Count int `json:"count"`
	}
	if json.Unmarshal(body, &tags) != nil {
		return ids
	}
	sort.SliceStable(tags, func(i, j int) bool { return tags[i].Count < tags[j].Count })
	var specific, all []int
	for _, t := range tags {
		all = append(all, t.ID)
		if t.Count < genericTagThreshold {
			specific = append(specific, t.ID)
		}
	}
	if len(specific) == 0 {
		return all
	}
	return specific
}

// postsByIDs fetches summaries keeping the given order.
func (c *Client) postsByIDs(ctx context.Context, ids []int) ([]PostSummary, error) {
	return c.postsWhere(ctx, url.Values{
		"include":  {joinInts(ids)},
		"orderby":  {"include"},
		"per_page": {strconv.Itoa(len(ids))},
	})
}

func (c *Client) postsWhere(ctx context.Context, q url.Values) ([]PostSummary, error) {
	q.Set("_embed", "author,wp:term")
	q.Set("_fields", postFields)
	body, _, _, _, err := c.fetch(ctx, "/wp/v2/posts", q)
	if err != nil {
		return nil, err
	}
	var raw []rawPost
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	out := make([]PostSummary, 0, len(raw))
	for _, r := range raw {
		out = append(out, r.summary())
	}
	return out, nil
}

func joinInts(ids []int) string {
	s := make([]string, len(ids))
	for i, id := range ids {
		s[i] = strconv.Itoa(id)
	}
	return strings.Join(s, ",")
}
