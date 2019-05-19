#! /bin/bash

cd ~/proj/cse-scripts-lambda
npm install --save-dev

go get -v -t -d ./..
make build
go test -v ./..
make deploy
