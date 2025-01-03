package fetch

import (
	"path"
	"path/filepath"

	"github.com/moby/buildkit/client/llb"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/wabenet/dodfile-syntax/pkg/state"
)

const (
	Type = "fetch"

	defaultBaseImage = "debian"
)

type Action struct {
	Source      string `mapstructure:"source"`
	Sha256      string `mapstructure:"sha256"`
	Unpack      string `mapstructure:"unpack"`
	Destination string `mapstructure:"destination"`
}

func (a *Action) Type() string {
	return Type
}

func (a *Action) Execute(base llb.State) (llb.State, error) {
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
		switch filepath.Ext(targetFile) {
		case ".tgz":
			downloader.Exec("/bin/tar", "-zxf", targetFile, "-C", path.Dir(a.Destination))
		case ".zip":
			downloader.Exec("/usr/bin/unzip", targetFile, "-d", path.Dir(a.Destination))
		default:
			// TODO: should this really be the default?
			downloader.Exec("/bin/tar", "-zxf", targetFile, "-C", path.Dir(a.Destination))
		}
		targetFile = path.Join(path.Dir(a.Destination), a.Unpack)
	} else {
		downloader.Exec("/bin/chmod", "+x", targetFile)
	}

	s := state.FromLLB(defaultBaseImage, base)
	s.Copy(downloader, targetFile, a.Destination)

	return s.Get(), nil
}

func (a *Action) UpdateImage(_ *oci.ImageConfig) {}
