package lambdabuilders

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/hashicorp/go-multierror"
	"github.com/powerman/rpc-codec/jsonrpc2"
)

// Client represents the entity that can launch lambda-builders and interface with it via JSON-RPC
type Client struct {
	lambdaBuildersPath string
}

type ClientOption func(client *Client) error

func WithLambdaBuilders(path string) ClientOption {
	return func(client *Client) error {
		client.lambdaBuildersPath = path
		return nil
	}
}

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

type ServiceMethod string

const (
	ServiceMethodBuild ServiceMethod = "LambdaBuilder.build"
)

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
