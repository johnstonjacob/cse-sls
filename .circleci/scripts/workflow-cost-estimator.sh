#! /bin/bash

cd ~/proj/workflow-cost-estimator
npm install --save-dev

go get -v -t -d ./..
make build
go test -v ./..
make deploy
