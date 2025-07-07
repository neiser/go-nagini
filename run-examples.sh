#!/usr/bin/env bash
set -xeuo pipefail
go run ./examples/simple --my-name Harry
go run ./examples/slice --some-ints 5,6,7
go run ./examples/subcommand wizard --use-magic
IS_EVIL=true FAVORITE_HOUSE=Slytherin go run ./examples/viper
go run ./examples/mark --i-am-voldemort
go run ./examples/verbose -vvv
