package main

import (
	"log"

	"github.com/neiser/go-nagini/command"
	"github.com/neiser/go-nagini/flag"
)

type (
	Wand string
)

func main() {
	var (
		myName       string
		favoriteWand Wand = "elder"
		iAmVoldemort bool
	)
	_ = command.New().
		Flag(flag.String(&myName, flag.NotEmpty), flag.RegisterOptions{
			Name:     "my-name",
			Required: true,
		}).
		Flag(flag.String(&favoriteWand, flag.NotEmptyTrimmed), flag.RegisterOptions{
			Name:  "favorite-wand",
			Usage: "Specify magic wand",
		}).
		Flag(flag.Bool(&iAmVoldemort), flag.RegisterOptions{
			Name: "i-am-voldemort",
		}).
		Run(func() error {
			if iAmVoldemort {
				return command.WithExitCodeError{ExitCode: 66}
			}
			log.Printf("I'm %s and my favorite wand is '%s'", myName, favoriteWand)
			return nil
		}).
		Execute()
}
