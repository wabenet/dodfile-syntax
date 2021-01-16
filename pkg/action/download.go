package action

import (
	"path"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
)

type DownloadAction struct {
	Source      string
	Sha256      string
	Unpack      string
	Destination string
}

func (a *DownloadAction) Execute(base llb.State) llb.State {
	targetFile := path.Base(a.Source)
	s := Sh(llb.Image(defaultBaseImage), "apt-get update && apt-get install -y --no-install-recommends --no-install-suggests apt-transport-https curl ca-certificates")

	s = Sh(s, "curl -Lo %s %s", targetFile, a.Source)

	if a.Sha256 != "" {
		s = Sh(s, "echo \"%s  %s\" | sha256sum -c -", a.Sha256, targetFile)
	}

	if a.Unpack != "" {
		s = Sh(s, "tar -zxf %s -C %s", targetFile, path.Dir(a.Destination))
		targetFile = path.Join(path.Dir(a.Destination), a.Unpack)
	} else {
		s = Sh(s, "chmod +x %s", targetFile)
	}

	return Copy(s, targetFile, base, a.Destination)
}

func (*DownloadAction) UpdateMetadata(_ *dockerfile2llb.Image) {
}
