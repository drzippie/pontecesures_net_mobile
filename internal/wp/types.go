// Package wp is a small, read-only client for the WordPress REST API of
// pontecesures.net. It converts raw WordPress payloads into compact,
// frontend-friendly models (decoded titles, plain-text excerpts, sanitised
// HTML and a best-effort thumbnail extracted from the post body).
package wp

// Category is a WordPress post category.
type Category struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Count  int    `json:"count"`
	Parent int    `json:"parent"`
	// Illustration is the section artwork used by the website theme.
	Illustration string `json:"illustration,omitempty"`
}

// Image is a remote image reference.
type Image struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Alt    string `json:"alt"`
	// Fallback is true when the image is the section illustration rather
	// than a photo taken from the post itself.
	Fallback bool `json:"fallback"`
	// Landscape mirrors the website rule (width >= 1.3 * height): landscape
	// photos are cropped to fill the card, everything else (posters,
	// obituaries...) is shown whole because the text matters.
	Landscape bool `json:"landscape"`
}

// PostSummary is the lightweight representation used in lists.
type PostSummary struct {
	ID         int        `json:"id"`
	Date       string     `json:"date"` // Local site time, ISO 8601 without zone (Europe/Madrid).
	Title      string     `json:"title"`
	Excerpt    string     `json:"excerpt"`
	Link       string     `json:"link"`
	Author     string     `json:"author"`
	Categories []Category `json:"categories"`
	Image      *Image     `json:"image"`
}

// Post is a full article ready to be rendered.
type Post struct {
	PostSummary
	// Content is sanitised HTML safe to inject into the WebView.
	Content string `json:"content"`
}

// PostPage is a page of post summaries.
type PostPage struct {
	Items      []PostSummary `json:"items"`
	Page       int           `json:"page"`
	TotalPages int           `json:"totalPages"`
	Total      int           `json:"total"`
	// Stale is true when the data comes from the offline cache because the
	// network request failed.
	Stale bool `json:"stale"`
}

// rawPost mirrors the subset of the /wp/v2/posts payload we consume.
type rawPost struct {
	ID    int    `json:"id"`
	Date  string `json:"date"`
	Link  string `json:"link"`
	Title struct {
		Rendered string `json:"rendered"`
	} `json:"title"`
	Excerpt struct {
		Rendered string `json:"rendered"`
	} `json:"excerpt"`
	Content struct {
		Rendered string `json:"rendered"`
	} `json:"content"`
	Embedded struct {
		Author []struct {
			Name string `json:"name"`
		} `json:"author"`
		Terms [][]struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Slug     string `json:"slug"`
			Taxonomy string `json:"taxonomy"`
		} `json:"wp:term"`
	} `json:"_embedded"`
}

// Page is a static WordPress page (e.g. "Como colaborar", "Normas de uso").
type Page struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"` // sanitised HTML
	Link     string `json:"link"`
	Modified string `json:"modified"`
}
