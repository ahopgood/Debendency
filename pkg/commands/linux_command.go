package commands

import (
	"errors"
	"fmt"
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
	// command.Stderr = os.Stderr
	command.Stdout = &output
	// command.Stderr = &output

	err := command.Start()

	var exit *exec.ExitError
	if errors.As(err, &exit) {
		output.Write(exit.Stderr)
		fmt.Printf("Standard Error: %s\n", output.String())
		return output.String(), exit.ProcessState.ExitCode(), err
	}

	err = command.Wait()
	if errors.As(err, &exit) {
		output.Write(exit.Stderr)
		fmt.Printf("Standard Error from Wait: %s\n", output.String())
		return output.String(), exit.ProcessState.ExitCode(), err
	}

	standardOut := output.String()
	return standardOut, command.ProcessState.ExitCode(), err
}
