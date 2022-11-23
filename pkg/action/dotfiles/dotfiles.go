package dotfiles

import (
	"fmt"

	"github.com/dodo-cli/dodfile-syntax/pkg/state"
	"github.com/moby/buildkit/client/llb"
)

const (
	defaultBaseImage = "debian"
	defaultUser      = "user"
	defaultUID       = 1000
	defaultShell     = "/bin/bash"
	superUser        = "root"
)

type Action struct {
	Name     string
	UID      int
	GID      int
	Shell    string
	Dotfiles string
}

func (a *Action) setDefaults() {
	if a.Name == "" {
		a.Name = defaultUser
	}

	if a.UID == 0 {
		a.UID = defaultUID
	}

	if a.GID == 0 {
		a.GID = defaultUID
	}

	if a.Shell == "" {
		a.Shell = defaultShell
	}
}

func (a *Action) Execute(base llb.State) llb.State {
	home := fmt.Sprintf("/home/%s", a.Name)

	s := state.FromLLB(defaultBaseImage, base)

	if a.Dotfiles != "" {
		source := state.FromLLB(defaultBaseImage, llb.Local("context"))
		s.Copy(source, a.Dotfiles, home)
		s.Exec("/bin/chown", "-R", fmt.Sprintf("%d:%d", a.UID, a.GID), home)
	}
	return s.Get()
}
