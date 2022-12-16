package build_test

import (
	"os"
	"testing"

	"github.com/dodo-cli/dodfile-syntax/pkg/action/base"
	"github.com/dodo-cli/dodfile-syntax/pkg/action/copy"
	"github.com/dodo-cli/dodfile-syntax/pkg/action/download"
	"github.com/dodo-cli/dodfile-syntax/pkg/action/env"
	"github.com/dodo-cli/dodfile-syntax/pkg/action/install"
	"github.com/dodo-cli/dodfile-syntax/pkg/action/script"
	"github.com/dodo-cli/dodfile-syntax/pkg/action/user"
	"github.com/dodo-cli/dodfile-syntax/pkg/build"
	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	t.Parallel()

	dockerfile, err := os.ReadFile("test/dockerfile.yaml")
	assert.Nil(t, err)

	image, err := build.ParseConfig(dockerfile)
	assert.Nil(t, err)

	assert.Equal(t, 7, len(image))

	assert.Equal(t, &base.Action{Config: "debian"}, image[0])

	assert.Equal(t, &env.Action{Env: map[string]string{"PATH": "/usr/local/bin:$PATH"}}, image[1])

	assert.Equal(t, &download.Action{Config: []download.ActionConfig{{
		Source:      "https://files.example.com/test.zip",
		Unpack:      "test",
		Destination: "/bin/test",
	}}}, image[2])

	assert.Equal(t, &copy.Action{Config: []copy.ActionConfig{{
		Image: "test",
		Path:  "/some/file",
	}}}, image[3])

	assert.Equal(t, &install.Action{
		Name: []string{"test"},
		Repo: []string{"deb [arch=amd64] https://repo.example.com/ buster main"},
		Gpg:  []string{"https://repo.example.com/keys/test.asc"},
	}, image[4])

	assert.Equal(t, &user.Action{Name: "dodo", Dotfiles: "path/to/files"}, image[5])

	assert.Equal(t, &script.Action{Config: []script.ActionConfig{{
		Script: "echo Hello World",
	}}}, image[6])
}
