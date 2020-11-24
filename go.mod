module github.com/dodo-cli/dodfile-syntax

go 1.15

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.3.1-0.20200512144102-f13ba8f2f2fd
	github.com/docker/docker => github.com/docker/docker v17.12.0-ce-rc1.0.20200310163718-4634ce647cf2+incompatible
	github.com/hashicorp/go-immutable-radix => github.com/tonistiigi/go-immutable-radix v0.0.0-20170803185627-826af9ccf0fe
	github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305
)

require (
	github.com/containerd/containerd v1.4.0-0
	github.com/docker/docker v0.0.0
	github.com/moby/buildkit v0.7.1
	github.com/opencontainers/image-spec v1.0.1
	gopkg.in/yaml.v2 v2.3.0
)
