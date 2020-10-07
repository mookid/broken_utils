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

PROGS = $(addprefix $(GO_DIR),$(addsuffix .exe, $(NAMES)))
SOURCES = $(addprefix $(GO_DIR),$(addsuffix .go, $(NAMES)))

go.all: $(PROGS)

.go.exe:
	go build -o $@ $<

go.install:
	cp $(PROGS) $(GO_BIN)

go.clean:
	rm -f $(PROGS)

all: go.all
install: go.install
clean: go.clean
