package main

import (
	"embed"
	"os"
	"strings"
	"time"

	"internal/build"
	"internal/emoji"
)

//go:embed "website"
var files embed.FS

func main() {
	site := build.Build{}
	site.WalkDir(files, "website", false)
	site.Transform(build.CollectFrontMatter{})

	// Add flash-cards tool to the site from git repository
	if err := site.FromGit(
		"https://github.com/NDoolan360/flash-cards",
		"main",
		"tools/flash-cards",
	); err != nil {
		log.Fatalf("Failed to clone repository: %v", err)
		os.Exit(1)
	}
	site.Filter(withParentDir("/tools/flash-cards")).AddToMeta("HideSocialLinks", "true").AddToMeta("IsDraft", "true")

	params := map[string]any{
		"PublishTime": time.Now(),
		"Projects":    getProjects().SortByCreatedDate(),
	}

	var components build.Assets
	components = site.Pop(withMeta("IsComponent"))
	components.Filter(withMeta("IsStatic")).Transform(build.MarkdownTransformer{})
	componentMap := components.ToMap("Name")

	// Pre-build articles as I want their meta for other pages
	articles := site.Filter(
		withParentDir("/articles"),
		withoutMeta("IsDraft"),
		withMeta("Title"),
		withMeta("Description"),
	).SetMetaFunc("URL", func(asset build.Asset) string {
		return strings.TrimSuffix(asset.Path, ".md") + ".html"
	})
	articles.Transform(build.TemplateTransformer{}, build.MarkdownTransformer{})
	params["Articles"] = articles

	// Pre-fill the emojis page
	site.Filter(withPath("/other/emojis.md")).Transform(
		build.TemplateTransformer{
			GlobalData: map[string]any{"Emojis": emoji.GetEmojis()},
		},
	)

	// Process markdown files
	site.Filter(withExtensions(".md"), withoutMeta("IsDraft")).Transform(
		build.TemplateTransformer{GlobalData: params, Components: componentMap},
		build.MarkdownTransformer{},
	)

	// Add reload script to all html files when in development
	if os.Getenv("ENV") == "dev" {
		site.Filter(withExtensions(".html")).Transform(
			&build.AddAutoReload{WebSocketPath: "/reload", Timeout: 2500},
		)
	}

	// Wrap all html files in base_template.html
	baseTemplate := site.Pop(withPath("/templates/base_template.html"))[0]
	site.Filter(withExtensions(".html")).Transform(
		build.TemplateTransformer{
			GlobalData: params,
			WrapperTemplate: &build.WrapperTemplate{
				Template:       baseTemplate,
				ChildBlockName: baseTemplate.Meta["ChildBlockName"].(string),
			},
			Components: componentMap,
		},
	)

	// Unescape escaped double curly braces
	site.Transform(&build.ReplacerTransformer{
		Replacements: map[string]string{"\\{\\{": "{{", "\\}\\}": "}}"},
	})

	// Minify
	site.Transform(&build.MinifyTransformer{})

	// Write to dir "build"
	os.RemoveAll("build")
	site.Filter(withoutMeta("IsDraft")).Write("build")

	if os.Getenv("ENV") == "dev" {
		devServer("8888", "build")
	}
}
