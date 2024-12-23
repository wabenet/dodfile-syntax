package env

import (
	"fmt"

	"github.com/moby/buildkit/client/llb"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
)

const Type = "env"

type Action struct {
	Env map[string]string `mapstructure:",remain"`
}

func (a *Action) Type() string {
	return Type
}

func (a *Action) Execute(base llb.State) (llb.State, error) {
	return base, nil
}

func (a *Action) UpdateImage(config *oci.ImageConfig) {
	env := map[string]string{}

	for key, value := range a.Env {
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

	envs = append(config.Env, envs...)
	config.Env = envs
}
