#!/bin/bash

cd $GOPATH/src/github.com/dingotiles/dingo-postgresql-agent

go test $(go list ./... | grep -v /vendor/)
