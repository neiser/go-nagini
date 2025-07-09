package flag

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Binding allows a command flag to be bound to other sources of configuration.
// This interface is detected while the flag is registered, see RegisterOptions.AfterRegistration.
// See for example [github.com/neiser/go-nagini/flag/binding.Viper].
type Binding interface {
	BindTo() Binder
}

// Binder is called during command execution to actually bind the flag.
// Binding is deferred to ensure that the flag value has been parsed properly.
type Binder func(flag *pflag.Flag) error

type cobraRunFuncPtr *func(cmd *cobra.Command, args []string) error

// addToPreRunE adds the given action to the command PreRunE phase.
// Using this phase is important as we need possibly flag values to be present.
// Otherwise, flags would not override values read in from the config file.
func addToCobraRun(cobraRunPtr cobraRunFuncPtr, action func(cmd *cobra.Command, args []string) error) {
	// important to catch existing as local variable,
	// as otherwise chaining the action callbacks leads to a stack overflow
	if existing := *cobraRunPtr; existing != nil {
		*cobraRunPtr = func(cmd *cobra.Command, args []string) error {
			if err := existing(cmd, args); err != nil {
				return err
			}
			return action(cmd, args)
		}
	} else {
		*cobraRunPtr = action
	}
}
