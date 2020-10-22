GO_BIN = $(GOPATH)/bin
GO_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

.SUFFIXES:  .go .exe

NAMES=                                                            \
	git-ff                                                    \
	l                                                         \
	namespaces                                                \
	rm_orig_files                                             \
	re_kill                                                   \
	s                                                         \
	tasks                                                     \
	token                                                     \
	vs                                                        \

TESTED_NAMES=                                                     \
	re_kill                                                   \

PROGS = $(addprefix $(GO_DIR),$(addsuffix .exe, $(NAMES)))
SOURCES = $(addprefix $(GO_DIR),$(addsuffix .go, $(NAMES)))

TESTS = $(addprefix $(GO_DIR),$(addsuffix .test, $(TESTED_NAMES)))

go.all: $(PROGS)

go.tests: $(TESTS)

.go.exe:
	go build -o $@ $<

go.install:
	cp $(PROGS) $(GO_BIN)

go.clean:
	rm -f $(PROGS)

all: go.all
install: go.install
clean: go.clean
tests: go.tests

%.test: %_test.go %.go
	go test -v $^
