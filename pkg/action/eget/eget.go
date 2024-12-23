package eget

import (
	"github.com/moby/buildkit/client/llb"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/wabenet/dodfile-syntax/pkg/state"
)

const (
	Type = "eget"

	defaultBaseImage  = "debian"
	egetInstallerURL  = "https://zyedidia.github.io/eget.sh"
	egetInstallerSHA  = "0e64b8a3c13f531da005096cc364ac77835bda54276fedef6c62f3dbdc1ee919"
	egetInstallerPath = "/tmp/eget.sh"
	egetDownloadDir   = "/tmp/eget/"
)

type Action struct {
	Repo string `mapstructure:"repo"`

	Tag          string   `mapstructure:"tag"`
	Prerelease   bool     `mapstructure:"pre-release"`
	ExtractFile  string   `mapstructure:"file"`
	All          bool     `mapstructure:"all"`
	Asset        []string `mapstructure:"asset"`
	VerifySHA256 string   `mapstructure:"verify-sha256"`
}

func (a *Action) Type() string {
	return Type
}

func (a *Action) Execute(base llb.State) (llb.State, error) {
	downloader := state.From(defaultBaseImage)
	downloader.Install("apt-transport-https", "curl", "ca-certificates", "tar")
	downloader.Exec("/usr/bin/curl", "-o", egetInstallerPath, egetInstallerURL)
	downloader.Sh("echo \"%s  %s\" | sha256sum -c -", egetInstallerSHA, egetInstallerPath)
	downloader.Exec("/bin/sh", egetInstallerPath)
	downloader.Exec("/bin/mkdir", "-p", egetDownloadDir)
	downloader.Cwd(egetDownloadDir)

	egetCmd := []string{"/eget", a.Repo}

	if a.Tag != "" {
		egetCmd = append(egetCmd, "--tag", a.Tag)
	}

	if a.Prerelease {
		egetCmd = append(egetCmd, "--pre-release")
	}

	if a.ExtractFile != "" {
		egetCmd = append(egetCmd, "--file", a.ExtractFile)
	}

	if a.All {
		egetCmd = append(egetCmd, "--all")
	}

	for _, asset := range a.Asset {
		egetCmd = append(egetCmd, "--asset", asset)
	}

	if a.VerifySHA256 != "" {
		egetCmd = append(egetCmd, "--verify-sha256", a.VerifySHA256)
	}

	downloader.Exec(egetCmd...)

	s := state.FromLLB(defaultBaseImage, base)
	s.CopyDir(downloader, egetDownloadDir, "/bin")

	return s.Get(), nil

}

func (a *Action) UpdateImage(_ *oci.ImageConfig) {}
