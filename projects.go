package main

import (
	"internal/projects"
	"log"
	"os"
)

func getProjects() projects.Projects {
	projects, err := projects.GetProjects(map[string]projects.Host{
		"github": projects.Github{
			BaseURL:     "https://api.github.com",
			User:        os.Getenv("GITHUB_USERNAME"),
			BearerToken: os.Getenv("GITHUB_TOKEN"),
		},
		"bgg": projects.Bgg{
			BaseURL:  "https://boardgamegeek.com/xmlapi",
			Geeklist: os.Getenv("BGG_GEEKLIST"),
		},
		"cults3d": projects.Cults3d{
			BaseURL: "https://cults3d.com",
			User:    os.Getenv("CULTS3D_USERNAME"),
			APIKey:  os.Getenv("CULTS3D_API_KEY"),
		},
	})
	if err != nil {
		log.Fatalf("failed to create projects: %v", err)
	}

	return projects
}
