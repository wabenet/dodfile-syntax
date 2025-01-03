package action

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/moby/buildkit/client/llb"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/wabenet/dodfile-syntax/pkg/action/base"
	cbase "github.com/wabenet/dodfile-syntax/pkg/action/compat/base"
	ccopy "github.com/wabenet/dodfile-syntax/pkg/action/compat/copy"
	cdownload "github.com/wabenet/dodfile-syntax/pkg/action/compat/download"
	cenv "github.com/wabenet/dodfile-syntax/pkg/action/compat/env"
	cinstall "github.com/wabenet/dodfile-syntax/pkg/action/compat/install"
	cscript "github.com/wabenet/dodfile-syntax/pkg/action/compat/script"
	"github.com/wabenet/dodfile-syntax/pkg/action/copy"
	"github.com/wabenet/dodfile-syntax/pkg/action/eget"
	"github.com/wabenet/dodfile-syntax/pkg/action/env"
	"github.com/wabenet/dodfile-syntax/pkg/action/fetch"
	"github.com/wabenet/dodfile-syntax/pkg/action/golang"
	"github.com/wabenet/dodfile-syntax/pkg/action/install"
	"github.com/wabenet/dodfile-syntax/pkg/action/nodejs"
	"github.com/wabenet/dodfile-syntax/pkg/action/python"
	"github.com/wabenet/dodfile-syntax/pkg/action/ruby"
	"github.com/wabenet/dodfile-syntax/pkg/action/rust"
	"github.com/wabenet/dodfile-syntax/pkg/action/script"
	"github.com/wabenet/dodfile-syntax/pkg/action/user"
)

type Action interface {
	Type() string
	Execute(llb.State) (llb.State, error)
	UpdateImage(*oci.ImageConfig)
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
	case cbase.Type:
		return &cbase.Action{}, nil
	case ccopy.Type:
		return &ccopy.Action{}, nil
	case cdownload.Type:
		return &cdownload.Action{}, nil
	case cenv.Type:
		return &cenv.Action{}, nil
	case cinstall.Type:
		return &cinstall.Action{}, nil
	case cscript.Type:
		return &cscript.Action{}, nil
	case base.Type:
		return &base.Action{}, nil
	case copy.Type:
		return &copy.Action{}, nil
	case eget.Type:
		return &eget.Action{}, nil
	case env.Type:
		return &env.Action{}, nil
	case fetch.Type:
		return &fetch.Action{}, nil
	case golang.Type:
		return &golang.Action{}, nil
	case install.Type:
		return &install.Action{}, nil
	case python.Type:
		return &python.Action{}, nil
	case ruby.Type:
		return &ruby.Action{}, nil
	case nodejs.Type:
		return &nodejs.Action{}, nil
	case rust.Type:
		return &rust.Action{}, nil
	case script.Type:
		return &script.Action{}, nil
	case user.Type:
		return &user.Action{}, nil
	default:
		return nil, errors.New("Unknown action")
	}
}
