package wp

import (
	"html"
	"regexp"
	"strconv"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var (
	tagRe      = regexp.MustCompile(`(?s)<[^>]*>`)
	spaceRe    = regexp.MustCompile(`\s+`)
	imgRe      = regexp.MustCompile(`(?is)<img\b[^>]*>`)
	attrRe     = regexp.MustCompile(`(?is)\b(src|srcset|width|height|alt)\s*=\s*("([^"]*)"|'([^']*)')`)
	readMoreRe = regexp.MustCompile(`(?i)\s*(\[(&hellip;|…|\.\.\.)\]|Continue reading.*|Ler máis.*|Leer más.*)$`)

	// Iframes are only allowed from well-known video providers.
	embedSrcRe = regexp.MustCompile(`^https://(www\.)?(youtube\.com|youtube-nocookie\.com|player\.vimeo\.com|www\.facebook\.com)/`)

	policy = newPolicy()
)

func newPolicy() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("srcset", "sizes", "loading", "decoding").OnElements("img")
	// Keep the orientation classes added by the theme (e-apaisada/e-vertical).
	p.AllowAttrs("class").Matching(regexp.MustCompile(`^[a-zA-Z0-9_\- ]+$`)).OnElements("img", "figure")
	p.AllowElements("figure", "figcaption", "video", "source", "audio")
	p.AllowAttrs("src", "controls", "poster", "preload").OnElements("video", "audio", "source")
	p.AllowAttrs("type").OnElements("source")
	p.AllowAttrs("src").Matching(embedSrcRe).OnElements("iframe")
	p.AllowAttrs("width", "height", "allowfullscreen", "allow", "title").OnElements("iframe")
	p.RequireNoReferrerOnFullyQualifiedLinks(false)
	p.AddTargetBlankToFullyQualifiedLinks(false)
	return p
}

// SanitizeHTML removes scripts, inline handlers, styles and any element not
// explicitly allowed. The result is safe to render with x-html.
func SanitizeHTML(s string) string {
	out := policy.Sanitize(s)
	// Old posts embed http:// media, which iOS blocks (ATS).
	out = strings.ReplaceAll(out, `src="http://`, `src="https://`)
	return strings.TrimSpace(out)
}

// PlainText strips tags, decodes entities and collapses whitespace.
func PlainText(s string) string {
	s = tagRe.ReplaceAllString(s, " ")
	s = html.UnescapeString(s)
	s = strings.ReplaceAll(s, " ", " ")
	return strings.TrimSpace(spaceRe.ReplaceAllString(s, " "))
}

// Excerpt returns a plain-text excerpt without "read more" trailers,
// truncated to at most maxRunes runes on a word boundary.
func Excerpt(s string, maxRunes int) string {
	t := readMoreRe.ReplaceAllString(PlainText(s), "")
	r := []rune(t)
	if len(r) <= maxRunes {
		return t
	}
	cut := string(r[:maxRunes])
	if i := strings.LastIndex(cut, " "); i > maxRunes/2 {
		cut = cut[:i]
	}
	return strings.TrimRight(cut, " ,.;:") + "…"
}

// skipImageRe matches avatars, emoji and tracking pixels, which are never
// used as a post thumbnail (same rule as the website theme).
var skipImageRe = regexp.MustCompile(`(?i)(gravatar\.com|/emoji/|\.gif$|^data:)`)

// FirstImage extracts the first usable <img> of an HTML fragment and picks
// the srcset candidate closest to (but not below) targetWidth, so list
// thumbnails do not download full-size originals.
func FirstImage(content string, targetWidth int) *Image {
	for _, tag := range imgRe.FindAllString(content, -1) {
		if img := parseImg(tag, targetWidth); img != nil {
			return img
		}
	}
	return nil
}

func parseImg(tag string, targetWidth int) *Image {
	attrs := map[string]string{}
	for _, m := range attrRe.FindAllStringSubmatch(tag, -1) {
		v := m[3]
		if v == "" {
			v = m[4]
		}
		attrs[strings.ToLower(m[1])] = html.UnescapeString(v)
	}
	src := secureURL(attrs["src"])
	if src == "" || skipImageRe.MatchString(src) {
		return nil
	}
	img := &Image{URL: src, Alt: attrs["alt"]}
	img.Width, _ = strconv.Atoi(attrs["width"])
	img.Height, _ = strconv.Atoi(attrs["height"])
	img.Landscape = img.Height > 0 && float64(img.Width) >= float64(img.Height)*1.3

	if best, w := pickSrcset(attrs["srcset"], targetWidth); best != "" {
		if img.Width > 0 && img.Height > 0 && w > 0 {
			img.Height = img.Height * w / img.Width
		}
		img.URL = best
		if w > 0 {
			img.Width = w
		}
	}
	return img
}

// secureURL upgrades plain-http URLs to https (iOS App Transport Security
// blocks cleartext loads) and rejects anything that is not a web URL.
func secureURL(u string) string {
	u = strings.TrimSpace(u)
	switch {
	case strings.HasPrefix(u, "https://"):
		return u
	case strings.HasPrefix(u, "http://"):
		return "https://" + strings.TrimPrefix(u, "http://")
	case strings.HasPrefix(u, "//"):
		return "https:" + u
	}
	return ""
}

func pickSrcset(srcset string, target int) (string, int) {
	bestURL, bestW := "", 0
	fallbackURL, fallbackW := "", 0
	for _, cand := range strings.Split(srcset, ",") {
		f := strings.Fields(strings.TrimSpace(cand))
		if len(f) != 2 || !strings.HasSuffix(f[1], "w") {
			continue
		}
		w, err := strconv.Atoi(strings.TrimSuffix(f[1], "w"))
		f[0] = secureURL(f[0])
		if err != nil || f[0] == "" {
			continue
		}
		if w >= target && (bestW == 0 || w < bestW) {
			bestURL, bestW = f[0], w
		}
		if w > fallbackW {
			fallbackURL, fallbackW = f[0], w
		}
	}
	if bestURL != "" {
		return bestURL, bestW
	}
	return fallbackURL, fallbackW
}
