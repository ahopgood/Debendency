.PHONY: build generate

build:
	go build -o build/debendency main.go

test:
	ginkgo -r -cover --skip-package integrationtests

# Integration tests
int:
	ginkgo -v -r --skip-package pkg

fmt:
	gofmt -s -w .
	golangci-lint run --fast --fix

lint:
	golangci-lint run

generate:
	go generate ./...
	#go generate pkg
# needs to be in pkg

coverage-html: coverage-clean test
	go tool cover -html=coverprofile.out -o coverage.html

coverage-cli: coverage-clean test
	go tool cover -func=coverprofile.out

coverage-clean:
	- rm 'coverprofile.out'

## Show this help message.
help:
	echo "usage: make [target] ..."
	echo ""
	echo "targets:"
	egrep '^(.+)\:(\ |.*) ##\ (.+)' ${MAKEFILE_LIST} | column -t -c 2 -s ':#'