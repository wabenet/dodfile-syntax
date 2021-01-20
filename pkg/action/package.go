package action

import (
	"github.com/dodo-cli/dodfile-syntax/pkg/state"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
)

type PackageAction struct {
	Name []string
	Repo []string
	Gpg  []string
}

func (a *PackageAction) Execute(base llb.State) llb.State {
	s := state.FromLLB(defaultBaseImage, base)

	for _, repo := range a.Repo {
		s.Sh("echo \"%s\" >> /etc/apt/sources.list", repo)
	}

	for _, url := range a.Gpg {
		curl := state.From(defaultBaseImage)
		curl.Install("apt-transport-https", "curl", "ca-certificates")
		curl.Exec("/usr/bin/curl", "-Lo", "/key.gpg", url)
		s.Copy(curl, "/key.gpg", "/key.gpg")
		s.Sh("apt-key add /key.gpg && rm /key.gpg")
	}

	if len(a.Name) > 0 {
		s.Install(a.Name...)
	}

	return s.Get()
}

func (*PackageAction) UpdateMetadata(_ *dockerfile2llb.Image) {
}
