package config

import (
	"github.com/moby/buildkit/client/llb"
)

const (
	baseImage = "debian"
)

type Image struct {
	Packages *Packages   `yaml:"packages"`
	Download []*Download `yaml:"download"`
	CopyFrom []*CopyFrom `yaml:"copy_from"`
}

func (img *Image) Build() llb.State {
	s := llb.Image(baseImage)

	if img.Packages != nil {
		s = Install(s, img.Packages)
	}

	for _, d := range img.Download {
		s = Fetch(s, d)
	}

	for _, c := range img.CopyFrom {
		s = CopyFromDockerfile(s, c)
	}

	return s
}
