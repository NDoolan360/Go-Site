---
Name: list-item
IsComponent: true
---

- [{{ if .Image.Src }}![{{ if .Image.Alt }}{{ .Image.Alt }}{{end}}]({{ .Image.Src }}){{ end }}]({{.URL}})
  ### {{ if .Logo }}![Logo]({{ .Logo }}){{ end }} {{ .Title }} {.list-item}
  {{ if .Created }}<time datetime="{{ .Created }}">{{ .Created }}</time>{{ end }}

  {{ if .Description }}{{ .Description }}{{ end }}

  {{ range .Tags }}
  - {{ . }}
  {{ end }}
