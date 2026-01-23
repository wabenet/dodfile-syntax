package install

import (
	"crypto/md5"
	"encoding/hex"
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
	Name    string `mapstructure:"name"`
	Repo    string `mapstructure:"repo"`
	Gpg     string `mapstructure:"gpg"`
	Release string `mapstructure:"release"`
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

			hasher := md5.New()
			hasher.Write([]byte(a.Repo))
			gpgFile := fmt.Sprintf("%s/%s.gpg", keyringsDir, hex.EncodeToString(hasher.Sum(nil)))

			s.Copy(curl, "/key.gpg", gpgFile)
			s.Sh("echo \"deb [signed-by=%s] %s %s main\" >> /etc/apt/sources.list", gpgFile, a.Repo, a.Release)
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
