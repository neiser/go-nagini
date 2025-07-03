package parameter

import (
	"errors"
	"fmt"
	"strings"
)

// Parser is used in New to convert from string to the generic type T.
// Use ErrParser when an error occurs.
type Parser[T any] func(s string) (T, error)

// ErrParser can be used when implementations of Parser fail.
var ErrParser = errors.New("cannot parse parameter")

// NotEmptyString is a commonly used Parser and requires that the given string is not only whitespace.
func NotEmptyString[T ~string](s string) (T, error) {
	if strings.TrimSpace(s) != "" {
		return T(s), nil
	}
	return "", fmt.Errorf("%w: value '%s' is not set", ErrParser, s)
}

// SliceParser is a Parser Slice, see NewSlice and ParseSliceOf.
type SliceParser[T any] func(ss []string) ([]T, error)

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
