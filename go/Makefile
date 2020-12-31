GO_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

go.all:
	cd $(GO_DIR) && go build ./...

go.install:
	cd $(GO_DIR) && go install ./...

go.clean:
	cd $(GO_DIR) && go clean ./...

go.tests:
	cd $(GO_DIR) && go test ./... -v

all: go.all
install: go.install
clean: go.clean
tests: go.tests
