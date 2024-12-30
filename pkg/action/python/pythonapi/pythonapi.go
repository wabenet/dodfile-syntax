package pythonapi

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"strconv"

	"github.com/wabenet/dodfile-syntax/pkg/simplerest"
)

const (
	APIRoot = "https://www.python.org/api/v2"
	Major   = "3" // We only support Python 3
	Latest  = "latest"
)

var (
	ErrNoRelease       = errors.New("no valid release found")
	ErrTooManyReleases = errors.New("too many valid releases found")
)

type OS int

const (
	_ = iota
	Windows
	MacOS
	Source
)

//nolint:tagliatelle
type Release struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Version     int    `json:"version"`
	IsPublished bool   `json:"is_published"`
	IsLatest    bool   `json:"is_latest"`
	ReleaseDate string `json:"release_date"`
	PreRelease  bool   `json:"pre_release"`
	ResourceURI string `json:"resource_uri"`
}

//nolint:tagliatelle
type ReleaseFile struct {
	Name             string `json:"name"`
	Slug             string `json:"slug"`
	OS               string `json:"os"`
	URL              string `json:"url"`
	GPGSignatureFile string `json:"gpg_signature_file"`
	MD5Sum           string `json:"md5_sum"`
	ResourceURI      string `json:"resource_uri"`
}

func GetDownload(version string, os OS) (ReleaseFile, error) {
	var result ReleaseFile

	release, err := GetReleaseForVersion(version)
	if err != nil {
		return result, err
	}

	getUrl, err := url.Parse(release.ResourceURI)
	if err != nil {
		return result, fmt.Errorf("invalid download URL for release: %s: %w", release.ResourceURI, err)
	}

	releaseID := path.Base(getUrl.Path)

	getUrl, err = url.Parse(APIRoot)
	if err != nil {
		return result, fmt.Errorf("invalid API endpoint URL: %s: %w", APIRoot, err)
	}

	getUrl = getUrl.JoinPath("downloads", "release_file")

	query := getUrl.Query()
	query.Set("release", releaseID)
	query.Set("os", strconv.Itoa(int(os)))
	getUrl.RawQuery = query.Encode()

	releaseFiles, err := simplerest.Get[[]ReleaseFile](getUrl.String())
	if err != nil {
		return result, err
	}

	if len(releaseFiles) < 1 {
		return result, ErrNoRelease
	}

	result = releaseFiles[0]

	return result, nil
}

func GetReleaseForVersion(version string) (Release, error) {
	var result Release

	getUrl, err := url.Parse(APIRoot)
	if err != nil {
		return result, fmt.Errorf("invalid API endpoint URL: %s: %w", APIRoot, err)
	}

	getUrl = getUrl.JoinPath("downloads", "release")

	query := getUrl.Query()
	query.Set("version", Major)

	if version != Latest {
		query.Set("name", "Python "+version)
	}

	getUrl.RawQuery = query.Encode()

	releases, err := simplerest.Get[[]Release](getUrl.String())
	if err != nil {
		return result, err
	}

	if version == Latest {
		for _, r := range releases {
			if r.IsLatest {
				return r, nil
			}
		}

		return result, ErrNoRelease
	}

	if len(releases) != 1 {
		return result, ErrTooManyReleases
	}

	result = releases[0]

	return result, nil
}
