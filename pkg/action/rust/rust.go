package rust

import (
	"fmt"
	"strings"

	"github.com/moby/buildkit/client/llb"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/wabenet/dodfile-syntax/pkg/state"
)

const (
	Type = "rust"

	defaultBaseImage = "debian"

	installPath = "/opt/rust"

	rustupURL  = "https://sh.rustup.rs"
	rustupFile = "/tmp/rustup-init.sh"
)

type Action struct {
	Crates            []string `mapstructure:"crates"`
	BuildDependencies []string `mapstructure:"build_dependencies"`
}

func (a *Action) Type() string {
	return Type
}

func (a *Action) Execute(base llb.State) (llb.State, error) {
	build := state.From(defaultBaseImage)

	build.Install("ca-certificates", "curl")

	build.Env("RUSTUP_HOME", fmt.Sprintf("%s/rustup", installPath))
	build.Env("CARGO_HOME", fmt.Sprintf("%s/cargo", installPath))

	build.Download(rustupURL, rustupFile)
	build.Exec("/bin/sh", rustupFile, "-y", "--no-modify-path")

	if len(a.Crates) > 0 {
		build.Install(a.BuildDependencies...)
		build.Exec(append([]string{fmt.Sprintf("%s/cargo/bin/cargo", installPath), "install"}, a.Crates...)...)
	}

	s := state.FromLLB(defaultBaseImage, base)

	s.Copy(build, installPath, installPath)

	return s.Get(), nil
}

func (a *Action) UpdateImage(config *oci.ImageConfig) {
	envs := config.Env

	for i, env := range envs {
		if parts := strings.SplitN(env, "=", 2); parts[0] == "PATH" {
			envs[i] = fmt.Sprintf("PATH=%s/cargo/bin:%s", installPath, parts[1])
		}
	}

	envs = append(envs, fmt.Sprintf("RUSTUP_HOME=%s/rustup", installPath))
	envs = append(envs, fmt.Sprintf("CARGO_HOME=%s/cargo", installPath))

	config.Env = envs
}
