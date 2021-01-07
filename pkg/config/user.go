package config

import (
	"fmt"

	"github.com/moby/buildkit/client/llb"
)

type User struct {
	Name     string `yaml:"name"`
	Uid      int    `yaml:"uid"`
	Gid      int    `yaml:"gid"`
	Shell    string `yaml:"shell"`
	Dotfiles string `yaml:"dotfiles"`
}

func (u *User) SetDefaults() {
	if u.Name == "" {
		u.Name = "user"
	}
	if u.Uid == 0 {
		u.Uid = 1000
	}
	if u.Gid == 0 {
		u.Gid = 1000
	}
	if u.Shell == "" {
		u.Shell = "/bin/bash"
	}
}

func SetupUser(base llb.State, u *User) llb.State {
	if u.Name == "root" {
		return base
	}

	u.SetDefaults()

	home := fmt.Sprintf("/home/%s", u.Name)
	base = Sh(base, "addgroup --gid %d %s", u.Gid, u.Name)
	base = Sh(base, "adduser --uid %d --gid %d --home %s --shell %s --disabled-password %s", u.Uid, u.Gid, home, u.Shell, u.Name)

	if u.Dotfiles != "" {
		source := llb.Local("context")
		base = Copy(source, u.Dotfiles, base, home)
		base = Sh(base, "chown -R %d:%d %s", u.Uid, u.Gid, home)
	}

	return base
}
