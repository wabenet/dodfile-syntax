package build_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wabenet/dodfile-syntax/pkg/action/base"
	"github.com/wabenet/dodfile-syntax/pkg/action/copy"
	"github.com/wabenet/dodfile-syntax/pkg/action/env"
	"github.com/wabenet/dodfile-syntax/pkg/action/fetch"
	"github.com/wabenet/dodfile-syntax/pkg/action/install"
	"github.com/wabenet/dodfile-syntax/pkg/action/script"
	"github.com/wabenet/dodfile-syntax/pkg/action/user"
	"github.com/wabenet/dodfile-syntax/pkg/build"
)

func TestParseConfigNew(t *testing.T) {
	t.Parallel()

	dockerfile, err := os.ReadFile("test/dockerfile_new.yaml")
	assert.Nil(t, err)

	image, err := build.ParseConfig(dockerfile)
	assert.Nil(t, err)

	assert.Equal(t, 7, len(image))

	assert.Equal(t, &base.Action{Name: "debian"}, image[0])

	assert.Equal(t, &env.Action{Variables: map[string]string{"PATH": "/usr/local/bin:$PATH"}}, image[1])

	assert.Equal(t, &user.Action{Name: "dodo", Dotfiles: "path/to/files"}, image[2])

	assert.Equal(t, &fetch.Action{
		Source:      "https://files.example.com/test.zip",
		Unpack:      "test",
		Destination: "/bin/test",
	}, image[3])

	assert.Equal(t, &copy.Action{
		Image: "test",
		Path:  "/some/file",
	}, image[4])

	assert.Equal(t, &install.Action{
		Name: "test",
		Repo: "deb [arch=amd64] https://repo.example.com/ buster main",
		Gpg:  "https://repo.example.com/keys/test.asc",
	}, image[5])

	assert.Equal(t, &script.Action{
		Script: "echo Hello World",
	}, image[6])
}
