package copy

import (
	"context"

	"github.com/dodo-cli/dodfile-syntax/pkg/state"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
)

const defaultBaseImage = "debian"

type CopyAction struct {
	Directory  string
	Image      string
	Dockerfile string
	Path       string
}

func (a *CopyAction) Execute(base llb.State) llb.State {
	s := state.FromLLB(defaultBaseImage, base)

	if len(a.Directory) > 0 {
		source := state.FromLLB(defaultBaseImage, llb.Local("context"))
		s.Copy(source, a.Directory, a.Path)

		return s.Get()
	}

	if len(a.Image) > 0 {
		source := state.From(a.Image)
		s.Copy(source, a.Path, a.Path)

		return s.Get()
	}

	// TODO: error handling here
	buildContext := llb.Local("context")
	dockerImg, _, _ := dockerfile2llb.Dockerfile2LLB(
		context.Background(),
		[]byte(a.Dockerfile),
		dockerfile2llb.ConvertOpt{
			BuildContext: &buildContext,
		},
	)

	source := state.FromLLB(defaultBaseImage, *dockerImg)
	s.Copy(source, a.Path, a.Path)

	return s.Get()
}

func (*CopyAction) UpdateMetadata(_ *dockerfile2llb.Image) {
}
