package command

import (
	"errors"
	"log"
	"os"
)

// Execute executes the command using cobra and takes care of error handling.
// By default, exits the application with proper exit code and never returns.
// By default, logs an error originating from Run callback execution using [log.Printf].
// See  WithExiter and WithErrorLogger to change this default behavior (which can be useful for testing).
//
//nolint:wrapcheck
func (c Command) Execute(options ...ExecuteOption) (err error) {
	opts := executeOptions{
		Exiter: os.Exit,
		ErrorLogger: func(err error) {
			log.Printf("Command failed: %s", err.Error())
		},
	}.apply(options)

	if opts.Args != nil {
		c.SetArgs(opts.Args)
	}
	err = c.Command.Execute()

	if err != nil {
		exitCode := 1
		var errFromRunCallback fromRunCallbackError
		if errors.As(err, &errFromRunCallback) {
			var errWithExitCode WithExitCodeError
			if errors.As(err, &errWithExitCode) {
				exitCode = errWithExitCode.ExitCode
			}
			opts.ErrorLogger(errFromRunCallback.Wrapped)
		} else {
			// see cobra.Execute implementation, this mimics the behavior as if
			// SilenceErrors and SilenceUsage were false.
			c.PrintErrln(c.ErrPrefix(), err.Error())
			c.Println(c.UsageString())
		}
		opts.Exiter(exitCode)
	} else {
		opts.Exiter(0)
	}
	return err
}

// ExecuteOption is given to Command.Execute for modifying the default executeOptions.
type ExecuteOption func(options *executeOptions)

// WithArgs sets the used arguments before the command is executed.
// Providing this ExecuteOption without any arguments has the effect that [os.Args] is not considered by Cobra,
// which is useful for tests.
func WithArgs(args ...string) ExecuteOption {
	return func(options *executeOptions) {
		// when called with zero arguments,
		// explicitly set empty args array to make Cobra ignore os.Args (see above).
		if args == nil {
			args = []string{}
		}
		options.Args = args
	}
}

// WithExiter sets a different exit function.
// The default is [os.Exit] which makes Command.Execute never return.
func WithExiter(exiter func(exitCode int)) ExecuteOption {
	return func(options *executeOptions) {
		options.Exiter = exiter
	}
}

// WithErrorLogger uses the given logger for errors originating from executing Command.Run callbacks.
// By default, logs using [log.Printf].
func WithErrorLogger(logger func(err error)) ExecuteOption {
	return func(options *executeOptions) {
		options.ErrorLogger = logger
	}
}

// executeOptions are options for running Command.Execute.
// Use WithArgs, WithExiter, WithErrorLogger.
type executeOptions struct {
	Args        []string
	Exiter      func(exitCode int)
	ErrorLogger func(err error)
}

func (o executeOptions) apply(opts []ExecuteOption) executeOptions {
	for _, opt := range opts {
		opt(&o)
	}
	return o
}
