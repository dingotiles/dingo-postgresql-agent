#!/bin/bash

set -e -x

export GOPATH=$PWD/gopath

cd $GOPATH/src/github.com/dingotiles/dingo-postgresql-agent

go test $(go list ./... | grep -v /vendor/)
