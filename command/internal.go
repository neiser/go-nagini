package command

import (
	"strings"

	"github.com/spf13/cobra"
)

// long is internally used by Long and LongParagraph.
func (c Command) long(add, sep string) Command {
	c.Command.Long = strings.TrimSpace(strings.Join([]string{c.Command.Long, add}, sep))
	return c
}

func wrapRunCallbackError(run func() error) func(*cobra.Command, []string) error {
	return func(*cobra.Command, []string) error {
		if err := run(); err != nil {
			return fromRunCallbackError{err}
		}
		return nil
	}
}

func (c Command) addToPersistentPreRunE(action func(*cobra.Command, []string) error) {
	// important to catch existing as local variable,
	// as otherwise chaining the action callbacks leads to a stack overflow
	if existing := c.PersistentPreRunE; existing != nil {
		c.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			if err := existing(cmd, args); err != nil {
				return err
			}
			return action(cmd, args)
		}
	} else {
		c.PersistentPreRunE = action
	}
}
