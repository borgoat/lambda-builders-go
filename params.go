package lambdabuilders

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

type ParamsCapability struct {
	Language             string `json:"language"`
	DependencyManager    string `json:"dependency_manager"`
	ApplicationFramework string `json:"application_framework"`
}

type ParamsProtocolVersion string

const (
	ProtocolVersion01 ParamsProtocolVersion = "0.1"
	ProtocolVersion02 ParamsProtocolVersion = "0.2"
	ProtocolVersion03 ParamsProtocolVersion = "0.3"
)
