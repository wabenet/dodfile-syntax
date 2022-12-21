package user

import (
	"fmt"
	"strconv"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/wabenet/dodfile-syntax/pkg/state"
)

const (
	Type = "user"

	defaultBaseImage = "debian"
	defaultUser      = "user"
	defaultUID       = 1000
	defaultShell     = "/bin/bash"
	superUser        = "root"
)

type Action struct {
	Name     string `mapstructure:"name"`
	UID      int    `mapstructure:"uid"`
	GID      int    `mapstructure:"gid"`
	Shell    string `mapstructure:"shell"`
	Dotfiles string `mapstructure:"dotfiles"`
}

func (a *Action) Type() string {
	return Type
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

func (a *Action) Execute(base llb.State) (llb.State, error) {
	if a.Name == superUser {
		return base, nil
	}

	a.setDefaults()

	s := state.FromLLB(defaultBaseImage, base)

	home := fmt.Sprintf("/home/%s", a.Name)
	s.Exec("/usr/sbin/addgroup", "--gid", strconv.Itoa(a.GID), a.Name)
	s.Exec("/usr/sbin/adduser", "--uid", strconv.Itoa(a.UID), "--gid", strconv.Itoa(a.GID), "--home", home, "--shell", a.Shell, "--disabled-password", a.Name)

	if a.Dotfiles != "" {
		source := state.FromLLB(defaultBaseImage, llb.Local("context"))
		s.Copy(source, a.Dotfiles, home)
		s.Exec("/bin/chown", "-R", fmt.Sprintf("%d:%d", a.UID, a.GID), home)
	}

	return s.Get(), nil
}

func (a *Action) UpdateImage(i *dockerfile2llb.Image) {
	i.Config.User = a.Name

	if a.Shell != "" {
		i.Config.Cmd = []string{a.Shell}
	}
}
