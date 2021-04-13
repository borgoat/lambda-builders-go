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

func ExampleBuilder_Build() error {
	c, err := NewClient()
	if err != nil {
		return err
	}

	b, err := c.NewBuilder("go", "modules", "")
	if err != nil {
		return err
	}

	err = b.Build(
		"/path/to/source",
		"/path/to/compile",
		"/path/to/scratchdir",
		"/path/to/source/go.mod",
		lambdabuilders.WithRuntime("go1.x"),
		lambdabuilders.WithOptions(map[string]interface{}{
			"artifact_executable_name": "my-handler",
		}),
	)
	return err
}
