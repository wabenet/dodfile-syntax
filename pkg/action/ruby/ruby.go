package ruby

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/moby/buildkit/client/llb"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/wabenet/dodfile-syntax/pkg/action/ruby/github"
	"github.com/wabenet/dodfile-syntax/pkg/state"
)

const (
	Type = "ruby"

	defaultBaseImage = "debian"

	buildPath   = "/src/ruby"
	installPath = "/opt/ruby"
	gitRepo     = "https://github.com/ruby/ruby.git"

	latest = "latest"

	gemRC = `
install: --no-document
update: --no-document
`
)

var ErrNoRelease = errors.New("no valid release found")

type Action struct {
	Version           string   `mapstructure:"version"`
	Gems              []string `mapstructure:"gems"`
	BuildDependencies []string `mapstructure:"build_dependencies"`
}

func (a *Action) Type() string {
	return Type
}

func (a *Action) Verify() error {
	if a.Version == "" {
		a.Version = latest
	}

	return nil
}

func buildDependencies() *state.State {
	s := state.From(defaultBaseImage)

	s.Install(
		"autoconf",
		"build-essential",
		"ca-certificates",
		"bison",
		"git",
		"libssl-dev",
		"libyaml-dev",
		"ruby", // yeah
		"zlib1g-dev",
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

	if err := a.Verify(); err != nil {
		return build, err
	}

	tag, err := getTagForVersion(a.Version)
	if err != nil {
		return build, fmt.Errorf("no valid git tag found: %w", err)
	}

	build.Exec("/usr/bin/git", "clone", "--branch", tag, "--depth", "1", gitRepo, buildPath)

	build.Cwd(buildPath)
	build.Sh("./autogen.sh")
	build.Sh("./configure --build=\"$(dpkg-architecture --query DEB_BUILD_GNU_TYPE)\" --disable-install-doc --prefix=%s", installPath)
	build.Sh("make -j \"$(nproc)\"")
	build.Sh("make install")

	build.CreateDirectory(path.Join(installPath, "etc"))
	build.Sh("echo '%s' > %s/etc/gemrc", gemRC, installPath)

	build.Env("GEM_HOME", GemHome())
	build.CreateDirectory(GemHome())

	if len(a.Gems) > 0 {
		build.Install(a.BuildDependencies...)
		GemInstall(build, a.Gems...)
	}

	return build, nil
}

func (a *Action) UpdateImage(config *oci.ImageConfig) {
	envs := config.Env

	for i, env := range envs {
		if parts := strings.SplitN(env, "=", 2); parts[0] == "PATH" {
			envs[i] = fmt.Sprintf("PATH=%s/bin:%s/bin:%s/gems/bin:%s", installPath, GemHome(), GemHome(), parts[1])
		}
	}

	envs = append(envs, "GEM_HOME="+GemHome())
	envs = append(envs, "BUNDLE_PATH="+GemHome())
	envs = append(envs, "BUNDLE_SILENCE_ROOT_WARNING=1")
	envs = append(envs, "BUNDLE_APP_CONFIG="+GemHome())

	config.Env = envs
}

func GemInstall(s *state.State, gems ...string) {
	s.Exec(append([]string{path.Join(installPath, "bin/gem"), "install"}, gems...)...)
}

func GemHome() string {
	return path.Join(installPath, "bundle")
}

func getTagForVersion(version string) (string, error) {
	if version == latest {
		r, err := github.GetLatestRelease("ruby", "ruby")
		if err != nil {
			return "", err
		}

		return r.TagName, nil
	}

	rs, err := github.GetReleases("ruby", "ruby")
	if err != nil {
		return "", err
	}

	for _, r := range rs {
		if r.Name == version {
			return r.TagName, nil
		}
	}

	return "", ErrNoRelease
}
