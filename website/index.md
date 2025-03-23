---
Title: n.doolan.dev
Description:
  Software Engineer, tinkering on side projects,
  shipping code at [Kaluza](https://kaluza.com),
  and learning a little along the way.
Styles:
  - /static/styles/list-item.css
ArticleLimit: 3
ProjectLimit: 3
---

## Hi, I'm Nathan 👋

I'm passionate about programming, technology, and learning new things I enjoy working on a variety of projects, tinkering with new technologies, and documenting some of that here.

[More about me](/about.html)

## Recent Articles

{{ range $i, $_ := .Global.Articles }}
  {{ if lt $i $.ArticleLimit }}
    {{ template "list-item" .Meta }}
  {{ end }}
{{ end }}

## Recent Projects

{{ range $i, $_ := .Global.Projects }}
  {{ if lt $i $.ProjectLimit }}
    {{ template "list-item" . }}
  {{ end }}
{{ end }}
