package script

import (
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/wabenet/dodfile-syntax/pkg/state"
)

const (
	Type = "script"

	defaultBaseImage = "debian"
)

type Action struct {
	Script string `mapstructure:"script"`
	User   string `mapstructure:"user"`
	Cwd    string `mapstructure:"cwd"`
}

func (a *Action) Type() string {
	return Type
}

func (a *Action) Execute(base llb.State) (llb.State, error) {
	s := state.FromLLB(defaultBaseImage, base)

	if len(a.User) > 0 {
		s.User(a.User)
	}

	if len(a.Cwd) > 0 {
		s.Cwd(a.Cwd)
	}

	s.Sh(a.Script)

	return s.Get(), nil
}

func (a *Action) UpdateImage(_ *dockerfile2llb.Image) {}
