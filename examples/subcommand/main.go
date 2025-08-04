package main

import (
	"errors"
	"log"

	"github.com/neiser/go-nagini/command"
	"github.com/neiser/go-nagini/flag"
)

var ErrCannotUseMagic = errors.New("cannot use magic")

func main() {
	var (
		useMagic bool
	)
	_ = command.New().
		Flag(flag.Bool(&useMagic), flag.RegisterOptions{
			Name:       "use-magic",
			Usage:      "Use some magic, c'mon",
			Persistent: true,
		}).
		AddCommands(
			command.New().
				Use("muggle").
				Short("A person which cannot use magic").
				Run(func() error {
					if useMagic {
						return command.WithExitCodeError{
							ExitCode: 21,
							Wrapped:  ErrCannotUseMagic,
						}
					}
					return nil
				}),
			command.New().
				Use("wizard").
				Short("A person which may use magic").
				Run(func() error {
					if useMagic {
						log.Printf("Abracadabra!")
					}
					return nil
				}),
		).
		AddPersistentPreRun(func() error {
			log.Printf("Will always run!")
			return nil
		}).
		AddPersistentPreRun(func() error {
			log.Printf("Will also run!")
			return nil
		}).
		Execute()
}
