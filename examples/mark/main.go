package main

import (
	"log"

	"github.com/neiser/go-nagini/command"
	"github.com/neiser/go-nagini/flag"
)

func main() {
	var (
		name         = "Harry"
		iAmVoldemort bool
	)
	command.New().
		Flag(flag.String(&name, flag.NotEmpty), flag.RegisterOptions{
			Name: "name",
		}).
		Flag(flag.String(&name, flag.NotEmpty), flag.RegisterOptions{
			Name: "nickname",
		}).
		Flag(flag.Bool(&iAmVoldemort), flag.RegisterOptions{
			Name: "i-am-voldemort",
		}).
		MarkFlagsMutuallyExclusive(&name, &iAmVoldemort).
		Run(func() error {
			switch {
			case iAmVoldemort:
				log.Print("My name is Voldemort!")
			case name != "":
				log.Printf("My name is %s", name)
			}
			return nil
		}).
		Execute()
}
