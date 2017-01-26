#!/usr/bin/env bash

go build -o $GOPATH/bin/stack

go test github.com/gophercloud/cli/acceptance_tests/...
