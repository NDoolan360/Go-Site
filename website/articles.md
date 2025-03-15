---
Title: Articles
Styles:
  - /static/styles/list-item.css
---

## {{.Title}}

{{ range .Global.Articles }}
{{ template "list-item" .Meta }}
{{ end }}
