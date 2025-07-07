package flag

import "strconv"

// Bool provides a boolean flag.
// It is a simple wrapper around String, and relies on the support by [Value.IsBoolFlag].
func Bool(target *bool) Value {
	return String(target, strconv.ParseBool)
}
