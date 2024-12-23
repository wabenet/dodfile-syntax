package base

import (
	"github.com/moby/buildkit/client/llb"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
)

const Type = "base"

type Action struct {
	Name string `mapstructure:"name"`
}

func (a *Action) Type() string {
	return Type
}

func (a *Action) Execute(_ llb.State) (llb.State, error) {
	return llb.Image(a.Name), nil
}

func (a *Action) UpdateImage(_ *oci.ImageConfig) {}
