package config

import (
	"context"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
)

type CopyFrom struct {
	Directory  string            `yaml:"directory"`
	Image      string            `yaml:"image"`
	Dockerfile string            `yaml:"dockerfile"`
	Path       string            `yaml:"path"`
	Env        map[string]string `yaml:"env"`
}

func CopyFromDockerfile(base llb.State, cfg *CopyFrom) llb.State {
	if len(cfg.Directory) > 0 {
		source := llb.Local("context")
		return Copy(source, cfg.Directory, base, cfg.Path)
	}

	if len(cfg.Image) > 0 {
		source := llb.Image(cfg.Image)
		return Copy(source, cfg.Path, base, cfg.Path)
	}

	buildContext := llb.Local("context")

	// TODO: error handling here
	source, _, _ := dockerfile2llb.Dockerfile2LLB(
		context.Background(),
		[]byte(cfg.Dockerfile),
		dockerfile2llb.ConvertOpt{
			BuildContext: &buildContext,
		},
	)

	return Copy(*source, cfg.Path, base, cfg.Path)
}
