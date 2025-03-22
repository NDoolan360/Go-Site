package projects

import "strings"

// Escape Markdown and HTML special characters.
func EscapeSpecialChars(s string) string {
	s = strings.ReplaceAll(s, "\u00A0", "&nbsp;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "[", "&#91;")
	s = strings.ReplaceAll(s, "/", "&#92;")
	s = strings.ReplaceAll(s, "]", "&#93;")
	s = strings.ReplaceAll(s, "^", "&#94;")
	s = strings.ReplaceAll(s, "_", "&#95;")
	s = strings.ReplaceAll(s, "`", "&#96;")
	s = strings.ReplaceAll(s, "{", "&#123;")
	s = strings.ReplaceAll(s, "|", "&#124;")
	s = strings.ReplaceAll(s, "}", "&#125;")
	s = strings.ReplaceAll(s, "~", "&#126;")
	return s
}
