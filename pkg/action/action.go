package action

import (
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
)

const (
	defaultBaseImage = "debian"
	defaultUser      = "user"
	defaultUID       = 1000
	defaultShell     = "/bin/bash"
	superUser        = "root"
)

type Action interface {
	Execute(llb.State) llb.State
	UpdateMetadata(*dockerfile2llb.Image)
}
