// Package command constructs cobra.Command instances fluently starting from New.
package command

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/neiser/go-nagini/parameter"
)

// A Command wraps cobra.Command and provides a fluent API to build CLI commands.
// Use New to construct one and starting fluently building.
// You can either use Run and register parameters using Parameter or Toggle.
// Or use SubCommands to build a command hierarchy.
// In any case, use Short and Long and LongParagraph to build the help message.
type Command struct {
	*cobra.Command
}

// New constructs a command with some use keyword.
func New(use string) Command {
	return Command{&cobra.Command{
		Use: use,
		// We do our own usage output in Command.Execute below.
		SilenceErrors: true,
		SilenceUsage:  true,
	}}
}

// Short sets the short command description.
func (c Command) Short(short string) Command {
	c.Command.Short = short
	return c
}

// Long adds the given sentence to the long command description (separated by newline).
func (c Command) Long(sentence string) Command {
	return c.long(sentence, "\n")
}

// LongParagraph adds the given sentence as a new paragraph to the long command description (preceded by two newlines).
func (c Command) LongParagraph(paragraph string) Command {
	return c.long(paragraph, "\n\n")
}

// long is internally used by Long and LongParagraph.
//
//nolint:funcorder
func (c Command) long(add, sep string) Command {
	c.Command.Long = strings.TrimSpace(strings.Join([]string{c.Command.Long, add}, sep))
	return c
}

// Parameter registers a new parameter.
// Use parameter.New to construct one and set appropriate parameter.RegisterOptions.
// The given param might implement parameter.Binding.
func (c Command) Parameter(param pflag.Value, options parameter.RegisterOptions) Command {
	flags := options.SelectFlags(c.Command)
	flag := flags.VarPF(param, options.Name, options.Shorthand, options.Usage)
	options.AfterRegistration(c.Command, flag, param)
	return c
}

// Toggle register a boolean parameter, which allows toggling a
// Set appropriate parameter.RegisterOptions!
func (c Command) Toggle(target *bool, defValue bool, options parameter.RegisterOptions) Command {
	flags := options.SelectFlags(c.Command)
	flags.BoolVarP(target, options.Name, options.Shorthand, defValue, options.Usage)
	options.AfterRegistration(c.Command, flags.Lookup(options.Name), nil)
	return c
}

// SubCommands registers children commands.
// This is used to build a hierarchy of commands.
func (c Command) SubCommands(commands ...Command) Command {
	for _, command := range commands {
		c.AddCommand(command.Command)
	}
	return c
}

// Run sets the given code to run during Execute and returned errors are logged.
// The error may implement WithExitCodeError and is wrapped in fromRunCallbackError.
func (c Command) Run(run func() error) Command {
	c.RunE = func(*cobra.Command, []string) error {
		if err := run(); err != nil {
			return fromRunCallbackError{err}
		}
		return nil
	}
	return c
}

// Execute executes the command using cobra and takes care of error handling.
// Note: This function never returns.
func (c Command) Execute() {
	_ = c.execute(os.Exit)
	panic("never returns")
}

func (c Command) execute(exiter func(exitCode int)) (err error) {
	err = c.Command.Execute()
	if err != nil {
		exitCode := 1
		var errFromRunCallback fromRunCallbackError
		if errors.As(err, &errFromRunCallback) {
			var errWithExitCode WithExitCodeError
			if errors.As(err, &errWithExitCode) {
				exitCode = errWithExitCode.ExitCode
			}
			ErrorLogger(errFromRunCallback.Wrapped)
		} else {
			// see cobra.Execute implementation, this mimics the behavior as if
			// SilenceErrors and SilenceUsage were false.
			c.PrintErrln(c.ErrPrefix(), err.Error())
			c.Println(c.UsageString())
		}
		exiter(exitCode)
	} else {
		exiter(0)
	}
	return
}
