package agent_bin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type AgentBinary interface {
	RetrieveUrl() (string, error)
}

type DevContainerAgentBinary struct{}

func NewDevContainerAgentBinary() *DevContainerAgentBinary {
	return &DevContainerAgentBinary{}
}

func (r *DevContainerAgentBinary) RetrieveUrl() (string, error) {
	return "http://minio:9000/mear-bin/mear-agent", nil
}

var GitHubRepoOwner = "KTachibanaM"
var GitHubRepoName = "mear"
var GitHubReleaseApiUrl = fmt.Sprintf("https://api.github.com/repos/%v/%v/releases/latest", GitHubRepoOwner, GitHubRepoName)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

type GitHubAgentBinary struct{}

func NewGithubAgentBinary() *GitHubAgentBinary {
	return &GitHubAgentBinary{}
}

func (r *GitHubAgentBinary) RetrieveUrl() (string, error) {
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
