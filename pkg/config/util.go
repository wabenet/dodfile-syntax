package config

import (
	"fmt"

	"github.com/moby/buildkit/client/llb"
)

func Sh(s llb.State, cmd string, v ...interface{}) llb.State {
	return s.Run(llb.Args([]string{"/bin/sh", "-c", fmt.Sprintf(cmd, v...)})).Root()
}

func Copy(src llb.State, srcPath string, dest llb.State, destPath string) llb.State {
	cp := llb.Image(defaultBaseImage).Run(llb.Shlexf("cp -a /src%s /dest%s", srcPath, destPath))
	cp.AddMount("/src", src, llb.Readonly)

	return cp.AddMount("/dest", dest)
}
