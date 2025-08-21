package build

import (
	"bytes"
	"path"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	fences "github.com/stefanfritsch/goldmark-fences"
	"github.com/yuin/goldmark"
	g_emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	goldmark_parser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"

	emoji "internal/emoji"
	inlinesvg "internal/inline_svg"
)

type MarkdownTransformer struct {
	root string
}

func (p MarkdownTransformer) newGoldmark() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			&fences.Extender{},
			inlinesvg.InlineSvg,
			extension.Typographer,
			g_emoji.New(g_emoji.WithEmojis(emoji.GetEmojisAsGoldmark())),
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithWrapperRenderer(codeWrapperRenderer),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
					chromahtml.WithLineNumbers(true),
					chromahtml.LineNumbersInTable(true),
				),
			),
		),
		goldmark.WithParserOptions(
			goldmark_parser.WithAttribute(),
			goldmark_parser.WithHeadingAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			inlinesvg.WithParentPath(p.root),
		),
	)
}

func (p MarkdownTransformer) WithRoot(root string) MarkdownTransformer {
	return MarkdownTransformer{root: root}
}

func (p MarkdownTransformer) Transform(asset *Asset) error {
	if path.Ext(asset.Path) != ".md" {
		return nil
	}

	html := &bytes.Buffer{}
	if err := p.newGoldmark().Convert(asset.Data, html); err != nil {
		return err
	}

	asset.Path = strings.TrimSuffix(asset.Path, ".md") + ".html"
	asset.Data = html.Bytes()

	return nil
}

func codeWrapperRenderer(w util.BufWriter, context highlighting.CodeBlockContext, entering bool) {
	language, ok := context.Language()
	lang := string(language)

	// code block with a language
	noLang := !ok || language == nil
	if entering {
		w.WriteString(`<figure class="codeblock"`)
		if !noLang {
			w.WriteString(` data-lang="` + lang + `"`)
		}
		w.WriteString(`>`)

		w.WriteString(`<figcaption>`)
		w.WriteString(`<button class="copycode" disabled>Copy</button>`)
		w.WriteString(`</figcaption>`)

		if noLang {
			w.WriteString(`<pre class="chroma">`)
			w.WriteString(`<code>`)
		}
	} else {
		if noLang {
			w.WriteString(`</code>`)
			w.WriteString(`</pre>`)
		}

		w.WriteString(`</figure>`)
	}
}
