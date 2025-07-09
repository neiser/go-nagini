package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/neiser/go-nagini/command"
	"github.com/neiser/go-nagini/flag"
)

type VerboseLevel int

func (v *VerboseLevel) String() string {
	return fmt.Sprintf("%d", *v)
}

func (v *VerboseLevel) Set(s string) error {
	if s == "true" {
		*v++
		return nil
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("cannot convert: %w", err)
	}
	*v = VerboseLevel(val)
	return nil
}

func (v *VerboseLevel) Type() string {
	return "int"
}

//nolint:ireturn
func (v *VerboseLevel) Target() any {
	return v
}

func (v *VerboseLevel) IsBoolFlag() bool {
	return true
}

func main() {
	var (
		verboseLevel VerboseLevel
		enableDebug  bool
	)
	command.New().
		Flag(&verboseLevel, flag.RegisterOptions{Name: "verbose", Shorthand: "v"}).
		Flag(flag.Bool(&enableDebug), flag.RegisterOptions{Name: "debug"}).
		MarkFlagsMutuallyExclusive(&enableDebug, &verboseLevel).
		Run(func() error {
			log.Printf("Verbose level is %d\n", verboseLevel)
			return nil
		}).
		Execute()
}
