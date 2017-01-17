#!/bin/bash

set -e -x

cd $GOPATH
tree .

cd $GOPATH/src/github.com/dingotiles/dingo-postgresql-agent

go test $(go list ./... | grep -v /vendor/)
