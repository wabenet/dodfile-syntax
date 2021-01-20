package action

import (
	"path"

	"github.com/dodo-cli/dodfile-syntax/pkg/state"
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

	downloader := state.From(defaultBaseImage)
	downloader.Install("apt-transport-https", "curl", "ca-certificates")

	if a.Unpack != "" {
		downloader.Install("tar", "unzip")
	}

	downloader.Exec("/usr/bin/curl", "-Lo", targetFile, a.Source)

	if a.Sha256 != "" {
		downloader.Sh("echo \"%s  %s\" | sha256sum -c -", a.Sha256, targetFile)
	}

	if a.Unpack != "" {
		downloader.Exec("/bin/tar", "-zxf", targetFile, "-C", path.Dir(a.Destination))
		targetFile = path.Join(path.Dir(a.Destination), a.Unpack)
	} else {
		downloader.Exec("/bin/chmod", "+x", targetFile)
	}

	s := state.FromLLB(defaultBaseImage, base)
	s.Copy(downloader, targetFile, a.Destination)

	return s.Get()
}

func (*DownloadAction) UpdateMetadata(_ *dockerfile2llb.Image) {
}
