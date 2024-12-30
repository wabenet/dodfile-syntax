package golangapi

import (
	"errors"
	"fmt"
	"net/url"
	"runtime"

	"github.com/wabenet/dodfile-syntax/pkg/simplerest"
)

const (
	APIRoot = "https://go.dev/dl"
	Latest  = "latest"
)

var ErrNoRelease = errors.New("no valid release found")

type Release struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []File `json:"files"`
}

type File struct {
	Filename       string `json:"filename"`
	OS             string `json:"os"`
	Arch           string `json:"arch"`
	Version        string `json:"version"`
	ChecksumSHA256 string `json:"sha256"`
	Size           int64  `json:"size"`
	Kind           string `json:"kind"` // "archive", "installer", "source"
}

func (f File) URL() (string, error) {
	getUrl, err := url.Parse(APIRoot)
	if err != nil {
		return "", fmt.Errorf("invalid API endpoint URL: %s: %w", APIRoot, err)
	}

	getUrl = getUrl.JoinPath(f.Filename)

	return getUrl.String(), nil
}

func GetDownload(version string) (File, error) {
	release, err := GetReleaseForVersion(version)
	if err != nil {
		return File{}, err
	}

	for _, f := range release.Files {
		if f.OS == "linux" && f.Arch == runtime.GOARCH && f.Kind == "archive" {
			return f, nil
		}
	}

	return File{}, ErrNoRelease
}

func GetReleaseForVersion(version string) (Release, error) {
	var result Release

	getUrl, err := url.Parse(APIRoot)
	if err != nil {
		return result, fmt.Errorf("invalid API endpoint URL: %s: %w", APIRoot, err)
	}

	getUrl = getUrl.JoinPath("/")

	query := getUrl.Query()
	query.Set("mode", "json")

	if version != Latest {
		query.Set("include", "all")
	}

	getUrl.RawQuery = query.Encode()

	releases, err := simplerest.Get[[]Release](getUrl.String())
	if err != nil {
		return result, err
	}

	for _, r := range releases {
		if version == Latest && r.Stable {
			return r, nil
		}

		if r.Version == "go"+version {
			return r, nil
		}
	}

	return result, ErrNoRelease
}
