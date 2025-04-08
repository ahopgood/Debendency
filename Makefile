.PHONY: build generate

build:
	go build -o build/debendency main.go

static-build:
	CGO_ENABLED=0 GOOS=$(goos) go build -ldflags "-X main.version=$(VERSION)" -o build/static-main main.go

test:
	ginkgo -r -cover --skip-package integrationtests -coverprofile=unit.coverprofile

# Integration tests
int:
	ginkgo -v -r -cover --skip-package pkg -coverprofile=int.coverprofile -coverpkg ./pkg/commands
#	ginkgo -v -r -cover --skip-package pkg -coverprofile=coverprofile.int -coverpkg ./pkg/commands

fmt:
	gofmt -s -w .
	golangci-lint run --fast --fix

lint:
	golangci-lint run

generate:
	go generate ./...
	#go generate pkg
# needs to be in pkg

coverage-html:  coverage-clean test int
	gover . coverprofile.aggregate
	go tool cover -html=coverprofile.aggregate -o coverage.html

coverage-cli: coverage-clean test int
	go tool cover -func=coverprofile.out

coverage-clean:
	- rm 'coverprofile.aggregate' 'coverprofile.unit' 'coverprofile.int'

## Show this help message.
help:
	echo "usage: make [target] ..."
	echo ""
	echo "targets:"
	egrep '^(.+)\:(\ |.*) ##\ (.+)' ${MAKEFILE_LIST} | column -t -c 2 -s ':#'