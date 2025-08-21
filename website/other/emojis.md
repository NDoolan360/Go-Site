---
Title: Emojis
Styles:
  - /static/styles/code.css
Scripts:
  - /static/scripts/hyperscript@0.9.14.min.js
---

<input type="text" style="width:100%;" placeholder="Search Emojis..." _="on input show <tbody>tr/> in next <table/> when its textContent.toLowerCase() contains my value">

| Emoji | Custom Shortcode | Github Shortcodes |
| ----- | ---------------- | ----------------- |
{{- range .Global.Emojis}}
| {{ .String }} | {{range index .ShortNames "Custom"}} `{{.}}`&nbsp; {{end}} | {{range index .ShortNames "Github"}} `{{.}}`&nbsp; {{else}} - {{end}} |
{{- end }}
