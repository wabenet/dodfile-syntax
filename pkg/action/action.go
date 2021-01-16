package action

import (
	"fmt"
	"path"

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

func Sh(s llb.State, cmd string, v ...interface{}) llb.State {
	return s.Run(llb.Args([]string{"/bin/sh", "-c", fmt.Sprintf(cmd, v...)})).Root()
}

func Copy(src llb.State, srcPath string, dest llb.State, destPath string) llb.State {
	cp := llb.Image(defaultBaseImage).Run(llb.Args([]string{"/bin/sh", "-c", fmt.Sprintf("mkdir -p /dest/%s && cp -aT /src/%s /dest/%s", path.Dir(destPath), srcPath, destPath)}))
	cp.AddMount("/src", src, llb.Readonly)

	return cp.AddMount("/dest", dest)
}
