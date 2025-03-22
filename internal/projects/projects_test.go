package projects

import (
	_ "embed"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//go:embed testdata/expected_github.json
var ExpectedGithubFetchResponse []byte

//go:embed testdata/expected_cults3d.json
var ExpectedCults3dFetchResponse []byte

//go:embed testdata/expected_bgg.json
var ExpectedBggFetchResponse []byte

//go:embed testdata/expected_all.json
var ExpectedAllFetchResponse []byte

func TestFetchAndParse(t *testing.T) {
	server := mockExternalServer(t)

	mockHosts := map[string]Host{
		"github":  Github{BaseURL: server.URL, User: "NDoolan360", BearerToken: "TEST"},
		"cults3d": Cults3d{BaseURL: server.URL, User: "TEST", APIKey: "TEST"},
		"bgg":     Bgg{BaseURL: server.URL, Geeklist: "332832"},
	}

	tests := []struct {
		name  string
		hosts map[string]Host
		want  []byte
	}{
		{"Github Host Test", map[string]Host{"github": mockHosts["github"]}, ExpectedGithubFetchResponse},
		{"Cults3D Host Test", map[string]Host{"cults3d": mockHosts["cults3d"]}, ExpectedCults3dFetchResponse},
		{"BGG Host Test", map[string]Host{"bgg": mockHosts["bgg"]}, ExpectedBggFetchResponse},
		{"All Hosts Test", mockHosts, ExpectedAllFetchResponse},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			projects, err := GetProjects(tc.hosts)
			if err != nil {
				t.Fatalf("Got error: %v", err)
			}

			json, err := json.MarshalIndent(projects, "", "  ")
			if err != nil {
				t.Fatalf("Got error: %v", err)
			}
			json = append(json, '\n')
			if diff := cmp.Diff(json, tc.want); diff != "" {
				t.Errorf("(+want -got):\n%s", diff)
			}
		})
	}

	server.Close()
}

//go:embed testdata/mock_github.json
var MockGithubHttpResponse []byte

//go:embed testdata/mock_cults3d.json
var MockCults3dHttpResponse []byte

//go:embed testdata/mock_bgg_geeklist.xml
var MockBggHttpResponse []byte

//go:embed testdata/mock_bgg_boardgame.xml
var MockBggXmlHttpResponse []byte

func mockExternalServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		switch r.URL.String() {
		case "/graphql?gh":
			_, err = w.Write(MockGithubHttpResponse)
		case "/graphql?cults":
			_, err = w.Write(MockCults3dHttpResponse)
		case "/geeklist/332832":
			_, err = w.Write(MockBggHttpResponse)
		case "/boardgame/330653":
			_, err = w.Write(MockBggXmlHttpResponse)
		default:
			err = errors.New("mock url not defined")
		}
		if err != nil {
			t.Error(err)
		}
	}))
}
