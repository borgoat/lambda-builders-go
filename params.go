package lambdabuilders

// Params represents the parameters in the JSON-RPC call to lambda-builders
// https://github.com/aws/aws-lambda-builders/blob/develop/DESIGN.md#command-line-interface-internal
type Params struct {
	ProtocolVersion ParamsProtocolVersion  `json:"__protocol_version"`
	Capability      *ParamsCapability      `json:"capability"`
	SourceDir       string                 `json:"source_dir"`
	ArtifactsDir    string                 `json:"artifacts_dir"`
	ScratchDir      string                 `json:"scratch_dir"`
	ManifestPath    string                 `json:"manifest_path"`
	Runtime         string                 `json:"runtime"`
	Optimizations   map[string]interface{} `json:"optimizations"`
	Options         map[string]interface{} `json:"options"`
}

// ParamsCapability is the 3-tuple configuring a certain workflow to build Lambda functions
type ParamsCapability struct {
	Language             string `json:"language"`
	DependencyManager    string `json:"dependency_manager"`
	ApplicationFramework string `json:"application_framework"`
}

// ParamsProtocolVersion defines the versions of the JSON-RPC API exposed by lambda-builders
type ParamsProtocolVersion string

const (
	// ProtocolVersion03 is the only one we support as of today
	ProtocolVersion03 ParamsProtocolVersion = "0.3"
)
