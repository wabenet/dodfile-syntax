package config

import (
	"fmt"

	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/moby/buildkit/util/system"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

func (img *Image) GetDockerMetadata() dockerfile2llb.Image {
	metadata := dockerfile2llb.Image{
		Image: specs.Image{
			Architecture: "amd64",
			OS:           "linux",
		},
	}

	env := map[string]string{"PATH": system.DefaultPathEnv}

	for _, from := range img.From {
		for key, value := range from.Env {
			switch key {
			case "PATH":
				env["PATH"] = fmt.Sprintf("%s:%s", env["PATH"], value)
			default:
				env[key] = value
			}
		}
	}

	for key, value := range img.Env {
		switch key {
		case "PATH":
			env["PATH"] = fmt.Sprintf("%s:%s", env["PATH"], value)
		default:
			env[key] = value
		}
	}

	envs := []string{}
	for key, value := range env {
		envs = append(envs, fmt.Sprintf("%s=%s", key, value))
	}

	if img.User != nil {
		metadata.Config.User = img.User.Name
		if img.User.Shell != "" {
			metadata.Config.Cmd = []string{img.User.Shell}
		}
	}

	metadata.RootFS.Type = "layers"
	metadata.Config.Env = envs

	return metadata
}
