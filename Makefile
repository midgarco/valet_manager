SHELL := /bin/bash

# The name of the executable (default is current directory name)
TARGET := $(shell echo $${PWD\#\#*/})
.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
VERSION := 1.0.0
BUILD := `git rev-parse --short HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build clean fmt lint test test-all

all: clean build run

# pb:
# 	@protoc -I/usr/local/include -I. \
# 		-I${GOPATH}/src \
# 		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
# 		--go_out=plugins=grpc:. \
# 		proto/app.proto
# 	@protoc -I/usr/local/include -I. \
# 		-I${GOPATH}/src \
# 		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
# 		--grpc-gateway_out=logtostderr=true:. \
# 		proto/app.proto
# 	@protoc -I/usr/local/include -I. \
# 		-I${GOPATH}/src \
# 		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
# 		--swagger_out=logtostderr=true:. \
# 		proto/app.proto	

$(TARGET): $(SRC)
	go build $(LDFLAGS) -o $(TARGET) -v ./cmd/api

build: test $(TARGET)
	@true

clean:
	@rm -f $(TARGET)
	@rm -f cmd/api/api
	@rm -f cmd/cli/cli
	# @rm -f proto/app.pb.go
	# @rm -f proto/app.pb.gw.go
	# @rm -f proto/app.swagger.json

fmt:
	# gofmt -l -w $(SRC)

test:
	# go test -short ./...

lint:
	# go vet ./...

test-all: lint test
	# go test -race ./...

run:
	@./$(TARGET)