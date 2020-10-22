CARGO_BIN = $(HOME)/.cargo/bin
RUST_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
MANIFEST = --manifest-path $(addprefix $(RUST_DIR),Cargo.toml)

rust.all:
	cargo build --release $(MANIFEST)

rust.install:
	cp $(addprefix $(RUST_DIR),target/release/*.exe) $(CARGO_BIN)

rust.clean:
	cargo clean $(MANIFEST)

all: rust.all
install: rust.install
clean: rust.clean
