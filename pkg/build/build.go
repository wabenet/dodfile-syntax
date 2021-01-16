package build

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/moby/buildkit/frontend/gateway/client"
	yaml "gopkg.in/yaml.v2"
)

func Build(ctx context.Context, c client.Client) (*client.Result, error) {
	img, err := GetConfig(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("getting config: %w", err)
	}

	st, metadata := img.Build()

	def, err := st.Marshal()
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

func GetConfig(ctx context.Context, c client.Client) (*Image, error) {
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
		dockerfile2llb.WithInternalName(name),
	)

	def, err := src.Marshal()
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

	cfg := &Image{}
	if err := yaml.UnmarshalStrict(dtDockerfile, cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}
