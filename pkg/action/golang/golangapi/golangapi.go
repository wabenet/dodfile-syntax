package golangapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
)

const (
	APIRoot = "https://go.dev/dl"
	Latest  = "latest"
)

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
	u, err := url.Parse(APIRoot)
	if err != nil {
		return "", err
	}

	u = u.JoinPath(f.Filename)

	return u.String(), nil
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

	return File{}, errors.New("no matching release file found")
}

func GetReleaseForVersion(version string) (Release, error) {
	u, err := url.Parse(APIRoot)
	if err != nil {
		return Release{}, err
	}

	u = u.JoinPath("/")

	q := u.Query()
	q.Set("mode", "json")
	if version != Latest {
		q.Set("include", "all")
	}
	u.RawQuery = q.Encode()

	releases, err := get[[]Release](u.String())
	if err != nil {
		return Release{}, err
	}

	for _, r := range releases {
		if version == Latest && r.Stable {
			return r, nil
		}

		if r.Version == fmt.Sprintf("go%s", version) {
			return r, nil
		}
	}

	return Release{}, errors.New("no matching version found")
}

func get[T any](url string) (T, error) {
	var result T

	resp, err := http.Get(url)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("HTTP error %d", resp.StatusCode)
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
