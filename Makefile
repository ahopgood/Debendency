.PHONY: build

build:
	go build -o build/debendency main.go

test:
	ginkgo -r --skip-package integrationtests -cover

# Integration tests
int:
	ginkgo -r --skip-package pkg

fmt:

lint:

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