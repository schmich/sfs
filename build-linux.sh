#!/bin/sh

set -eufx

version=`git tag | tail -n1`
commit=`git rev-parse HEAD`

go get -v
go build -ldflags "-w -s -X main.version=$version -X main.commit=$commit" -o sfs
