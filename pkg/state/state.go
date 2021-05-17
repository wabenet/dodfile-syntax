package state

import (
	"fmt"
	"path"

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

func (s *State) Install(pkgs ...string) {
	execState := s.current.Run(llb.Args(updateCmd()))
	execState = execState.Run(llb.Args(append(installCmd(), pkgs...)))
	s.current = execState.Root()
}

func (s *State) Copy(src *State, srcPath string, destPath string) {
	execState := llb.Image(s.baseImage).Run(llb.Args([]string{"/bin/mkdir", "-p", path.Join("/dest", path.Dir(destPath))}))
	s.current = execState.AddMount("/dest", s.current)

	execState = s.current.Run(llb.Args([]string{"/bin/cp", "-a", "-T", path.Join("/src", srcPath), path.Join("/dest", destPath)}))
	execState.AddMount("/src", src.current, llb.Readonly)
	s.current = execState.AddMount("/dest", s.current)
}

func updateCmd() []string {
	return []string{"/usr/bin/apt-get", "update"}
}

func installCmd() []string {
	return []string{
		"/usr/bin/apt-get",
		"install",
		"-y",
		"--no-install-recommends",
		"--no-install-suggests",
	}
}
