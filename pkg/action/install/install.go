package install

import (
	"fmt"
	"strings"

	"github.com/moby/buildkit/client/llb"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/wabenet/dodfile-syntax/pkg/state"
)

const (
	Type = "install"

	defaultBaseImage = "debian"
	keyringsDir      = "/etc/apt/keyrings"
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
		if len(a.Gpg) > 0 {
			curl := state.From(defaultBaseImage)
			curl.Install("apt-transport-https", "curl", "ca-certificates")
			curl.Exec("/usr/bin/curl", "-Lo", "/key.gpg", a.Gpg)
			s.CreateDirectory(keyringsDir)
			s.Copy(curl, "/key.gpg", fmt.Sprintf("%s/%s.gpg", keyringsDir, a.Name))

			s.Sh("echo \"deb [signed-by=%s/%s.gpg] %s any main\" >> /etc/apt/sources.list", keyringsDir, a.Name, a.Repo)
		} else {
			s.Sh("echo \"deb %s trixie main\" >> /etc/apt/sources.list", a.Repo)
		}

	}

	if len(a.Name) > 0 {
		s.Install(strings.Fields(a.Name)...)
	}

	return s.Get(), nil
}

func (a *Action) UpdateImage(_ *oci.ImageConfig) {}
