package main

import (
	"log"
	"strconv"

	"github.com/neiser/go-nagini/command"
	"github.com/neiser/go-nagini/flag"
)

func main() {
	var (
		someInts []int
	)
	command.New().
		Flag(flag.NewSlice(&someInts, flag.ParseSliceOf(strconv.Atoi)), flag.RegisterOptions{
			Name:     "some-ints",
			Required: true,
		}).
		Run(func() error {
			log.Printf("Got integers: '%v'", someInts)
			return nil
		}).
		Execute()
}
