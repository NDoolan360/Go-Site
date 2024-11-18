package codewrapper

import (
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/util"
)

func WrapperRenderer(w util.BufWriter, context highlighting.CodeBlockContext, entering bool) {
	language, ok := context.Language()
	lang := string(language)
	if entering {

	} else {

	}
	// code block with a language
	if ok && lang != "" {
		if entering {
			w.WriteString(`<figure class="codeblock" data-lang="` + lang + `">`)
			w.WriteString(`<figcaption>`)
			w.WriteString(`<button class="copycode" disabled>Copy</button>`)
			w.WriteString(`</figcaption>`)
		} else {
			w.WriteString(`</figure>`)
		}
		return
	}

	// code block with no language specified
	if language == nil {
		if entering {
			w.WriteString(`<figure class="codeblock">`)
			w.WriteString(`<figcaption>`)
			w.WriteString(`<button class="copycode" disabled>Copy</button>`)
			w.WriteString(`</figcaption>`)
			w.WriteString(`<pre class="chroma">`)
			w.WriteString(`<code>`)
		} else {
			w.WriteString(`</code>`)
			w.WriteString(`</pre>`)
			w.WriteString(`</figure>`)
		}
	}
}
