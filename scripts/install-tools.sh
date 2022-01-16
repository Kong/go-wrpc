#!/bin/bash

go install ./cmd/protoc-gen-go-wrpc
go install google.golang.org/protobuf/cmd/protoc-gen-go

curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.43.0
