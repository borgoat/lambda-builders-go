# Lambda Builders Go

[![Go Reference](https://pkg.go.dev/badge/github.com/borgoat/lambda-builders-go.svg)](https://pkg.go.dev/github.com/borgoat/lambda-builders-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/borgoat/lambda-builders-go)](https://goreportcard.com/report/github.com/borgoat/lambda-builders-go)
[![GitHub license](https://img.shields.io/github/license/borgoat/lambda-builders-go?color=yellow)](https://github.com/borgoat/lambda-builders-go/blob/development/LICENSE)

Go wrapper for [AWS Lambda Builders](https://github.com/aws/aws-lambda-builders),
using its JSON-RPC API.

> Lambda Builders is a Python library to compile, build and package AWS Lambda functions for several runtimes & frameworks.

## Getting Started

First, the lambda-builders executable is required:

```shell
$ pip install --user aws-lambda-builders
```

Now, download this module:

```shell
$ go get -u github.com/borgoat/lambda-builders-go
```

Here is an example how to use this library:

```go
package main

import lambdabuilders "github.com/borgoat/lambda-builders-go"

func main() {
	// Client executes lambda-builders and uses JSON-RPC to communicate with it
	client, err := lambdabuilders.NewClient()
	if err != nil {
		panic(err)
	}

	// Here we create a new builder for a Lambda written in Go, using Go Modules 
	// Other workflows may be found in the aws-lambda-builders Python library:
	// https://github.com/aws/aws-lambda-builders/tree/develop/aws_lambda_builders/workflows
	b, err := client.NewBuilder("go", "modules", "")
	if err != nil {
		panic(err)
	}

	// Finally we actually call LambdaBuilders.build
	// Check out how to configure it from the Python library:
	// https://github.com/aws/aws-lambda-builders/blob/165f92f35753d87e4abe1115fd2399826b371e1f/aws_lambda_builders/builder.py#L56-L67
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
	if err != nil {
		panic(err)
	}
}
```
