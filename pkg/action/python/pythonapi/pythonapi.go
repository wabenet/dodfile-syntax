package pythonapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	APIRoot = "https://www.python.org/api/v2"
	Major   = "3" // We only support Python 3
	Latest  = "latest"
)

type OS int

const (
	_ = iota
	Windows
	MacOS
	Source
)

type Release struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Version     int    `json:"version"`
	IsPublished bool   `json:"is_published"`
	IsLatest    bool   `json:"is_latest"`
	ReleaseDate string `json:"release_date"`
	PreRelease  bool   `json:"pre_release"`
	ResourceUri string `json:"resource_uri"`
}

type ReleaseFile struct {
	Name             string `json:"name"`
	Slug             string `json:"slug"`
	OS               string `json:"os"`
	URL              string `json:"url"`
	GPGSignatureFile string `json:"gpg_signature_file"`
	MD5Sum           string `json:"md5_sum"`
	ResourceUri      string `json:"resource_uri"`
}

func GetDownload(version string, os OS) (ReleaseFile, error) {
	release, err := GetReleaseForVersion(version)
	if err != nil {
		return ReleaseFile{}, err
	}

	u, err := url.Parse(release.ResourceUri)
	if err != nil {
		return ReleaseFile{}, err
	}

	releaseID := path.Base(u.Path)

	u, err = url.Parse(APIRoot)
	if err != nil {
		return ReleaseFile{}, err
	}

	u = u.JoinPath("downloads", "release_file")

	q := u.Query()
	q.Set("release", releaseID)
	q.Set("os", strconv.Itoa(int(os)))
	u.RawQuery = q.Encode()

	releaseFiles, err := get[[]ReleaseFile](u.String())
	if err != nil {
		return ReleaseFile{}, err
	}

	if len(releaseFiles) < 1 {
		return ReleaseFile{}, errors.New("no releases")
	}

	return releaseFiles[0], nil
}

func GetReleaseForVersion(version string) (Release, error) {
	u, err := url.Parse(APIRoot)
	if err != nil {
		return Release{}, err
	}

	u = u.JoinPath("downloads", "release")

	q := u.Query()
	q.Set("version", Major)
	if version != Latest {
		q.Set("name", fmt.Sprintf("Python %s", version))
	}
	u.RawQuery = q.Encode()

	releases, err := get[[]Release](u.String())
	if err != nil {
		return Release{}, err
	}

	if version == Latest {
		for _, r := range releases {
			if r.IsLatest {
				return r, nil
			}
		}
		return Release{}, errors.New("no latest found")
	}

	if len(releases) != 1 {
		return Release{}, errors.New("too many releases")
	}

	return releases[0], nil
}

func get[T any](url string) (T, error) {
	var result T

	resp, err := http.Get(url)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, err
	}

	return result, nil
}
