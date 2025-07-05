package main

import (
	"log"

	"github.com/neiser/go-nagini/command"
	"github.com/neiser/go-nagini/flag"
	"github.com/neiser/go-nagini/flag/binding"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv() // tell Viper to read env
	var (
		gitlabToken string
	)
	command.New().
		Flag(
			binding.Viper{
				Value:     flag.New(&gitlabToken, flag.NotEmptyTrimmed),
				ConfigKey: "GITLAB_TOKEN",
			},
			flag.RegisterOptions{
				Name:  "gitlab-token",
				Usage: "A secret GitLab Token",
			},
		).
		Run(func() error {
			log.Printf("Gitlab Token length '%d'", len(gitlabToken))
			return nil
		}).
		Execute()
}
