#! /bin/bash

go get -v -t -d ./..
make build
go test -v ./..
make deploy
