package flag

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// RegisterOptions are used when registering a flag with a cobra.Command and the [pflag.Flag].
// Some properties here are straightforwardly set, but some others need support in AfterRegistration below.
type RegisterOptions struct {
	// Name is the flag name (double dash prefix). This should always be set!
	Name string
	// Shorthand is an optional short flag (single dash prefix).
	Shorthand string
	// Usage describes how to use that flag.
	Usage string
	// Deprecated is shown as an alternative for this deprecated flag.
	Deprecated string
	// Hidden hides the flag from the usage help output.
	Hidden bool
	// Required forces this flag to be present.
	Required bool
	// Persistent makes the flag to be registered as a persistent flag.
	// The flag is then inherited to sub commands.
	// See RegisterOptions.SelectFlags
	Persistent bool
}

// RegisterModifier tweak RegisterOptions.
type RegisterModifier func(options *RegisterOptions)

// WithUsage sets the usage string as a RegisterModifier.
func WithUsage(usage string, args ...any) RegisterModifier {
	return func(options *RegisterOptions) {
		options.Usage = fmt.Sprintf(usage, args...)
	}
}

// Persistent makes the parameter to be registered as a persistent flag.
// The flag is then inherited to sub commands.
func Persistent() RegisterModifier {
	return func(options *RegisterOptions) {
		options.Persistent = true
	}
}

// SelectFlags is used when registering parameters. See command package.
func (o RegisterOptions) SelectFlags(cmd *cobra.Command) *pflag.FlagSet {
	if o.Persistent {
		return cmd.PersistentFlags()
	}
	return cmd.Flags()
}

// AfterRegistration is called after the parameter was registered.
// Note that 'value' parameter can be nil (happens when registering a simple FlagBool parameter).
// See command.Command.
func (o RegisterOptions) AfterRegistration(cmd *cobra.Command, flag *pflag.Flag, value Value) {
	if binding, ok := value.(Binding); ok {
		if binder := binding.BindTo(); binder != nil {
			action := func(*cobra.Command, []string) error {
				return binder(flag)
			}
			if o.Persistent {
				addToCobraRun(&cmd.PersistentPreRunE, action)
			} else {
				addToCobraRun(&cmd.PreRunE, action)
			}
		}
	}
	flag.Deprecated = o.Deprecated
	flag.Hidden = o.Hidden
	if o.Required {
		_ = cmd.MarkFlagRequired(flag.Name)
	}
	if value.IsBoolFlag() {
		flag.NoOptDefVal = "true"
	}
}

// Apply applies the given RegisterModifier's to this instance of RegisterOptions.
func (o RegisterOptions) Apply(modifiers ...RegisterModifier) RegisterOptions {
	for _, modifier := range modifiers {
		modifier(&o)
	}
	return o
}
