package flag

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/pflag"
)

// Slice is constructed from NewSlice.
// This interface explicitly declares that this is a pflag.SliceValue,
// but it must also implement Value to make it registrable as a command flag.
// See also String and Bool.
type Slice interface {
	Value
	pflag.SliceValue
}

// NewSlice construct a new Slice flag having multiple values, parsed by given SliceParser.
// See also ParseSliceOf.
// If SliceParser is nil, falls back to an implementation of SliceTargetParser on type *E,
// panics if nothing is found.
func NewSlice[T any, E ~[]T](target *E, parser SliceParser[T]) Slice {
	result := anySliceValue[T, E]{target: target, parser: parser}
	if parser == nil {
		var ok bool
		if result.targetParser, ok = any(target).(SliceTargetParser); !ok {
			panic(fmt.Sprintf("flag for slice target %p (value '%v') must specify non-nil parser, "+
				"as flag.SliceTargetParser interface is also not implemented", target, *target))
		}
	}
	return result
}

type anySliceValue[T any, E ~[]T] struct {
	target       *E
	parser       SliceParser[T]
	targetParser SliceTargetParser
}

//nolint:ireturn
func (v anySliceValue[T, E]) Target() any {
	return v.target
}

func (v anySliceValue[T, E]) IsBoolFlag() bool {
	return false
}

func (v anySliceValue[T, E]) String() string {
	switch {
	case *v.target == nil:

		return "<nil>"
	case len(*v.target) == 0:
		return "<empty>"
	default:
		return v.asCsv()
	}
}

func readAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	result, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("cannot read value '%s' as comma-separated values: %w", val, err)
	}
	return result, nil
}

func (v anySliceValue[T, E]) asCsv() string {
	var buffer bytes.Buffer
	csvWriter := csv.NewWriter(&buffer)
	_ = csvWriter.Write(v.GetSlice())
	csvWriter.Flush()
	return strings.TrimSpace(buffer.String())
}

func (v anySliceValue[T, E]) Set(s string) error {
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")
	values, err := readAsCSV(s)
	if err != nil {
		return fmt.Errorf("cannot parse '%s' as comma-separated values: %w", s, err)
	}
	return v.Replace(values)
}

func (v anySliceValue[T, E]) Type() string {
	var t T
	return "[]" + reflect.TypeOf(t).Name()
}

//nolint:wrapcheck
func (v anySliceValue[T, E]) Append(s string) (err error) {
	var added []T
	if v.targetParser != nil {
		return v.targetParser.ParseAndAppend(s)
	}
	added, err = v.parser([]string{s})
	if err != nil {
		return err
	}
	*v.target = append(*v.target, added...)
	return nil
}

//nolint:wrapcheck
func (v anySliceValue[T, E]) Replace(ss []string) (err error) {
	if v.targetParser != nil {
		return v.targetParser.ParseAndReplace(ss)
	}
	*v.target, err = v.parser(ss)
	return
}

func (v anySliceValue[T, E]) GetSlice() (result []string) {
	for _, item := range *v.target {
		result = append(result, convertToString(item))
	}
	return
}
