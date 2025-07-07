// Package command constructs cobra.Command instances fluently starting from New.
package command

import (
	"errors"
	"os"
	"strings"

	"github.com/neiser/go-nagini/flag"
	"github.com/spf13/cobra"
)

// A Command wraps cobra.Command and provides a fluent API to build CLI commands.
// Use New to construct one and starting fluently building.
// You can either use Run and register parameters using Flag or FlagBool.
// Or use AddCommands to build a command hierarchy.
// In any case, use Short and Long and LongParagraph to build the help message.
type Command struct {
	*cobra.Command

	// flagNames holds the registered flag names for a command,
	// identified by the target pointer as the map key.
	flagNames map[uintptr][]string
}

// New constructs a command.
func New() Command {
	return Command{
		&cobra.Command{
			// We do our own usage output in Command.Execute below.
			SilenceErrors: true,
			SilenceUsage:  true,
		},
		map[uintptr][]string{},
	}
}

// Use specify the command usage, see cobra.Command#Use.
func (c Command) Use(use string) Command {
	c.Command.Use = use
	return c
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

// Flag registers a new flag.
// Use [flag.New], [flag.Bool] or [flag.Slice] to construct one and set appropriate [flag.RegisterOptions].
// The given param might implement [flag.Binding].
func (c Command) Flag(flagValue flag.Value, options flag.RegisterOptions) Command {
	flags := options.SelectFlags(c.Command)
	newFlag := flags.VarPF(flagValue, options.Name, options.Shorthand, options.Usage)
	c.addFlagName(flagValue.Target(), options.Name)
	options.AfterRegistration(c.Command, newFlag, flagValue)
	return c
}

// AddCommands registers children commands.
// This is used to build a hierarchy of commands.
func (c Command) AddCommands(commands ...Command) Command {
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
