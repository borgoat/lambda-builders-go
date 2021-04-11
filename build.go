package lambdabuilders

import (
	"fmt"
)

type Builder struct {
	client     *Client
	capability *ParamsCapability
}

func (c *Client) NewBuilder(language, dependencyManager, applicationFramework string) (*Builder, error) {
	b := &Builder{
		client: c,
		capability: &ParamsCapability{
			Language:             language,
			DependencyManager:    dependencyManager,
			ApplicationFramework: applicationFramework,
		},
	}

	return b, nil
}

type BuildOption func(params *Params) error

func WithRuntime(runtime string) BuildOption {
	return func(params *Params) error {
		params.Runtime = runtime
		return nil
	}
}

func WithOptions(options map[string]interface{}) BuildOption {
	return func(params *Params) error {
		params.Options = options
		return nil
	}
}

func (b *Builder) Build(sourceDir, artifactsDir, scratchDir, manifestPath string, opts ...BuildOption) error {
	params := &Params{
		ProtocolVersion: ProtocolVersion03,
		SourceDir:       sourceDir,
		ArtifactsDir:    artifactsDir,
		ScratchDir:      scratchDir,
		ManifestPath:    manifestPath,
		Capability:      b.capability,
	}

	for _, opt := range opts {
		err := opt(params)
		if err != nil {
			return fmt.Errorf("failed to set custom option: %w", err)
		}
	}

	err := b.client.GenericCall(ServiceMethodBuild, params, nil)
	if err != nil {
		return fmt.Errorf("failed to perform call: %w", err)
	}

	return nil
}
