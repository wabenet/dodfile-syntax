package config

import (
	"strings"

	"github.com/moby/buildkit/client/llb"
)

type Download struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
	Sha256      string `yaml:"sha256"`
}

func Fetch(base llb.State, d *Download) llb.State {
	s := Install(llb.Image(defaultBaseImage), &Packages{Name: []string{"curl", "ca-certificates"}})

	if isTgz(d.Source) {
		s = Sh(s, "curl -Lo tmp.tar.gz %s", d.Source)

		if d.Sha256 != "" {
			s = Sh(s, "echo \"%s  tmp.tar.gz\" | sha256sum -c -", d.Sha256)
		}

		s = Sh(s, "mkdir -p %[1]s && tar -zxvf tmp.tar.gz -C %[1]s && rm tmp.tar.gz", d.Destination)
	} else {
		s = Sh(s, "curl -Lo %s %s && chmod +x %s", d.Destination, d.Source, d.Destination)

		if d.Sha256 != "" {
			s = Sh(s, "echo \"%s  %s\" | sha256sum -c -", d.Sha256, d.Destination)
		}
	}

	return Copy(s, d.Destination, base, d.Destination)
}

func isTgz(file string) bool {
	return strings.HasSuffix(file, ".tar.gz")
}
