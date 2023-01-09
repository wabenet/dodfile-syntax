package base

import (
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
)

const Type = "base_image"

type Action struct {
	Config string `mapstructure:"config"`
}

func (a *Action) Type() string {
	return Type
}

func (a *Action) Execute(_ llb.State) (llb.State, error) {
	return llb.Image(a.Config), nil
}

func (a *Action) UpdateImage(_ *dockerfile2llb.Image) {}
