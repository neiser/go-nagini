// Package binding binds flags to external configuration systems, such as Viper.
package binding

import (
	"github.com/spf13/cobra"
)

// addToPreRunE adds the given action to the command PreRunE phase.
// Using this phase is important as we need possibly flag values to be present.
// Otherwise, flags would not override values read in from the config file.
func addToPreRunE(cmd *cobra.Command, action func(cmd *cobra.Command, args []string) error) {
	// important to catch existingPreRunE as local variable,
	// as otherwise chaining the action callbacks leads to a stack overflow
	if existingPreRunE := cmd.PreRunE; existingPreRunE != nil {
		cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
			if err := existingPreRunE(cmd, args); err != nil {
				return err
			}
			return action(cmd, args)
		}
	} else {
		cmd.PreRunE = action
	}
}
