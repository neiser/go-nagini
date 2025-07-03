package parameter

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
// but it must also implement pflag.Value to make it registrable as a command parameter.
type Slice interface {
	pflag.SliceValue
	pflag.Value
}

// NewSlice construct a new Slice parameter.
func NewSlice[T any, E ~[]T](target *E, defValue E, parser SliceParser[T]) Slice {
	*target = defValue
	return anySliceValue[T, E]{target, parser}
}

type anySliceValue[T any, E ~[]T] struct {
	target *E
	parser SliceParser[T]
}

func (v anySliceValue[T, E]) String() string {
	return v.asCsv()
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

func (v anySliceValue[T, E]) Append(s string) error {
	added, err := v.parser([]string{s})
	if err != nil {
		return err
	}
	*v.target = append(*v.target, added...)
	return nil
}

func (v anySliceValue[T, E]) Replace(ss []string) (err error) {
	*v.target, err = v.parser(ss)
	return
}

func (v anySliceValue[T, E]) GetSlice() (result []string) {
	for _, item := range *v.target {
		result = append(result, convertToString(item))
	}
	return
}
