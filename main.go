package main

import (
	"embed"
	"internal/build"
	"internal/projects"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

//go:embed static
var staticFiles embed.FS

//go:embed pages
var pageFiles embed.FS

func main() {
	outDir := "assets"

	assets := build.NewBuild(outDir)
	if err := assets.StartWith(staticFiles); err != nil {
		log.Fatalf("failed to recreate assets dir: %v", err)
	}

	if err := assets.CollectFiles(pageFiles, "pages"); err != nil {
		log.Fatalf("failed to collect files: %v", err)
	}

	projects, err := projects.GetProjects([]string{"github", "bgg", "cults3d"})
	if err != nil {
		log.Fatalf("failed to create projects: %v", err)
	}

	if err := assets.ApplyTemplates(
		map[string]any{
			"Env":         os.Getenv("ENV"),
			"PublishTime": time.Now(),
			"Projects":    projects,
			"BlogPosts":   assets.Pages.GetBlogPosts().SortByDate(),
		},
	); err != nil {
		log.Fatalf("failed to apply template: %v", err)
	}

	if err := assets.WriteFiles(); err != nil {
		log.Fatalf("failed to write files: %v", err)
	}

	if os.Getenv("ENV") == "dev" {
		devServer("8888", outDir)
	}
}
