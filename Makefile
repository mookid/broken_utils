all:
	$(MAKE) -C rust $@
	$(MAKE) -C go $@

install:
	$(MAKE) -C rust $@
	$(MAKE) -C go $@

clean:
	$(MAKE) -C rust $@
	$(MAKE) -C go $@
