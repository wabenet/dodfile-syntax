package pythonapi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wabenet/dodfile-syntax/pkg/action/python/pythonapi"
)

func TestGetRelease(t *testing.T) {
	t.Parallel()

	release, err := pythonapi.GetDownload("3.11.5", pythonapi.Source)
	assert.Nil(t, err)

	assert.Equal(t, release.URL, "https://www.python.org/ftp/python/3.11.5/Python-3.11.5.tgz")
	assert.Equal(t, release.MD5Sum, "b628f21aae5e2c3006a12380905bb640")
}
