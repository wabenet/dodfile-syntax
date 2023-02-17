package env

import (
	"fmt"
	"strings"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
)

const Type = "environment"

type Action struct {
	Variables map[string]string `mapstructure:"variables"`
}

func (a *Action) Type() string {
	return Type
}

func (a *Action) Execute(base llb.State) (llb.State, error) {
	return base, nil
}

func (a *Action) UpdateImage(i *dockerfile2llb.Image) {
	env := listToMap(i.Config.Env)

	for key, value := range a.Variables {
		switch key {
		case "PATH":
			env["PATH"] = fmt.Sprintf("%s:%s", value, env["PATH"])
		default:
			env[key] = value
		}
	}

	i.Config.Env = mapToList(env)
}

func listToMap(l []string) map[string]string {
	m := map[string]string{}

	for _, s := range l {
		vs := strings.SplitN(s, "=", 2)
		m[vs[0]] = vs[1]
	}

	return m
}

func mapToList(m map[string]string) []string {
	l := []string{}

	for k, v := range m {
		l = append(l, fmt.Sprintf("%s=%s", k, v))
	}

	return l
}
