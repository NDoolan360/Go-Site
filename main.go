package main

import (
	"embed"
	"log"
	"os"
	"strings"
	"time"

	"internal/emoji"
	"internal/inlinesvg"

	. "github.com/NDoolan360/site-tools"
)

//go:embed "website"
var files embed.FS

func main() {
	// Verify environment variables
	for _, env := range []string{
		"ENV",
		"GITHUB_USERNAME",
		"GITHUB_TOKEN",
		"CULTS3D_USERNAME",
		"CULTS3D_API_KEY",
		"BGG_GEEKLIST",
	} {
		if _, ok := os.LookupEnv(env); !ok {
			log.Fatalf("Missing environment variable: %s", env)
			os.Exit(1)
		}
	}

	site := Build{}
	if err := site.FromDir(files, "website"); err != nil {
		log.Fatalf("Failed to build site from directory: %v", err)
	}

	site.Transform(CollectFrontMatter{})

	// Add flash-cards tool to the site from github
	path := "tools/flash-cards"
	if err := site.FromGit("https://github.com/NDoolan360/flash-cards", "main", path); err != nil {
		log.Fatalf("Failed to clone repository: %v", err)
	}
	flashCards := site.Filter(WithParentDir("/" + path))
	flashCards.AddToMeta("HideSocialLinks", "true")
	flashCards.AddToMeta("AltHeading", "flash-cards")

	// Reusable markdown transformer
	mdTransformer := MarkdownTransformer{
		Extensions:    Extensions{inlinesvg.InlineSvg, emoji.GoldMarkCustomEmojiExtension()},
		RenderOptions: RenderOptions{inlinesvg.WithParentPath("website")},
	}

	params := map[string]any{
		"PublishTime": time.Now(),
		"Projects":    getProjects(),
	}

	var components Assets
	components = site.Pop(WithMeta("IsComponent"))
	components.Filter(WithMeta("IsStatic")).Transform(mdTransformer)
	componentMap := components.ToMap("Name")

	// Pre-build articles as I want their meta for other pages
	articles := site.Filter(
		WithParentDir("/articles"),
		WithoutMeta("IsDraft"),
		WithMeta("Title"),
		WithMeta("Description"),
	).SetMetaFunc("URL", func(asset Asset) string {
		return strings.TrimSuffix(asset.Path, ".md") + ".html"
	})
	articles.Transform(TemplateTransformer{}, mdTransformer)
	params["Articles"] = articles

	// Pre-fill the emojis page
	site.Filter(WithPath("/other/emojis.md")).Transform(
		TemplateTransformer{
			GlobalData: map[string]any{"Emojis": emoji.GetEmojis()},
		},
	)

	// Process markdown files
	site.Filter(WithExtensions(".md"), WithoutMeta("IsDraft")).Transform(
		TemplateTransformer{GlobalData: params, Components: componentMap},
		mdTransformer,
	)

	// Add reload script to all html files when in development
	if os.Getenv("ENV") == "dev" {
		site.Filter(WithExtensions(".html")).Transform(
			&AddAutoReload{WebSocketPath: "/reload", Timeout: 2500},
		)
	}

	// Wrap all html files in base_template.html
	baseTemplate := site.Pop(WithPath("/templates/base_template.html"))[0]
	site.Filter(WithExtensions(".html")).Transform(
		TemplateTransformer{
			GlobalData: params,
			WrapperTemplate: &WrapperTemplate{
				Template:       baseTemplate,
				ChildBlockName: baseTemplate.Meta["ChildBlockName"].(string),
			},
			Components: componentMap,
		},
	)

	// Unescape escaped double curly braces
	site.Transform(&ReplacerTransformer{
		Replacements: map[string]string{"\\{\\{": "{{", "\\}\\}": "}}"},
	})

	// Minify
	site.Transform(&MinifyTransformer{})

	// Write to dir "build"
	os.RemoveAll("build")
	site.Filter(WithoutMeta("IsDraft")).Write("build")

	if os.Getenv("ENV") == "dev" {
		devServer("8888", "build")
	}
}
