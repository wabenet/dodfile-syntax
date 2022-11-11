package action

import (
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
)

type Action interface {
	Execute(llb.State) llb.State
	UpdateMetadata(*dockerfile2llb.Image)
}
