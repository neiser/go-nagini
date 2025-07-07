package flag

import (
	"errors"
	"fmt"
	"strings"
)

// Parser is used in New to convert from string to the generic type T.
// Note that this signature matches standard conversion functions, such as [strconv.Atoi].
// Implementations of Parser may use ErrParser when an error occurs.
type Parser[T any] func(s string) (T, error)

// TargetParser can be implemented by given target pointers
// Implementations of TargetParser may use ErrParser when an error occurs.
type TargetParser interface {
	Parse(s string) error
}

// ErrParser can be used when implementations of Parser fail.
var ErrParser = errors.New("cannot parse parameter")

// NotEmptyTrimmed is a commonly used Parser and requires that the given string
// is not empty after whitespace was trimmed.
func NotEmptyTrimmed[T ~string](s string) (T, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed != "" {
		return AnyString[T](trimmed)
	}
	return "", fmt.Errorf("%w: value '%s' after trimming is empty", ErrParser, s)
}

// NotEmpty is a commonly used Parser and requires that the given string
// is not empty without trimming whitespace.
func NotEmpty[T ~string](s string) (T, error) {
	if s != "" {
		return AnyString[T](s)
	}
	return "", fmt.Errorf("%w: value '%s' is empty", ErrParser, s)
}

// AnyString is a Parser which returns any string as type T.
// Prefer using NotEmptyTrimmed or NotEmpty to prevent users from accidentally passing empty or non-trimmed flag values.
func AnyString[T ~string](s string) (T, error) {
	return T(s), nil
}

// SliceParser is a Parser Slice, see NewSlice and ParseSliceOf.
type SliceParser[T any] func(ss []string) ([]T, error)

// SliceTargetParser can be implemented by given target pointers
// Implementations of SliceTargetParser may use ErrParser when an error occurs.
type SliceTargetParser interface {
	Parse(ss []string) error
}

// ParseSliceOf turns a Parser into a SliceParser, calling the given single element parser for each slice element.
// It propagates parsing failures of the single element parser and telling at which element the error happened.
func ParseSliceOf[T any](p Parser[T]) SliceParser[T] {
	return func(ss []string) ([]T, error) {
		var result []T
		for i, s := range ss {
			parsed, err := p(s)
			if err != nil {
				return nil, fmt.Errorf("cannot parse slice element %d: %w", i, err)
			}
			result = append(result, parsed)
		}
		return result, nil
	}
}
