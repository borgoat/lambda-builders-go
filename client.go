package lambdabuilders

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/hashicorp/go-multierror"
	"github.com/powerman/rpc-codec/jsonrpc2"
)

type rwCloserCmd struct {
	io.WriteCloser
	io.ReadCloser
}

func NewReadWriteCloserFromCmd(cmd *exec.Cmd) (*rwCloserCmd, error) {
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

func GenericCall(serviceMethod string, args, reply interface{}) error {
	lambdaBuilders := exec.Command("lambda-builders") // TODO: Support env variable to customise cmd path

	rwc, err := NewReadWriteCloserFromCmd(lambdaBuilders)
	if err != nil {
		return fmt.Errorf("failed to set up stdin/stdout: %w", err)
	}

	cli := jsonrpc2.NewClient(rwc)

	call := cli.Go("LambdaBuilder.build", args, reply, nil)
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
