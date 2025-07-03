package parameter

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// RegisterOptions are used when registering a parameter with a cobra.Command.
// Some properties here are straightforwardly set, but some others need the AfterRegistration below.
type RegisterOptions struct {
	Name, Shorthand, Usage string
	// Required forces this parameter to be present.
	Required bool
	// Persistent makes the parameter to be registered as a persistent flag.
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
// Note that 'value' parameter can be nil (happens when registering a simple Toggle parameter).
// See command.Command.
func (o RegisterOptions) AfterRegistration(cmd *cobra.Command, flag *pflag.Flag, value pflag.Value) {
	if binding, ok := value.(Binding); ok {
		binding.BindTo(cmd, flag)
	}
	if o.Required {
		_ = cmd.MarkFlagRequired(flag.Name)
	}
}

// Apply applies the given RegisterModifier's to this instance of RegisterOptions.
func (o RegisterOptions) Apply(modifiers ...RegisterModifier) RegisterOptions {
	for _, modifier := range modifiers {
		modifier(&o)
	}
	return o
}
