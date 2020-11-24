GO_BIN = $(GOPATH)/bin
GO_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

.SUFFIXES:  .go .exe

NAMES=                                                            \
	git-codereview                                            \
	git-ff                                                    \
	git-where-is                                              \
	git-zip                                                   \
	l                                                         \
	namespaces                                                \
	re-kill                                                   \
	rgw                                                       \
	rm-orig-files                                             \
	s                                                         \
	sha256sum                                                 \
	tasks                                                     \
	token                                                     \
	vs                                                        \

TESTED_NAMES=                                                     \
	re-kill                                                   \
	rgw                                                       \
	tasks                                                     \

PROGS = $(addprefix $(GO_DIR),$(addsuffix .exe, $(NAMES)))
TESTS = $(addprefix $(GO_DIR),$(addsuffix .test, $(TESTED_NAMES)))
TEST_PROGS = $(addprefix $(GO_DIR),$(addsuffix .testexe, $(TESTED_NAMES)))

go.all: $(PROGS)
go.all: $(TEST_PROGS)

go.tests: $(TESTS)

%.exe: %.go
	go build -o $@ $<

%.test: %_test.go %.go
	go test -v $^

%.testexe: %_test.go %.go
	go test -c -o $@ $^

go.install:
	cp $(PROGS) $(GO_BIN)

go.clean:
	rm -f $(PROGS) $(TEST_PROGS)

all: go.all
install: go.install
clean: go.clean
tests: go.tests
