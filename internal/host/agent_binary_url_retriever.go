package host

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type AgentBinaryURLRetriever interface {
	RetrieveUrl() (string, error)
}

type DevContainerAgentBinaryURLRetriever struct{}

func NewDevContainerAgentBinaryURLRetriever() *DevContainerAgentBinaryURLRetriever {
	return &DevContainerAgentBinaryURLRetriever{}
}

func (r *DevContainerAgentBinaryURLRetriever) RetrieveUrl() (string, error) {
	return "http://minio:9000/bin/mear-agent", nil
}

var GitHubRepoOwner = "KTachibanaM"
var GitHubRepoName = "mear"
var GitHubReleaseApiUrl = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", GitHubRepoOwner, GitHubRepoName)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

type GitHubAgentBinaryURLRetriever struct{}

func NewGithubAgentBinaryURLRetriever() *GitHubAgentBinaryURLRetriever {
	return &GitHubAgentBinaryURLRetriever{}
}

func (r *GitHubAgentBinaryURLRetriever) RetrieveUrl() (string, error) {
	req, err := http.NewRequest("GET", GitHubReleaseApiUrl, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request for github release API: %v", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to request github release API: %v", err)
	}
	defer resp.Body.Close()

	var release GitHubRelease
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return "", fmt.Errorf("failed to parse github release API response: %v", err)
	}

	version := strings.TrimPrefix(release.TagName, "v")

	return fmt.Sprintf("https://github.com/KTachibanaM/mear/releases/download/v%v/mear-agent_%v_linux_amd64", version, version), nil
}
