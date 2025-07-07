package flag

import "strconv"

// Bool provides a boolean flag.
func Bool(target *bool) Value {
	return New(target, strconv.ParseBool)
}
