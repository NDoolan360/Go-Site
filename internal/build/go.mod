module github.com/NDoolan360/go-site/internal/build

go 1.23.1

replace internal/code_wrapper => ../code_wrapper

replace internal/inline_svg => ../inline_svg

require (
	github.com/alecthomas/chroma/v2 v2.14.0
	github.com/yuin/goldmark v1.7.8
	github.com/yuin/goldmark-emoji v1.0.4
	github.com/yuin/goldmark-highlighting/v2 v2.0.0-20230729083705-37449abec8cc
	github.com/yuin/goldmark-meta v1.1.0
	internal/code_wrapper v0.0.0
	internal/inline_svg v0.0.0
)

require (
	github.com/dlclark/regexp2 v1.11.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.6 // indirect
	golang.org/x/net v0.30.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
