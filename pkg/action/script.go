package action

import (
	"github.com/dodo-cli/dodfile-syntax/pkg/state"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
)

type ScriptAction struct {
	Script string
}

func (a *ScriptAction) Execute(base llb.State) llb.State {
	s := state.FromLLB(defaultBaseImage, base)

        s.Sh(a.Script)

        return s.Get()
}

func (*ScriptAction) UpdateMetadata(_ *dockerfile2llb.Image) {
}
