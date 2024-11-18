---
Title: Projects
StyleSheets:
  - /static/styles/utils/projects.css
---

# {{.Page.Title}}  {#projects}

{{ range .Site.Projects }}
- <div class="project-card-content">

  ### {{ .Title }} ![Logo]({{ .Logo }})
  {{if .Description}}{{ .Description }}{{end}}
  {{ range .Tags }}
  - {{ . }}
  {{ end }}

  </div>
  <a href="{{ .URL }}"></a>
  <img src="{{ .Image.Src }}" alt="{{ .Image.Alt }}">
{{ end }}
