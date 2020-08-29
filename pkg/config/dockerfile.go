package config

import (
	"context"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
)

type CopyFrom struct {
	Dockerfile  string `yaml:"dockerfile"`
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

func CopyFromDockerfile(base llb.State, cfg *CopyFrom) llb.State {
	// TODO: error handling here. convert options?
	source, _, _ := dockerfile2llb.Dockerfile2LLB(
		context.Background(),
		[]byte(cfg.Dockerfile),
		dockerfile2llb.ConvertOpt{},
	)

	return Copy(*source, cfg.Source, base, cfg.Destination)
}
