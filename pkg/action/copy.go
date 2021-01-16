package action

import (
	"context"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
)

type CopyAction struct {
	Directory  string
	Image      string
	Dockerfile string
	Path       string
}

func (a *CopyAction) Execute(base llb.State) llb.State {
	if len(a.Directory) > 0 {
		source := llb.Local("context")
		return Copy(source, a.Directory, base, a.Path)
	}

	if len(a.Image) > 0 {
		source := llb.Image(a.Image)
		return Copy(source, a.Path, base, a.Path)
	}

	buildContext := llb.Local("context")

	// TODO: error handling here
	source, _, _ := dockerfile2llb.Dockerfile2LLB(
		context.Background(),
		[]byte(a.Dockerfile),
		dockerfile2llb.ConvertOpt{
			BuildContext: &buildContext,
		},
	)

	return Copy(*source, a.Path, base, a.Path)
}

func (*CopyAction) UpdateMetadata(_ *dockerfile2llb.Image) {
}
