package wp

// Section illustrations, ported from the website theme
// (wp-content/themes/pontecesures/inc/imaxes.php). pontecesures.net does not
// use featured images: a post shows the first image of its body or, when it
// has none, the illustration of its section.

const illustrationBaseURL = "https://pontecesures.net/wp-content/themes/pontecesures/assets/img/categorias/"

// Categories that have their own illustration.
var illustratedSlugs = map[string]bool{
	"sociedade": true, "politica": true, "cultura": true, "festas": true,
	"obras": true, "deportes": true, "sucesos": true, "emprego": true,
	"medio-ambiente": true, "esquelas": true, "turismo": true, "avisos": true,
}

// Categories without an illustration that borrow another one.
var illustrationEquivalents = map[string]string{
	"concello":         "politica",
	"pontecesures-cf":  "deportes",
	"formacionemprego": "emprego",
}

// Illustrations that were replaced under a new file name (cache busting).
var illustrationFiles = map[string]string{
	"politica":  "politica-v2",
	"sociedade": "sociedade-v2",
	"turismo":   "turismo-v2",
}

// IllustrationSlug returns the illustration that corresponds to a list of
// post categories, following the theme's precedence rules.
func IllustrationSlug(cats []Category) string {
	for _, c := range cats {
		if illustratedSlugs[c.Slug] {
			return c.Slug
		}
	}
	for _, c := range cats {
		if eq, ok := illustrationEquivalents[c.Slug]; ok {
			return eq
		}
	}
	return "xeral"
}

// IllustrationURL returns the absolute URL of an illustration slug.
func IllustrationURL(slug string) string {
	if f, ok := illustrationFiles[slug]; ok {
		slug = f
	}
	if slug != "xeral" && !illustratedSlugs[slug] && !isReplacedFile(slug) {
		slug = "xeral"
	}
	return illustrationBaseURL + slug + ".jpg"
}

func isReplacedFile(name string) bool {
	for _, f := range illustrationFiles {
		if f == name {
			return true
		}
	}
	return false
}

// PostImage resolves the image of a post: first body image, otherwise the
// section illustration.
func PostImage(content string, cats []Category, targetWidth int) *Image {
	if img := FirstImage(content, targetWidth); img != nil {
		return img
	}
	return &Image{
		URL:       IllustrationURL(IllustrationSlug(cats)),
		Fallback:  true,
		Landscape: true,
	}
}
