package flag

import "github.com/spf13/pflag"

// Value extends pflag.Value with a method to
// obtain the target pointer of the registered flag and
// support for boolean-like behavior.
type Value interface {
	pflag.Value
	// Target returns the target pointer, used as key for looking up flags registered to a command.
	// See [github.com/neiser/go-nagini/command.Command.Flag].
	Target() any
	// IsBoolFlag returns true if the target points to a type with kind boolean.
	// See also Bool.
	IsBoolFlag() bool
}
