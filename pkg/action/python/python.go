package python

import (
	"fmt"
	"path"
	"strings"

	"github.com/moby/buildkit/client/llb"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/wabenet/dodfile-syntax/pkg/state"
)

const (
	Type = "python"

	defaultBaseImage = "debian"

	tarFile     = "/python.tar.xz"
	buildPath   = "/src/python"
	installPath = "/opt/python"

	pythonGPGKey       = "E3FF2839C048B25C084DEBE9B26995E310250568"
	pythonPipVersion   = "19.2.3"
	pythonGetPipURL    = "https://github.com/pypa/get-pip/raw/309a56c5fd94bd1134053a541cb4657a4e47e09d/get-pip.py"
	pythonGetPipSHA256 = "57e3643ff19f018f8a00dfaa6b7e4620e3c1a7a2171fd218425366ec006b3bfe"
)

type Action struct {
	Version string
}

func (a *Action) Type() string {
	return Type
}

func (a *Action) Execute(base llb.State) (llb.State, error) {
	build := state.BuildBase(defaultBaseImage)
	build.Install("libsqlite3-dev", "tk-dev", "uuid-dev")

	url := fmt.Sprintf("https://www.python.org/ftp/python/%s/Python-%s.tar.xz", a.Version, a.Version)
	ascUrl := fmt.Sprintf("%s.asc", url)

	build.Download(url, tarFile)
	build.Download(ascUrl, fmt.Sprintf("%s.asc", tarFile))
	build.GPGVerify(tarFile, fmt.Sprintf("%s.asc", tarFile), []string{pythonGPGKey})

	build.CreateDirectory(buildPath)
	build.Exec("/bin/tar", "-xJf", tarFile, "-C", buildPath, "--strip-components=1")
	build.Cwd(buildPath) // TODO: CWD does not work??
	build.Sh("./configure --build=\"$(dpkg-architecture --query DEB_BUILD_GNU_TYPE)\" --enable-loadable-sqlite-extensions --prefix=%s", installPath)
	build.Sh("make -j \"$(nproc)\"")
	build.Sh("make install")

	build.SymLink("idle3", path.Join(installPath, "bin/idle"))
	build.SymLink("pydoc3", path.Join(installPath, "bin/pydoc"))
	build.SymLink("python3", path.Join(installPath, "bin/python"))
	build.SymLink("python3-config", path.Join(installPath, "bin/python-config"))

	build.Download(pythonGetPipURL, "/get-pip.py")
	build.Sh("echo \"%s *get-pip.py\" | sha256sum --check --strict -", pythonGetPipSHA256)
	build.Sh("%s/bin/python get-pip.py --disable-pip-version-check --no-cache-dir \"pip==%s\"", installPath, pythonPipVersion)

	build.Sh(`
find %s -depth \
  \( \
    \( -type d -a \( -name test -o -name tests \) \) \
    -o \
    \( -type f -a \( -name '*.pyc' -o -name '*.pyo' \) \) \
  \) -exec rm -rf '{}' +
`, installPath)

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

	config.Env = envs
}
