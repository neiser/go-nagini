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
    favoriteHouse = "Hufflepuff"
    isEvil        = false
  )
  command.New().
    Flag(
      binding.Viper{
        Value:     flag.String(&favoriteHouse, flag.NotEmptyTrimmed),
        ConfigKey: "FAVORITE_HOUSE",
      },
      flag.RegisterOptions{
        Name: "house",
      },
    ).
    Flag(
      binding.Viper{
        Value:     flag.Bool(&isEvil),
        ConfigKey: "IS_EVIL",
      },
      flag.RegisterOptions{
        Shorthand: "e",
      },
    ).
    Run(func() error {
      prefix := "Favorite"
      if isEvil {
        prefix = "Evil favorite"
      }
      log.Printf("%s house is %s", prefix, favoriteHouse)
      return nil
    }).
    Execute()
}
