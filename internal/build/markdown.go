package build

import (
	"bytes"
	"text/template"

	codewrapper "internal/code_wrapper"
	"internal/inline_svg"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var mdParser = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		meta.Meta,
		inline_svg.InlineSvg,
		extension.Typographer,
		emoji.Emoji,
		highlighting.NewHighlighting(
			highlighting.WithWrapperRenderer(codewrapper.WrapperRenderer),
			highlighting.WithFormatOptions(
				chromahtml.WithClasses(true),
				chromahtml.WithLineNumbers(true),
				chromahtml.LineNumbersInTable(true),
			),
		),
	),
	goldmark.WithParserOptions(
		parser.WithAttribute(),
		parser.WithHeadingAttribute(),
	),
	goldmark.WithRendererOptions(html.WithUnsafe(), inline_svg.WithParentPath(".")),
)

func GetMeta(source []byte) (map[string]any, error) {
	context := parser.NewContext()
	if err := mdParser.Convert(source, &bytes.Buffer{}, parser.WithContext(context)); err != nil {
		return nil, err
	}
	return meta.Get(context), nil
}

func ApplyTemplate(source *[]byte, data any) error {
	buf := &bytes.Buffer{}
	tmpl, err := template.New("markdown").Parse(string(*source))
	if err != nil {
		return err
	}
	tmpl.Execute(buf, data)
	*source = buf.Bytes()
	return nil
}

func ToHtml(source *[]byte) error {
	buf := &bytes.Buffer{}
	if err := mdParser.Convert(*source, buf); err != nil {
		return err
	}
	*source = buf.Bytes()
	return nil
}
