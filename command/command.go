// Package command constructs cobra.Command instances fluently starting from New.
package command

import (
	"iter"

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

	// commands holds the subcommands added via AddCommands below.
	// Uses a pointer to slice to enable Command value modification during AddCommands.
	// New ensures that the pointer-to-pointer is never nil.
	commands *[]Command
	// parent points to pointer of the parent if added via AddCommands below.
	// Otherwise, points to a nil value.
	// Uses a pointer to slice to enable Command value modification during AddCommands.
	// New ensures that the pointer-to-pointer is never nil.
	parent **Command

	// flagNames holds the registered flag names for a command,
	// identified by the target pointer as the map key.
	flagNames map[uintptr][]string
}

// New constructs a command.
func New() Command {
	var noParent *Command
	var noCommands []Command
	return Command{
		&cobra.Command{
			// We do our own usage output in Command.Execute below.
			SilenceErrors: true,
			SilenceUsage:  true,
		},
		&noCommands,
		&noParent,
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

// Flag registers a new flag.
// Use [flag.String], [flag.Bool] or [flag.Slice] to construct one and set appropriate [flag.RegisterOptions].
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
		*command.parent = &c
		*c.commands = append(*c.commands, command)
	}
	return c
}

// Run sets the given code to run during Execute and returned errors are logged.
// The error may implement WithExitCodeError and is wrapped in fromRunCallbackError.
func (c Command) Run(run func() error) Command {
	c.RunE = wrapRunCallbackError(run)
	return c
}

// AddPersistentPreRun adds the given code to run persistently, that means it will be executed
// also for sub commands added with AddCommands.
// Pre-runs before PreRun of command itself and the Run of the command.
func (c Command) AddPersistentPreRun(run func() error) Command {
	c.addToPersistentPreRunE(wrapRunCallbackError(run))
	return c
}

// All returns this command and all sub-commands added via AddCommands recursively as an iterator.
func (c Command) All() iter.Seq[Command] {
	return c.all
}

// Parents returns all parent commands up and including the root Command as an iterator.
func (c Command) Parents() iter.Seq[Command] {
	return c.parents
}
