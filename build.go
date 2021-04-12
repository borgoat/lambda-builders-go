package lambdabuilders

import (
	"fmt"
)

// Builder helps you build Lambda functions. It tries to mimic closely the Python implementation:
// https://github.com/aws/aws-lambda-builders/blob/develop/aws_lambda_builders/builder.py
// adapting it to Go idioms where needed.
//
// Create it NewBuilder
type Builder struct {
	client     *Client
	capability *ParamsCapability
}

// NewBuilder creates a new Builder using language, dependencyManager, and applicationFramework to define its capability.
// Workflows supported by lambda-builders may be found here:
// https://github.com/aws/aws-lambda-builders/tree/develop/aws_lambda_builders/workflows
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

// BuildOption is the functional interface to configure a Builder
type BuildOption func(params *Params) error

// WithRuntime defines the Lambda runtime to pass to the Builder
func WithRuntime(runtime string) BuildOption {
	return func(params *Params) error {
		params.Runtime = runtime
		return nil
	}
}

// WithOptions defines the custom options to pass to the Builder
func WithOptions(options map[string]interface{}) BuildOption {
	return func(params *Params) error {
		params.Options = options
		return nil
	}
}

// Build a Lambda function.
// Look at the Python implementation in lambda-builders for more info:
// https://github.com/aws/aws-lambda-builders/blob/165f92f35753d87e4abe1115fd2399826b371e1f/aws_lambda_builders/builder.py#L56-L67
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
