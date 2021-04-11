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

	c, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	b, err := c.NewBuilder("go", "modules", "")
	if err != nil {
		t.Fatal(err)
	}

	err = b.Build(
		".",
		"./out",
		scratchDir,
		"./go.mod",
		WithRuntime("go1.x"),
		WithOptions(map[string]interface{}{
			"artifact_executable_name": "aws-lambda-builders-go",
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
}
