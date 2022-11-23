package user

import (
	"fmt"
	"strconv"

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

	s := state.FromLLB(defaultBaseImage, base)

	home := fmt.Sprintf("/home/%s", a.Name)
	s.Exec("/usr/sbin/addgroup", "--gid", strconv.Itoa(a.GID), a.Name)
	s.Exec("/usr/sbin/adduser", "--uid", strconv.Itoa(a.UID), "--gid", strconv.Itoa(a.GID), "--home", home, "--shell", a.Shell, "--disabled-password", a.Name)

	return s.Get()
}
