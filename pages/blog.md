---
Title: Blog Posts
StyleSheets:
  - /static/styles/utils/posts.css
---

# {{.Page.Title}}  {#posts}

{{ range .Site.BlogPosts }}
- <a class="post-card" href="{{ .URL }}">

  ### {{ .Title }} <time datetime="{{ .Created }}">{{ .Created }}</time>
  {{ .Description }}

  </a>
{{ end }}
