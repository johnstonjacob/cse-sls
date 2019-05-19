#! /bin/bash

cd ~/proj/workflow-cost-estimator

go get -v -t -d ./..
make build
go test -v ./..
make deploy
