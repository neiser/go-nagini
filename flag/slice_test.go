package flag

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_anySliceValue_Type(t *testing.T) {
	t.Run("built-in type", func(t *testing.T) {
		var (
			target []string
		)
		sut := Slice(&target, ParseSliceOf[string](NotEmptyTrimmed[string]))
		assert.Equal(t, "[]string", sut.Type())
	})
	t.Run("custom type", func(t *testing.T) {
		type someType string
		type someTypes []someType
		var (
			target someTypes
		)
		sut := Slice[someType](&target, ParseSliceOf[someType](NotEmptyTrimmed[someType]))
		assert.Equal(t, "[]someType", sut.Type())
	})
}

type someSliceType []string

func (s *someSliceType) ParseAndReplace(ss []string) error {
	*s = ss
	return nil
}

func (s *someSliceType) ParseAndAppend(ss ...string) error {
	*s = append(*s, ss...)
	return nil
}

func Test_anySliceValue_Append(t *testing.T) {
	t.Run("using ParseSliceOf", func(t *testing.T) {
		var (
			target = []string{"item0"}
		)
		sut := Slice(&target, ParseSliceOf[string](AnyString))
		require.NoError(t, sut.Append("item1"))
		assert.Equal(t, []string{"item0", "item1"}, sut.GetSlice())
	})
	t.Run("using ParseSliceOf but fails", func(t *testing.T) {
		var (
			target = []string{"item0"}
		)
		sut := Slice(&target, ParseSliceOf[string](func(string) (string, error) {
			return "", errors.New("some error")
		}))
		require.ErrorContains(t, sut.Append("item1"), "some error")
	})
	t.Run("using SliceTargetParser", func(t *testing.T) {
		var (
			target = someSliceType{"item0"}
		)
		sut := Slice(&target, nil)
		require.NoError(t, sut.Append("item1"))
		assert.Equal(t, []string{"item0", "item1"}, sut.GetSlice())
	})
}

func Test_anySliceValue_Replace(t *testing.T) {
	t.Run("using ParseSliceOf", func(t *testing.T) {
		var (
			target = []string{"item0"}
		)
		sut := Slice(&target, ParseSliceOf[string](AnyString))
		require.NoError(t, sut.Replace([]string{"item0", "item1"}))
		assert.Equal(t, []string{"item0", "item1"}, sut.GetSlice())
	})
	t.Run("using ParseSliceOf but fails", func(t *testing.T) {
		var (
			target = []string{"item0"}
		)
		sut := Slice(&target, ParseSliceOf[string](func(string) (string, error) {
			return "", errors.New("some error")
		}))
		require.ErrorContains(t, sut.Replace([]string{"item0", "item1"}), "some error")
	})
	t.Run("using SliceTargetParser", func(t *testing.T) {
		var (
			target = someSliceType{"item0"}
		)
		sut := Slice(&target, nil)
		require.NoError(t, sut.Replace([]string{"item0", "item1"}))
		assert.Equal(t, []string{"item0", "item1"}, sut.GetSlice())
	})
}

func Test_readAsCSV(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr assert.ErrorAssertionFunc
	}{
		{"empty string", "", []string{}, assert.NoError},
		{"one value", "some value with spaces", []string{"some value with spaces"}, assert.NoError},
		{"two values", "val2,val1", []string{"val2", "val1"}, assert.NoError},
		{"parsing fails", `",`, nil, func(t assert.TestingT, err error, args ...any) bool {
			return assert.ErrorContains(t, err, `cannot read value '",' as comma-separated values`, args...)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readAsCSV(tt.input)
			if !tt.wantErr(t, err, fmt.Sprintf("readAsCSV(%v)", tt.input)) {
				return
			}
			assert.Equalf(t, tt.want, got, "readAsCSV(%v)", tt.input)
		})
	}
}
