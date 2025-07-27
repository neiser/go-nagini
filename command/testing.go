package command

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertExitCode is a ExecuteOption which uses WithExiter to assert the given exit code upon execution.
//
// This function is only useful for testing.
func AssertExitCode(t *testing.T, expectedExitCode int) ExecuteOption {
	t.Helper()
	hasRun := false
	t.Cleanup(func() {
		// defer check as we execute later (after ExecuteOption was applied)
		assert.Truef(t, hasRun, "exiter was not run")
	})
	return WithExiter(func(exitCode int) {
		hasRun = true
		assert.Equalf(t, expectedExitCode, exitCode, "exit code should be %d", expectedExitCode)
	})
}

// AssertWithRun installs a Run callback into the command to assert that it was called.
// Additional runs provided onRun if not nil.
// A possibly existing Run function is overwritten.
// Restores the previous possibly nil RunE function using [testing.T.Cleanup].
//
// This function is only useful for testing.
func AssertWithRun(t *testing.T, onRun func()) ExecuteOption {
	t.Helper()
	return ApplyToCommand(func(command Command) {
		// command.Run below modifies RunE, so restore it after test
		previousRunE := command.RunE
		t.Cleanup(func() {
			command.RunE = previousRunE
		})
		runCalled := false
		command.Run(func() (err error) {
			if onRun != nil {
				onRun()
			}
			runCalled = true
			return
		})
		t.Cleanup(func() {
			// defer check as Run is only executed upon Execute, which runs after this ExecuteOption callback
			assert.Truef(t, runCalled, "Command.Run callback not called")
		})
	})
}

// CaptureCobraOutput install buffers for Stdout and Stderr and
// provides suppliers once the command has been executed to obtain the output.
//
// This function is only useful for testing.
func (c Command) CaptureCobraOutput(t *testing.T) (getStdout, getStderr func() string) {
	t.Helper()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	c.SetOut(&stdout)
	c.SetErr(&stderr)
	t.Cleanup(func() {
		c.SetOut(os.Stdout)
		c.SetErr(os.Stderr)
	})
	return stdout.String, stderr.String
}
