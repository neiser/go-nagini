package flag

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Binding allows a command flag to be bound to other sources of configuration.
// This interface is detected while the flag is registered, see RegisterOptions.AfterRegistration.
// See for example ViperBinding.
type Binding interface {
	BindTo(cmd *cobra.Command, flag *pflag.Flag)
}

// ViperBinding is a helper to bind a command flag (represented as a pflag.Value) to a config file ConfigKey.
type ViperBinding struct {
	pflag.Value

	ConfigKey string
}

// BindTo binds the flag in the PreRunE phase of the command to viper.
// Using this phase is important as we need possibly flag values to be present.
// Otherwise, flags would not override values read in from the config file.
func (v ViperBinding) BindTo(cmd *cobra.Command, flag *pflag.Flag) {
	addToPreRunE(cmd, func(*cobra.Command, []string) error {
		// Check if viper has a config value before binding the flag,
		// as otherwise the config value would always be reported as present
		// (value source would then always be the bound flag)
		configValuePresent := false
		if configValue := viper.Get(v.ConfigKey); configValue != nil {
			configValuePresent = true
		}
		if err := viper.BindPFlag(v.ConfigKey, flag); err != nil {
			return fmt.Errorf("cannot bind value to viper: %w", err)
		}
		// Only set value from viper if the value is actually present
		// which makes a flag-set value precedent over viper values
		if !configValuePresent {
			return nil
		}
		return v.setValueFromViper()
	})
}

func (v ViperBinding) setValueFromViper() error {
	// If the current flag value, and we have something set from Viper,
	// we set the current value to the viper config value.
	if sliceValue, ok := v.Value.(pflag.SliceValue); ok {
		// checking for nil only allows defining an empty slice [] as a valid (default) value for the flag
		if sliceValue.GetSlice() == nil {
			var configValue []string
			if err := viper.UnmarshalKey(v.ConfigKey, &configValue); err != nil {
				return fmt.Errorf("cannot unmarshal slice value of viper config key '%s': %w", v.ConfigKey, err)
			}
			if err := sliceValue.Replace(configValue); err != nil {
				return fmt.Errorf("cannot set slice value to viper config %s='%s': %w", v.ConfigKey, configValue, err)
			}
		}
	} else if strings.TrimSpace(v.String()) == "" {
		var configValue string
		if err := viper.UnmarshalKey(v.ConfigKey, &configValue); err != nil {
			return fmt.Errorf("cannot unmarshal value of viper config key '%s': %w", v.ConfigKey, err)
		}
		if err := v.Set(configValue); err != nil {
			return fmt.Errorf("cannot set value to viper config value %s='%s': %w", v.ConfigKey, configValue, err)
		}
	}
	return nil
}

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
