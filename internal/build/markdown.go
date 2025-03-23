package build

import (
	"bytes"
	"path"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	fences "github.com/stefanfritsch/goldmark-fences"
	"github.com/yuin/goldmark"
	g_emoji "github.com/yuin/goldmark-emoji"
	g_emoji_def "github.com/yuin/goldmark-emoji/definition"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	goldmark_parser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"

	emoji "internal/emoji"
	inlinesvg "internal/inline_svg"
)

type MarkdownTransformer struct{}

func (MarkdownTransformer) newGoldmark(root string) goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			&fences.Extender{},
			inlinesvg.InlineSvg,
			extension.Typographer,
			g_emoji.New(g_emoji.WithEmojis(emojisAsGoldmark(emoji.GetEmojis()))),
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
			inlinesvg.WithParentPath(root),
		),
	)
}

func (p MarkdownTransformer) Transform(asset *Asset) error {
	if path.Ext(asset.Path) != ".md" {
		return nil
	}

	html := &bytes.Buffer{}
	if err := p.newGoldmark(asset.SourceRoot).Convert(asset.Data, html); err != nil {
		return err
	}

	asset.Path = strings.TrimSuffix(asset.Path, ".md") + ".html"
	asset.Data = html.Bytes()

	return nil
}

func emojisAsGoldmark(emojis []emoji.Emoji) g_emoji_def.Emojis {
	goldmarkEmojis := []g_emoji_def.Emoji{}
	for _, e := range emojis {
		goldmarkEmojis = append(goldmarkEmojis, g_emoji_def.Emoji{
			Name:       e.Name,
			Unicode:    e.Unicode,
			ShortNames: append(e.ShortNames["Custom"], e.ShortNames["Github"]...),
		})
	}

	return g_emoji_def.NewEmojis(goldmarkEmojis...)
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
