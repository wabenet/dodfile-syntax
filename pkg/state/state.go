package state

import (
	"fmt"

	"github.com/moby/buildkit/client/llb"
)

const (
	contextKey = "context"
)

type State struct {
	baseImage string
	current   llb.State
}

func From(base string) *State {
	return &State{
		baseImage: base,
		current:   llb.Image(base),
	}
}

func FromLLB(base string, state llb.State) *State {
	return &State{
		baseImage: base,
		current:   state,
	}
}

func (s *State) Get() llb.State {
	return s.current
}

func (s *State) User(u string) {
	s.current = s.current.User(u)
}

func (s *State) Cwd(dir string) {
	s.current = s.current.Dir(dir)
}

func (s *State) Exec(args ...string) {
	s.current = s.current.Run(llb.Args(args)).Root()
}

func (s *State) Sh(cmd string, v ...interface{}) {
	s.Exec("/bin/sh", "-c", fmt.Sprintf(cmd, v...))
}
