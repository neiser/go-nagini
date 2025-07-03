package command

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/neiser/go-nagini/parameter"
)

func TestNew(t *testing.T) {
	type someType string

	t.Run("build simple command and run it", func(t *testing.T) {
		require.NoError(t, New("simple").runTest(t, []string{"simple"}, nil))
	})

	t.Run("bind optional parameter", func(t *testing.T) {
		var (
			someParameter someType
		)
		cmd := New("simple").
			Parameter(parameter.New(&someParameter, "", parameter.NotEmptyString), parameter.RegisterOptions{
				Name:      "some-param",
				Shorthand: "p",
			})

		// sub testcases modify state of someParameter, so run the "not set" case first

		t.Run("optional string param not set", func(t *testing.T) {
			require.NoError(t, cmd.runTest(t, []string{"simple"}, func() {
				require.Empty(t, someParameter)
			}))
		})

		t.Run("optional string param set", func(t *testing.T) {
			require.NoError(t, cmd.runTest(t, []string{"simple", "--some-param", "bla"}, func() {
				require.Equal(t, someType("bla"), someParameter)
			}))
		})

		t.Run("optional string param set with shorthand", func(t *testing.T) {
			require.NoError(t, cmd.runTest(t, []string{"simple", "-p", "blabla"}, func() {
				require.Equal(t, someType("blabla"), someParameter)
			}))
		})
	})

	t.Run("bind required parameter", func(t *testing.T) {
		var (
			someRequiredParam someType
		)
		cmd := New("simple").
			Parameter(parameter.New(&someRequiredParam, "", parameter.NotEmptyString), parameter.RegisterOptions{
				Name:     "some-required",
				Required: true,
			}).
			Run(func() error {
				return nil // dummy to make this cmd runnable
			})

		// sub testcases modify state of someRequiredParam, so run the "not set" case first

		t.Run("param not set", func(t *testing.T) {
			getStdout, getStderr := cmd.captureCobraOutput(t)

			err := cmd.runTestWithArgs(t, []string{"simple"}, func(exitCode int) {
				assert.Equal(t, 1, exitCode)
			})

			require.ErrorContains(t, err, `required flag(s) "some-required" not set`)
			assert.Contains(t, getStdout(), "Usage:")
			assert.Equal(t, `Error: required flag(s) "some-required" not set`+"\n", getStderr())
		})

		t.Run("param set", func(t *testing.T) {
			require.NoError(t, cmd.runTest(t, []string{"simple", "--some-required", "bla"}, func() {
				require.Equal(t, someType("bla"), someRequiredParam)
			}))
		})
	})

	t.Run("error handling and propagation", func(t *testing.T) {
		someError := errors.New("some error")
		t.Run("without exit code", func(t *testing.T) {
			cmd := New("simple").Run(func() error {
				return someError
			})
			getCobraStdout, getCobraStderr := cmd.captureCobraOutput(t)
			previousErrorLogger := ErrorLogger
			var capturedError error
			ErrorLogger = func(err error) {
				capturedError = err
			}
			t.Cleanup(func() {
				ErrorLogger = previousErrorLogger
			})

			err := cmd.runTestWithArgs(t, []string{"simple"}, func(exitCode int) {
				assert.Equal(t, 1, exitCode)
			})

			require.ErrorIs(t, err, someError)
			assert.Empty(t, getCobraStdout())
			assert.Empty(t, getCobraStderr())
			assert.Equal(t, someError, capturedError)
		})
		t.Run("with exit code", func(t *testing.T) {
			cmd := New("simple").Run(func() error {
				return WithExitCodeError{42, someError}
			})
			err := cmd.runTestWithArgs(t, []string{"simple"}, func(exitCode int) {
				assert.Equal(t, 42, exitCode)
			})
			require.ErrorIs(t, err, someError)
		})
	})
}

func (c Command) runTest(t *testing.T, args []string, onRun func()) error {
	t.Helper()
	haveRun := false
	c.Run(func() error {
		haveRun = true
		if onRun != nil {
			onRun()
		}
		return nil
	})
	err := c.runTestWithArgs(t, args, func(exitCode int) {
		assert.Equal(t, 0, exitCode, "exit code should be zero")
	})
	assert.True(t, haveRun)
	return err
}

func (c Command) runTestWithArgs(t *testing.T, args []string, exiter func(exitCode int)) error {
	t.Helper()
	c.SetArgs(args)
	return c.execute(exiter)
}

func (c Command) captureCobraOutput(t *testing.T) (getStdout, getStderr func() string) {
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
