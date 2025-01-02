package yarnapi

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/wabenet/dodfile-syntax/pkg/simplerest"
)

const (
	APIRoot = "https://registry.yarnpkg.com"
	Latest  = "latest"
)

var ErrNoRelease = errors.New("no valid release found")

type Package struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Versions    map[string]Version `json:"versions"`
	DistTags    map[string]string  `json:"dist-tags"`
}

type Version struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Dist    Dist   `json:"dist"`
}

type Dist struct {
	SHASum  string `json:"shasum"`
	Tarball string `json:"tarball"`
}

type File struct {
	URL    string
	SHASum string
}

func GetDownload(version string) (File, error) {
	var result File

	pkg, err := GetPackage("yarn")
	if err != nil {
		return result, err
	}

	ver, err := pkg.Version(version)
	if err != nil {
		return result, err
	}

	return File{
		URL:    ver.Dist.Tarball,
		SHASum: ver.Dist.SHASum,
	}, nil
}

func (p Package) Version(version string) (Version, error) {
	if version == Latest {
		version = p.DistTags["latest"]
	}
	if v, ok := p.Versions[version]; ok {
		return v, nil
	}

	return Version{}, ErrNoRelease
}

func GetPackage(name string) (Package, error) {
	var result Package

	getUrl, err := url.Parse(APIRoot)
	if err != nil {
		return result, fmt.Errorf("invalid API endpoint URL: %s: %w", APIRoot, err)
	}

	getUrl = getUrl.JoinPath(name)

	result, err = simplerest.Get[Package](getUrl.String())
	if err != nil {
		return result, err
	}

	return result, nil
}
