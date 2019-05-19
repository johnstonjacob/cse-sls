.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux CGO_ENABLED=0 go build -v -ldflags="-s -w" -o bin/workflow-cost-estimator src/workflow-cost-estimator/main.go src/workflow-cost-estimator/types.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

ci-deploy: clean build
	sls deploy --verbose --stage ${SLS_STAGE}

dev: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
