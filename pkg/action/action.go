package action

import (
	"errors"

	"github.com/go-viper/mapstructure/v2"
	"github.com/moby/buildkit/client/llb"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/wabenet/dodfile-syntax/pkg/action/base"
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
	"github.com/wabenet/dodfile-syntax/pkg/config"
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

func New(name string, cfg interface{}) (Action, error) {
	at, ac := decode(name, cfg)

	act, err := getByType(at)
	if err != nil {
		return nil, err
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:     &act,
		DecodeHook: config.TemplatingDecodeHook(),
	})
	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(ac); err != nil {
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
