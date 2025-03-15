package projects

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Github struct {
	BaseURL     string
	User        string
	BearerToken string
}

var _ Host = (*Github)(nil)

type GithubData struct {
	Data struct {
		User struct {
			Repositories struct {
				Nodes []struct {
					Name              string `json:"name"`
					Description       string `json:"description"`
					URL               string `json:"url"`
					OpenGraphImageURL string `json:"openGraphImageUrl"`
					CreatedAt         string `json:"createdAt"`
					PrimaryLanguage   struct {
						Name  string `json:"name"`
						Color string `json:"color"`
					} `json:"primaryLanguage"`
					RepositoryTopics struct {
						Nodes []struct {
							Topic struct {
								Name string `json:"name"`
							} `json:"topic"`
						} `json:"nodes"`
					} `json:"repositoryTopics"`
				} `json:"nodes"`
			} `json:"repositories"`
		} `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (gh Github) Fetch() ([]byte, error) {
	client := &http.Client{}
	body := fmt.Sprintf(`{"query":"{user(login:\"%s\"){repositories(first:100,isFork:false,visibility:PUBLIC){nodes{name,description,url,openGraphImageUrl,createdAt,primaryLanguage{name,color},repositoryTopics(first:10){nodes{topic{name}}}}}}}"}`, gh.User)
	request, err := http.NewRequest(http.MethodPost, gh.BaseURL+"/graphql?gh", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+gh.BearerToken)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request to github failed with status code: %d", response.StatusCode)
	}

	if data, err := io.ReadAll(response.Body); err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func (Github) Parse(data []byte) (projects Projects, err error) {
	var githubProjects GithubData

	if unmarshalErr := json.Unmarshal(data, &githubProjects); unmarshalErr != nil {
		return nil, errors.Join(errors.New("error parsing GitHub projects"), unmarshalErr)
	}
	if len(githubProjects.Errors) > 0 {
		return nil, errors.New(githubProjects.Errors[0].Message)
	}

	for _, project := range githubProjects.Data.User.Repositories.Nodes {
		topics := []string{}
		// Add primary language to tags if it exists
		if project.PrimaryLanguage.Name != "" && project.PrimaryLanguage.Color != "" {
			topics = append(topics, fmt.Sprintf(`<p><i class="language-dot" style="background-color: %s"></i>%s</p>`, project.PrimaryLanguage.Color, project.PrimaryLanguage.Name))
		}
		for _, topic := range project.RepositoryTopics.Nodes {
			topics = append(topics, topic.Topic.Name)
		}

		createdAtTime, err := time.Parse(time.RFC3339, project.CreatedAt)
		if err != nil {
			log.Printf("error parsing Cults3D date value '%s'", project.CreatedAt)
			continue
		}

		projects = append(projects, Project{
			Title: EscapeSpecialChars(project.Name),
			URL:   project.URL,
			Image: Image{
				Src: project.OpenGraphImageURL,
				Alt: project.Name + " social preview",
			},
			Logo:        "static/images/logos/github.svg",
			Description: EscapeSpecialChars(project.Description),
			Created:     createdAtTime.Format(time.DateOnly),
			Tags:        topics,
		})
	}

	return projects, err
}
