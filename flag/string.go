package flag

import (
	"fmt"
	"path"
	"reflect"
)

// String constructs a new flag with a generic type as a target.
// The Parser converts from the given string value to the target type.
// If the given Parser is nil, falls back to implementation of TargetParser on type *T,
// panics if nothing is found.
// See also Bool and Slice.
func String[T any](target *T, parser Parser[T]) Value {
	result := anyValue[T]{target: target, parser: parser}
	if parser == nil {
		var ok bool
		if result.targetParser, ok = any(target).(TargetParser); !ok {
			panic(fmt.Sprintf("flag for target %p (value '%v') must specify non-nil parser, "+
				"as flag.TargetParser interface is also not implemented", target, *target))
		}
	}
	return result
}

type anyValue[T any] struct {
	target       *T
	parser       Parser[T]
	targetParser TargetParser
}

//nolint:ireturn
func (v anyValue[T]) Target() any {
	return v.target
}

func (v anyValue[T]) String() string {
	return convertToString(*v.target)
}

//nolint:wrapcheck
func (v anyValue[T]) Set(s string) (err error) {
	if v.targetParser != nil {
		return v.targetParser.Parse(s)
	}
	*v.target, err = v.parser(s)
	return
}

func (v anyValue[T]) Type() string {
	if v.IsBoolFlag() {
		return ""
	}
	typeOf := reflect.TypeOf(*v.target)
	if _, pkgName := path.Split(typeOf.PkgPath()); len(pkgName) > 0 {
		return fmt.Sprintf("%s.%s", pkgName, typeOf.Name())
	}
	return typeOf.Name()
}

func (v anyValue[T]) IsBoolFlag() bool {
	typeOf := reflect.TypeOf(*v.target)
	return typeOf.Kind() == reflect.Bool
}

func convertToString[T any](t T) string {
	if stringer, ok := any(t).(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%v", t)
}
