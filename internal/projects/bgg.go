package projects

import (
	_ "embed"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Bgg struct {
	BaseURL  string
	Geeklist string
}

var _ Host = (*Bgg)(nil)

type BggProject struct {
	Item struct {
		Id string `xml:"objectid,attr"`
	} `xml:"item"`
}

type BggItem struct {
	Title         string   `xml:"boardgame>name"`
	Description   string   `xml:"boardgame>description"`
	YearPublished string   `xml:"boardgame>yearpublished"`
	ImageSrc      string   `xml:"boardgame>image"`
	Tags          []string `xml:"boardgame>boardgamemechanic"`
}

func (bgg Bgg) Fetch() ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("%s/geeklist/%s", bgg.BaseURL, bgg.Geeklist))
	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

func (bgg Bgg) Parse(data []byte) (projects Projects, err error) {
	var projectItems []BggProject
	if unmarshalErr := xml.Unmarshal(data, &projectItems); unmarshalErr != nil {
		return nil, errors.Join(errors.New("error parsing BGG projects"), unmarshalErr)
	}

	for _, item := range projectItems {
		resp, err := http.Get(fmt.Sprintf("%s/boardgame/%s", bgg.BaseURL, item.Item.Id))
		if err != nil {
			return nil, err
		}

		projectData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var bggProject BggItem
		if unmarshalErr := xml.Unmarshal(projectData, &bggProject); unmarshalErr != nil {
			log.Print("error parsing BGG project: ", item.Item.Id, ": ", unmarshalErr)
			continue
		}

		projects = append(projects, Project{
			Title: EscapeSpecialChars(bggProject.Title),
			URL:   fmt.Sprintf("https://boardgamegeek.com/boardgame/%s", item.Item.Id),
			Image: Image{
				Src: bggProject.ImageSrc,
				Alt: fmt.Sprintf("Board Game: %s", bggProject.Title),
			},
			Description: EscapeSpecialChars(bggProject.Description),
			Logo:        "static/images/logos/bgg.svg",
			Created:     bggProject.YearPublished,
			Tags:        bggProject.Tags,
		})
	}

	return projects, err
}
