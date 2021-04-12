package lambdabuilders

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/hashicorp/go-multierror"
	"github.com/powerman/rpc-codec/jsonrpc2"
)

// Client represents the entity that can launch lambda-builders and interface with it via JSON-RPC.
// Create it via NewClient
type Client struct {
	lambdaBuildersPath string
}

// ClientOption represents the functiona interface to customize Client
type ClientOption func(client *Client) error

// WithLambdaBuilders is a ClientOption to specify the path to the lambda-builders executable.
// This is optional: if it's not provided, the default PATH will be used to look it up.
func WithLambdaBuilders(path string) ClientOption {
	return func(client *Client) error {
		client.lambdaBuildersPath = path
		return nil
	}
}

// NewClient is used to create a new instance of a Client - to interface with lambda-builders
func NewClient(opts ...ClientOption) (*Client, error) {
	var c Client
	var err error

	for _, opt := range opts {
		err = opt(&c)
		if err != nil {
			return nil, fmt.Errorf("failure in setting option: %w", err)
		}
	}

	if c.lambdaBuildersPath == "" {
		c.lambdaBuildersPath, err = exec.LookPath("lambda-builders")
		if err != nil {
			return nil, fmt.Errorf("lambda-builders executable could not be found")
		}
	}

	return &c, nil
}

// rwCloserCmd is used to wrap exec.Cmd to expose a single io.ReadWriteCloser interface
type rwCloserCmd struct {
	io.WriteCloser
	io.ReadCloser
}

func newReadWriteCloserFromCmd(cmd *exec.Cmd) (*rwCloserCmd, error) {
	var err, result error

	stdin, err := cmd.StdinPipe()
	if err != nil {
		result = multierror.Append(result, err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		result = multierror.Append(result, err)
	}

	rwc := &rwCloserCmd{stdin, stdout}

	return rwc, result
}

func (rwc *rwCloserCmd) Close() error {
	var result error

	result = multierror.Append(result, rwc.WriteCloser.Close())
	result = multierror.Append(result, rwc.ReadCloser.Close())

	return result
}

// ServiceMethod is used to define the lambda-builders method to be called.
type ServiceMethod string

const (
	// ServiceMethodBuild may be used to call "LambdaBuilder.build"
	// As of today, this is the only available method exposed by lambda-builders:
	// https://github.com/aws/aws-lambda-builders/blob/165f92f35753d87e4abe1115fd2399826b371e1f/aws_lambda_builders/__main__.py#L90-L92
	ServiceMethodBuild ServiceMethod = "LambdaBuilder.build"
)

// GenericCall may be used for any call to lambda-builders
func (c *Client) GenericCall(serviceMethod ServiceMethod, args, reply interface{}) error {
	lambdaBuilders := exec.Command(c.lambdaBuildersPath)

	rwc, err := newReadWriteCloserFromCmd(lambdaBuilders)
	if err != nil {
		return fmt.Errorf("failed to set up stdin/stdout: %w", err)
	}

	cli := jsonrpc2.NewClient(rwc)

	call := cli.Go(string(serviceMethod), args, reply, nil)
	err = lambdaBuilders.Start()
	if err != nil {
		return fmt.Errorf("failed to launch lambda-builders: %w", err)
	}

	err = rwc.WriteCloser.Close()
	if err != nil {
		return fmt.Errorf("failed to close stdin: %w", err)
	}

	// Wait for lambda-builders to return
	<-call.Done

	if call.Error != nil {
		return fmt.Errorf("lambda-builders RPC returned an error: %w", call.Error)
	}

	return nil
}
