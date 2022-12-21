package action

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/wabenet/dodfile-syntax/pkg/action/base"
	"github.com/wabenet/dodfile-syntax/pkg/action/copy"
	"github.com/wabenet/dodfile-syntax/pkg/action/download"
	"github.com/wabenet/dodfile-syntax/pkg/action/env"
	"github.com/wabenet/dodfile-syntax/pkg/action/install"
	"github.com/wabenet/dodfile-syntax/pkg/action/script"
	"github.com/wabenet/dodfile-syntax/pkg/action/user"
)

type Action interface {
	Type() string
	Execute(llb.State) (llb.State, error)
	UpdateImage(*dockerfile2llb.Image)
}

type actionConfig struct {
	ID     string                 `mapstructure:"id"`
	Type   string                 `mapstructure:"type"`
	Config map[string]interface{} `mapstructure:",remain"`
}

func New(name string, config interface{}) (Action, error) {
	at, ac := decode(name, config)

	act, err := getByType(at)
	if err != nil {
		return nil, err
	}

	if err := mapstructure.Decode(ac, &act); err != nil {
		return nil, err
	}

	return act, nil
}

func decode(key string, value interface{}) (string, map[string]interface{}) {
	ac := actionConfig{}
	if err := mapstructure.Decode(value, &ac); err != nil {
		return key, map[string]interface{}{
			"config": value,
		}
	}

	if t := ac.Type; t != "" {
		return t, ac.Config
	}

	if t := ac.ID; t != "" {
		return t, ac.Config
	}

	return key, ac.Config
}

func getByType(t string) (Action, error) {
	switch t {
	case base.Type:
		return &base.Action{}, nil
	case user.Type:
		return &user.Action{}, nil
	case env.Type:
		return &env.Action{}, nil
	case install.Type:
		return &install.Action{}, nil
	case download.Type:
		return &download.Action{}, nil
	case copy.Type:
		return &copy.Action{}, nil
	case script.Type:
		return &script.Action{}, nil
	default:
		return nil, errors.New("Unknown action")
	}
}
