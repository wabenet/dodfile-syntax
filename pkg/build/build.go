package build

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/util/system"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	defaultBaseImage = "debian"
)

func Build(ctx context.Context, c client.Client) (*client.Result, error) {
	img, err := GetConfig(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("getting config: %w", err)
	}

	st := llb.Image(defaultBaseImage)

	metadata := specs.Image{Platform: specs.Platform{Architecture: runtime.GOARCH, OS: "linux"}}
	metadata.RootFS.Type = "layers"
	metadata.Config.Env = []string{fmt.Sprintf("PATH=%s", system.DefaultPathEnv("unix"))}

	for _, a := range img {
		st, err = a.Execute(st)
		if err != nil {
			return nil, err
		}

		a.UpdateImage(&metadata.Config)
	}

	def, err := st.Marshal(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal local source: %w", err)
	}

	res, err := c.Solve(ctx, client.SolveRequest{
		Definition: def.ToPB(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to resolve dockerfile: %w", err)
	}

	ref, err := res.SingleRef()
	if err != nil {
		return nil, err
	}

	config, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal image config: %w", err)
	}

	res.AddMeta(exptypes.ExporterImageConfigKey, config)
	res.SetRef(ref)

	return res, nil
}

func GetConfig(ctx context.Context, c client.Client) (Image, error) {
	opts := c.BuildOpts().Opts

	filename := opts["filename"]
	if filename == "" {
		filename = "Dockerfile"
	}

	name := "load Dockerfile"
	if filename != "Dockerfile" {
		name += " from " + filename
	}

	src := llb.Local("dockerfile",
		llb.IncludePatterns([]string{filename}),
		llb.SessionID(c.BuildOpts().SessionID),
		llb.SharedKeyHint("Dockerfile"),
		//dockerfile2llb.WithInternalName(name),
	)

	def, err := src.Marshal(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal local source: %w", err)
	}

	res, err := c.Solve(ctx, client.SolveRequest{Definition: def.ToPB()})
	if err != nil {
		return nil, fmt.Errorf("failed to resolve dockerfile: %w", err)
	}

	ref, err := res.SingleRef()
	if err != nil {
		return nil, err
	}

	dtDockerfile, err := ref.ReadFile(ctx, client.ReadRequest{Filename: filename})
	if err != nil {
		return nil, fmt.Errorf("failed to read dockerfile: %w", err)
	}

	return ParseConfig(dtDockerfile)
}
