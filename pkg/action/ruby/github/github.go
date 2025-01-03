package github

import (
	"fmt"
	"net/url"

	"github.com/wabenet/dodfile-syntax/pkg/simplerest"
)

const (
	APIRoot = "https://api.github.com"
)

//nolint:tagliatelle
type Release struct {
	ID      int64  `json:"id"`
	URL     string `json:"url"`
	Name    string `json:"name"`
	TagName string `json:"tag_name"`
}

func GetLatestRelease(org, repo string) (Release, error) {
	getUrl, err := url.Parse(APIRoot)
	if err != nil {
		return Release{}, fmt.Errorf("invalid API endpoint URL: %s: %w", APIRoot, err)
	}

	getUrl = getUrl.JoinPath("repos", org, repo, "releases", "latest")

	release, err := simplerest.Get[Release](getUrl.String())
	if err != nil {
		return Release{}, err
	}

	return release, nil
}

func GetReleases(org, repo string) ([]Release, error) {
	getUrl, err := url.Parse(APIRoot)
	if err != nil {
		return []Release{}, fmt.Errorf("invalid API endpoint URL: %s: %w", APIRoot, err)
	}

	getUrl = getUrl.JoinPath("repos", org, repo, "releases")

	releases, err := simplerest.Get[[]Release](getUrl.String())
	if err != nil {
		return []Release{}, err
	}

	return releases, nil
}
