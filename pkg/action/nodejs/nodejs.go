package nodejs

import (
	"fmt"
	"path"
	"strings"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/util/system"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/wabenet/dodfile-syntax/pkg/action/nodejs/nodeapi"
	"github.com/wabenet/dodfile-syntax/pkg/action/nodejs/yarnapi"
	"github.com/wabenet/dodfile-syntax/pkg/state"
)

const (
	Type = "nodejs"

	defaultBaseImage = "debian"

	yarnTarFile = "/yarn.tgz"
	installPath = "/opt/nodejs"
)

type Action struct {
	Version           string   `mapstructure:"version"`
	YarnVersion       string   `mapstructure:"yarn_version"`
	Modules           []string `mapstructure:"modules"`
	BuildDependencies []string `mapstructure:"build_dependencies"`
}

func (a *Action) Type() string {
	return Type
}

func (a *Action) Verify() error {
	if a.Version == "" {
		a.Version = nodeapi.Latest
	}

	if a.YarnVersion == "" {
		a.YarnVersion = yarnapi.Latest
	}

	return nil
}

func buildDependencies() *state.State {
	s := state.From(defaultBaseImage)

	s.Install(
		"build-essential",
		"ca-certificates",
		"curl",
		"gnupg2",
	)

	return s
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
	build := buildDependencies()
	build.Install(a.BuildDependencies...)

	if err := a.Verify(); err != nil {
		return build, err
	}

	release, err := nodeapi.GetDownload(a.Version)
	if err != nil {
		return build, fmt.Errorf("no valid download URL found: %w", err)
	}

	build.Download(release.URL, release.Filename)

	sumsFile := "SHASUMS256.txt.asc"
	build.Download(release.SignedHashFile, sumsFile)

	for _, key := range validGPGKeys() {
		build.Exec("/usr/bin/gpg", "--batch", "--keyserver", "keys.openpgp.org", "--recv-keys", key)
	}

	build.Exec("/usr/bin/gpg", "--batch", "--verify", sumsFile)
	build.Sh("grep %s SHASUMS256.txt.asc | sha256sum -c -", release.Filename)

	build.CreateDirectory(installPath)

	switch {
	case strings.HasSuffix(release.Filename, ".tar.xz"):
		build.Exec("/bin/tar", "-xJf", release.Filename, "-C", installPath, "--strip-components=1", "--no-same-owner")
	case strings.HasSuffix(release.URL, ".tar.gz"), strings.HasSuffix(release.URL, ".tgz"):
		build.Exec("/bin/tar", "-xzf", release.Filename, "-C", installPath, "--strip-components=1", "--no-same-owner")
	}

	// TODO: this feels unneccesary
	build.Env("PATH", fmt.Sprintf("%s/bin:%s", installPath, system.DefaultPathEnv("unix")))

	yarnRelease, err := yarnapi.GetDownload(a.YarnVersion)
	if err != nil {
		return build, fmt.Errorf("no valid download URL found: %w", err)
	}

	build.Download(yarnRelease.URL, yarnTarFile)
	build.Sh("echo '%s  %s' | shasum -c -", yarnRelease.SHASum, yarnTarFile)
	build.Exec("/bin/tar", "-xzf", yarnTarFile, "-C", installPath, "--strip-components=1", "--no-same-owner")

	build.Exec(path.Join(installPath, "/bin/yarn"), "config", "set", "global-folder", installPath)

	for _, module := range a.Modules {
		build.Exec(path.Join(installPath, "/bin/yarn"), "global", "add", module, "--prefix", installPath)
	}

	return build, nil
}

func (a *Action) UpdateImage(config *oci.ImageConfig) {
	envs := config.Env

	for i, env := range envs {
		if parts := strings.SplitN(env, "=", 2); parts[0] == "PATH" {
			envs[i] = fmt.Sprintf("PATH=%s/bin:%s", installPath, parts[1])
		}
	}

	config.Env = envs
}

func validGPGKeys() []string {
	return []string{
		"5BE8A3F6C8A5C01D106C0AD820B1A390B168D356", // Antoine du Hamel
		"DD792F5973C6DE52C432CBDAC77ABFA00DDBF2B7", // Juan José Arboleda
		"CC68F5A3106FF448322E48ED27F5E38D5B0A215F", // Marco Ippolito
		"8FCCA13FEF1D0C2E91008E09770F7A9A5AE15600", // Michaël Zasso
		"890C08DB8579162FEE0DF9DB8BEAB4DFCF555EF4", // Rafael Gonzaga
		"C82FA3AE1CBEDC6BE46B9360C43CEC45C17AB93C", // Richard Lau
		"108F52B48DB57BB0CC439B2997B01419BD92F80A", // Ruy Adorno
		"A363A499291CBBC940DD62E41F10027AF002F8B0", // Ulises Gascón
	}
}
