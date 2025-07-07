package command

import (
	"fmt"
	"reflect"
)

func getPointerValue(target any) uintptr {
	valueOf := reflect.ValueOf(target)
	if valueOf.Kind() != reflect.Pointer {
		panic(fmt.Sprintf("given target must be of type pointer, but is of type %s (value '%v')",
			reflect.TypeOf(target), target))
	}
	return uintptr(valueOf.UnsafePointer())
}

func (c Command) addFlagName(target any, flagName string) {
	key := getPointerValue(target)
	c.flagNames[key] = append(c.flagNames[key], flagName)
}

func (c Command) getFlagNames(targets []any) (result []string) {
	for _, target := range targets {
		flagNames, found := c.flagNames[getPointerValue(target)]
		if !found {
			panic(fmt.Sprintf("cannot find flag names for target %p=%+v", target, target))
		}
		result = append(result, flagNames...)
	}
	return
}

// MarkFlagsRequiredTogether exposes [github.com/spf13/cobra.Command.MarkFlagsRequiredTogether] fluently,
// accepting pointers to already registered flag values via Flag.
func (c Command) MarkFlagsRequiredTogether(targets ...any) Command {
	c.Command.MarkFlagsRequiredTogether(c.getFlagNames(targets)...)
	return c
}

// MarkFlagsOneRequired exposes [github.com/spf13/cobra.Command.MarkFlagsOneRequired] fluently,
// accepting pointers to already registered flag values via Flag.
func (c Command) MarkFlagsOneRequired(targets ...any) Command {
	c.Command.MarkFlagsOneRequired(c.getFlagNames(targets)...)
	return c
}

// MarkFlagsMutuallyExclusive exposes [github.com/spf13/cobra.Command.MarkFlagsMutuallyExclusive] fluently,
// accepting pointers to already registered flag values via Flag.
func (c Command) MarkFlagsMutuallyExclusive(targets ...any) Command {
	c.Command.MarkFlagsMutuallyExclusive(c.getFlagNames(targets)...)
	return c
}
