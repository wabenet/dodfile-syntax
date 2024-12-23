package copy

import (
	"context"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/client/llb/imagemetaresolver"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/moby/buildkit/solver/pb"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/wabenet/dodfile-syntax/pkg/state"
)

const (
	Type = "copy"

	defaultBaseImage = "debian"
)

type Action struct {
	Directory  string `mapstructure:"directory"`
	Image      string `mapstructure:"image"`
	Dockerfile string `mapstructure:"dockerfile"`
	Path       string `mapstructure:"path"`
}

func (a *Action) Type() string {
	return Type
}

func (a *Action) Execute(base llb.State) (llb.State, error) {
	s := state.FromLLB(defaultBaseImage, base)

	if len(a.Directory) > 0 {
		source := state.FromLLB(defaultBaseImage, llb.Local("context"))
		s.Copy(source, a.Directory, a.Path)

		return s.Get(), nil
	}

	if len(a.Image) > 0 {
		source := state.From(a.Image)
		s.Copy(source, a.Path, a.Path)

		return s.Get(), nil
	}

	buildContext := llb.Local("context")
	caps := pb.Caps.CapSet(pb.Caps.All())
	dockerImg, _, _, _, err := dockerfile2llb.Dockerfile2LLB(
		context.Background(),
		[]byte(a.Dockerfile),
		dockerfile2llb.ConvertOpt{
			MetaResolver: imagemetaresolver.Default(),
			LLBCaps:      &caps,
			MainContext:  &buildContext,
		},
	)
	if err != nil {
		return s.Get(), err
	}

	source := state.FromLLB(defaultBaseImage, *dockerImg)
	s.Copy(source, a.Path, a.Path)

	return s.Get(), nil
}

func (a *Action) UpdateImage(_ *oci.ImageConfig) {}
