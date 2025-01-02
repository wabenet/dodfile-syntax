package nodeapi

import (
	"errors"
	"fmt"
	"net/url"
	"runtime"
	"strings"

	"github.com/wabenet/dodfile-syntax/pkg/simplerest"
)

const (
	APIRoot = "https://nodejs.org/dist"
	Latest  = "latest"
)

var ErrNoRelease = errors.New("no valid release found")

type Release struct {
	Version string   `json:"version"`
	Files   []string `json:"files"`
}

type ReleaseFile struct {
	Filename       string
	URL            string
	SignedHashFile string
}

func GetDownload(version string) (ReleaseFile, error) {
	var result ReleaseFile

	release, err := GetReleaseForVersion(version)
	if err != nil {
		return result, err
	}

	getUrl, err := url.Parse(APIRoot)
	if err != nil {
		return result, fmt.Errorf("invalid API endpoint URL: %s: %w", APIRoot, err)
	}

	// TODO: You would expect that the `files` in the release information
	// somehow tell us which files there are to download. But oddly enough
	// there seems to be no relation whatsoever. We can just hope for now
	// that this pattern always works.
	result.Filename = fmt.Sprintf("node-%s-%s-%s.tar.gz", release.Version, runtime.GOOS, runtime.GOARCH)
	result.URL = getUrl.JoinPath(release.Version, result.Filename).String()
	result.SignedHashFile = getUrl.JoinPath(release.Version, "SHASUMS256.txt.asc").String()

	return result, nil
}

func GetReleaseForVersion(version string) (Release, error) {
	var result Release

	getUrl, err := url.Parse(APIRoot)
	if err != nil {
		return result, fmt.Errorf("invalid API endpoint URL: %s: %w", APIRoot, err)
	}

	getUrl = getUrl.JoinPath("index.json")

	releases, err := simplerest.Get[[]Release](getUrl.String())
	if err != nil {
		return result, err
	}

	if version == Latest {
		return releases[0], nil
	}

	for _, release := range releases {
		if release.Version == version {
			return release, nil
		}
	}

	return result, ErrNoRelease
}

func validFileForOS(name string) bool {
	return strings.Contains(name, runtime.GOOS) && strings.Contains(name, runtime.GOARCH)
}
