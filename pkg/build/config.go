package build

import (
	"fmt"

	"github.com/wabenet/dodfile-syntax/pkg/action"
	"github.com/wabenet/dodfile-syntax/pkg/action/base"
	"github.com/wabenet/dodfile-syntax/pkg/action/copy"
	"github.com/wabenet/dodfile-syntax/pkg/action/download"
	"github.com/wabenet/dodfile-syntax/pkg/action/env"
	"github.com/wabenet/dodfile-syntax/pkg/action/install"
	"github.com/wabenet/dodfile-syntax/pkg/action/script"
	"github.com/wabenet/dodfile-syntax/pkg/action/user"
	yaml "gopkg.in/yaml.v2"
)

type Image []action.Action

func ParseConfig(input []byte) (Image, error) {
	actionsByType := map[string][]action.Action{
		base.Type:     {},
		copy.Type:     {},
		download.Type: {},
		env.Type:      {},
		install.Type:  {},
		script.Type:   {},
		user.Type:     {},
	}

	var cfg map[string]interface{}
	if err := yaml.Unmarshal(input, &cfg); err != nil {
		return nil, fmt.Errorf("invalid yaml syntax: %w", err)
	}

	for name, value := range cfg {
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
	sorted = append(sorted, actionsByType[download.Type]...)
	sorted = append(sorted, actionsByType[copy.Type]...)
	sorted = append(sorted, actionsByType[install.Type]...)
	sorted = append(sorted, actionsByType[user.Type]...)
	sorted = append(sorted, actionsByType[script.Type]...)

	return sorted, nil
}
