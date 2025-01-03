package build

import (
	"fmt"

	"github.com/wabenet/dodfile-syntax/pkg/action"
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
	yaml "gopkg.in/yaml.v2"
)

type Image []action.Action

type Config struct {
	Actions map[string]interface{} `yaml:"actions"`
}

func ParseConfig(input []byte) (Image, error) {
	var cfg Config
	if err := yaml.Unmarshal(input, &cfg); err != nil {
		return nil, fmt.Errorf("invalid yaml syntax: %w", err)
	}

	actionsByType := map[string][]action.Action{
		base.Type:    {},
		copy.Type:    {},
		eget.Type:    {},
		env.Type:     {},
		fetch.Type:   {},
		install.Type: {},
		python.Type:  {},
		ruby.Type:    {},
		nodejs.Type:  {},
		golang.Type:  {},
		rust.Type:    {},
		script.Type:  {},
		user.Type:    {},
	}

	for name, value := range cfg.Actions {
		act, err := action.New(name, value)
		if err != nil {
			return nil, fmt.Errorf("could not decode action: %w", err)
		}

		acts := actionsByType[act.Type()]
		acts = append(acts, act)
		actionsByType[act.Type()] = acts
	}

	// TODO: implement something smarter to put the actions in the correct order
	// This list is currently hardcoded, so we have the exact same behavior
	// as before
	sorted := []action.Action{}
	sorted = append(sorted, actionsByType[base.Type]...)
	sorted = append(sorted, actionsByType[env.Type]...)
	sorted = append(sorted, actionsByType[user.Type]...)
	sorted = append(sorted, actionsByType[fetch.Type]...)
	sorted = append(sorted, actionsByType[eget.Type]...)
	sorted = append(sorted, actionsByType[copy.Type]...)
	sorted = append(sorted, actionsByType[python.Type]...)
	sorted = append(sorted, actionsByType[ruby.Type]...)
	sorted = append(sorted, actionsByType[nodejs.Type]...)
	sorted = append(sorted, actionsByType[golang.Type]...)
	sorted = append(sorted, actionsByType[rust.Type]...)
	sorted = append(sorted, actionsByType[install.Type]...)
	sorted = append(sorted, actionsByType[script.Type]...)

	return sorted, nil
}
