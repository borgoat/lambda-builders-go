package lambdabuilders

import (
	"os"
	"testing"
)

func TestBuild_go(t *testing.T) {
	scratchDir, err := os.MkdirTemp("", "aws-lambda-builders-go")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(scratchDir)

	params := &Params{
		ProtocolVersion: ProtocolVersion03,
		SourceDir:       ".",
		ArtifactsDir:    "./out",
		ScratchDir:      scratchDir,
		ManifestPath:    "./manifest",
		Runtime:         "go1.x",
		Capability:      ParamsCapability{
			Language:             "go",
			DependencyManager:    "modules",
			ApplicationFramework: "",
		},
		Options: map[string]interface{}{
			"artifact_executable_name": "aws-lambda-builders-go",
		},
	}

	err = Build(params)
	if err != nil {
		t.Fatal(err)
	}
}
