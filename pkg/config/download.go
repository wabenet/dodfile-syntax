package config

import (
	"path"

	"github.com/moby/buildkit/client/llb"
)

type Download struct {
	Source      string `yaml:"source"`
	Sha256      string `yaml:"sha256"`
	Unpack      string `yaml:"unpack"`
	Destination string `yaml:"destination"`
}

func Fetch(base llb.State, d *Download) llb.State {
	targetFile := path.Base(d.Source)
	s := Install(llb.Image(defaultBaseImage), &Packages{Name: []string{"curl", "ca-certificates"}})

	s = Sh(s, "curl -Lo %s %s", targetFile, d.Source)

	if d.Sha256 != "" {
		s = Sh(s, "echo \"%s  %s\" | sha256sum -c -", d.Sha256, targetFile)
	}

	if d.Unpack != "" {
		s = Sh(s, "tar -zxf %s -C %s", targetFile, path.Dir(d.Destination))
                targetFile = path.Join(path.Dir(d.Destination), d.Unpack)
	} else {
		s = Sh(s, "chmod +x %s", targetFile)
	}

	return Copy(s, targetFile, base, d.Destination)
}
