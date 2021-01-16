package action

import (
	"fmt"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
)

type UserAction struct {
	Name     string
	UID      int
	GID      int
	Shell    string
	Dotfiles string
}

func (a *UserAction) setDefaults() {
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

func (a *UserAction) Execute(base llb.State) llb.State {
	if a.Name == superUser {
		return base
	}

	a.setDefaults()

	home := fmt.Sprintf("/home/%s", a.Name)
	base = Sh(base, "addgroup --gid %d %s", a.GID, a.Name)
	base = Sh(base, "adduser --uid %d --gid %d --home %s --shell %s --disabled-password %s", a.UID, a.GID, home, a.Shell, a.Name)

	if a.Dotfiles != "" {
		source := llb.Local("context")
		base = Copy(source, a.Dotfiles, base, home)
		base = Sh(base, "chown -R %d:%d %s", a.UID, a.GID, home)
	}

	return base
}

func (a *UserAction) UpdateMetadata(metadata *dockerfile2llb.Image) {
	metadata.Config.User = a.Name

	if a.Shell != "" {
		metadata.Config.Cmd = []string{a.Shell}
	}
}
