---
Title: Projects
Styles:
  - /static/styles/list-item.css
---

## {{.Title}}

<!-- TODO: Link each project to an internal page that has links to the external -->

{{ range .Global.Projects }}
{{ template "list-item" . }}
{{ end }}
