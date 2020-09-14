GO_BIN=$(GOPATH)/bin

.SUFFIXES:  .go .exe

NAMES=                                                            \
	git-ff                                                    \
	namespaces                                                \
	rm_orig_files                                             \
	s                                                         \
	tasks                                                     \
	vs                                                        \

PROGS=$(addsuffix .exe, $(NAMES))
SOURCES=$(addsuffix .go, $(NAMES))

all: $(PROGS)

.go.exe:
	go build -o $@ $<

install:
	cp $(PROGS) $(GO_BIN)

clean:
	rm -f $(PROGS)
