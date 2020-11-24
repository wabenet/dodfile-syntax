package config

import (
	"strings"

	"github.com/moby/buildkit/client/llb"
)

type Packages struct {
	Name []string `yaml:"name"`
	Repo []string `yaml:"repo"`
	Gpg  []string `yaml:"gpg"`
}

func Install(base llb.State, p *Packages) llb.State {
	base = Sh(base, "apt-get update && apt-get install apt-transport-https -y")

	for _, repo := range p.Repo {
		base = Sh(base, "echo \"%s\" >> /etc/apt/sources.list", repo)
	}

	for _, url := range p.Gpg {
		curl := Install(llb.Image(defaultBaseImage), &Packages{Name: []string{"curl"}})
		downloadSt := Sh(curl, "curl -Lo /key.gpg %s", url)
		base = Copy(downloadSt, "/key.gpg", base, "/key.gpg")
		base = Sh(base, "apt-key add /key.gpg && rm /key.gpg")
	}

	if len(p.Name) > 0 {
		packages := strings.Join(p.Name, " ")
		base = Sh(base, "apt-get update && apt-get install --no-install-recommends --no-install-suggests -y %s", packages)
	}

	return base
}
