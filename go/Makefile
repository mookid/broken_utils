GO_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
LDFLAGS := "-s -w"

go.all:
	cd $(GO_DIR) && go build -ldflags=$(LDFLAGS) ./...

go.install:
	cd $(GO_DIR) && go install -ldflags=$(LDFLAGS) ./...

go.clean:
	cd $(GO_DIR) && go clean ./...

go.tests:
	cd $(GO_DIR) && go test ./... -v

all: go.all
install: go.install
clean: go.clean
tests: go.tests
