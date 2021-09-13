.DEFAULT_GOAL := test-all

test-all: lint test

.PHONY: test
test:
	go test -race ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: verify-codegen
verify-codegen:
	./scripts/verify-codegen.sh

.PHONY: update-codegen
update-codegen:
	buf generate

