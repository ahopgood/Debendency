
build: go build main.go -o build/main

test: ginkgo --skip-package integrationtests -cover

# Integration tests
int: ginkgo --skip-package pkg

fmt:

lint:

coverage-html: coverage-clean test

coverage-cli: coverage-clean test

coverage-clean:

help: ## Show this help message.
	echo "usage: make [target] ..."
	echo ""
	echo "targets:"
	egrep '^(.+)\:(\ |.*) ##\ (.+)' ${MAKEFILE_LIST} | column -t -c 2 -s ':#'