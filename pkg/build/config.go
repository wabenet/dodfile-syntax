package build

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/wabenet/dodfile-syntax/pkg/action"
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
	"github.com/wabenet/dodfile-syntax/pkg/action/install"
	"github.com/wabenet/dodfile-syntax/pkg/action/python"
	"github.com/wabenet/dodfile-syntax/pkg/action/script"
	"github.com/wabenet/dodfile-syntax/pkg/action/user"
	yaml "gopkg.in/yaml.v2"
)

type Image []action.Action

func ParseConfig(input []byte) (Image, error) {
	initial := map[string][]action.Action{
		cbase.Type:     {},
		ccopy.Type:     {},
		cdownload.Type: {},
		cenv.Type:      {},
		cinstall.Type:  {},
		cscript.Type:   {},
		base.Type:      {},
		copy.Type:      {},
		eget.Type:      {},
		env.Type:       {},
		fetch.Type:     {},
		install.Type:   {},
		python.Type:    {},
		script.Type:    {},
		user.Type:      {},
	}

	var cfg map[string]interface{}
	if err := yaml.Unmarshal(input, &cfg); err != nil {
		return nil, fmt.Errorf("invalid yaml syntax: %w", err)
	}

	actionsByType, err := sortActions(cfg, initial)
	if err != nil {
		return nil, err
	}

	// TODO: implement something smarter to put the actions in the correct order
	// This list is currently hardcoded, so we have the exact same behavior
	// as before
	sorted := []action.Action{}
	sorted = append(sorted, actionsByType[cbase.Type]...)
	sorted = append(sorted, actionsByType[base.Type]...)
	sorted = append(sorted, actionsByType[cenv.Type]...)
	sorted = append(sorted, actionsByType[env.Type]...)
	sorted = append(sorted, actionsByType[user.Type]...)
	sorted = append(sorted, actionsByType[cdownload.Type]...)
	sorted = append(sorted, actionsByType[fetch.Type]...)
	sorted = append(sorted, actionsByType[eget.Type]...)
	sorted = append(sorted, actionsByType[ccopy.Type]...)
	sorted = append(sorted, actionsByType[copy.Type]...)
	sorted = append(sorted, actionsByType[python.Type]...)
	sorted = append(sorted, actionsByType[cinstall.Type]...)
	sorted = append(sorted, actionsByType[install.Type]...)
	sorted = append(sorted, actionsByType[cscript.Type]...)
	sorted = append(sorted, actionsByType[script.Type]...)

	return sorted, nil
}

func sortActions(unsorted map[string]interface{}, actionsByType map[string][]action.Action) (map[string][]action.Action, error) {
	for name, value := range unsorted {
		if name == "actions" {
			subActions := map[string]interface{}{}
			if err := mapstructure.Decode(value, &subActions); err != nil {
				return nil, err
			}

			subSorted, err := sortActions(subActions, actionsByType)
			if err != nil {
				return nil, err
			}

			actionsByType = subSorted
		} else {
			act, err := action.New(name, value)
			if err != nil {
				return nil, fmt.Errorf("could not decode action: %w", err)
			}

			acts := actionsByType[act.Type()]
			acts = append(acts, act)
			actionsByType[act.Type()] = acts
		}
	}

	return actionsByType, nil
}
