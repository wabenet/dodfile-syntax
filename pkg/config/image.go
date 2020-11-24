package config

import (
	"github.com/moby/buildkit/client/llb"
)

const (
	defaultBaseImage = "debian"
)

type Image struct {
	BaseImage string      `yaml:"base_image"`
	User      string      `yaml:"user"`
	Paths     []string    `yaml:"paths"`
	Packages  *Packages   `yaml:"packages"`
	Download  []*Download `yaml:"download"`
	From      []*CopyFrom `yaml:"from"`
}

func (img *Image) base() llb.State {
	baseImage := img.BaseImage
	if len(baseImage) == 0 {
		baseImage = defaultBaseImage
	}

	return llb.Image(baseImage)
}

func (img *Image) Build() llb.State {
	s := img.base()

	if img.Packages != nil {
		s = Install(s, img.Packages)
	}

	for _, d := range img.Download {
		s = Fetch(s, d)
	}

	for _, c := range img.From {
		s = CopyFromDockerfile(s, c)
	}

	return s
}
