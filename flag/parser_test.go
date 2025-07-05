package flag

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotEmptyTrimmed(t *testing.T) {
	type testCase struct {
		name    string
		input   string
		want    string
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase{
		{"empty", "", "", func(t assert.TestingT, err error, msgAndArgs ...any) bool {
			return assert.ErrorContains(t, err, "cannot parse parameter: value '' after trimming is empty", msgAndArgs...)
		}},
		{"whitespace only", "\t\n   ", "", func(t assert.TestingT, err error, msgAndArgs ...any) bool {
			return assert.ErrorContains(t, err, "cannot parse parameter: value '\t\n   ' after trimming is empty", msgAndArgs...)
		}},
		{"some value", "some-value", "some-value", assert.NoError},
		{"some value with whitespace around", " some  value\t", "some  value", assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NotEmptyTrimmed[string](tt.input)
			if !tt.wantErr(t, err, fmt.Sprintf("NotEmptyTrimmed(%v)", tt.input)) {
				return
			}
			assert.Equalf(t, tt.want, got, "NotEmptyTrimmed(%v)", tt.input)
		})
	}
}

func TestNotEmpty(t *testing.T) {
	type testCase struct {
		name    string
		input   string
		want    string
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase{
		{"empty", "", "", func(t assert.TestingT, err error, msgAndArgs ...any) bool {
			return assert.ErrorContains(t, err, "cannot parse parameter: value '' is empty", msgAndArgs...)
		}},
		{"whitespace only", "\t\n   ", "\t\n   ", assert.NoError},
		{"some value", "some-value", "some-value", assert.NoError},
		{"some value with whitespace around", " some  value\t", " some  value\t", assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NotEmpty[string](tt.input)
			if !tt.wantErr(t, err, fmt.Sprintf("NotEmpty(%v)", tt.input)) {
				return
			}
			assert.Equalf(t, tt.want, got, "NotEmpty(%v)", tt.input)
		})
	}
}

func TestParseSliceOf(t *testing.T) {
	type testCase struct {
		name    string
		input   []string
		want    []int
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase{
		{"empty", nil, nil, assert.NoError},
		{"some ints", []string{"6", "8"}, []int{6, 8}, assert.NoError},
		{"not parsable", []string{"6", "x8x"}, nil, func(t assert.TestingT, err error, msgAndArgs ...any) bool {
			return assert.ErrorContains(t, err, `cannot parse slice element 1: strconv.Atoi: parsing "x8x": invalid syntax`, msgAndArgs...)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSliceOf(strconv.Atoi)(tt.input)
			if !tt.wantErr(t, err, fmt.Sprintf("NotEmpty(%v)", tt.input)) {
				return
			}
			assert.Equalf(t, tt.want, got, "NotEmpty(%v)", tt.input)
		})
	}
}
