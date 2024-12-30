package golang

import (
	"fmt"
	"strings"

	"github.com/moby/buildkit/client/llb"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/wabenet/dodfile-syntax/pkg/action/golang/golangapi"
	"github.com/wabenet/dodfile-syntax/pkg/state"
)

const (
	Type = "golang"

	defaultBaseImage = "debian"

	tarFile     = "/golang.tar.gz"
	installPath = "/opt/go"
)

type Action struct {
	Version string   `mapstructure:"version"`
	Install []string `mapstructure:"install"`
}

func (a *Action) Type() string {
	return Type
}

// TODO: make verify part of the action interface
func (a *Action) Verify() error {
	if a.Version == "" {
		a.Version = golangapi.Latest
	}

	return nil
}

func (a *Action) Execute(base llb.State) (llb.State, error) {
	build := state.From(defaultBaseImage)

	if err := a.Verify(); err != nil {
		return build.Get(), err
	}

	build.Install("ca-certificates", "gzip", "tar")

	release, err := golangapi.GetDownload(a.Version)
	if err != nil {
		return build.Get(), err
	}
	url, err := release.URL()
	if err != nil {
		return build.Get(), err
	}

	build.Download(url, tarFile)
	build.Sh("echo \"%s %s\" | sha256sum --check --strict -", release.ChecksumSHA256, tarFile)
	build.Exec("/bin/mkdir", "-p", installPath)
	build.Exec("/bin/tar", "-xzf", tarFile, "-C", installPath, "--strip-components=1")

	build.Env("GOPATH", installPath)

	for _, i := range a.Install {
		build.Exec(fmt.Sprintf("%s/bin/go", installPath), "install", i)
	}

	s := state.FromLLB(defaultBaseImage, base)

	s.Copy(build, installPath, installPath)

	return s.Get(), nil
}

func (a *Action) UpdateImage(config *oci.ImageConfig) {
	envs := config.Env

	for i, env := range envs {
		if parts := strings.SplitN(env, "=", 2); parts[0] == "PATH" {
			envs[i] = fmt.Sprintf("PATH=%s/bin:%s", installPath, parts[1])
		}
	}

	envs = append(envs, fmt.Sprintf("GOPATH=%s", installPath))

	config.Env = envs
}
