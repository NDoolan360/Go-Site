package projects

import (
	_ "embed"
	"log"
	"slices"
	"strings"
	"sync"
)

//go:embed testdata/expected_all.json
var AllProjectsMock []byte

type Project struct {
	Title       string   // Required
	URL         string   // Required
	Description string   // Required
	Image       Image    // Optional
	Logo        string   // Optional
	Created     string   // Optional
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

func GetProjects(hosts map[string]Host) (Projects, error) {
	var wg sync.WaitGroup
	var projects Projects

	wg.Add(len(hosts))
	for hostName, host := range hosts {
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

	// filter out any projects with missing titles, URLs, or descriptions
	filteredProjects := make(Projects, 0, len(projects))
	for _, project := range projects {
		if project.Title != "" &&
			project.URL != "" &&
			project.Description != "" {
			filteredProjects = append(filteredProjects, project)
		}
	}

	return filteredProjects, nil
}

func (projects Projects) SortByCreatedDate() Projects {
	return slices.SortedFunc(
		slices.Values(projects),
		func(project1, project2 Project) int {
			return strings.Compare(project2.Created, project1.Created)
		},
	)
}
