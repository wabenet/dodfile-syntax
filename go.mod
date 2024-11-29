module github.com/wabenet/dodfile-syntax

go 1.23

replace github.com/docker/docker => github.com/docker/docker v17.12.0-ce-rc1.0.20200310163718-4634ce647cf2+incompatible

replace github.com/containerd/containerd => github.com/containerd/containerd v1.3.1-0.20200512144102-f13ba8f2f2fd

require (
	github.com/mitchellh/mapstructure v1.5.0
	github.com/moby/buildkit v0.7.2
	github.com/opencontainers/image-spec v1.1.0-rc2
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/containerd/containerd v1.6.9 // indirect
	github.com/containerd/ttrpc v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.21+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/tonistiigi/fsutil v0.0.0-20220315205639-9ed612626da3 // indirect
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220502173005-c8bf987b8c21 // indirect
	google.golang.org/grpc v1.47.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
