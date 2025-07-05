# Nagini - Fluent API for Cobra

![coverage](https://raw.githubusercontent.com/neiser/go-nagini/badges/.badges/main/coverage.svg)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)

Nagini wraps the famous 
[Cobra CLI library](https://github.com/spf13/cobra) 
with a fluent API.

## Installation

```shell
go get github.com/neiser/go-nagini
```

## Example Usage

### Simple command with generic, string-like flag and boolean flag

Showing [`examples/simple/main.go`](./examples/simple/main.go):

```go:examples/simple/main.go

```

Run with
```shell
go run github.com/neiser/go-nagini/examples/simple --my-name Harry
```

### Binding a flag to Viper, flag value takes precedence over Viper

Showing [`examples/viper/main.go`](examples/viper/main.go):

```go:examples/viper/main.go

```

Run with
```shell
GITLAB_TOKEN=bla go run github.com/neiser/go-nagini/examples/viper
```
or 
```shell
GITLAB_TOKEN=bla go run github.com/neiser/go-nagini/examples/viper --gitlab-token blub
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
