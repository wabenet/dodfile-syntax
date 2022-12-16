package copy

import (
	"context"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/wabenet/dodfile-syntax/pkg/state"
)

const (
	Type = "from"

	defaultBaseImage = "debian"
)

type Action struct {
	Config []ActionConfig `mapstructure:"config"`
}

type ActionConfig struct {
	Directory  string `mapstructure:"directory"`
	Image      string `mapstructure:"image"`
	Dockerfile string `mapstructure:"dockerfile"`
	Path       string `mapstructure:"path"`
}

func (a *Action) Type() string {
	return Type
}

func (a *Action) Execute(base llb.State) (llb.State, error) {
	var err error

	s := base

	for _, ac := range a.Config {
		s, err = ac.Execute(s)
		if err != nil {
			return s, err
		}
	}

	return s, nil
}

func (a *ActionConfig) Execute(base llb.State) (llb.State, error) {
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
	dockerImg, _, err := dockerfile2llb.Dockerfile2LLB(
		context.Background(),
		[]byte(a.Dockerfile),
		dockerfile2llb.ConvertOpt{
			BuildContext: &buildContext,
		},
	)
	if err != nil {
		return s.Get(), err
	}

	source := state.FromLLB(defaultBaseImage, *dockerImg)
	s.Copy(source, a.Path, a.Path)

	return s.Get(), nil
}

func (a *Action) UpdateImage(_ dockerfile2llb.Image) {}
