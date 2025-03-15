---
Styles:
  - /static/styles/code.css
Scripts:
  - /static/scripts/hyperscript@0.9.14.min.js
---

## Emojis

<input type="text" placeholder="Search Emojis..."
  _="on input
     show <tr/> in next <table>tbody/>
     when its textContent.toLowerCase() contains my value.toLowerCase()">

| Emoji | Custom Shortcode | Github Shortcodes |
| ----- | ---------------- | ----------------- |
{{- range .Global.Emojis}}
| {{ .String }} | {{range index .ShortNames "Custom"}} `{{.}}` {{end}} | {{range index .ShortNames "Github"}} `{{.}}` {{else}} - {{end}} |
{{- end }}
