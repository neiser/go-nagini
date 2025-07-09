package binding

import (
	"fmt"

	"github.com/neiser/go-nagini/flag"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Viper binds a command flag (given by a  [flag.Value] instance) to a ConfigKey for Viper.
// Implements flag.Binding.
type Viper struct {
	flag.Value

	ConfigKey string
}

// BindTo binds the flag of the command to a viper configuration value.
func (v Viper) BindTo() flag.Binder {
	if v.ConfigKey == "" {
		return nil
	}
	return func(flag *pflag.Flag) error {
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
		// the flag has not been changed on the cmd lint.
		// This makes a flag-set value precedent over viper values
		if configValuePresent && !flag.Changed {
			return v.setValueFromViper()
		}
		return nil
	}
}

func (v Viper) setValueFromViper() error {
	// If the current flag value, and we have something set from Viper,
	// we set the current value to the viper config value.
	if sliceValue, ok := v.Value.(pflag.SliceValue); ok {
		var configValue []string
		if err := viper.UnmarshalKey(v.ConfigKey, &configValue); err != nil {
			return fmt.Errorf("cannot unmarshal slice value of viper config key '%s': %w", v.ConfigKey, err)
		}
		if err := sliceValue.Replace(configValue); err != nil {
			return fmt.Errorf("cannot replace slice value to viper config %s='%s': %w", v.ConfigKey, configValue, err)
		}
	} else {
		var configValue string
		if err := viper.UnmarshalKey(v.ConfigKey, &configValue); err != nil {
			return fmt.Errorf("cannot unmarshal value of viper config key '%s': %w", v.ConfigKey, err)
		}
		if err := v.Set(configValue); err != nil {
			return fmt.Errorf("cannot set value to viper config %s='%s': %w", v.ConfigKey, configValue, err)
		}
	}
	return nil
}
