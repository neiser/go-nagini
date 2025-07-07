package command

import (
	"bytes"
	"errors"
	"os"
	"strconv"
	"testing"

	"github.com/neiser/go-nagini/flag"
	"github.com/neiser/go-nagini/flag/binding"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	type someType string

	t.Run("build simple command and run it", func(t *testing.T) {
		require.NoError(t, New().runTest(t, []string{}, nil))
	})

	t.Run("add optional flag", func(t *testing.T) {
		var (
			someVal someType
		)
		cmd := New().
			Flag(flag.New(&someVal, flag.NotEmptyTrimmed), flag.RegisterOptions{
				Name:      "some-val",
				Shorthand: "p",
			})

		// sub testcases modify state of someVal, so run the "not set" case first

		t.Run("optional string flag not set", func(t *testing.T) {
			require.NoError(t, cmd.runTest(t, []string{}, func() {
				require.Empty(t, someVal)
			}))
		})

		t.Run("optional string flag set", func(t *testing.T) {
			require.NoError(t, cmd.runTest(t, []string{"--some-val", "bla"}, func() {
				require.Equal(t, someType("bla"), someVal)
			}))
		})

		t.Run("optional string flag set with shorthand", func(t *testing.T) {
			require.NoError(t, cmd.runTest(t, []string{"-p", "blabla"}, func() {
				require.Equal(t, someType("blabla"), someVal)
			}))
		})
	})

	t.Run("add required flag", func(t *testing.T) {
		var (
			someRequiredVal someType
		)
		cmd := New().
			Flag(flag.New(&someRequiredVal, flag.NotEmptyTrimmed), flag.RegisterOptions{
				Name:     "some-required",
				Required: true,
			}).
			Run(func() error {
				return nil // dummy to make this cmd runnable
			})

		// sub testcases modify state of someRequiredVal, so run the "not set" case first

		t.Run("param not set", func(t *testing.T) {
			getStdout, getStderr := cmd.captureCobraOutput(t)

			err := cmd.runTestWithArgs(t, []string{}, func(exitCode int) {
				assert.Equal(t, 1, exitCode)
			})

			require.ErrorContains(t, err, `required flag(s) "some-required" not set`)
			assert.Contains(t, getStdout(), "Usage:")
			assert.Equal(t, `Error: required flag(s) "some-required" not set`+"\n", getStderr())
		})

		t.Run("param set", func(t *testing.T) {
			require.NoError(t, cmd.runTest(t, []string{"--some-required", "bla"}, func() {
				require.Equal(t, someType("bla"), someRequiredVal)
			}))
		})
	})

	t.Run("bool flag", func(t *testing.T) {
		var (
			someBool bool
		)
		cmd := New().
			Flag(flag.Bool(&someBool), flag.RegisterOptions{
				Name: "some-bool",
			}).
			Run(func() error {
				return nil // dummy to make this cmd runnable
			})
		require.NoError(t, cmd.runTest(t, []string{"--some-bool"}, func() {
			require.True(t, someBool)
		}))
	})

	t.Run("error handling and propagation", func(t *testing.T) {
		someError := errors.New("some error")
		t.Run("without exit code", func(t *testing.T) {
			cmd := New().Run(func() error {
				return someError
			})
			getCobraStdout, getCobraStderr := cmd.captureCobraOutput(t)
			getError := captureError(t)

			err := cmd.runTestWithArgs(t, []string{}, func(exitCode int) {
				assert.Equal(t, 1, exitCode)
			})

			require.ErrorIs(t, err, someError)
			assert.Empty(t, getCobraStdout())
			assert.Empty(t, getCobraStderr())
			assert.Equal(t, someError, getError())
		})
		t.Run("with exit code", func(t *testing.T) {
			cmd := New().Run(func() error {
				return WithExitCodeError{42, someError}
			})
			err := cmd.runTestWithArgs(t, []string{}, func(exitCode int) {
				assert.Equal(t, 42, exitCode)
			})
			require.ErrorIs(t, err, someError)
		})
	})

	t.Run("command hierarchy", func(t *testing.T) {
		cmd := New().
			AddCommands(
				New().Use("sub1").Short("First subcommand").Run(func() error {
					return WithExitCodeError{ExitCode: 42}
				}),
				New().Use("sub2").Short("Second subcommand").Run(func() error {
					return WithExitCodeError{ExitCode: 43}
				}),
			)
		t.Run("no subcommand", func(t *testing.T) {
			getStdout, getStderr := cmd.captureCobraOutput(t)
			require.NoError(t, cmd.runTestWithArgs(t, []string{}, func(exitCode int) {
				assert.Equal(t, 0, exitCode)
			}))
			stdout := getStdout()
			assert.Contains(t, stdout, "Usage:")
			assert.Contains(t, stdout, "First subcommand")
			assert.Contains(t, stdout, "Second subcommand")
			assert.Empty(t, getStderr())
		})

		t.Run("run sub1", func(t *testing.T) {
			getStdout, getStderr := cmd.captureCobraOutput(t)
			getError := captureError(t)
			require.ErrorContains(t, cmd.runTestWithArgs(t, []string{"sub1"}, func(exitCode int) {
				assert.Equal(t, 42, exitCode)
			}), "exit code 42")
			assert.Empty(t, getStdout())
			assert.Empty(t, getStderr())
			assert.Error(t, getError())
		})

		t.Run("run sub2", func(t *testing.T) {
			getStdout, getStderr := cmd.captureCobraOutput(t)
			getError := captureError(t)
			require.ErrorContains(t, cmd.runTestWithArgs(t, []string{"sub2"}, func(exitCode int) {
				assert.Equal(t, 43, exitCode)
			}), "exit code 43")
			assert.Empty(t, getStdout())
			assert.Empty(t, getStderr())
			assert.Error(t, getError())
		})
	})

	t.Run("viper binding", func(t *testing.T) {
		viper.AutomaticEnv()

		t.Run("scalar value", func(t *testing.T) {
			var (
				someVal string
			)
			cmd := New().Flag(
				binding.Viper{
					Value:     flag.New(&someVal, flag.NotEmptyTrimmed),
					ConfigKey: "SOME_VAL",
				},
				flag.RegisterOptions{
					Name: "some-val",
				},
			)

			// sub testcases modify state of someVal, so run the "not set" case first

			t.Run("flag not set, env not set", func(t *testing.T) {
				require.NoError(t, cmd.runTest(t, []string{}, func() {
					require.Empty(t, someVal)
				}))
			})

			t.Setenv("SOME_VAL", "value-from-env")

			t.Run("flag not set, but env set", func(t *testing.T) {
				require.NoError(t, cmd.runTest(t, []string{}, func() {
					require.Equal(t, "value-from-env", someVal)
				}))
			})

			t.Run("flag is preferred over env", func(t *testing.T) {
				require.NoError(t, cmd.runTest(t, []string{"--some-val", "blabla"}, func() {
					require.Equal(t, "blabla", someVal)
				}))
			})
		})

		t.Run("slice of ints", func(t *testing.T) {
			type (
				sliceOfInts []int
			)

			var (
				someInts sliceOfInts
				someBool yesOrNoType = true
			)
			cmd := New().
				Flag(
					binding.Viper{
						Value:     flag.NewSlice(&someInts, flag.ParseSliceOf[int](strconv.Atoi)),
						ConfigKey: "SOME_INTEGERS",
					},
					flag.RegisterOptions{Name: "some-ints"},
				).
				Flag(
					binding.Viper{
						Value:     flag.New(&someBool, nil),
						ConfigKey: "SOME_BOOL",
					},
					flag.RegisterOptions{Name: "some-bool"},
				).
				Run(func() error {
					// have some dummy to always run through PreRunE
					return nil
				})

			// sub testcases modify state of someInts, so run the "not set" case first

			t.Run("flag not set, env not set", func(t *testing.T) {
				require.NoError(t, cmd.runTest(t, []string{}, func() {
					require.Empty(t, someInts)
					require.True(t, bool(someBool))
				}))
			})

			t.Setenv("SOME_INTEGERS", "2,x3x,4")

			t.Run("env parsing fails, help message properly shown", func(t *testing.T) {
				getStdout, getStderr := cmd.captureCobraOutput(t)
				require.ErrorContains(t, cmd.runTestWithArgs(t, []string{}, func(exitCode int) {
					require.Equal(t, 1, exitCode)
				}), `cannot replace slice value to viper config SOME_INTEGERS='[2 x3x 4]': cannot parse slice element 1: strconv.Atoi: parsing "x3x": invalid syntax`)
				assert.Contains(t, getStderr(), "Error: cannot replace slice value to viper config")
				stdout := getStdout()
				assert.Contains(t, stdout, "--some-bool[=true]    (default true)")
				assert.Contains(t, stdout, "--some-ints []int     (default <nil>)")
			})

			t.Setenv("SOME_INTEGERS", "2,3,4")
			t.Setenv("SOME_BOOL", "no")

			t.Run("flag not set, but env set", func(t *testing.T) {
				require.NoError(t, cmd.runTest(t, []string{}, func() {
					require.Equal(t, sliceOfInts{2, 3, 4}, someInts)
					require.False(t, bool(someBool))
				}))
			})

			t.Run("flag is preferred over env", func(t *testing.T) {
				require.NoError(t, cmd.runTest(t, []string{"--some-ints", "5,6,7", "--some-bool=yes"}, func() {
					require.Equal(t, sliceOfInts{5, 6, 7}, someInts)
					require.True(t, bool(someBool))
				}))
			})

			t.Run("flag parsing fails", func(t *testing.T) {
				cmd.captureCobraOutput(t) // just to silence confusing error output during tests
				require.ErrorContains(t, cmd.runTestWithArgs(t, []string{"--some-ints", "5,x6x,7"}, func(exitCode int) {
					require.Equal(t, 1, exitCode)
				}), `invalid argument "5,x6x,7" for "--some-ints" flag: cannot parse slice element 1: strconv.Atoi: parsing "x6x": invalid syntax`)
			})
		})
	})
}

func (c Command) runTest(t *testing.T, args []string, onRun func()) error {
	t.Helper()
	haveRun := false
	previousRunE := c.RunE
	t.Cleanup(func() {
		c.RunE = previousRunE
	})
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

func captureError(t *testing.T) (getError func() error) {
	t.Helper()
	previousErrorLogger := ErrorLogger
	var capturedError error
	ErrorLogger = func(err error) {
		capturedError = err
	}
	t.Cleanup(func() {
		ErrorLogger = previousErrorLogger
	})
	return func() error {
		return capturedError
	}
}

type yesOrNoType bool

func (v *yesOrNoType) Parse(s string) error {
	switch s {
	case "yes":
		*v = true
	case "no":
		*v = false
	default:
		vAsBool, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		*v = yesOrNoType(vAsBool)
	}
	return nil
}
