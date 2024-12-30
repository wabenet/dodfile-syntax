package rust

import (
	"fmt"
	"path"
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
	s := state.FromLLB(defaultBaseImage, base)

	build, err := a.Build()
	if err != nil {
		return s.Get(), err
	}

	s.Copy(build, installPath, installPath)

	return s.Get(), nil
}

func (a *Action) Build() (*state.State, error) {
	build := state.From(defaultBaseImage)

	build.Install("ca-certificates", "curl")

	build.Env("RUSTUP_HOME", RustupHome())
	build.Env("CARGO_HOME", CargoHome())

	build.Download(rustupURL, rustupFile)
	build.Exec("/bin/sh", rustupFile, "-y", "--no-modify-path")

	if len(a.Crates) > 0 {
		build.Install(a.BuildDependencies...)
		CargoInstall(build, a.Crates...)
	}

	return build, nil
}

func (a *Action) UpdateImage(config *oci.ImageConfig) {
	envs := config.Env

	for i, env := range envs {
		if parts := strings.SplitN(env, "=", 2); parts[0] == "PATH" {
			envs[i] = fmt.Sprintf("PATH=%s/cargo/bin:%s", installPath, parts[1])
		}
	}

	envs = append(envs, "RUSTUP_HOME="+RustupHome())
	envs = append(envs, "CARGO_HOME="+CargoHome())

	config.Env = envs
}

func CargoInstall(s *state.State, crates ...string) {
	s.Exec(append([]string{path.Join(installPath, "cargo/bin/cargo"), "install"}, crates...)...)
}

func RustupHome() string {
	return path.Join(installPath, "rustup")
}

func CargoHome() string {
	return path.Join(installPath, "cargo")
}
