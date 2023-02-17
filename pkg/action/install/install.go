package install

import (
	"strings"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/wabenet/dodfile-syntax/pkg/state"
)

const (
	Type = "install"

	defaultBaseImage = "debian"
)

type Action struct {
	Name string `mapstructure:"name"`
	Repo string `mapstructure:"repo"`
	Gpg  string `mapstructure:"gpg"`
}

func (a *Action) Type() string {
	return Type
}

func (a *Action) Execute(base llb.State) (llb.State, error) {
	s := state.FromLLB(defaultBaseImage, base)

	if len(a.Gpg) > 0 {
		s.Install("gnupg")
	}

	if len(a.Repo) > 0 {
		s.Sh("echo \"%s\" >> /etc/apt/sources.list", a.Repo)
	}

	if len(a.Gpg) > 0 {
		curl := state.From(defaultBaseImage)
		curl.Install("apt-transport-https", "curl", "ca-certificates")
		curl.Exec("/usr/bin/curl", "-Lo", "/key.gpg", a.Gpg)
		s.Copy(curl, "/key.gpg", "/key.gpg")
		s.Sh("apt-key add /key.gpg && rm /key.gpg")
	}

	if len(a.Name) > 0 {
		s.Install(strings.Fields(a.Name)...)
	}

	return s.Get(), nil
}

func (a *Action) UpdateImage(_ *dockerfile2llb.Image) {}
