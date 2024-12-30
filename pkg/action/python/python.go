package python

import (
	"fmt"
	"path"
	"strings"

	"github.com/moby/buildkit/client/llb"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/wabenet/dodfile-syntax/pkg/action/python/pythonapi"
	"github.com/wabenet/dodfile-syntax/pkg/state"
)

const (
	Type = "python"

	defaultBaseImage = "debian"

	tarFile     = "/python.tar"
	buildPath   = "/src/python"
	installPath = "/opt/python"

	pythonGPGKey    = "E3FF2839C048B25C084DEBE9B26995E310250568"
	pythonGetPipURL = "https://bootstrap.pypa.io/get-pip.py"
)

type Action struct {
	Version           string   `mapstructure:"version"`
	PipVersion        string   `mapstructure:"pip_version"`
	PipPackages       []string `mapstructure:"pip_packages"`
	BuildDependencies []string `mapstructure:"build_dependencies"`
}

func (a *Action) Type() string {
	return Type
}

func (a *Action) Verify() error {
	if a.Version == "" {
		a.Version = pythonapi.Latest
	}

	return nil
}

func buildDependencies() *state.State {
	s := state.From(defaultBaseImage)

	s.Install(
		"autoconf",
		"build-essential",
		"curl",
		"git",
		"gnupg2",
		"libffi-dev",
		"libgdbm-dev",
		"libncurses5-dev",
		"libreadline6-dev",
		"libsqlite3-dev",
		"libssl-dev",
		"libyaml-dev",
		"tk-dev",
		"uuid-dev",
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
	build.Install(a.BuildDependencies...)

	if err := a.Verify(); err != nil {
		return build, err
	}

	release, err := pythonapi.GetDownload(a.Version, pythonapi.Source)
	if err != nil {
		return build, fmt.Errorf("no valid download URL found: %w", err)
	}

	build.Download(release.URL, tarFile)

	signatureFile := tarFile + ".asc"
	build.Download(release.GPGSignatureFile, signatureFile)
	build.GPGVerify(tarFile, signatureFile, []string{pythonGPGKey})

	build.CreateDirectory(buildPath)

	switch {
	case strings.HasSuffix(release.URL, ".tar.xz"):
		build.Exec("/bin/tar", "-xJf", tarFile, "-C", buildPath, "--strip-components=1")
	case strings.HasSuffix(release.URL, ".tar.gz"), strings.HasSuffix(release.URL, ".tgz"):
		build.Exec("/bin/tar", "-xzf", tarFile, "-C", buildPath, "--strip-components=1")
	}

	build.Cwd(buildPath)
	build.Sh("./configure --build=\"$(dpkg-architecture --query DEB_BUILD_GNU_TYPE)\" --enable-loadable-sqlite-extensions --prefix=%s", installPath)
	build.Sh("make -j \"$(nproc)\"")
	build.Sh("make install")

	build.Cwd("/")
	build.SymLink("idle3", path.Join(installPath, "bin/idle"))
	build.SymLink("pydoc3", path.Join(installPath, "bin/pydoc"))
	build.SymLink("python3", path.Join(installPath, "bin/python"))
	build.SymLink("python3-config", path.Join(installPath, "bin/python-config"))

	build.Download(pythonGetPipURL, "/get-pip.py") // TODO: verify checksum

	if a.PipVersion == "" {
		build.Sh("%s/bin/python get-pip.py --no-cache-dir", installPath)
	} else {
		build.Sh("%s/bin/python get-pip.py --disable-pip-version-check --no-cache-dir \"pip==%s\"", installPath, a.PipVersion)
	}

	if len(a.PipPackages) > 0 {
		build.Sh("%s/bin/pip3 install %s", installPath, strings.Join(a.PipPackages, " "))
	}

	build.Sh(`
find %s -depth \
  \( \
    \( -type d -a \( -name test -o -name tests \) \) \
    -o \
    \( -type f -a \( -name '*.pyc' -o -name '*.pyo' \) \) \
  \) -exec rm -rf '{}' +
`, installPath)

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
