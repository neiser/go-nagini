# Nagini - Fluent and generic API for Cobra

![coverage](https://raw.githubusercontent.com/neiser/go-nagini/badges/.badges/main/coverage.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/neiser/go-nagini)](https://goreportcard.com/report/github.com/neiser/go-nagini)
[![Go Reference](https://pkg.go.dev/badge/github.com/neiser/go-nagini.svg)](https://pkg.go.dev/github.com/neiser/go-nagini)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)

Nagini wraps the famous 
[Cobra CLI library](https://github.com/spf13/cobra) 
with a fluent API and Go generics.

It supports slice values (comma-separated values) and arbitrary parsing.

It optionally binds flags to the
[Viper](https://github.com/spf13/viper) configuration library.

## Installation

```shell
go get github.com/neiser/go-nagini
```

## Example Usage

### Simple command with generic, string-like flag and boolean flag

Showing [`examples/simple/main.go`](./examples/simple/main.go):

```go:examples/simple/main.go
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
	command.New().
		Flag(flag.New(&myName, flag.NotEmpty), flag.RegisterOptions{
			Name:     "my-name",
			Required: true,
		}).
		Flag(flag.New(&favoriteWand, flag.NotEmptyTrimmed), flag.RegisterOptions{
			Name:  "favorite-wand",
			Usage: "Specify magic wand",
		}).
		FlagBool(&iAmVoldemort, flag.RegisterOptions{
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

```

Run with
```shell
go run ./examples/simple --my-name Harry
```

### Command with slice flag

Showing [`examples/slice/main.go`](./examples/slice/main.go):

```go:examples/slice/main.go
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

```

Run with
```shell
go run ./examples/slice --some-ints 5,6,7
```

### Adding subcommands with fluent description

Showing [`examples/subcommand/main.go`](./examples/subcommand/main.go):

```go:examples/subcommand/main.go
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
	command.New().
		FlagBool(&useMagic, flag.RegisterOptions{
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
						return ErrCannotUseMagic
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
		Execute()
}

```

Run with
```shell
go run ./examples/subcommand wizard --use-magic
go run ./examples/subcommand muggle --use-magic
```

### Binding a flag to Viper, flag value takes precedence over Viper

Showing [`examples/viper/main.go`](examples/viper/main.go):

```go:examples/viper/main.go
package main

import (
	"log"

	"github.com/neiser/go-nagini/command"
	"github.com/neiser/go-nagini/flag"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv() // tell Viper to read env
	var (
		gitlabToken string
	)
	command.New().
		Flag(
			flag.ViperBinding{
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

```

Run with
```shell
GITLAB_TOKEN=bla go run ./examples/viper
```
or 
```shell
GITLAB_TOKEN=bla go run ./examples/viper --gitlab-token blub
```


## Development and Contributions

Install the provided 
[pre-commit](https://pre-commit.com)
hooks with
```shell
pre-commit install
```

This library only exposes some limited feature set of Cobra.
Please open an issue if you really miss something which fits well into this library.
Otherwise, you can always modify the embedded `*cobra.Command` directly as a workaround.
