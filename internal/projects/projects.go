package projects

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"sync"
)

//go:embed testdata/expected_all.json
var AllProjectsMock []byte

type Project struct {
	Title       string   // Required
	URL         string   // Required
	Image       Image    // Required
	Logo        string   // Required
	Description string   // Optional
	CreatedDate string   // Optional
	Tags        []string // Optional
}

type Projects []Project

type Image struct {
	Src string
	Alt string
}

type Language struct {
	Name   string
	Colour string
}

type Host interface {
	Fetch() ([]byte, error)
	Parse([]byte) (Projects, error)
}

var hostInterfaces = map[string]Host{
	"github":  GithubHost{BaseURL: "https://api.github.com", User: os.Getenv("GITHUB_USERNAME")},
	"bgg":     BggHost{BaseURL: "https://boardgamegeek.com/xmlapi", Geeklist: os.Getenv("BGG_GEEKLIST")},
	"cults3d": Cults3dHost{BaseURL: "https://cults3d.com", User: os.Getenv("CULTS3D_USERNAME")},
}

func GetProjects(hostNames []string) (Projects, error) {
	var wg sync.WaitGroup
	var projects Projects

	for _, hostName := range hostNames {
		host, ok := hostInterfaces[hostName]
		if !ok {
			return nil, fmt.Errorf("Interface for host '%s' not found.", host)
		}

		wg.Add(1)
		go func(hostName string, host Host, projects *Projects, wg *sync.WaitGroup) {
			defer wg.Done()

			data, err := host.Fetch()
			if err != nil {
				log.Print(hostName, ": ", err)
				return
			}

			moreProjects, err := host.Parse(data)
			if err != nil {
				log.Print(hostName, ": ", err)
				return
			}

			*projects = append(*projects, moreProjects...)
		}(hostName, host, &projects, &wg)
	}
	wg.Wait()

	// filter out any projects with missing titles, URLs, logos, or images
	filteredProjects := make(Projects, 0, len(projects))
	for _, project := range projects {
		if project.Title != "" && project.URL != "" && project.Logo != "" && project.Image.Src != "" {
			filteredProjects = append(filteredProjects, project)
		}
	}

	// sort projects by creation date
	slices.SortStableFunc(filteredProjects, func(project1, project2 Project) int {
		return -strings.Compare(project1.CreatedDate, project2.CreatedDate)
	})

	return filteredProjects, nil
}
