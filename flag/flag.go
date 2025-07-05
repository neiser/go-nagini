// Package flag adds flags to commands.
package flag

import (
	"fmt"
	"path"
	"reflect"

	"github.com/spf13/pflag"
)

// New constructs a new flag with a generic type as a target.
// The Parser converts from the given string value to the target type.
func New[T any](target *T, parser Parser[T]) pflag.Value {
	return anyValue[T]{target, parser}
}

type anyValue[T any] struct {
	target *T
	parser Parser[T]
}

func (v anyValue[T]) String() string {
	return convertToString(*v.target)
}

func (v anyValue[T]) Set(s string) (err error) {
	*v.target, err = v.parser(s)
	return
}

func (v anyValue[T]) Type() string {
	typeOf := reflect.TypeOf(*v.target)
	if _, pkgName := path.Split(typeOf.PkgPath()); len(pkgName) > 0 {
		return fmt.Sprintf("%s.%s", pkgName, typeOf.Name())
	}
	return typeOf.Name()
}

func convertToString[T any](t T) string {
	if stringer, ok := any(t).(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%v", t)
}
