package build

import (
	"fmt"

	"github.com/dodo-cli/dodfile-syntax/pkg/action"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/moby/buildkit/util/system"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	defaultBaseImage = "debian"
)

type Image struct {
	BaseImage string            `yaml:"base_image"`
	User      *User             `yaml:"user"`
	Env       map[string]string `yaml:"env"`
	Packages  *Packages         `yaml:"packages"`
	Download  []*Download       `yaml:"download"`
	From      []*CopyFrom       `yaml:"from"`
	Run       []*Run            `yaml:"run"`
}

type User struct {
	Name     string `yaml:"name"`
	UID      int    `yaml:"uid"`
	GID      int    `yaml:"gid"`
	Shell    string `yaml:"shell"`
	Dotfiles string `yaml:"dotfiles"`
}

type Packages struct {
	Name []string `yaml:"name"`
	Repo []string `yaml:"repo"`
	Gpg  []string `yaml:"gpg"`
}

type Download struct {
	Source      string `yaml:"source"`
	Sha256      string `yaml:"sha256"`
	Unpack      string `yaml:"unpack"`
	Destination string `yaml:"destination"`
}

type CopyFrom struct {
	Directory  string            `yaml:"directory"`
	Image      string            `yaml:"image"`
	Dockerfile string            `yaml:"dockerfile"`
	Path       string            `yaml:"path"`
	Env        map[string]string `yaml:"env"`
}

type Run struct {
	Script string `yaml:"script"`
}

func (img *Image) base() llb.State {
	baseImage := img.BaseImage
	if len(baseImage) == 0 {
		baseImage = defaultBaseImage
	}

	return llb.Image(baseImage)
}

func (img *Image) metadata() dockerfile2llb.Image {
	metadata := dockerfile2llb.Image{
		Image: specs.Image{
			Architecture: "amd64",
			OS:           "linux",
		},
	}

	env := map[string]string{"PATH": system.DefaultPathEnv}

	for key, value := range img.Env {
		switch key {
		case "PATH":
			env["PATH"] = fmt.Sprintf("%s:%s", env["PATH"], value)
		default:
			env[key] = value
		}
	}

	envs := []string{}
	for key, value := range env {
		envs = append(envs, fmt.Sprintf("%s=%s", key, value))
	}

	metadata.RootFS.Type = "layers"
	metadata.Config.Env = envs

	return metadata
}

func (img *Image) Build() (llb.State, dockerfile2llb.Image) {
	s := img.base()
	metadata := img.metadata()

	if img.User != nil {
		a := &action.UserAction{
			Name:     img.User.Name,
			UID:      img.User.UID,
			GID:      img.User.GID,
			Shell:    img.User.Shell,
			Dotfiles: img.User.Dotfiles,
		}

		s = a.Execute(s)
		a.UpdateMetadata(&metadata)
	}

	if img.Packages != nil {
		a := &action.PackageAction{
			Name: img.Packages.Name,
			Repo: img.Packages.Repo,
			Gpg:  img.Packages.Gpg,
		}

		s = a.Execute(s)
		a.UpdateMetadata(&metadata)
	}

	for _, d := range img.Download {
		a := &action.DownloadAction{
			Source:      d.Source,
			Sha256:      d.Sha256,
			Unpack:      d.Unpack,
			Destination: d.Destination,
		}

		s = a.Execute(s)
		a.UpdateMetadata(&metadata)
	}

	for _, c := range img.From {
		a := &action.CopyAction{
			Directory:  c.Directory,
			Image:      c.Image,
			Dockerfile: c.Dockerfile,
			Path:       c.Path,
		}

		s = a.Execute(s)
		a.UpdateMetadata(&metadata)
	}

	for _, r := range img.Run {
		a := &action.ScriptAction{
			Script: r.Script,
		}

		s = a.Execute(s)
		a.UpdateMetadata(&metadata)
	}

	return s, metadata
}
