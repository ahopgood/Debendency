package commands

import (
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

// Interface for executing a command
//
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o internal/fake_command.go . Command
type Command interface {
	Command(programName string, args ...string) (string, int, error)
}

type LinuxCommand struct{}

func (lc LinuxCommand) Command(programName string, args ...string) (string, int, error) {
	var output strings.Builder
	command := exec.Command(programName, args...)
	// command.Stdout = os.Stdout
	//command.Stderr = os.Stderr
	command.Stdout = &output
	command.Stderr = &output

	err := command.Start()

	var exit *exec.ExitError
	if errors.As(err, &exit) {
		output.Write(exit.Stderr)
		slog.Error(fmt.Sprintf("Standard Error: %s\n", output.String()))
		return output.String(), exit.ProcessState.ExitCode(), err
	}

	err = command.Wait()
	slog.Error(fmt.Sprintf("Wait error: %#v", err))
	if errors.As(err, &exit) {
		output.Write(exit.Stderr)
		slog.Error(fmt.Sprintf("Standard Error (exit code %d) from Wait: %s\n", exit.ExitCode(), output.String()))
		return output.String(), exit.ProcessState.ExitCode(), err
	}

	standardOut := output.String()
	return standardOut, command.ProcessState.ExitCode(), err
}
