module github.com/NDoolan360/go-site/internal/build

go 1.23.4

replace internal/inline_svg => ../md_inline_svg

replace internal/emoji => ../emoji

require (
	github.com/adrg/frontmatter v0.2.0
	github.com/alecthomas/chroma/v2 v2.15.0
	github.com/stefanfritsch/goldmark-fences v1.0.0
	github.com/tdewolff/minify/v2 v2.21.3
	github.com/yuin/goldmark v1.7.8
	github.com/yuin/goldmark-emoji v1.0.5
	github.com/yuin/goldmark-highlighting/v2 v2.0.0-20230729083705-37449abec8cc
	github.com/yuin/goldmark-meta v1.1.0
	internal/emoji v0.0.0
	internal/inline_svg v0.0.0
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/dlclark/regexp2 v1.11.4 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/tdewolff/parse/v2 v2.7.19 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
