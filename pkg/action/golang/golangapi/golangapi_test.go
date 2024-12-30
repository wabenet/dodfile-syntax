package golangapi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wabenet/dodfile-syntax/pkg/action/golang/golangapi"
)

// TODO: test is not reproducible (esp. goarch)
func TestGetRelease(t *testing.T) {
	t.Parallel()

	release, err := golangapi.GetDownload("latest")
	assert.Nil(t, err)

	url, err := release.URL()
	assert.Nil(t, err)

	assert.Equal(t, "https://go.dev/dl/go1.23.4.linux-arm64.tar.gz", url)
	assert.Equal(t, "16e5017863a7f6071363782b1b8042eb12c6ca4f4cd71528b2123f0a1275b13e", release.ChecksumSHA256)
}
