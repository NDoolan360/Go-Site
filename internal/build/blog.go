package build

import (
	"slices"
	"strings"
)

type BlogPost map[string]string
type BlogPosts []BlogPost

func NewBlogPost(path string, meta map[string]any) (BlogPost, bool) {
	blogPost := BlogPost{
		"URL":      strings.ReplaceAll(path, " ", "%20"),
		"Created":  "",
		"Modified": "",
	}
	for key, value := range meta {
		if stringValue, ok := value.(string); ok {
			blogPost[key] = stringValue
		}
	}

	if meta == nil {
		return nil, false
	}

	if url, ok := blogPost["URL"]; ok && !strings.HasPrefix(url, "/blog/") {
		return nil, false
	}

	if isDraft, ok := blogPost["IsDraft"]; ok && isDraft == "true" {
		return nil, false
	}

	return blogPost, true
}

func (blogPosts BlogPosts) SortByDate() BlogPosts {
	slices.SortStableFunc(blogPosts, func(i, j BlogPost) int {
		return -strings.Compare(i["Created"], j["Created"])
	})
	return blogPosts
}
