package projects

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Cults3dHost struct {
	BaseURL string
	User    string
}

var _ Host = (*Cults3dHost)(nil)

type Cults3dData struct {
	Data struct {
		User struct {
			Creations []struct {
				Title       string   `json:"name"`
				Description string   `json:"description"`
				Url         string   `json:"url"`
				PublishedAt string   `json:"publishedAt"`
				Downloads   int      `json:"downloadsCount"`
				ImageSrc    string   `json:"illustrationImageUrl"`
				Topics      []string `json:"tags"`
			} `json:"creations"`
		} `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (cults Cults3dHost) Fetch() ([]byte, error) {
	client := &http.Client{}
	body := fmt.Sprintf(`{"query":"{user(nick:\"%s\"){creations(limit:100,sort:BY_DOWNLOADS){name,url,description,publishedAt,downloadsCount,illustrationImageUrl,tags}}}"}`, cults.User)
	request, err := http.NewRequest(http.MethodPost, cults.BaseURL+"/graphql?cults", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.SetBasicAuth(os.Getenv("CULTS3D_USERNAME"), os.Getenv("CULTS3D_API_KEY"))

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request to cults3d failed with status code: %d", response.StatusCode)
	}

	if data, err := io.ReadAll(response.Body); err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func (Cults3dHost) Parse(data []byte) (projects Projects, err error) {
	var cults3dData Cults3dData
	if unmarshalErr := json.Unmarshal(data, &cults3dData); unmarshalErr != nil {
		return nil, errors.Join(errors.New("error parsing Cults3D projects"), unmarshalErr)
	}
	if len(cults3dData.Errors) > 0 {
		return nil, errors.New(cults3dData.Errors[0].Message)
	}

	for _, project := range cults3dData.Data.User.Creations {
		// prepend downloads count to topics
		project.Topics = append([]string{fmt.Sprintf("Downloads: %d", project.Downloads)}, project.Topics...)

		project.Description = strings.ReplaceAll(project.Description, "\r\n", " ")

		projects = append(projects, Project{
			Title: project.Title,
			URL:   project.Url,
			Image: Image{
				Src: project.ImageSrc,
				Alt: fmt.Sprintf("3D Model: %s", project.Title),
			},
			Description: project.Description,
			CreatedDate: project.PublishedAt,
			Logo:        "static/images/logos/cults3d.svg",
			Tags:        project.Topics,
		})
	}

	return projects, err
}
