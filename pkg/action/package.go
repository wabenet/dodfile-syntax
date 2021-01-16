package action

import (
	"strings"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
)

type PackageAction struct {
	Name []string
	Repo []string
	Gpg  []string
}

func (a *PackageAction) Execute(base llb.State) llb.State {
	base = Sh(base, "apt-get update && apt-get install apt-transport-https -y")

	for _, repo := range a.Repo {
		base = Sh(base, "echo \"%s\" >> /etc/apt/sources.list", repo)
	}

	for _, url := range a.Gpg {
		curl := Sh(llb.Image(defaultBaseImage), "apt-get update && apt-get install -y --no-install-recommends --no-install -suggests apt-transport-https curl ca-certificates")
		downloadSt := Sh(curl, "curl -Lo /key.gpg %s", url)
		base = Copy(downloadSt, "/key.gpg", base, "/key.gpg")
		base = Sh(base, "apt-key add /key.gpg && rm /key.gpg")
	}

	if len(a.Name) > 0 {
		packages := strings.Join(a.Name, " ")
		base = Sh(base, "apt-get update && apt-get install --no-install-recommends --no-install-suggests -y %s", packages)
	}

	return base
}

func (*PackageAction) UpdateMetadata(_ *dockerfile2llb.Image) {
}
